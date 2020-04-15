set GOBIN=C:\gomodules\modules\gochess

set GOOS=linux
set GOARCH=amd64

go install uciengine/zurimain.go

set GOOS=windows
set GOARCH=amd64

go install uciengine/zurimain.go

copy zurimain.exe %EASYCHESS_PATH%\resources\server\bin\gochess.exe
copy zurimain %EASYCHESS_PATH%\resources\server\bin\gochess
