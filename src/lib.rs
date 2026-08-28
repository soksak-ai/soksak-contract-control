use sha2::{Digest, Sha256};
use std::path::Path;

pub const PROCESS_LABEL_ENVIRONMENT: &str = "SOKSAK_PROCESS_LABEL";

pub fn parse_process_label(_value: &str) -> Option<&str> {
    None
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
