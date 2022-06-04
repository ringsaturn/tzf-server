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
curl "http://localhost:8080/tz?lng=-0.1276&lat=51.5073"
```

Check timezone info page with link to <http://geojson.io> to view polygon:

```bash
open "http://localhost:8080/info?lng=-0.1276&lat=51.5073"
```
