Cymbal Supplements has a mountain of data that needs processing to perfect the Inferno Elixir. To handle this massive workload, you need to deploy a high-performance computing cluster using [Google Cloud's Cluster Toolkit](https://cloud.google.com/cluster-toolkit/docs/overview). We've saved you a job and pre-installed this onto your workstation 💪.

This cluster will be the foundation for your **AI Hypercomputer**, enabling you to analyse chili pepper varieties, fermentation times, customer feedback, and genetic markers.

### Task

1. Open Cloud Workstations and start your instance: [Cloud Workstations](https://console.cloud.google.com/workstations/overview?project=%%CLIENT_PROJECT_ID%%).
2. Inside workstations, select **Open Folder** and select **cluster-toolkit**
3. Select the menu in the top-left and select **Terminal --> New Terminal**
4. You can now use the `gcluster` command to deploy a cluster! Use the below deployment file and [the documentation](https://cloud.google.com/cluster-toolkit/docs/overview) to get started!
5. Read teh tips below if you get stuck

#### Deployment File

Expand the below to view your deployment specification to get you started

:::collapse{title="Show deployment.yaml"}

::rawfile{file=deployment.yaml type=code language=yaml}

:::

### Verification

* Provide a screenshot of the Slurm login screen. You should be able to navigate to your instance in the GCE console and ssh into it.
* Login to your cluster with the following command:

```bash
gcloud compute ssh hpcslurm-slurm-login-001 --tunnel-through-iap --zone "%%LOCATION%%-b" --project %%CLIENT_PROJECT_ID%%
```

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

