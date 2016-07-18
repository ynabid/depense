#!/bin/sh


cp $GOPATH/bin/depense res
tar -czf res.tar.gz res/
scp -i "/home/yassine/Dev/keystore/AmsKEY.pem" res.tar.gz  ec2-user@ec2-52-33-196-167.us-west-2.compute.amazonaws.com:.

rm res/depense
rm res.tar.gz
ssh -i "/home/yassine/Dev/keystore/AmsKEY.pem" ec2-user@ec2-52-33-196-167.us-west-2.compute.amazonaws.com


