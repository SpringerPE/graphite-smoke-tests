# Smoke tests for graphite

## How to install this project
```
go get github.com/SpringerPE/graphite-smoke-tests
```

## How to run these tests

Create a configuration file, let's suppose `./smoke_tests_config.json` with the following contents:

```
{
    "api": "<graphite_api_endpoint>", #api.graphite.example
    "api_port": 80 # (optional, default: 80)
    "host": "<graphite_host>", #graphite.host.example
    "port": 2003 # (optional, default: 80)
    "enable_tcp": true,
    "enable_udp": true
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