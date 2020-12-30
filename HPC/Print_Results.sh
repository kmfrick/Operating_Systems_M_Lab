#!/bin/bash

files=$(find -type f -path *out)
if [[ ${files} == "" ]]; then
	echo "No files to process";
	exit
fi
echo "aname, n, m, t"
for file in ${files}
do
	cat ${file}
done

