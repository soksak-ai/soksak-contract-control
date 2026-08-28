use sha2::{Digest, Sha256};
use serde::{Deserialize, Serialize};
use std::path::Path;

pub const PROTOCOL: u32 = 2;
pub const PROCESS_LABEL_ENVIRONMENT: &str = "SOKSAK_PROCESS_LABEL";
pub const DEFAULT_PROCESS_LABEL: &str = "soksak";

pub fn parse_process_label(value: &str) -> Option<&str> {
    let bytes = value.as_bytes();
    if bytes.is_empty() || bytes.len() > 31 || !bytes[0].is_ascii_alphanumeric() {
        return None;
    }
    bytes[1..]
        .iter()
        .all(|byte| byte.is_ascii_alphanumeric() || matches!(byte, b'.' | b'_' | b'-'))
        .then_some(value)
}

pub fn format_process_name(label: &str, role: &str) -> Option<String> {
    let label = parse_process_label(label)?;
    let bytes = role.as_bytes();
    if bytes.is_empty()
        || bytes.len() > 128
        || !bytes[0].is_ascii_lowercase() && !bytes[0].is_ascii_digit()
        || !bytes[1..]
            .iter()
            .all(|byte| byte.is_ascii_lowercase() || byte.is_ascii_digit() || *byte == b'-')
    {
        return None;
    }
    Some(format!("{label}-{role}"))
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Announcement {
    pub protocol: u32,
    pub socket: String,
    pub process_label: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub token: Option<String>,
}

pub fn announcement(socket: &str, process_label: &str, token: Option<&str>) -> Option<Announcement> {
    Some(Announcement {
        protocol: PROTOCOL,
        socket: socket.to_owned(),
        process_label: parse_process_label(process_label)?.to_owned(),
        token: token.map(str::to_owned),
    })
}

pub fn address(runtime_root: &Path, identifier: &str, windows: bool) -> String {
    if !windows {
        return runtime_root
            .join(format!("{identifier}.sock"))
            .to_string_lossy()
            .into_owned();
    }
    let normalized = normalize_windows(runtime_root);
    let mut hash = Sha256::new();
    hash.update(normalized.as_bytes());
    hash.update([0]);
    hash.update(identifier.as_bytes());
    let digest = hash.finalize();
    let suffix: String = digest[..16]
        .iter()
        .map(|byte| format!("{byte:02x}"))
        .collect();
    format!(r"\\.\pipe\soksak-control-{suffix}")
}

fn normalize_windows(path: &Path) -> String {
    let value = path.to_string_lossy().replace('/', "\\");
    let mut parts = Vec::new();
    for part in value.split('\\') {
        match part {
            "" | "." => {}
            ".." => {
                parts.pop();
            }
            _ => parts.push(part),
        }
    }
    parts.join("\\").to_lowercase()
}

#[cfg(test)]
mod tests {
    use super::*;
    use serde::Deserialize;

    #[derive(Deserialize)]
    struct Vector {
        runtime: String,
        identifier: String,
        windows: bool,
        address: String,
    }

    #[derive(Deserialize)]
    struct ProcessLabelVector {
        input: String,
        valid: bool,
    }

    #[derive(Deserialize)]
    struct ProcessNameVector {
        label: String,
        role: String,
        name: String,
        valid: bool,
    }

    #[test]
    fn process_label_vectors_are_the_rust_contract() {
        let vectors: Vec<ProcessLabelVector> =
            serde_json::from_str(include_str!("../process-label-vectors.json")).unwrap();
        for vector in vectors {
            let parsed = parse_process_label(&vector.input);
            assert_eq!(parsed, vector.valid.then_some(vector.input.as_str()), "{}", vector.input);
        }
    }

    #[test]
    fn process_name_vectors_are_the_rust_contract() {
        let vectors: Vec<ProcessNameVector> =
            serde_json::from_str(include_str!("../process-name-vectors.json")).unwrap();
        for vector in vectors {
            let formatted = format_process_name(&vector.label, &vector.role);
            assert_eq!(formatted.as_deref(), vector.valid.then_some(vector.name.as_str()));
        }
    }

    #[test]
    fn announcement_requires_and_carries_the_process_label() {
        let value = announcement("/runtime/sidecar.sock", "soksakv3", Some("token")).unwrap();
        assert_eq!(value.protocol, 2);
        assert_eq!(value.process_label, "soksakv3");
        let json = serde_json::to_value(value).unwrap();
        assert_eq!(json["processLabel"], "soksakv3");
    }

    #[test]
    fn matches_shared_address_vectors() {
        let vectors: Vec<Vector> =
            serde_json::from_str(include_str!("../address-vectors.json")).unwrap();
        for vector in vectors {
            assert_eq!(
                address(
                    Path::new(&vector.runtime),
                    &vector.identifier,
                    vector.windows
                ),
                vector.address
            );
        }
    }

    #[test]
    fn windows_normalization_is_case_and_separator_stable() {
        assert_eq!(
            address(Path::new(r"C:\RUNTIME/one"), "com.soksak.test", true),
            address(Path::new(r"c:\runtime\one"), "com.soksak.test", true)
        );
    }
}
