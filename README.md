# api-load-test

> simple and elegant load testing tool.

## Development

```
git clone https://github.com/Marvin9/api-load-test
cd api-load-test
make build
./bin/api-load-test --endpoint "YOUR_API_ENDPOINT" --method "GET" --rate 10 --until 10
```

- To list available flags

```
./bin/api-load-test --help

Usage:
  loadtest [flags]

Flags:
  -e, --endpoint string   target endpoint. eg: http://locahost:8000/
  -h, --help              help for loadtest
  -m, --method string     method of target endpoint [GET/POST/PUT...] (default "GET")
  -r, --rate int          load of requets per second (default 100)
  -u, --until int         duration of load in seconds (default 10)
```

## Limitations

`--rate` of requests will be dependent on the available file descriptors `ulimit -n`, processes `ulimit -u` and the round trip time of request. In short number of `rate` should not exceed certain limit.