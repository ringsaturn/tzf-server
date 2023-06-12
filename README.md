# Simple server convert longitude&latitude to timezone name

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

## HTTP API

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

## Redis Protocol Commands

### `redis-cli`

```
$ redis-cli -p 6380
127.0.0.1:6380> GET_TZ 116.3883 39.9289
Asia/Shanghai
127.0.0.1:6380> GET_TZS 87.4160 44.0400
1) "Asia/Shanghai"
2) "Asia/Urumqi"
```

## `redis-py`

```python
>>> from redis import Redis
>>> rc = Redis.from_url("redis://localhost:6380")
>>> rc.ping()
True
>>> rc.execute_command("get_tz", 116.3883, 39.9289).decode()
'Asia/Shanghai'
>>> rc.execute_command("get_tzs", 87.4160, 44.0400)
[b'Asia/Shanghai', b'Asia/Urumqi']
```
