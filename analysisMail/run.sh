#!/bin/bash

go build .

NAME="./analysisMail"
ID=`ps -ef | grep "$NAME" | grep -v "grep" | awk '{print $2}'`
echo $ID

for id in $ID
do
kill -9 $id
echo "killed $id"
done

nohup $NAME prod > ./run.log &
