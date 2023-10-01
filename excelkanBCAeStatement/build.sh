GOOS=linux GOARCH=amd64 go build -o excelkanBCAeStatement-Linux-amd64 -v
GOOS=darwin GOARCH=amd64 go build -o excelkanBCAeStatement-Mac-amd64 -v
GOOS=darwin GOARCH=amd64 go build -o excelkanBCAeStatement-Mac-arm64 -v
GOOS=windows GOARCH=amd64 go build -o excelkanBCAeStatement.exe -v
