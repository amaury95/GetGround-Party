# GetGround-Party

This repository is the solution for the practical task if GetGround company for a Golang Backend Developer

## Requirements

As the application is written in golang, system requirements are quite straight forward:

### System

- OS: Linux Based, Windows 7+, Mac OS X
- RAM: 128MB
- HDD: 1GB

### Software

- [Golang]() (v1.10 - latest)
- [MySQL Server]() (v5.7)
- [Docker]() (v1.4 - latest) [Optional]
- [VSCode]() (latest) [Optional]

## Quick setup

1 - Download the project respository into your host machine

```sh
git clone https://github.com/amaury95/GetGround-Party
```

2 - Create a database instance using `docker-compose`:

```sh
docker-compose -f deployments/mysql.docker-compose.yaml up -d
```

\*you can use `docker-compose -f deployments/mysql.docker-compose.yaml stop|down` to stop the testing instance.

3 - Open browser on `http://localhost:8080` and enter the testing credentials (username: `root`, password: `example`)

4 - Create a database called `party`.

5 - Run the application with the default config:

```sh
go run cmd/server/main.go
```

\*you can run `go run cmd/server/main.go --help` and configure the database connection

6 - API is exposed on `http://localhost:3033`. Now you can use any tool such as [Postman]() or [REST Client]() to manipulate the API.

## Configuration options

You can run `server.exe --help` to pop up the server configuration.

```sh
usage: party [-h|--help] [-p|--port <integer>] [-n|--username "<value>"]
             [-k|--password "<value>"] [-u|--url "<value>"] [-d|--database
             "<value>"]

             Party is the webserver to manage GetGround invitations and guests.

Arguments:

  -h  --help      Print help information
  -p  --port      server port to listen for requests. Default: 3033
  -n  --username  mysql connection username. Default: root
  -k  --password  mysql connection password. Default: example
  -u  --url       mysql connection server url. Default: 127.0.0.1:3306
  -d  --database  mysql connection database name. Default: party
```

For security reasons sensible credentials must be stored as environment variables, should not be included inside of the project repository and should not be promped directly to the command line. We can make use of the system variables by using them directly from the command line:

```sh
go run cmd/server/main.go --username=$MYSQL_ROOT --password=$MYSQL_ROOT_PASSWORD
```

## API Docs

Api Docs can be described in [OpenAPI Specs](). Since there are not the objective of this assignment, I toke the liberty of provide a detailed `routes.rest` file that can be used with the [REST Client]() extension of Visual Studio Code.

## Building solution

We can compile the source code of this project for any golang compatible platform. This are two examples of compiling configurations:

```sh
GOOS=linux GOARCH=amd64 go build -o bin/main-linux-amd64 ./cmd/server
```

```sh
GOOS=windows GOARCH=amd64 go build -o bin/main-windows-amd64.exe ./cmd/server
```

## Testing

Testing was carried on using state of the art testing environment tools as:

 - [Ginkgo](https://github.com/onsi/ginkgo)
 - [Gomega](https://github.com/onsi/gomega)
 - [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)
 - [httpexpect](https://github.com/gavv/httpexpect)

We can test the application by going into the following steps:

```sh
cd ./tests
go test
```

we can also make use of the `ginkgo` binary by running:
```sh
ginkgo -r
```
