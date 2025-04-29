Now it's time to put your AI Hypercomputer to work! You'll be running a Python program on the Slurm cluster to analyse the vast dataset and determine the optimal formula for the Inferno Elixir.

### Task

1. **Copy the Python Script:**
   - Copy the provided Python script (`optimising_elixir.py`) to the Slurm Login Node
2. **Copy the Slurm Submission Script:**
   - Copy the provided slurm script that executes the `optimising_elixir.py` to the Slurm Login Node
3. **Submit the Job:**
   - Submit the Slurm job to the cluster (check the docs or Gemini to find out how!)
4. **Verification:**
   - Once the job has run, run the results collection tool to determine the best pepper combination for optimal youthfulness!

### Tips

- Ensure Python is installed on the cluster nodes.
- Check the Slurm job output for any errors.
- The output of the script will contain the optimal elixir formula.
- Use `nano` to easily copy files to the login node
  - If you prefer to use `vim` then ensure to use the `:set paste` command if pasting multi-line python code
- Use the `cat` command to obtain the contents of the slurm output (e.g. `cat output.txt`)
- You can use `watch squeue` to see the state of your Slurm job

### Provided Code

Expand the below to view your code and deployment script.

:::collapse{title="Show optimising_elixir.py"}

Filename:

```plaintext
optimising_elixir.py
```

::rawfile{file=optimising_elixir.py type=code language=python}

:::

:::collapse{title="submit.sh"}

Filename:

```plaintext
submit.sh
```

::rawfile{file=submit.sh type=code language=bash}

:::

### Running

Once you have submitted your job, you can monitor it with the `squeue` command:

```bash
watch squeue
```

### Verification

Once complete, run the following command to collect your results:

First set your job ID (e.g. `export JOB_ID=2`)

```bash
export JOB_ID=xx
```

Now run the collector program

```bash
python3 optimising_elixir.py --output-dir "results_job_${JOB_ID}" collect
```

Copy the output and paste it below to get your score!
