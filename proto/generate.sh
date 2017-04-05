#! /bin/bash
protoc -I=. --go_out=../hlcmsg ./*.proto