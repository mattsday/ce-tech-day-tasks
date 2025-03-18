#!/bin/bash
#SBATCH --array=1-10
#SBATCH --cpus-per-task=1
#SBATCH --mem=100M
#SBATCH --time=00:10:00
#SBATCH --output=output_%a.txt

python optimising_elixir.py
