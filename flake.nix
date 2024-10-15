{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  description = "Valkyrie";

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      rec {

        docsDependencies = with pkgs; [ python312Packages.mkdocs-material redocly ];
        devDependencies = with pkgs; [ 
          sqlc 
          go-migrate 
          go_1_22
          # caddy
          pkg-config ] ++ lib.optionals stdenv.isLinux [ 
            gpgme 
            libgpg-error 
            libassuan
            btrfs-progs
          ];

        packages = {
          odin = pkgs.callPackage ./build/package/nix/odin.nix { inherit pkgs; };
        };
		    defaultPackage = packages.odin;

        devShells = {
          default = pkgs.mkShell {
            buildInputs = devDependencies;
          };
          load-test = pkgs.mkShell {
            buildInputs = with pkgs; [ k6 openapi-generator-cli go_1_22 ];
          };
        };
    }
  );
}
