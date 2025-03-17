{
  inputs = {
    nixpkgs.url = "/var/cache/nixpkgs/NixOS-nixpkgs-b27ba4e";
    flake-parts.url = "/var/cache/flake-parts-main";
    process-compose-flake.url = "/var/cache/process-compose-flake-main";
    services-flake.url = "/var/cache/services-flake-main";
  };

  outputs = inputs: inputs.flake-parts.lib.mkFlake { inherit inputs; } {
    systems = [ "x86_64-linux" "aarch64-darwin" "x86_64-darwin" "aarch64-linux" ];
    imports = [
      inputs.process-compose-flake.flakeModule
    ];
    perSystem = { self', pkgs, lib, system, ... }: {
      _module.args.pkgs = import inputs.nixpkgs {
        inherit system;
        config.allowUnfree = true;
      };

      packages.default = pkgs.buildEnv {
        name = "Odin-Sandbox-Env";
        paths = with pkgs; [
        ];
      };

      devShells.default = pkgs.mkShell {
        buildInputs = with pkgs; [
        ];
      };
      process-compose."odin" = pc: {
        imports = [
          inputs.services-flake.processComposeModules.default
        ];
        services = {
        };
      };
    };
    flake = {
    };
  };
}
