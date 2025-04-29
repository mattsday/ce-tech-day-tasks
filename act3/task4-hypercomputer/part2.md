Before running the complex AI optimisation program, you need to verify that your Slurm cluster is functioning correctly. A simple 'Hello World' workload will serve as a test run.

### Task

1. **Login to your controller:**
    * You completed this in the previous step. You can run this command to get back in:

```bash
gcloud compute ssh hpcslurm-slurm-login-001 --tunnel-through-iap --zone "%%LOCATION%%-b" --project %%CLIENT_PROJECT_ID%%
```

2. **Create a Simple Slurm Script:**
    * Write a Slurm script that prints the hostname of the compute node it runs on.
    * The [Quickstart guide](https://cloud.google.com/cluster-toolkit/docs/quickstarts/slurm-cluster#run_a_job_on_the_hpc_cluster) and/or Gemini will definitely help you here!
3. **Submit the Job:**
    * Submit the script to the Slurm cluster.
4. **Verification:**
    * Provide a screenshot of the Slurm job output, showing the hostname.

### Tips

* Use basic Slurm commands for job submission.
* Ensure your script is executable on the cluster nodes.
* Check the standard output of the Slurm job.
