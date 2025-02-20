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

        odinDependencies = with pkgs; [
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
          odin = pkgs.callPackage ./builds/nix/odin.nix { inherit pkgs; };
          agent = pkgs.callPackage ./builds/nix/agent.nix { inherit pkgs; };
          odinDockerImage = nix2containerPkgs.nix2container.buildImage {
            name = "odin";
            tag = "0.0.1";
            fromImage = alpine;
            config = {
              Entrypoint = [ "${odin}/bin/odin" ];
            };
          };
        };

        defaultPackage = packages.odin;

        docsDependencies = with pkgs; [
          python312Packages.mkdocs-material
          redocly
          quicktype
          uv
        ];

        devShells = {
          default = pkgs.mkShell {
            buildInputs = odinDependencies ++ import ./frontend/deps.nix { inherit pkgs; };
          };
          odin = pkgs.mkShell {
            buildInputs = odinDependencies;
          };
          python-sdk = import ./sdk/odinpy/shell.nix { inherit pkgs; };
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
