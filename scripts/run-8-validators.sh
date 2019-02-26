#!/bin/sh

DATA_PATH=/tmp/data
PASSWORD_PATH=$DATA_PATH/password.txt
PASSWORD="password"

SESSION=8validators


bazel build //validator

for i in `seq 1 8`;
do
  KEYSTORE=$DATA_PATH/keystore$i
  DATADIR=$DATA_PATH/data$i
  
  CMD="bazel-bin/validator/linux_amd64_pure_stripped/validator --demo-config --password $PASSWORD --keystore-path $KEYSTORE"

  nohup $CMD $> /tmp/validator$i.log &
done

echo "8 validators are running in the background. You can follow their logs at /tmp/validator#.log where # is replaced by the validator index of 1 through 8."

echo "To stop the processes, use 'pkill validator'"
