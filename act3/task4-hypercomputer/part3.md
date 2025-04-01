Now it's time to put your AI Hypercomputer to work! You'll be running a Python program on the Slurm cluster to analyse the vast dataset and determine the optimal formula for the Inferno Elixir.

### Task

1. **Copy the Python Script:**
   - Copy the provided Python script (`optimising_elixir.py`) to the Slurm Login Node
2. **Copy the Slurm Submission Script:**
   - Copy the provided slurm script that executes the `optimising_elixir.py` to the Slurm Login Node
3. **Submit the Job:**
   - Submit the Slurm job to the cluster.
4. **Verification:**
   - Once the job has run, copy the contents of your `output.txt` file and paste it in the verification box below.

### Tips

- Ensure Python is installed on the cluster nodes.
- Check the Slurm job output for any errors.
- The output of the script will contain the optimal elixir formula.
- Use `nano` to easily copy files to the login node
  - If you prefer to use `vim` then ensure to use the `:set paste` command if pasting multi-line python code
- Use the `cat` command to obtain the contents of the slurm output (e.g. `cat output.txt`)

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

### Verification
