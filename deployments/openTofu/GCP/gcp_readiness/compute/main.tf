module "persistant_disk" {
  source = "../../modules/persistant_disk"

  disk_name = "nix-disk"
  disk_size_gb = 200
  disk_type = "pd-balanced"
  project_name = "Valkyrie-project"
  zone_name = "asia-south1-c"
  disk_iops = 3500
  disk_throughput = 300
}

module "server_vm" {
  source = "../../modules/compute_engine"

  vm_name = "valkyrie-server"
  os_image = "projects/ubuntu-os-cloud/global/images/ubuntu-2410-oracular-amd64-v20241115"
  os_disk_size = 10
  os_disk_type = "pd-standard"

  disk_id = module.persistant_disk.disk_id
  disk_name = module.persistant_disk.disk_name

  machine_type = "n2-custom-4-4096"
  snet_id = data.google_compute_subnetwork.server-subnet.id

}

module "spot_vm_template" {
  source = "../../modules/spot_vm_template"

	template_name = "valkyrie-worker-template"
	vpc_name = data.google_compute_network.vpc.id #module.vpc.id
	snet_name = data.google_compute_subnetwork.worker-subnet.id #module.vpc.snet_id
	machine_type = "n2-custom-4-8192"
	os = "projects/ubuntu-os-cloud/global/images/ubuntu-2410-oracular-amd64-v20241115"
	os_disk_size = 10
  os_disk_access_mode = "READ_WRITE"
	data_disk_name = module.persistant_disk.disk_name
  data_disk_access_mode = "READ_ONLY"
}

module "spot_mig" {
  source = "../../modules/spot_MIG"

  mig_name = "valkyrie-worker"
  mig_target_size = 1
  vm_template_id = module.spot_vm_template.spot_vm_template_id
  data_disk_name = module.persistant_disk.disk_name
  zone_name = "asia-south1-c"
}