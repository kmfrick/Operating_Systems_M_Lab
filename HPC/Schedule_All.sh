for i in $(find -name weak); do
	cd $i
	for j in $(ls .); do
		sbatch $j
	done
	cd ~/frick
done
for i in $(find -name strong); do
	cd $i
	for j in $(ls .); do
		sbatch $j
	done
	cd ~/frick
done

