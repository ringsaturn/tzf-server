# Simple HTTP API convert longitude&latitude to timezone name

```bash
go install github.com/ringsaturn/tzf-server@latest

# check args
tzf-server --help

# start server
tzf-server
```

Example links:

- Query by location [`http://localhost:8080/tz?lng=-0.1276&lat=51.5073`](http://localhost:8080/tz?lng=-0.1276&lat=51.5073)
- Query by offset [`http://localhost:8080/tz/offset?offset=0`](http://localhost:8080/tz/offset?offset=0)
- Query by name [`http://localhost:8080/info?name=Etc/GMT-12`](http://localhost:8080/info?name=Etc/GMT-12)
