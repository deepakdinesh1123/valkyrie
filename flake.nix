{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.11";
    flake-utils.url = "github:numtide/flake-utils";
    nix2container.url = "github:nlewo/nix2container";
  };

  description = "Valkyrie";

  outputs = { self, nixpkgs, flake-utils, nix2container, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        nix2containerPkgs = nix2container.packages.${system};
        arch = builtins.head (builtins.match "^([^-]+)-.*" system);
      in
      rec {
        alpine = nix2containerPkgs.nix2container.pullImage {
          imageName = "alpine";
          imageDigest = "sha256:31687a2fdd021f85955bf2d0c2682e9c0949827560e1db546358ea094f740f12";
          arch = if arch == "x86_64" then "amd64" else "arm64";
          sha256 = "sha256-USv9xcTnGqG78ep3wEPPuidyL27nehNqjRioZDx+iQo=";
        };
        ubuntu = nix2containerPkgs.nix2container.pullImage {
          imageName = "ubuntu";
          imageDigest = "sha256:72297848456d5d37d1262630108ab308d3e9ec7ed1c3286a32fe09856619a782"; 
          arch = if arch == "x86_64" then "amd64" else "arm64";
          sha256 = "sha256-H2ddt+ZxnnzrGBoTyAVMs/qkQuUHG+HelIgcqzVcjS4="; 
        };
        valkyrieDependencies = with pkgs; [
          sqlc
          go-migrate
          go_1_22
          caddy
          pkg-config 
          just
        ] ++ lib.optionals stdenv.isLinux [
          nsjail
          gpgme
          libgpg-error
          libassuan
          btrfs-progs
          fuse-overlayfs
          gvisor
          crun
        ];

        packages = rec {
          valkyrie = pkgs.callPackage ./builds/nix/valkyrie.nix { inherit pkgs; };
          agent = pkgs.callPackage ./builds/nix/agent.nix { inherit pkgs; };
          valkyrieDockerImage = nix2containerPkgs.nix2container.buildImage {
            name = "valkyrie";
            tag = "0.0.1";
            fromImage = ubuntu;
            config = {
              Entrypoint = [ "${valkyrie}/bin/valkyrie" ];
            };
          };
        };
        defaultPackage = packages.valkyrie;

        devShells = {
          default = pkgs.mkShell {
            buildInputs = valkyrieDependencies;
          };
          valkyrie = pkgs.mkShell {
            buildInputs = valkyrieDependencies;
          };
          python-sdk = import ./sdk/pyvalkyrie/shell.nix { inherit pkgs; };
          js-sdk = import ./sdk/ts/shell.nix { inherit pkgs; };
          docs = import ./docs/shell.nix { inherit pkgs; };
          schemas = import ./schemas/shell.nix { inherit pkgs; };
          agent = pkgs.mkShell {
            buildInputs = [ pkgs.go_1_22 ];
          };
        };
      }
    );
}
