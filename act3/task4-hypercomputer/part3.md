Now it's time to put your AI Hypercomputer to work! You'll be running a Python program on the Slurm cluster to analyse the vast dataset and determine the optimal formula for the Inferno Elixir.

### Task

1. **Upload the Python Script:**
    * Upload the provided Python script (`optimising_elixir.py`) to the Slurm cluster.
2. **Create a Slurm Submission Script:**
    * Create a slurm script that executes the `optimising_elixir.py` script.
3. **Submit the Job:**
    * Submit the Slurm job to the cluster.
4. **Verification:**
    * Provide a screenshot of the Slurm job output, showing the optimal pepper variety, fermentation time, and youthfulness score.
    * Verify that the script ran successfully and produced the expected output.

### Tips

* Ensure Python is installed on the cluster nodes.
* Check the Slurm job output for any errors.
* The output of the script will contain the optimal elixir formula.

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
