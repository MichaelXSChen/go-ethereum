#!/bin/bash
if [ ! -f bootnode.key ];then
    ./bootnode -genkey bootnode.key
fi
nohup ./bootnode -nodekey=bootnode.key > bootnode.log&
