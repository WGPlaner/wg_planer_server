# WGPlaner Server

[![Build status](https://ci.appveyor.com/api/projects/status/ok5rq84eh6sx8lxd/branch/master?svg=true)](https://ci.appveyor.com/project/archer96/wg-planer-server/branch/master)

## Setup
To create the go API, install `go-swagger` and run:

```bash
rm -rf gen # Delete old generated files
mkdir gen  # Create gen directory
swagger generate server -t gen -f ./swagger.yml --exclude-main -A wgplaner
```
