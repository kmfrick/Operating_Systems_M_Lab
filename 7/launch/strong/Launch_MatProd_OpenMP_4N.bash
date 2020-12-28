#!/bin/bash
#SBATCH --account=tra20_IngInfBo
#SBATCH --partition=skl_usr_dbg
#SBATCH -t 00:15:00
#SBATCH --nodes=1
#SBATCH --ntasks-per-node=1
#SBATCH -c 4
#SBATCH -o ../../../out/MatrixProductOpenMP_Frick_4K_4N.out
# env. variables and modules
module load autoload intelmpi
# execution lines
srun ../../../bin/MatrixProductOpenMP 4000 4
