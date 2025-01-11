{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  description = "Valkyrie";

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      rec {
        devDependencies = with pkgs; [
          sqlc
          go-migrate
          go_1_22
          nodejs_20
          # podman-compose
          # caddy
          postgresql_16
          pkg-config ] ++ lib.optionals stdenv.isLinux [
            nsjail
            gpgme
            libgpg-error
            libassuan
            btrfs-progs
            fuse-overlayfs
          ];

        packages = {
          odin = pkgs.callPackage ./build/package/nix/odin.nix { inherit pkgs; };
        };
		    defaultPackage = packages.odin;

        devShells = {
          default = pkgs.mkShell {
            buildInputs = devDependencies;
          };
          dev = pkgs.mkShell {
            buildInputs = devDependencies;
          };
          docs = import ./docs/shell.nix { inherit pkgs; };
        };
    }
  );
}
