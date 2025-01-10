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

## Installing Odin
Odin can be installed using either of these methods:
Method 1: Using Nix Flakes
!!! example "Steps for Nix Flakes Installation"
    1. Enable flakes in your Nix configuration
    2. Add Odin flake to your inputs
    3. Import and use Odin in your configuration
Method 2: Using Nix Shell
!!! example "Steps for Nix Shell Installation"
    1. Create a shell.nix file
    2. Add Odin dependencies
    3. Enter the Nix shell environment
!!! info "Version Information"
    Check the official documentation for specific version compatibility details.

## Next Steps
After installation, you might want to:

Configure your environment
Set up your first Odin project
Explore Odin's features

!!! success "Need Help?"
    If you encounter any issues during installation, please refer to our troubleshooting guide or reach out to the community for support.