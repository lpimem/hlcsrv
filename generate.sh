#! /bin/bash
cd hlcproto
protoc -I=. --go_out=../hlcmsg ./*.proto