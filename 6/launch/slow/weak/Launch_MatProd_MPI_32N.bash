#!/bin/bash
#SBATCH --account=tra20_IngInfBo
#SBATCH --partition=skl_usr_dbg
#SBATCH -t 00:30:00
#SBATCH --nodes=4
#SBATCH --ntasks-per-node=8
#SBATCH -o ../../../out/MatrixProductMPI_Slow_Weak.out
#SBATCH --open-mode=append

# env. variables and modules
module load autoload intelmpi
# execution lines
srun ../../../bin/MatrixProductMPI_Slow 3200 mpi.s.weak
