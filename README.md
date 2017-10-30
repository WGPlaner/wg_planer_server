# WGPlaner Server

[![Build status](https://ci.appveyor.com/api/projects/status/ok5rq84eh6sx8lxd/branch/master?svg=true)](https://ci.appveyor.com/project/archer96/wg-planer-server/branch/master)

## Setup
To create the go API, install `go-swagger` and run:

```bash
swagger generate server -t gen -f ./swagger/swagger.yml --exclude-main -A wgplaner
```
