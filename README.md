# Multifinance-Apps
This repository contains multifinance app for test and learning purposes.

## Preqrequisites
The Multifinance Apps system requires Go 1.19 or higher and Docker installed on the local machine in order to run the binary.

Docker
You need to have docker installed in your machine. Follow this step if you don't have docker on your machine :

+ Download the Docker CE (Community Edition) package from the Docker website (https://www.docker.com/products/docker-desktop).
+ Install the package by following the instructions provided during the installation process.
+ Once the installation is complete, verify that Docker has been installed correctly by running the following command in your terminal: "docker run hello-world".


Go Programming Language
You need to have golang 1.19 installed in your machine. Follow this step if you don't have golang 1.19 on your machine :
+ Download the Go 1.19 binary package from the official Go website (https://golang.org/dl/).
+ Install the package by following the instructions provided during the installation process.
+ Once the installation is complete, verify that Go has been installed correctly by running the following command in your terminal: "go version".

## How to run locally
Once you have all the prerequisites properly installed, you can start by cloning this repository.
+ Clone the repo to your local folder
```
git clone https://github.com/jaysm12/multifinance-apps.git
```
+ Navigate to cloned repo
```
cd multifinance-apps
```

Note: These instructions assume that you have Git installed on your machine. If you don't have Git installed, you can follow the instructions on the Git website to install it.

## Docker Setup:
To run the Multifinance Apps system binary correctly, it is necessary to connect it with the related dependencies. This can be done simply by executing the following command:
```
make deps-init
```
The deps-init command will perform the following actions:

+ Build RabbitMQ and verify that it is running
+ Build Mysql and verify that it is running
  
To stop the dependencies, run :
```
make deps-tear
```
### Running Binary:
Once you have cloned the repository and set up the docker dependencies, you can run the binary using either of the following methods:

+ Run vendor to download package dependencies
```
go mod vendor
```
+ Change the config in /config/config.yaml accordingly

And run using :
```
make run-local
```
or
```
go run ./cmd/multifinance-apps/main.go
```
or
```
go build ./cmd/multifinance-apps/
./multifinance-apps
```

to run go test you can use :
```
make test
```

## Postman Collection
You can check out the Endpoint by importing the postman collection on `./Multifinance-apps.postman_collection`
