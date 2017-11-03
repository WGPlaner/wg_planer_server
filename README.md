# WGPlaner Server

[![Build Status AppVeyor](https://ci.appveyor.com/api/projects/status/ok5rq84eh6sx8lxd/branch/master?svg=true)](https://ci.appveyor.com/project/archer96/wg-planer-server/branch/master)
[![Build Status Travis](https://travis-ci.org/WGPlaner/wg_planer_server.svg?branch=master)](https://travis-ci.org/WGPlaner/wg_planer_server)

## Setup
To create the go API, install `go-swagger` and run:

```bash
rm -rf gen # Delete old generated files
mkdir gen  # Create gen directory
swagger generate server -t gen -f swagger.yml --exclude-main -A wgplaner
```

### Client Library
To create the Java Android Library, download `swagger-codegen`.

```bash
wget -o swagger-codegen-cli.jar https://oss.sonatype.org/content/repositories/releases/io/swagger/swagger-codegen-cli/2.2.3/swagger-codegen-cli-2.2.3.jar
java -jar swagger-codegen-cli.jar generate -i swagger.yml -l java --library=okhttp-gson -o android_client
```

### API Documentation
To create the API documentation, download `swagger-codegen`.

```bash
wget -o swagger-codegen-cli.jar https://oss.sonatype.org/content/repositories/releases/io/swagger/swagger-codegen-cli/2.2.3/swagger-codegen-cli-2.2.3.jar
java -jar swagger-codegen-cli.jar generate -i swagger.yml -l html2 -o html_api_doc
```

