---
date: 2024-08-14
authors:
  - manoj
  - deepak
categories:
  - Tech
description:
  - How to save money using spot instances
title: Saving money using spot instances
---

# Saving money using spot instances

**You can read about why and how we started using spot instances for Odin workers [here](spot_instances.md)**

## Requirements:
- AWS/GCP Account

<!-- more -->

## Reason

Odin workers operate independently of servers, allowing for the use of spot instances for worker tasks. Even if a spot instance worker is stopped, the task state remains unchanged in the database. Consequently, unfinished tasks can be resumed and completed by other workers. This setup helps maintain continuity and efficiency in task execution despite the potential instability of spot instances.

## Implementation

### AWS Spot Fleet Instances

This fleet ensures a minimum number of workers are always operational and can autoscale based on demand, each worker will have the shared nix store and other dependencies installed and ready to execute tasks.

### GCP Instance Groups

Similarly the instance group ensures a minimum number of workers are always operational and can autoscale based on demand.