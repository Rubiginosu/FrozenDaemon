#!/bin/sh
mkdir FrozenGo
echo 'Building FrozenGo x86'
GOOS=linux GOARCH=386 go build frozen.go
mv frozen FrozenGo/frozen.386
echo 'Building FrozenGo x64'
GOARCH=amd64 go build frozen.go
echo 'Done,Packing...'
mv frozen FrozenGo/frozen.amd64
tar -czvf FrozenGo.Exec.tar.gz FrozenGo
rm -rf FrozenGo
