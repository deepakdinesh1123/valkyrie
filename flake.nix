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
        docsDependencies = with pkgs; [ python312Packages.mkdocs-material ];
        devDependencies = with pkgs; [ sqlc go-migrate go_1_22 ] ++ docsDependencies ;
        packages = {
          odin = pkgs.callPackage ./build/package/nix/odin.nix { inherit pkgs; };
          geri = pkgs.callPackage ./build/package/nix/geri.nix { inherit pkgs; };
        };
		    defaultPackage = packages.odin;

        devShells = {
          default = pkgs.mkShell {
            buildInputs = [ packages.odin packages.geri ] ++ devDependencies ++ docsDependencies;
            shellHook = ''
              docker-compose up rabbitmq postgres -d
            '';
          };
          odin = pkgs.mkShell {
            buildInputs = [ packages.odin ];
          };
          geri = pkgs.mkShell {
            buildInputs = [ packages.geri ];
          };
          docs = pkgs.mkShell {
            buildInputs = docsDependencies;
            shellHook = ''
              mkdocs serve
            '';
          };
      };
    }
  );
}
