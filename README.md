# nginx2influxdb

Tool to import/stream Nginx logs directly into InfluxDB via wire protocol

## Usage

```
Usage of ./bin/nginx2influxdb:
  -d string
      InfluxDB database name
  -h string
      InfluxDB server url
  -p int
      Number of seconds between writes (default 5)
  -s  Stream
```

Streaming logs:

```
tail -f /var/log/nginx/access.log | nginx2influxdb -s -h http://myip:8086 -d mydb
```

Or just plain import:

```
cat /var/log/nginx/access.log | nginx2influxdb -h http://myip:8086 -d mydb
```

## Install

Clone repo and build the tool:

```
make build
```

Build for multiple platforms:

```
make all
```

## License

MIT