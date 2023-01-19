build: 
	env GOOS=linux go build -ldflags="-s -w" -o bin/main main.go

clean:
	rm -rf ./bin

start build:
	sls offline 
# --useDocker start --host 0.0.0.0
deploy_dev: clean build   # this means run build command first before deploy_prod
	serverless deploy 