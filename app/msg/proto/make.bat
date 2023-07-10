@echo off

"./protoc.exe" -I=./ --plugin="./protoc-gen-go.exe" --go_out=../ ./*.proto


