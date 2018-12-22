#/bin/bash

grep -h score StandardSetConcurrent* | cut -d " " -f "1 8" | sed 's/StandardSetConcurrent//g' | sed 's/40m0sC40_//g' | sed 's/20m0sC40_//g' | sort| sed 's/,//g' | awk ' { t = $1; $1 = $2; $2 = t; print; } ' | awk 'NR%2==0' | cut -d " " -f 1 | paste -sd+ - | bc
