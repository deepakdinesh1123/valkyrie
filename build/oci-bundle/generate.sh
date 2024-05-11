#!/bin/bash
set -e

OCI_BUNDLES_DIR="oci-bundles"
mkdir -p "$OCI_BUNDLES_DIR"

# Initialize PLATFORMS array if not already set
if [ -z "$PLATFORMS" ]; then
    PLATFORMS=("nix")
else
    # Split PLATFORMS string into an array using ':' as delimiter
    IFS=':' read -ra PLATFORMS <<< "$PLATFORMS"
fi

# Loop through each platform
for platform in "${PLATFORMS[@]}"; do
    # Create rootfs directory for the platform
    mkdir -p "$OCI_BUNDLES_DIR/$platform/rootfs"

    dockerImage="odin_platform_$platform"

    docker build -t $dockerImage -f "build/platforms/$platform.dockerfile" .
    
    # Run the Docker container and capture its ID
    containerId=$(docker run -d $dockerImage)
    
    # Export the container's filesystem to the platform's rootfs directory
    docker export "$containerId" | tar -C "$OCI_BUNDLES_DIR/$platform/rootfs" -xf -
    
    # Change directory to the platform's rootfs
    cd "$OCI_BUNDLES_DIR/$platform"
    
    # Run runc spec (assuming runc is available in the container)
    runc spec
    
    # Modify the config.json file
    jq '.process.terminal = false' config.json > temp.json && mv temp.json config.json
    
    # Stop the Docker container
    docker stop "$containerId"
    docker rm "$containerId"
    
    # Return to the parent directory
    cd ../..
done
