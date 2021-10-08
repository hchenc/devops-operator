IMG = 364554757/devops:latest

docker-build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o ./bin/controller main.go
	docker build . -t ${IMG} -f deploy/Dockerfile