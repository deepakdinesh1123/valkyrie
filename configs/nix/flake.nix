{
  inputs = {
    nixpkgs.url = "/var/cache/nixpkgs/NixOS-nixpkgs-b27ba4e";
    flake-utils.url = "/var/cache/flake-utils/numtide-flake-utils-11707dc";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
  } @ inputs:
    flake-utils.lib.eachDefaultSystem
    (
      system: let
        pkgs = import nixpkgs {
          inherit system;
        };
      in {
        packages.default = pkgs.buildEnv {
          name = "Odin-Sandbox-Environment";
          paths = with pkgs; [
            vim
            gnupatch
          ];
        };
      }
    );
}
