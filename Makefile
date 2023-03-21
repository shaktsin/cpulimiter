build:
	env GOOS=linux GOARCH=amd64 go build -o cpulim  main.go limiter.go 