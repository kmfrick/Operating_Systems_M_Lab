#!/bin/bash
#SBATCH --account=tra20_IngInfBo
#SBATCH --partition=skl_usr_dbg
#SBATCH -t 00:15:00
#SBATCH --nodes=4
#SBATCH --ntasks-per-node=2
#SBATCH -o ../../out/MatrixProductMPI_Slow_Frick_4K_8N.out
# env. variables and modules
module load autoload intelmpi
# execution lines
srun ../../bin/MatrixProductMPI_Slow 4000
