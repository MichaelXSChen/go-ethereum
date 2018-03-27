#!/bin/bash
./clear.sh
echo "clearing workspace"
sleep 1
./bootnode.sh
echo "starting bootnode"
for i in `seq 1 3`; do
    sleep 1
    echo -n "."
done 
echo " "
echo "starting geth node"
./start.sh 00
