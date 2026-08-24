#[test]
fn rust_package_uses_edition_2024() {
    let manifest = std::fs::read_to_string("Cargo.toml").expect("read Cargo.toml");
    assert!(manifest.lines().any(|line| line == r#"edition = "2024""#));
}

#[test]
fn verification_runs_once_on_main_with_exact_tools() {
    let toolchain = std::fs::read_to_string("rust-toolchain.toml").expect("read Rust owner");
    let workflow = std::fs::read_to_string(".github/workflows/verify.yml")
        .expect("read verification workflow");
    assert!(
        toolchain
            .lines()
            .any(|line| line == r#"channel = "1.96.0""#)
    );
    assert!(workflow.contains("push:\n    branches: [main]"));
    assert!(workflow.contains("actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1"));
    assert!(workflow.contains("rust-toolchain.toml"));
    assert!(workflow.contains("toolchain: ${{ steps.rust-toolchain.outputs.channel }}"));
    assert!(!workflow.contains("toolchain: \"1.96.0\""));
}
