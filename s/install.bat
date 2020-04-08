set GOBIN=C:\gomodules\modules\gochess

set GOOS=linux
set GOARCH=amd64

go install uciengine/gochess.go

set GOOS=windows
set GOARCH=amd64

go install uciengine/gochess.go

copy gochess.exe %EASYCHESS_PATH%\resources\server\bin
copy gochess %EASYCHESS_PATH%\resources\server\bin
