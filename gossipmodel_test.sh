#!/bin/bash

NUM_TESTS=10

SIZE=100 # Need to change
FANOUT=15
NUMEXP=10
DEFAULT_PROB=1
SCENARIO_NUM=6

LIMIT=10

for ((scenario=0; scenario < $SCENARIO_NUM; scenario++))
do
	for ((fanout=1; fanout < $FANOUT; fanout++))
	do
		counter=0
		for ((test_id=0; test_id < $NUM_TESTS; test_id++))
		do
			OUT=$(./gossipmodel -c $NUMEXP -p $DEFAULT_PROB -f $fanout -k 0.7/10 -k 0.9/20 -s $SIZE -e $scenario)
			iter=$(echo "$OUT" | tail -n3 | head -n1 | awk '{print $1}' | awk -F ':' '{print$1}')
			finexp=$(echo "$OUT" | tail -n3 | head -n1 | awk '{print $1}' | awk -F ':' '{print$2}')
			restime=$(echo "$OUT" | tail -n1)
			if [ $iter == "inf" ] || [ $iter -ne 3 ]
			then
				if [ $counter -eq $LIMIT ]
				then
					counter=0
					echo "scenario=$scenario fanout=$fanout iter=$iter finexp=$finexp time=$restime - FAIL!"
				else
					counter=$(($counter+1))
					test_id=$(($test_id-1))
				fi
			else
				echo "scenario=$scenario fanout=$fanout iter=$iter finexp=$finexp time=$restime"
			fi
		done
		echo "====================="
	done
	echo "++++++++++++++++++++++++++++++++++++++++++++++++"
	echo "++++++++++++++++++++++++++++++++++++++++++++++++"
	echo "++++++++++++++++++++++++++++++++++++++++++++++++"
done
