Cymbal Supplements has a mountain of data that needs processing to perfect the Inferno Elixir. To handle this massive workload, you need to deploy a high-performance computing cluster using [Google Cloud's Cluster Toolkit](https://cloud.google.com/cluster-toolkit/docs/overview).

This cluster will be the foundation for your **AI Hypercomputer**, enabling you to analyze chili pepper varieties, fermentation times, customer feedback, and genetic markers.

### Task

1. **Deploy a Slurm Cluster:**
    * Use the provided Cluster Toolkit and YAML blueprint to deploy a Slurm cluster.
    * Ensure the cluster has sufficient compute resources to handle the upcoming AI workload.
2. **Verification:**
    * Provide a screenshot of the deployed Slurm cluster in the Google Cloud Console, showing the nodes and their status.
    * Verify that the cluster is in a running state.

### Tips

* Pay close attention to the YAML blueprint.
* Ensure all necessary Google Cloud APIs are enabled.
* Cluster Toolkit is your friend!
* Check out [the documentation](https://cloud.google.com/cluster-toolkit/docs/quickstarts/slurm-cluster) for step-by-step instructions
* You can run this task either from the Cloud Shell or your [provisioned workstation](https://console.cloud.google.com/workstations/overview?project=%%CLIENT_PROJECT_ID%%)

:::collapse{title="Show deployment.yaml"}
```yaml
---
blueprint_name: hpc-slurm

vars:
  project_id: %%CLIENT_PROJECT_ID%%
  deployment_name: hpc-slurm
  region: %%LOCATION%%
  zone: %%LOCATION%%-a

# Documentation for each of the modules used below can be found at
# https://github.com/GoogleCloudPlatform/hpc-toolkit/blob/main/modules/README.md

deployment_groups:
- group: primary
  modules:
  # Source is an embedded module, denoted by "modules/*" without ./, ../, /
  # as a prefix. To refer to a local module, prefix with ./, ../ or /
  - id: network
    source: modules/network/vpc

  # Private Service Access (PSA) requires the compute.networkAdmin role which is
  # included in the Owner role, but not Editor.
  # PSA is a best practice for Filestore instances, but can be optionally
  # removed by deleting the private_service_access module and any references to
  # the module by Filestore modules.
  # https://cloud.google.com/vpc/docs/configure-private-services-access#permissions
  - id: private_service_access
    source: community/modules/network/private-service-access
    use: [network]

  - id: homefs
    source: modules/file-system/filestore
    use: [network, private_service_access]
    settings:
      local_mount: /home

  - id: debug_nodeset
    source: community/modules/compute/schedmd-slurm-gcp-v6-nodeset
    use: [network]
    settings:
      node_count_dynamic_max: 4
      machine_type: n2-standard-2
      allow_automatic_updates: false

  - id: debug_partition
    source: community/modules/compute/schedmd-slurm-gcp-v6-partition
    use:
    - debug_nodeset
    settings:
      partition_name: debug
      exclusive: false # allows nodes to stay up after jobs are done
      is_default: true

  - id: compute_nodeset
    source: community/modules/compute/schedmd-slurm-gcp-v6-nodeset
    use: [network]
    settings:
      node_count_dynamic_max: 20
      bandwidth_tier: gvnic_enabled
      allow_automatic_updates: false

  - id: compute_partition
    source: community/modules/compute/schedmd-slurm-gcp-v6-partition
    use:
    - compute_nodeset
    settings:
      partition_name: compute

  - id: h3_nodeset
    source: community/modules/compute/schedmd-slurm-gcp-v6-nodeset
    use: [network]
    settings:
      node_count_dynamic_max: 20
      # Note that H3 is available in only specific zones. https://cloud.google.com/compute/docs/regions-zones
      machine_type: h3-standard-88
      # H3 does not support pd-ssd and pd-standard
      # https://cloud.google.com/compute/docs/compute-optimized-machines#h3_disks
      disk_type: pd-balanced
      bandwidth_tier: gvnic_enabled
      allow_automatic_updates: false

  - id: h3_partition
    source: community/modules/compute/schedmd-slurm-gcp-v6-partition
    use:
    - h3_nodeset
    settings:
      partition_name: h3

  - id: slurm_login
    source: community/modules/scheduler/schedmd-slurm-gcp-v6-login
    use: [network]
    settings:
      machine_type: n2-standard-4
      enable_login_public_ips: true

  - id: slurm_controller
    source: community/modules/scheduler/schedmd-slurm-gcp-v6-controller
    use:
    - network
    - debug_partition
    - compute_partition
    - h3_partition
    - homefs
    - slurm_login
    settings:
      enable_controller_public_ips: true
```
:::