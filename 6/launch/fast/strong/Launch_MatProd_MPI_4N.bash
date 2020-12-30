#!/bin/bash
#SBATCH --account=tra20_IngInfBo
#SBATCH --partition=skl_usr_dbg
#SBATCH -t 00:15:00
#SBATCH --nodes=4
#SBATCH --ntasks-per-node=1
#SBATCH -o ../../../out/MatrixProductMPI_Fast_Strong.out
#SBATCH --open-mode=append
# env. variables and modules
module load autoload intelmpi
# execution lines
srun ../../../bin/MatrixProductMPI 400 mpi.f.strong
