# Simple HTTP API convert longitude&latitude to timezone name

> **Note**: It's designed as a debugger tool for package
> [tzf](https://github.com/ringsaturn/tzf), not production ready.

```bash
go install github.com/ringsaturn/tzf-server@latest

# check args
tzf-server --help

# start general server
./tzf-server

# start FuzzyFinder based server
./tzf-server -type 1
```

## Web Pages

### All supported timezone names

[`http://localhost:8080/web/tzs/all`](http://localhost:8080/web/tzs/all)

## API

### Lookup Location's timezone

```bash
curl "http://localhost:8080/api/v1/tz?lng=116.3883&lat=39.9289"
```

Output:

```json
{
  "timezone": "Asia/Shanghai"
}
```

### Lookup Location's timezones

```bash
curl "http://localhost:8080/api/v1/tzs?lng=87.6168&lat=43.8254"
```

Output:

```json
{
  "timezones": ["Asia/Shanghai", "Asia/Urumqi"]
}
```

### All supported timezone names

```bash
curl "http://localhost:8080/api/v1/tzs"
```

Output:

```json
{
  "timezones": [
    "Africa/Abidjan",
    // ...
    "Pacific/Tongatapu"
  ]
}
```
