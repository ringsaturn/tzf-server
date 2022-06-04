# Simple HTTP API convert longitude&latitude to timezone name

```bash
go install github.com/ringsaturn/tzf-server@latest

# check args
tzf-server --help

# start server
tzf-server
```

Check timezone name:

```bash
curl "http://localhost:8080/tz?lng=139.8753&lat=36.2330"
# Asia/Tokyo
```

Check timezone info page with link to <http://geojson.io> to view polygon:

```bash
curl "http://localhost:8080/info?lng=139.8753&lat=36.2330"
```
