<h1>Simple server convert longitude&latitude to timezone name</h1>

<img width="1920" src="https://github.com/ringsaturn/tzf-server/assets/13536789/9a820db1-1de9-49e6-9047-4e72c2fe41a8">

> [!NOTE]
>
> It's designed as a debugger tool for package
> [tzf](https://github.com/ringsaturn/tzf), not production ready.

- [Quick Start](#quick-start)
  - [Install](#install)
    - [Install via `go install`](#install-via-go-install)
    - [Download from release page](#download-from-release-page)
    - [Install from Docker Hub](#install-from-docker-hub)
  - [Usage](#usage)
- [Web Pages](#web-pages)
  - [All supported timezone names](#all-supported-timezone-names)
  - [\[Experiment\] Clickable debugger](#experiment-clickable-debugger)
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

#### Install from Docker Hub

```bash
docker pull ringsaturn/tzf-server
```

### Usage

```console
Usage of tzf-server:
  -disable-print-route
        Disable Print Route
  -hertz-prometheus-host-port string
        Hertz Prometheus Host&Port (default "0.0.0.0:8090")
  -hertz-prometheus-path string
        Hertz Prometheus Path (default "/hertz")
  -http-addr string
        HTTP Host&Port (default "0.0.0.0:8080")
  -path string
        custom data
  -prometheus-enable-go-coll
        Enable Go Collector (default true)
  -prometheus-host-port string
        Prometheus Host&Port (default "0.0.0.0:2112")
  -prometheus-path string
        Prometheus Path (default "/metrics")
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

### [Experiment] Clickable debugger

[`http://localhost:8080/web/click`](http://localhost:8080/web/click)

I have little knowledge about frontend development so it's just a experiment,
most codes are written by ChatGPT 3.5. You can access the prompts from
[gist](https://gist.github.com/ringsaturn/12b5509d80f69e7bca13da8745f7ce15).

## HTTP API

A swagger UI can be found at
[`http://localhost:8080/swagger/index.html`](http://localhost:8080/swagger/index.html).

### Lookup Location's timezone

```bash
curl "http://localhost:8080/api/v1/tz?longitude=116.3883&latitude=39.9289"
```

Output:

```json
{
  "timezone": "Asia/Shanghai",
  "abbreviation": "CST",
  "offset": "28800s"
}
```

### Lookup Location's timezones

```bash
curl "http://localhost:8080/api/v1/tzs?longitude=87.6168&latitude=43.8254"
```

Output:

```json
{
  "timezones": [
    {
      "timezone": "Asia/Shanghai",
      "abbreviation": "CST",
      "offset": "28800s"
    },
    {
      "timezone": "Asia/Urumqi",
      "abbreviation": "+06",
      "offset": "21600s"
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
      "offset": "0s"
    },
    // ...
    {
      "timezone": "Etc/GMT+12",
      "abbreviation": "-12",
      "offset": "-43200s"
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
