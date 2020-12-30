#include <stdio.h>
#include <math.h>
#include <sys/time.h>
#include <stdlib.h>
#include <stddef.h>
#include <mpi.h>
#include <time.h>


int main(int argc, char *argv[])
{
	MPI_Init(&argc, &argv);
	if (argc < 3) {
		printf("Usage: %s MATRIX_SIZE NAME\n", argv[0]);
		MPI_Abort(MPI_COMM_WORLD, EXIT_FAILURE);
	}
	int m = atoi(argv[1]);
	char* name = argv[2];

	int size, rank;
	double begin,end, local_elaps, global_elaps;// tempi
	int  d;
	MPI_Comm_size(MPI_COMM_WORLD, &size);
	MPI_Comm_rank(MPI_COMM_WORLD, &rank);
	if(size > m)
	{
		printf("No. of processes %d is greater than matrix size %d.\n",size, m);
		MPI_Abort(MPI_COMM_WORLD, EXIT_FAILURE);
	}
	if (m % size) {
		printf("Matrix size %d is not a multiple of process count %d.\n", m, size);
		MPI_Abort(MPI_COMM_WORLD, EXIT_FAILURE);
	}

	MPI_Barrier(MPI_COMM_WORLD);
	begin=MPI_Wtime();


	int A[m][m],B[m][m],C[m][m], i,j,k;
	int my_A[m*m], my_C[m][m];


	if (rank == 0){ // master:
		for(i=0;i<m;i++){
			for(j=0;j<m;j++){
				A[i][j]=i;
				B[i][j]=j;
			}
		}

		// invio i dati ad ogni processo

		//invio ad ogni processo Pk d=m/size righe
		d=m/size; // quante righe di A inviare ad ogni processo
		for (k=1; k<size; k++)
		{   //indice prima riga: k*d; indice prima colonna: 0
			i=k*d; 
			MPI_Send(&A[i][0], d*m, MPI_INT, k,0, MPI_COMM_WORLD);//invio d righe di A

		}
		for (i=0; i<d; i++)
			for(j=0; j<m; j++)
				my_A[m*i+j]=A[i][j];

	}
	else {
		d=m/size;
		MPI_Recv(&my_A,d*m, MPI_INT, 0, 0, MPI_COMM_WORLD, MPI_STATUS_IGNORE);

	}

	//broadcast matrix B a tutti
	MPI_Bcast(&B, m*m, MPI_INT, 0, MPI_COMM_WORLD);


	// a questo punto ogni processo ha my_A e B --> calcolo
	for(i=0;i<m/size;i++) {
		for(j=0;j<m;j++){
			my_C[i][j]=0;
			for( k=0;k<m;k++){
				my_C[i][j]+=my_A[m*i+k]*B[k][j];
			}
		}
	}


	if (rank==0) { 
		for (i=0; i<d; i++) // copio i risultati di rank 0 in C
			for(j=0; j<m; j++)
				C[i][j]=my_C[i][j];

		//ricevo da ogni processo Pk d=m/size righe di C

		for (k=1; k<size; k++)
		{   //indice prima riga: k*d
			i=k*d; 
			MPI_Recv(&C[i], d*m, MPI_INT, MPI_ANY_SOURCE,0, MPI_COMM_WORLD,MPI_STATUS_IGNORE );

		}
	}
	else // slave
		MPI_Send(&my_C[0][0], d*m, MPI_INT, 0,0, MPI_COMM_WORLD);

	//misura tempo:
	MPI_Barrier(MPI_COMM_WORLD);
	end=MPI_Wtime();
	local_elaps= end-begin;
	MPI_Reduce(&local_elaps, &global_elaps,1,MPI_DOUBLE,MPI_MAX,0,MPI_COMM_WORLD);
	if (rank == 0) {
		printf("%s, %d, %d, %lf\n", name, size, m, global_elaps);
	}




	MPI_Finalize();

}

