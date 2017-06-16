# Smoke tests for graphite

## How to install this project
```
go get github.com/SpringerPE/graphite-smoke-tests
go get github.com/onsi/ginkgo/ginkgo
```

## How to run these tests

Copy and modify the example configuration file under `./examples/test.config.json` or 
create a new one with the following contents:

```
{
    "api": "<graphite_api_endpoint>", #api.graphite.example
    "apiPort": 80 # (optional, default: 80)
    "host": "<graphite_host>", #graphite.host.example
    "port": 2003 # (optional, default: 80)
    "tcpEnabled": true,
    "udpEnabled": true
}
```

export the `SMOKE_TEST_CONFIG` variable with the location of the smoke tests:

```
export SMOKE_TEST_CONFIG=./smoke_tests_config.json
```

run the tests:

```
ginkgo -r
```
