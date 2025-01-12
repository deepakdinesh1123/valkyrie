# GCP Deployment
Detailed manual setup of [Odin server](./server.md), [Odin worker](./worker.md) and [shared nix store](./shared_nix_store.md).  

For convenience, we have provided Tofu configuration to experiment with deploying Odin to GCP.

### Step1: Install tofu
Download [tofu](https://opentofu.org/docs/intro/install/) and follow the installation instructions for you OS.

### Step2: Authenticate with GCP
Generate keys for service account having appropriate access and use it in provider.tf  
Different [gcp authentiation](https://cloud.google.com/docs/terraform/authentication) ways
```
provider "google" {
  credentials = "path/to/keys.json"
  project     = "project-name"
  region      = "gcp-region"
}
```

### Step3: Configure your gcp settings
Take sample odin.tfvars file. Use it to define GCP resources like VM size, region. Note that this template creates a new resource group for your Odin deployment.
```
```

### Step 4: Initialize and deploy with Tofu
Then run the following commands to deploy your Odin stack.

**Initialize Terraform:**  
```
tofu init
```

**Plan the deployment, and review it to ensure it matches your expectations:**  
```
tofu plan -var-file odin.tfvars
```

**Finally, apply the deployment:**  
```
tofu apply -var-file odin.tfvars
```

**After a few minutes, you can get the IP address of your instance with**
```
tofu output -raw public_ip_address
```
Add the ip as A record to your domain example.com given in caddyfile