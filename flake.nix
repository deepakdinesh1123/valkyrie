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

        k8sDependencies = with pkgs; [ 
          skaffold 
          k3d 
          skaffold 
          kubectl 
          kubectx 
          kubens 
          helm ];
        devDependencies = with pkgs; [ 
          sqlc 
          go-migrate 
          go_1_22  
          pkg-config ] ++ lib.optionals stdenv.isLinux [ 
            gpgme 
            libgpg-error 
            libassuanbtrfs-progs 
          ] ;

        packages = {
          odin = pkgs.callPackage ./build/package/nix/odin.nix { inherit pkgs; };
          nardump = pkgs.callPackage ./build/package/nix/nardump.nix { inherit pkgs; };
        };
		    defaultPackage = packages.odin;

        devShells = {
          default = pkgs.mkShell {
            buildInputs = devDependencies;
          };
          docs = pkgs.mkShell {
            packages = [
              pkgs.python3
              (pkgs.python3.withPackages ( python-pkgs: [
                python-pkgs.mkdocs
                python-pkgs.mkdocs-material
                python-pkgs.mike
                python-pkgs.mkdocs-minify-plugin
              ]))
              pkgs.mkdocs
            ];
          };
        };
    }
  );
}
