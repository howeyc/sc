set GOPATH=%BuildFolder%
go get -v
set GOARCH=386
go build -a -v -o sc_386.exe
set GOARCH=amd46
go build -a -v -o sc_amd64.exe

