build: 
	env GOOS=linux go build -ldflags="-s -w" -o bin/main main.go

clean:
	rm -rf ./bin

# start build:
# 	sls offline

offline: clean build
	sls offline --useDocker start --host 0.0.0.0

deploy_dev: clean build
	serverless deploy