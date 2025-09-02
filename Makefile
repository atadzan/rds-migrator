build:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o migrator  ./migrator.go