#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <mpi.h>
#include <time.h>

const size_t DIM = 10;
const int MAX = 10;

int main(int argc, char* argv[])
{
	MPI_Init(&argc, &argv);

	// Get number of processes and check that 4 processes are used
	int size;
	int rank;
	MPI_Comm_size(MPI_COMM_WORLD, &size); // size = num of components of MPI_COMM_WORLD
	MPI_Comm_rank(MPI_COMM_WORLD, &rank); // rank = process id inside MPI_COMM_WORLD
	if(size > DIM)
	{
		printf("No. of processes %d is greater than matrix size %d.\n",size, DIM);
		MPI_Abort(MPI_COMM_WORLD, EXIT_FAILURE);
	}
	if (DIM % size) {
		printf("Matrix size %d is not a multiple of process count %d.\n", DIM, size);
		MPI_Abort(MPI_COMM_WORLD, EXIT_FAILURE);
	}

	int rows = DIM / size;

	// Define my value
	int my_A[DIM * rows], my_C[DIM * rows];
	int B[DIM][DIM];

	if (rank == 0)
	{
		int **A = (int**) malloc(sizeof(int*) * DIM);
		int *A_data = (int*) malloc(sizeof(int) * DIM * DIM); // Ensure matrix is contiguous
		memset(A_data, 0, DIM * DIM * sizeof(*A_data));
		for (int i = 0; i < DIM; i++) {
			A[i] = &(A_data[DIM * i]);
		}
		srand((unsigned int)time(NULL)); 
		for(int i = 0; i < DIM; i++) {
			for (int j = 0; j < DIM; j++) {
				A[i][j]=rand() % MAX + 1;
				B[i][j]=rand() % MAX + 1;
			}
		}
		printf("A = \n");
		for(int i=0; i<DIM; i++) {
			for (int j = 0; j < DIM; j++) {
				printf("\t%d", B[i][j]);
			}
			printf("\n");
		}

		printf("B = \n");
		for(int i=0; i<DIM; i++) {
			for (int j = 0; j < DIM; j++) {
				printf("\t%d", B[i][j]);
			}
			printf("\n");
		}

		// Scatter matrix to processes
		MPI_Scatter(A_data, DIM * rows, MPI_INT, &my_A, DIM * rows, MPI_INT, 0, MPI_COMM_WORLD);
	} else  {
		// Scatter
		MPI_Scatter(NULL, DIM * rows, MPI_INT, &my_A, DIM * rows, MPI_INT, 0, MPI_COMM_WORLD);
	}
	MPI_Barrier(MPI_COMM_WORLD);

	MPI_Bcast(&B, DIM*DIM, MPI_INT, 0, MPI_COMM_WORLD);

	for(int i = 0; i < rows; i++) {
		for (int j = 0; j < DIM; j++) { 
			my_C[i * DIM + j] = 0;
		}
	}

	for(int i = 0; i < rows; i++) {
		for (int j = 0; j < DIM; j++) { 
			for (int k = 0; k < DIM; k++) {
				my_C[i * DIM + j] += my_A[i * DIM + k] * B[k][j]; 
			}
		}
	}

	if (rank==0) //collettore
	{   
		int  C[DIM][DIM];
		MPI_Gather(&my_C, DIM * rows, MPI_INT, C, DIM * rows, MPI_INT, 0, MPI_COMM_WORLD);
		printf("Risultato C = AB:\n");
		for(int i=0; i<DIM; i++) {
			for (int j = 0; j < DIM; j++) {
				printf("\t%d", C[i][j]);
			}
			printf("\n");
		}
	} else  {
		MPI_Gather(&my_C, DIM * rows, MPI_INT, NULL, 1, MPI_INT, 0, MPI_COMM_WORLD);
	}

	MPI_Finalize();

	return EXIT_SUCCESS;
}

