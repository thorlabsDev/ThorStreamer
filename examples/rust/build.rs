//! Build script for the thorStreamerExample project.
//!
//! This script uses `tonic_build` to compile the protocol buffer files located in the `proto` directory.
//! It also sets up the file descriptor and instructs Cargo to rerun the build if the proto files change.

fn main() {
    tonic_build::configure()
        .compile_well_known_types(true)
        .out_dir(std::env::var("OUT_DIR").unwrap())
        .protoc_arg("--experimental_allow_proto3_optional")
        .build_server(false)
        .file_descriptor_set_path(std::env::var("OUT_DIR").unwrap() + "/descriptor.bin")
        .compile(&["proto/events.proto", "proto/publisher.proto"], &["proto"])
        .expect("Failed to compile proto");

    println!("cargo:rerun-if-changed=proto/events.proto");
    println!("cargo:rerun-if-changed=proto/publisher.proto");
}
