resource "google_compute_instance_template" "new_template" {
  name         = var.template_name
  machine_type = var.machine_type

  network_interface {
    network        = var.vpc_name # Change this if using a custom VPC
    subnetwork     = var.snet_name
    # network_tier   = "STANDARD"
    stack_type     = "IPV4_ONLY"
  }

  service_account {
    scopes = [ "cloud-platform" ]
  }

  scheduling {
    automatic_restart   = false
    on_host_maintenance = "TERMINATE"
    provisioning_model  = "SPOT"
    preemptible         = true
    instance_termination_action = "STOP"
  }

  disk {
    auto_delete  = true
    boot         = true
    device_name  = var.template_name
    source_image = var.os #"projects/ubuntu-os-cloud/global/images/ubuntu-2410-oracular-amd64-v20241115"
    mode         = var.os_disk_access_mode
    type         = "pd-balanced"
    disk_size_gb = var.os_disk_size
  }

  disk {
    auto_delete = false
    boot        = false
    device_name = var.data_disk_name
    mode        = var.data_disk_access_mode
    source      = var.data_disk_name
  }

  shielded_instance_config {
    enable_secure_boot          = false
    enable_vtpm                 = true
    enable_integrity_monitoring = true
  }

  reservation_affinity {
    type = "ANY_RESERVATION"
  }
}