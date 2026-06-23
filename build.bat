:build_linux_amd64
echo ±‡“ÎLinux∞Ê±æ64Œª
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -v -a -o shopadmin .
@pause