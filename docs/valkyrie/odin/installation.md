# Installing Odin
This guide walks you through the process of installing Odin and its prerequisites.
Prerequisites
## Installing Nix
Nix can be installed in two modes: Single User or Multi User. Choose the appropriate installation method based on your needs.
### Single User Installation
```shell
sh <(curl -L https://nixos.org/nix/install) --no-daemon
```
!!! warning "Warning"
    Single User installation is required if you plan to use Podman as your container engine for Odin.
### Multi User Installation

```
sh <(curl -L https://nixos.org/nix/install) --daemon
```

### [Determinate Nix Installer](https://github.com/DeterminateSystems/nix-installer)

You can also use determinate nix installer to install nix easily, it allows you to easily roll back all the changes and uninstall nix if something goes wrong

## Installing Odin

### With Flakes

!!! note "Note"
    This is a temporary method until we start releasing versions of the odin package in our own binary cache

create a flake.nix with the following content

```nix
{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.11";
    flake-utils.url = "github:numtide/flake-utils";
    valkyrie.url = "github:deepakdinesh1123/valkyrie/development";
  };

  outputs = { self, nixpkgs, flake-utils, valkyrie, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        valkyriePkgs = valkyrie.packages.${system}; # Access packages from the Valkyrie flake
      in {
        devShells = {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [
              valkyriePkgs.odin # Use the odin package from Valkyrie
            ];
          };
        };
      }
    );
}

```

!!! tip "Tip"
    You can use the same flake input to add this package to your home-manager or nix-darwin or NixOS config

Then run the following command to create a shell with odin installed

```shell
nix develop
```
