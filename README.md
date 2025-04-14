# pickupPoint
______
# Table of Contents
* [About](#about)
* [Clone repo](#clone-repo)
* [Technologies](#technologies)
* [Quick start](#quick-start)
* [Additions](#additions)

## About
This is server part made for managing pickup points

## Clone repo
You can clone this repo to observe source code and run tests localy using this command (golang version 1.24.0 or higher):
```
git clone https://github.com/Razzle131/pickupPoint.git
```

## Technologies
Project is created with:
* Golang version: 1.24.0
* oapi-codegen
* squirrel
* jwt
* testcontainers-go (integration testing with postgres)

## Quick start
### Build from source
* Ensure that Go is installed on your machine and it`s version is equal or higther than 1.24.0
```
go version
```
* Clone repo
```
git clone https://github.com/Razzle131/ComputerClub.git
```
* Install dependecies
```
go mod tidy
```
* Run unit tests:
```
go test ./... -cover -short
```
or
```
make unit-test
```
* Run integration test:
```
go test -run Integration
```
or
```
make integr-test
```
* Configure enviroment variables
* Run programm
```
go run cmd/main.go
```
I recommend using [pplog](https://github.com/michurin/human-readable-json-logging), config file is located in the root dir of the project
```
pplog go run cmd/main.go
```
* You can use Makefile shortcut commands for running/testing/building docker image/starting docker container. For more info inspect Makefile
______
### Docker
This app could be builded and executed using docker:
* Clone repo
```
git clone https://github.com/Razzle131/ComputerClub.git
```
* use this command to build docker image  
```
docker build -t IMAGE-NAME .
```
replace the IMAGE-NAME with the name you want

or (configure docker image name in Makefile)
```
make docker-build
```
* run docker container:
```
docker run --rm -p 8080:8080 --env-file .env IMAGE-NAME
```
or (configure docker image name in Makefile)
```
make docker-run
```
______
### Docker Compose
* Clone repo:
```
git clone https://github.com/Razzle131/ComputerClub.git
```
* Configure enviroment variables in the docker-compose file
* Run:
```
docker compose up
```
______
## Additions
### Testing
* Business logic was covered for 80%-90% by unit tests.
* Unit tests was developed for cache repositories and services
______
### Get pvz response
* For some reason, oapi-codegen did not generate response struct for get pvz request, so I placed it in separate file in api directory
* I made "get pvz" response to hide pickup points without receptions
* I made "get pvz" response to show receptions without products