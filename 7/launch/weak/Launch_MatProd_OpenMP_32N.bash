#!/bin/bash
#SBATCH --account=tra20_IngInfBo
#SBATCH --partition=skl_usr_dbg
#SBATCH -t 00:30:00
#SBATCH --nodes=1
#SBATCH --ntasks-per-node=1
#SBATCH -c 32
#SBATCH -o ../../../out/MatrixProductOpenMP_Frick_32K_32N.out
# env. variables and modules
module load autoload intelmpi
# execution lines
srun ../../../bin/MatrixProductOpenMP 32000 32
