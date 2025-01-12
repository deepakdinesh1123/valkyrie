module "vpc" {
  source = "../../modules/vpc"

  vpc_name = "valkyrie-vpc"
}

module "worker_subnet" {
  source = "../../modules/vpc/subnet"

  snet_name = "valkyrie-worker-snet"
  snet_cidr = "10.0.0.0/22"
  region = var.location
  vpc_id = module.vpc.vpc_id
}

module "server_subnet" {
  source = "../../modules/vpc/subnet"

  snet_name = "valkyrie-server-snet"
  snet_cidr = "10.0.4.0/22"
  region = var.location
  vpc_id = module.vpc.vpc_id
}

module "db_subnet" {
  source = "../../modules/vpc/subnet"

  snet_name = "valkyrie-db-snet"
  snet_cidr = "10.0.8.0/22"
  region = var.location
  vpc_id = module.vpc.vpc_id
}

module "ssh_firewall" {
  source = "../../modules/vpc/firewall"

  vpc_id = module.vpc.vpc_id
  rule_name = "allow-ssh"
  ports = [ "22" ]
  protocol_name = "tcp"
  source_ips = ["0.0.0.0/0"]
}

module "db_firewall" {
  source = "../../modules/vpc/firewall"

  vpc_id = module.vpc.vpc_id
  rule_name = "allow-db"
  ports = [ "5432" ]
  protocol_name = "tcp"
  source_ips = ["0.0.0.0/0"]
}

module "persistant_disk" {
  source = "../../modules/persistant_disk"

  disk_name = "shared-nix-store"
  disk_size_gb = var.disk_size_gb
  disk_type = "pd-balanced"
  project_name = var.project_name
  zone_name = var.availability_zone
  disk_iops = var.disk_iops
  disk_throughput = var.disk_throughput
}

module "server_vm" {
  source = "../../modules/compute_engine"

  vm_name = "odin-server"
  os_image = "projects/ubuntu-os-cloud/global/images/ubuntu-2410-oracular-amd64-v20241115"
  os_disk_size = 10
  os_disk_type = "pd-standard"

  disk_id = module.persistant_disk.disk_id
  disk_name = module.persistant_disk.disk_name

  machine_type = var.server_machine_type
  snet_id = module.server_subnet.snet_id

}

module "spot_vm_template" {
  source = "../../modules/spot_vm_template"

	template_name = "odin-worker-template"
	vpc_name = module.vpc.vpc_id
	snet_name = module.worker_subnet.snet_id
	machine_type = var.worker_machine_type
	os = "projects/ubuntu-os-cloud/global/images/ubuntu-2410-oracular-amd64-v20241115"
	os_disk_size = 10
  os_disk_access_mode = "READ_WRITE"
	data_disk_name = module.persistant_disk.disk_name
  data_disk_access_mode = "READ_ONLY"
  service_account_email = data.google_service_account.svc_acc.email
}

module "spot_mig" {
  source = "../../modules/spot_MIG"

  mig_name = "odin-worker"
  mig_target_size = 1
  vm_template_id = module.spot_vm_template.spot_vm_template_id
  data_disk_name = module.persistant_disk.disk_name
  zone_name = var.availability_zone
}