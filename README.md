# WGPlaner Server

## Setup
To create the go API, install `swagger` and run:

```bash
swagger generate server -t gen -f ./swagger/swagger.yml --exclude-main -A wgplaner
```