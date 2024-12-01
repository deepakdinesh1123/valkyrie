resource "google_compute_instance_group_manager" "instance_group" {
  name               = var.mig_name
  base_instance_name = "${var.mig_name}-1"
  target_size        = var.mig_target_size
  zone               = var.zone_name
  
   version {
    instance_template = var.vm_template_id
    name              = "v1" # Optional but helps track versions
  }

  update_policy {
    type            = "PROACTIVE"
    minimal_action  = "RESTART"
    replacement_method = "RECREATE"
    max_unavailable_fixed = 1
  }

  instance_lifecycle_policy {
    force_update_on_repair    = "YES"
    # default_action_on_failure = "REPAIR"
  }
}

resource "google_compute_autoscaler" "autoscaler" {
  name    = "${var.mig_name}-autoscaler"
  zone    = var.zone_name
  target  = google_compute_instance_group_manager.instance_group.self_link

  autoscaling_policy {
    
    min_replicas    = 1
    max_replicas    = 10
    cooldown_period = 60

    cpu_utilization {
      target = 0.95
    }
  }
}