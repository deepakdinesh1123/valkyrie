# AWS Deployment
Detailed manual setup of [Odin server](./server.md), [Odin worker](./worker.md) and [shared nix store](./shared_nix_store.md).    
For convenience, we have provided Tofu configuration to experiment with deploying Odin to AWS.

### Step1: Install tofu
Download [tofu](https://opentofu.org/docs/intro/install/) and follow the installation instructions for you OS.

### Step2: Authenticate with AWS
[AWS authentiation](https://search.opentofu.org/provider/opentofu/aws/latest) uses access and secret key in tofu provider
```
provider "aws" {
  region     = "us-west-2"
  access_key = "my-access-key"
  secret_key = "my-secret-key"
}
```

### Step3: Configure your aws settings
Take sample odin.tfvars file. Use it to define AWS resources like VM size, region. Note that this template creates a new resource group for your Odin deployment.
```
location                = "us-east-1"
snet_availability_zone1 = "us-east-1a"
snet_availability_zone2 = "us-east-1b"
key_pair_name           = "ec2_key_pair"
ebs_size             = 80
multi_attach_enabled = true
ebs_iops             = 1000
ebs_type             = "io1"
ec2_instance_type = "t3.micro"
rds_compute_type  = "db.t3.micro"
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
