# WGPlaner Server

[![Build Status AppVeyor](https://ci.appveyor.com/api/projects/status/ok5rq84eh6sx8lxd/branch/master?svg=true)](https://ci.appveyor.com/project/archer96/wg-planer-server/branch/master)
[![Build Status Travis](https://travis-ci.org/WGPlaner/wg_planer_server.svg?branch=master)](https://travis-ci.org/WGPlaner/wg_planer_server)

## Setup
To generate the go API, install `go-swagger`.

```bash
go get -u github.com/go-swagger/go-swagger/cmd/swagger
go install github.com/go-swagger/go-swagger/cmd/swagger
```

Then run:

```bash
rm -rf gen # Delete old generated files
mkdir gen  # Create gen directory
swagger generate server -t . -f swagger.yml --exclude-main --skip-models -A wgplaner
```

To build `wg_planer_server` run:

```bash
go build -v -o "build/wg_planer_server" ./cmd/wgplaner-api/wgplaner-api.go
```

### Client Library
To create the Java Android Library, download `swagger-codegen`.

```bash
wget -O swagger-codegen-cli.jar https://oss.sonatype.org/content/repositories/releases/io/swagger/swagger-codegen-cli/2.2.3/swagger-codegen-cli-2.2.3.jar
java -jar swagger-codegen-cli.jar generate -i swagger.yml -l java --library=okhttp-gson -o build/android_client
```

### API Documentation
To create the API documentation, download `swagger-codegen`.

```bash
wget -O swagger-codegen-cli.jar https://oss.sonatype.org/content/repositories/releases/io/swagger/swagger-codegen-cli/2.2.3/swagger-codegen-cli-2.2.3.jar
java -jar swagger-codegen-cli.jar generate -i swagger.yml -l html2 -o build/api_doc_html
```

