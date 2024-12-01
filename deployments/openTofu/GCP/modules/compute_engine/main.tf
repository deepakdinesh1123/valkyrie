
resource "google_compute_instance" "compute-engine" {
  attached_disk {
    device_name = var.disk_name
    mode        = "READ_ONLY"
    source      = var.disk_id
  }

  zone = "asia-south1-c"

  boot_disk {
    auto_delete = true
    device_name = var.vm_name

    initialize_params {
      image = var.os_image
      size  = var.os_disk_size
      type  = var.os_disk_type
    }

    mode = "READ_WRITE"
  }

  can_ip_forward      = false
  deletion_protection = false
  enable_display      = false


  machine_type = var.machine_type

  metadata = {
    startup-script = "echo 'hello' > test.txt"
  }

  name = var.vm_name

  network_interface {
    access_config {
      network_tier = "STANDARD"
    }

    queue_count = 0
    stack_type  = "IPV4_ONLY"
    subnetwork  = var.snet_id
  }

  scheduling {
    # automatic_restart   = false
    # on_host_maintenance = "TERMINATE"
    # preemptible         = false
    provisioning_model  = "STANDARD"
  }

  service_account {
    scopes = [ "cloud-platform" ]
  }

  shielded_instance_config {
    enable_integrity_monitoring = true
    enable_secure_boot          = false
    enable_vtpm                 = true
  }

}
