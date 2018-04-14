# WGPlaner Server

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/492b2e1e8ce9415fa826016952baaa15)](https://www.codacy.com/app/archer96/wg_planer_server?utm_source=github.com&utm_medium=referral&utm_content=WGPlaner/wg_planer_server&utm_campaign=badger)
[![Go Report Card](https://goreportcard.com/badge/github.com/wgplaner/wg_planer_server)](https://goreportcard.com/report/github.com/wgplaner/wg_planer_server)
[![Build Status AppVeyor](https://ci.appveyor.com/api/projects/status/ok5rq84eh6sx8lxd/branch/master?svg=true)](https://ci.appveyor.com/project/archer96/wg-planer-server/branch/master)
[![Build Status Travis](https://travis-ci.org/WGPlaner/wg_planer_server.svg?branch=master)](https://travis-ci.org/WGPlaner/wg_planer_server)
[![codecov](https://codecov.io/gh/WGPlaner/wg_planer_server/branch/master/graph/badge.svg)](https://codecov.io/gh/WGPlaner/wg_planer_server)
[![GoDoc](https://godoc.org/github.com/WGPlaner/wg_planer_server?status.svg)](https://godoc.org/github.com/WGPlaner/wg_planer_server)

## Setup
To generate the go API, install `go-swagger`.

```bash
go get -u github.com/go-swagger/go-swagger/cmd/swagger
go install github.com/go-swagger/go-swagger/cmd/swagger
```

Then run:

```bash
rm -rf restapi
swagger generate server -t . -f swagger.yml --exclude-main --skip-models -P models.User -A wgplaner
```

To build `wg_planer_server` run:

```bash
go build -v -o "build/wg_planer_server" ./cmd/wgplaner-api/wgplaner-api.go
```

### Create Android Library
First download `swagger-codegen`:

```bash
wget http://central.maven.org/maven2/io/swagger/swagger-codegen-cli/2.3.1/swagger-codegen-cli-2.3.1.jar -O swagger-codegen-cli.jar
```

To create the Java Android Library run:

```bash
java -jar swagger-codegen-cli.jar generate -i swagger.yml -l java --library=okhttp-gson -o build/android_client
```

### API Documentation
Available here: https://doc.wgplaner.ameyering.de/docs

Using `go-swagger`, simply run:
```bash
swagger serve swagger.yml
```

