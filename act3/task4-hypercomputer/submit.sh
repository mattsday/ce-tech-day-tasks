#!/bin/bash

#SBATCH --job-name=elixir_quest
#SBATCH --array=1-10          # Run 10 tasks in parallel
#SBATCH --cpus-per-task=1     # Each task uses 1 CPU core
#SBATCH --mem=150M            # Increased slightly for safety, still low
#SBATCH --time=00:10:00       # VERY SHORT time limit
#SBATCH --output=output.txt   # Show output in output.txt
#SBATCH --error=error.txt     # Show errors in error.txt
#SBATCH --open-mode=append    # Append logs to one file

# --- Configuration ---
PYTHON_SCRIPT="optimising_elixir.py"
RESULTS_DIR="results_job_${SLURM_ARRAY_JOB_ID}" # Unique results dir per job array
ITERATIONS=1000                                 # <<< CRITICAL: Adjust this based on time limit and cluster speed! Must be low for 10min.

# --- Setup ---
# Load Python module if needed (adjust for your cluster environment)
# module load python/3.9

# Create directories for logs and results if they don't exist
mkdir -p logs
mkdir -p "${RESULTS_DIR}"

# Echo job information
echo "Starting Slurm Task ${SLURM_ARRAY_TASK_ID} on $(hostname)"
echo "Job ID: ${SLURM_ARRAY_JOB_ID}"
echo "Script: ${PYTHON_SCRIPT}"
echo "Iterations: ${ITERATIONS}"
echo "Results Dir: ${RESULTS_DIR}"
echo "-----------------------------------------"

# --- Run the Python script in 'run' mode ---
# Pass the iterations and output dir to the script
python3 ${PYTHON_SCRIPT} run --iterations ${ITERATIONS} --output-dir ${RESULTS_DIR}

echo "-----------------------------------------"
echo "Finished Slurm Task ${SLURM_ARRAY_TASK_ID}"

# NOTE: The 'collect' step is run MANUALLY after the array job finishes.
