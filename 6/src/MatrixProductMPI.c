#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <mpi.h>
#include <time.h>

const int MAX = 10;

int main(int argc, char* argv[])
{
	MPI_Init(&argc, &argv);
	if (argc < 3) {
		printf("Usage: %s MATRIX_SIZE NAME\n", argv[0]);
		MPI_Abort(MPI_COMM_WORLD, EXIT_FAILURE);
	}
	int m = atoi(argv[1]);
	char* name = argv[2];

	// Get number of processes and check that 4 processes are used
	int size;
	int rank;
	MPI_Comm_size(MPI_COMM_WORLD, &size); // size = num of components of MPI_COMM_WORLD
	MPI_Comm_rank(MPI_COMM_WORLD, &rank); // rank = process id inside MPI_COMM_WORLD
	if(size > m)
	{
		printf("No. of processes %d is greater than matrix size %d.\n",size, m);
		MPI_Abort(MPI_COMM_WORLD, EXIT_FAILURE);
	}
	if (m % size) {
		printf("Matrix size %d is not a multiple of process count %d.\n", m, size);
		MPI_Abort(MPI_COMM_WORLD, EXIT_FAILURE);
	}

	int rows = m / size;
	MPI_Bcast(&rows, 1, MPI_INT, 0, MPI_COMM_WORLD);

	// Define my value
	int *my_A = NULL, *my_C = NULL;
	int *B = (int*) malloc(m * m * sizeof(int));
	double begin;

	if (rank == 0)
	{
		int *A = (int*) malloc(sizeof(int) * m * m); // Ensure matrix is contiguous
		srand((unsigned int)time(NULL)); 
		for(int i = 0; i < m; i++) {
			for (int j = 0; j < m; j++) {
				A[i * m + j]=rand() % MAX + 1;
				B[i * m + j]=rand() % MAX + 1;
			}
		}
		/*
			 printf("A = \n");
			 for(int i=0; i<m; i++) {
			 for (int j = 0; j < m; j++) {
			 printf("\t%d", B[i][j]);
			 }
			 printf("\n");
			 }

			 printf("B = \n");
			 for(int i=0; i<m; i++) {
			 for (int j = 0; j < m; j++) {
			 printf("\t%d", B[i][j]);
			 }
			 printf("\n");
			 }
			 */
		// Scatter matrix to processes
		my_A = (int*)malloc(m * rows * sizeof(int));
		if (my_A == NULL) {
			printf("my_A is NULL n process %d\n", rank);
			MPI_Abort(MPI_COMM_WORLD, EXIT_FAILURE);
		}
		MPI_Barrier(MPI_COMM_WORLD);
		begin = MPI_Wtime();
		MPI_Scatter(A, m * rows, MPI_INT, my_A, m * rows, MPI_INT, 0, MPI_COMM_WORLD);
	} else  {
		// Scatter
		my_A = (int*)malloc(m * rows * sizeof(int));
		if (my_A == NULL) {
			printf("my_A is NULL n process %d\n", rank);
			MPI_Abort(MPI_COMM_WORLD, EXIT_FAILURE);
		}
		MPI_Barrier(MPI_COMM_WORLD);
		begin = MPI_Wtime();
		MPI_Scatter(NULL, m * rows, MPI_INT, my_A, m * rows, MPI_INT, 0, MPI_COMM_WORLD);
	}
	MPI_Bcast(B, m*m, MPI_INT, 0, MPI_COMM_WORLD);

	my_C = (int*)malloc(m * rows * sizeof(int));
	if (my_C == NULL) {
		printf("my_C is NULL n process %d\n", rank);
	}
	for(int i = 0; i < rows; i++) {
		for (int j = 0; j < m; j++) { 
			my_C[i * m + j] = 0;
		}
	}

	for(int i = 0; i < rows; i++) {
		for (int j = 0; j < m; j++) { 
			for (int k = 0; k < m; k++) {
				my_C[i * m + j] += my_A[i * m + k] * B[k * m + j]; 
			}
		}
	}

	if (rank==0) //collettore
	{   
		int  C[m][m];
		MPI_Gather(my_C, m * rows, MPI_INT, C, m * rows, MPI_INT, 0, MPI_COMM_WORLD);
		/*
			 printf("Risultato C = AB:\n");
			 for(int i=0; i<m; i++) {
			 for (int j = 0; j < m; j++) {
			 printf("\t%d", C[i][j]);
			 }
			 printf("\n");
			 }
			 */
	} else  {
		MPI_Gather(my_C, m * rows, MPI_INT, NULL, 1, MPI_INT, 0, MPI_COMM_WORLD);
	}
	MPI_Barrier(MPI_COMM_WORLD);
	double end=MPI_Wtime();
	double local_elaps= end-begin;
	double global_elaps;
	MPI_Reduce(&local_elaps, &global_elaps,1,MPI_DOUBLE,MPI_MAX,0,MPI_COMM_WORLD);
	if (rank == 0) {			
		printf("%s, %d, %d, %lf\n", name, size, m, global_elaps);
	}
	MPI_Finalize();

	return EXIT_SUCCESS;
}

