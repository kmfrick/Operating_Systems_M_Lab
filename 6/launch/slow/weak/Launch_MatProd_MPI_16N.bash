#!/bin/bash
#SBATCH --account=tra20_IngInfBo
#SBATCH --partition=skl_usr_dbg
#SBATCH -t 00:30:00
#SBATCH --nodes=4
#SBATCH --ntasks-per-node=4
#SBATCH -o ../../../out/MatrixProductMPI_Slow_Frick_16K_16N.out
# env. variables and modules
module load autoload intelmpi
# execution lines
srun ../../../bin/MatrixProductMPI_Slow 16000
