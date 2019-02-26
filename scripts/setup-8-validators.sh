#!/bin/sh

PRIVATE_KEY_PATH="/home/preston/priv0x8882042b8e93c85312f623f058ef252c8025a7ae"

echo "clearing data"
DATA_PATH=/tmp/data
rm -rf $DATA_PATH
mkdir -p $DATA_PATH

CONTRACT="0xa5f48Bd414fdF5143FE037a7ed399B95C39a6e6e"
PASSWORD="password"
PASSWORD_PATH=$DATA_PATH/password.txt

echo $PASSWORD > $PASSWORD_PATH 

for i in `seq 1 8`;
do
  echo "Generating validator $i"
  
  KEYSTORE=$DATA_PATH/keystore$i

  bazel run //validator -- accounts create \
    --password=$PASSWORD \
    --keystore-path=$KEYSTORE

  bazel run //contracts/deposit-contract/sendDepositTx -- \
    --httpPath=https://goerli.prylabs.net \
    --passwordFile=$PASSWORD_PATH \
    --depositContract=$CONTRACT \
    --numberOfDeposits=1 \
    --privKey=$(cat $PRIVATE_KEY_PATH) \
    --prysm-keystore=$KEYSTORE \
    --depositAmount=3200000
done
