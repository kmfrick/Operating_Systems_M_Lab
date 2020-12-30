#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <omp.h>
#include <time.h>

const int MAX = 10;

void run(int m, int size, char *name)
{
	// Define my value
	int B[m][m], C[m * m];

	int **A = (int**) malloc(sizeof(int*) * m);
	int *A_data = (int*) malloc(sizeof(int) * m * m); // Ensure matrix is contiguous
	memset(A_data, 0, m * m * sizeof(*A_data));
	for (int i = 0; i < m; i++) {
		A[i] = &(A_data[m * i]);
	}
	srand((unsigned int)time(NULL)); 
	// Generate A and B
	for(int i = 0; i < m; i++) {
		for (int j = 0; j < m; j++) {
			A[i][j] = rand() % MAX + 1;
			B[i][j] = rand() % MAX + 1;
		}
	}
	/*
		 printf("A = \n");
		 for(int i = 0; i<m; i++) {
		 for (int j = 0; j < m; j++) {
		 printf("\t%d", B[i][j]);
		 }
		 printf("\n");
		 }

		 printf("B = \n");
		 for(int i = 0; i<m; i++) {
		 for (int j = 0; j < m; j++) {
		 printf("\t%d", B[i][j]);
		 }
		 printf("\n");
		 }
	 */

	// Zero out C
	for(int i = 0; i < m; i++) {
		for (int j = 0; j < m; j++) { 
			C[i * m + j] = 0;
		}
	}
#pragma omp barrier
	double start = omp_get_wtime(); 
#pragma omp parallel num_threads(size) shared(A, B, C, m)
#pragma omp for schedule(static)
	for(int i = 0; i < m; i++) {
		for (int j = 0; j < m; j++) { 
			for (int k = 0; k < m; k++) {
				C[i * m + j] = C[i * m + j] + (A[i][k] * B[k][j]); 
			}
		}
	}
#pragma omp barrier
	double end = omp_get_wtime();
	printf("%s, %d, %d, %lf\n", name, size, m, end - start);
	/*
		 printf("C = AB\n");
		 for(int i = 0; i< m; i++) {
		 for (int j = 0; j < m; j++) {
		 printf("\t%d", C[i * m + j]);
		 }
		 printf("\n");
		 }
	 */
}

int main(int argc, char* argv[])
{
	if (argc < 3) {
		printf("Usage: %s MATRIX_SIZE NUM_PROCESSES NAME\n", argv[0]);
		exit(1);
	}
	int m = atoi(argv[1]);
	int size = atoi(argv[2]);
	char* name = argv[3];
	int rank;
	if (size > m) {
		printf("No. of processes %d is greater than matrix size %d.\n",size, m);
	} else if (m % size) {
		printf("Matrix size %d is not a multiple of process count %d.\n", m, size);
	} else {
		run(m, size, name);
	}
	return 0;

}

