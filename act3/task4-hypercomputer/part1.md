Cymbal Supplements has a mountain of data that needs processing to perfect the Inferno Elixir. To handle this massive workload, you need to deploy a high-performance computing cluster using [Google Cloud's Cluster Toolkit](https://cloud.google.com/cluster-toolkit/docs/overview).

This cluster will be the foundation for your **AI Hypercomputer**, enabling you to analyse chili pepper varieties, fermentation times, customer feedback, and genetic markers.

### Task

1. **Deploy a Slurm Cluster:**
    * Use the provided Cluster Toolkit and YAML blueprint to deploy a Slurm cluster.
    * Ensure the cluster has sufficient compute resources to handle the upcoming AI workload.
2. **Verification:**
    * Provide a screenshot of the Slurm login screen. You should be able to navigate to your instance in the GCE console and ssh into it.

### Tips

* Check out [the documentation](https://cloud.google.com/cluster-toolkit/docs/quickstarts/slurm-cluster) for step-by-step instructions on how to deploy a cluster using the YAML provided below
* The [Cluster Toolkit](https://cloud.google.com/cluster-toolkit/docs/setup/configure-environment) has been built and compiled for you inside your [cloud workstation](https://console.cloud.google.com/workstations/overview?project=%%CLIENT_PROJECT_ID%%) (inside the `cluster-toolkit` folder). Use this or the Cloud Shell
* Ensure you have logged into your workstation with both user and application-default credentials:

 ```bash
gcloud auth login
```

```bash
gcloud auth application-default login
```

#### Deployment File

Expand the below to view your deployment specification to get you started

:::collapse{title="Show deployment.yaml"}

::rawfile{file=deployment.yaml type=code language=yaml}

:::
