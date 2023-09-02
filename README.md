<h1>Simple server convert longitude&latitude to timezone name</h1>

> **Note**
>
> It's designed as a debugger tool for package
> [tzf](https://github.com/ringsaturn/tzf), not production ready.

- [Quick Start](#quick-start)
  - [Install](#install)
    - [Install via `go install`](#install-via-go-install)
    - [Download from release page](#download-from-release-page)
  - [Usage](#usage)
- [Web Pages](#web-pages)
  - [All supported timezone names](#all-supported-timezone-names)
- [HTTP API](#http-api)
  - [Lookup Location's timezone](#lookup-locations-timezone)
  - [Lookup Location's timezones](#lookup-locations-timezones)
  - [All supported timezone names](#all-supported-timezone-names-1)
- [Redis Protocol Commands](#redis-protocol-commands)
  - [`redis-cli`](#redis-cli)
  - [`redis-py`](#redis-py)

## Quick Start

### Install

#### Install via `go install`

```bash
go install github.com/ringsaturn/tzf-server@latest
```

#### Download from release page

Please visit <https://github.com/ringsaturn/tzf-server/releases> to get latest
release.

### Usage

```console
Usage of tzf-server:
  -debug
        Enable debug mode
  -disable-print-route
        Disable Print Route
  -http-addr string
        HTTP Host&Port (default "localhost:8080")
  -path string
        custom data
  -prometheus-enable-go-coll
        Enable Go Collector (default true)
  -prometheus-host-port string
        Prometheus Host&Port (default "localhost:8090")
  -prometheus-path string
        Prometheus Path (default "/hertz")
  -redis-addr string
        Redis Server Host&Port (default "localhost:6380")
  -type int
        which finder to use Polygon(0) or Fuzzy(1)
```

For example, start
[DefaultFinder](https://pkg.go.dev/github.com/ringsaturn/tzf#DefaultFinder)
server:

```bash
tzf-server
```

Or start [FuzzyFinder](https://pkg.go.dev/github.com/ringsaturn/tzf#FuzzyFinder)
based server:

```bash
tzf-server -type 1
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
  "timezone": "Asia/Shanghai",
  "abbreviation": "CST",
  "offset": 28800
}
```

### Lookup Location's timezones

```bash
curl "http://localhost:8080/api/v1/tzs?lng=87.6168&lat=43.8254"
```

Output:

```json
{
  "timezones": [
    {
      "timezone": "Asia/Shanghai",
      "abbreviation": "CST",
      "offset": 28800
    },
    {
      "timezone": "Asia/Urumqi",
      "abbreviation": "+06",
      "offset": 21600
    }
  ]
}
```

### All supported timezone names

```bash
curl "http://localhost:8080/api/v1/tzs/all"
```

Output:

```jsonc
{
  "timezones": [
    {
      "timezone": "Africa/Abidjan",
      "abbreviation": "GMT",
      "offset": 0
    },
    // ...
    {
      "timezone": "Etc/GMT+12",
      "abbreviation": "-12",
      "offset": -43200
    }
  ]
}
```

## Redis Protocol Commands

### `redis-cli`

```console
$ redis-cli -p 6380
127.0.0.1:6380> GET_TZ 116.3883 39.9289
Asia/Shanghai
127.0.0.1:6380> GET_TZS 87.4160 44.0400
1) "Asia/Shanghai"
2) "Asia/Urumqi"
```

### `redis-py`

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
