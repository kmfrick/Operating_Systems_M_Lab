#!/bin/bash
#SBATCH --account=tra20_IngInfBo
#SBATCH --partition=skl_usr_dbg
#SBATCH -t 00:30:00
#SBATCH --nodes=2
#SBATCH --ntasks-per-node=1
#SBATCH -o ../../../out/MatrixProductMPI_Slow_Frick_2K_2N.out
# env. variables and modules
module load autoload intelmpi
# execution lines
srun ../../../bin/MatrixProductMPI_Slow 2000
