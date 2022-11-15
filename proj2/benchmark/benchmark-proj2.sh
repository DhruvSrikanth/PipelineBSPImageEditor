#!/bin/bash
#
#SBATCH --mail-user=dhruvsrikanth@uchicago.edu
#SBATCH --mail-type=ALL
#SBATCH --job-name=proj2_benchmark 
#SBATCH --output=./slurm/out/%j.%N.stdout
#SBATCH --error=./slurm/out/%j.%N.stderr
#SBATCH --chdir=/home/dhruvsrikanth/temp/parallel_proj_2/project_2/proj2/grader
#SBATCH --partition=debug 
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=16
#SBATCH --mem-per-cpu=900
#SBATCH --exclusive
#SBATCH --time=30:00

module load golang/1.16.2 
python3 benchmark_graph.py