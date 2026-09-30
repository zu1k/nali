<h1 align="center">
  <br>Nali<br>
</h1>

<h4 align="center">An offline tool for querying IP geographic information and CDN provider.</h4>

<p align="center">
  <a href="https://github.com/zu1k/nali/actions">
    <img src="https://img.shields.io/github/actions/workflow/status/zu1k/nali/go.yml?branch=master&style=flat-square" alt="Github Actions">
  </a>
  <a href="https://goreportcard.com/report/github.com/zu1k/nali">
    <img src="https://goreportcard.com/badge/github.com/zu1k/nali?style=flat-square">
  </a>
  <a href="https://github.com/zu1k/nali/releases">
    <img src="https://img.shields.io/github/release/zu1k/nali/all.svg?style=flat-square">
  </a>
  <a href="https://github.com/zu1k/nali/releases">
    <img src="https://img.shields.io/github/downloads/zu1k/nali/total?style=flat-square">
  </a>
</p>

#### [中文文档](https://github.com/zu1k/nali/blob/master/README.md)

## Feature

- Multi database support
  - Chunzhen qqip database
  - ZX ipv6 database
  - ip2region IPv4 / IPv6 database
  - IPinfo Lite database
  - Geoip2 city database
  - IPIP free database
  - DB-IP database
  - IP2Location DB3 LITE database
- CDN provider query
- Pipeline support
- Interactive query
- Both ipv4 and ipv6 supported
- Multilingual support
- Offline query
- Full platform support
- Color print

## Install

### Install from source

Nali requires Go >= 1.26. You can build it from source:

```sh
$ go install github.com/zu1k/nali@latest
```

### Install pre-build binary

Pre-built binaries are available here: [release](https://github.com/zu1k/nali/releases)

Download the binary compatible with your platform, unpack and copy to the directory in path

### Arch Linux

We have published 3 packages in Aur:

- `nali-go`: release version, compile when installing
- `nali-go-bin`: release version, pre-compiled binary
- `nali-go-git`: the latest master branch version, compile when installing

## Usage

### Query a simple IP address

```
$ nali 1.2.3.4
1.2.3.4 [澳大利亚 APNIC Debogon-prefix网络]
```

#### or use `pipe`

```
$ echo IP 6.6.6.6 | nali
IP 6.6.6.6 [美国 亚利桑那州华楚卡堡市美国国防部网络中心]
```

### Query multiple IP addresses

```
$ nali 1.2.3.4 4.3.2.1 123.23.3.0
1.2.3.4 [澳大利亚 APNIC Debogon-prefix网络]
4.3.2.1 [美国 新泽西州纽瓦克市Level3Communications]
123.23.3.0 [越南 越南邮电集团公司]
```

### Interactive query

use `exit` or  `quit` to quit

```
$ nali
123.23.23.23
123.23.23.23 [越南 越南邮电集团公司]
1.0.0.1
1.0.0.1 [美国 APNIC&CloudFlare公共DNS服务器]
8.8.8.8
8.8.8.8 [美国 加利福尼亚州圣克拉拉县山景市谷歌公司DNS服务器]
quit
```

### Use with `dig`

```
$ dig nali.zu1k.com +short | nali
104.28.2.115 [美国 CloudFlare公司CDN节点]
104.28.3.115 [美国 CloudFlare公司CDN节点]
172.67.135.48 [美国 CloudFlare节点]
```

### Use with `nslookup`

```
$ nslookup nali.zu1k.com 8.8.8.8 | nali
Server:         8.8.8.8 [美国 加利福尼亚州圣克拉拉县山景市谷歌公司DNS服务器]
Address:        8.8.8.8 [美国 加利福尼亚州圣克拉拉县山景市谷歌公司DNS服务器]#53

Non-authoritative answer:
Name:   nali.zu1k.com
Address: 104.28.3.115 [美国 CloudFlare公司CDN节点]
Name:   nali.zu1k.com
Address: 104.28.2.115 [美国 CloudFlare公司CDN节点]
Name:   nali.zu1k.com
Address: 172.67.135.48 [美国 CloudFlare节点]
```

### Use with any other program

Because nali can read the contents of the `stdin` pipeline, it can be used with any program.

```
bash abc.sh | nali
```

Nali will insert IP information after IP address.

### IPv6 support

Use like IPv4

```
$ nslookup google.com | nali
Server:         127.0.0.53 [局域网 IP]
Address:        127.0.0.53 [局域网 IP]#53

Non-authoritative answer:
Name:   google.com
Address: 216.58.211.110 [美国 Google全球边缘网络]
Name:   google.com
Address: 2a00:1450:400e:809::200e [荷兰Amsterdam Google Inc. 服务器网段]
```

### Query CDN provider

```
$ nslookup www.gov.cn | nali
Server:         127.0.0.53 [局域网 IP]
Address:        127.0.0.53 [局域网 IP]#53

Non-authoritative answer:
www.gov.cn      canonical name = www.gov.cn.bsgslb.cn [白山云 CDN].
www.gov.cn.bsgslb.cn [白山云 CDN]       canonical name = zgovweb.v.bsgslb.cn [白山云 CDN].
Name:   zgovweb.v.bsgslb.cn [白山云 CDN]
Address: 103.104.170.25 [新加坡 ]
Name:   zgovweb.v.bsgslb.cn [白山云 CDN]
Address: 2001:428:6402:21b::5 [美国Louisiana州Monroe Qwest Communications Company, LLC (CenturyLink)]
Name:   zgovweb.v.bsgslb.cn [白山云 CDN]
Address: 2001:428:6402:21b::6 [美国Louisiana州Monroe Qwest Communications Company, LLC (CenturyLink)]
```

## Interface

After nali runs for the first time, a configuration file `config.yaml` will be generated in the config directory (use `nali info` to extract info), the configuration file defines the database information.

A database is defined as follows:

```yaml
- name: geoip
  name-alias:
  - geolite
  - geolite2
  format: mmdb
  file: GeoLite2-City.mmdb
  languages:
  - ALL
  types:
  - IPv4
  - IPv6
```

Configuration files written by older versions keep working: databases missing from the config fall back to the built-in defaults.

### Databases

| Database | Name / aliases | Query type | Default file | Auto download |
| --- | --- | --- | --- | --- |
| Chunzhen IPv4 | `qqwry`, `chunzhen` | IPv4 (default) | `qqwry.dat` | ✓ |
| ZX IPv6 | `zxipv6wry`, `zxipv6`, `zx` | IPv6 (default) | `zxipv6wry.db` | ✓ |
| ip2region IPv4 | `ip2region`, `i2r` | IPv4 | `ip2region.xdb` | ✓ |
| ip2region IPv6 | `ip2region-ipv6`, `i2r-ipv6` | IPv6 | `ip2region_v6.xdb` | ✓ |
| IPinfo Lite | `ipinfo` | IPv4, IPv6 | `ipinfo_lite.mmdb` | ✓ |
| GeoIP2 / GeoLite2 City | `geoip`, `geoip2`, `geolite`, `geolite2` | IPv4, IPv6 | `GeoLite2-City.mmdb` | manual |
| DB-IP | `dbip`, `db-ip` | IPv4, IPv6 | `dbip.mmdb` | manual |
| IPIP | `ipip` | IPv4, IPv6 | `ipipfree.ipdb` | manual |
| IP2Location DB3 LITE | `ip2location` | IPv4, IPv6 | `IP2LOCATION-LITE-DB3.IPV6.BIN` | manual |
| CDN | `cdn` | domain (default) | `cdn.yml` | ✓ |

Databases with automatic download are fetched on first use when the file is missing. For the others, put the file in the data directory shown by `nali info`, or set an absolute `file` path in the config.

### Use IPinfo Lite

[IPinfo Lite](https://ipinfo.io/lite) provides free country and ASN data for both IPv4 and IPv6. It is downloaded on first use:

```
$ NALI_DB_IP4=ipinfo NALI_DB_IP6=ipinfo nali 8.8.8.8 2001:4860:4860::8888
8.8.8.8 [United States AS15169 Google LLC]
2001:4860:4860::8888 [United States AS15169 Google LLC]

$ NALI_DB_IP4=ipinfo nali --json 8.8.8.8
{"type":0,"ip":"8.8.8.8","text":"United States AS15169 Google LLC","source":"ipinfo","info":{"network":"8.8.8.0/24","country":"United States","country_code":"US","continent":"North America","continent_code":"NA","asn":"AS15169","as_name":"Google LLC","as_domain":"google.com"}}
```

To use it permanently, set `selected.ipv4` and `selected.ipv6` to `ipinfo` in the config file.

Text output includes country, ASN and AS name. JSON includes `network`, `country`, `country_code`, `continent`, `continent_code`, `asn`, `as_name` and `as_domain`. When a record omits `network`, the matched CIDR is used. Names are the database's English values.

Databases with `format: mmdb` are detected as IPinfo from their metadata, so an `ipinfo_lite.mmdb` downloaded from your [IPinfo account](https://ipinfo.io/dashboard/downloads) works as well. The default download URL is the community mirror [NetworkCats/IPinfoLite-Download](https://github.com/NetworkCats/IPinfoLite-Download); change `download-urls` in the config to use another source:

```yaml
- name: ipinfo
  format: mmdb
  file: ipinfo_lite.mmdb
  download-urls:
  - https://github.com/NetworkCats/IPinfoLite-Download/releases/latest/download/ipinfo_lite.mmdb
  languages: [en]
  types: [IPv4, IPv6]
```

IPinfo Lite data is licensed under [CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/); please attribute IPinfo when you publish data derived from it.

### Use ip2region IPv6

ip2region ships separate IPv4 (`ip2region`) and IPv6 (`ip2region-ipv6`) database files, both downloaded on first use:

```
$ NALI_DB_IP4=ip2region NALI_DB_IP6=ip2region-ipv6 nali 223.5.5.5 240e::1
223.5.5.5 [中国|浙江省|杭州市|阿里|CN]
240e::1 [中国|北京市|北京市|电信|CN]
```

An existing `ip2region.xdb` from older versions keeps working, and the outdated download URL in old config files is migrated to `ip2region_v4.xdb` automatically.

### Help

```
$ nali --help
Usage:
  nali [flags]
  nali [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  info        get the necessary information of nali
  update      update ip databases and cdn, update nali to latest version if -v

Flags:
      --gbk       Use GBK decoder
  -h, --help      help for nali
  -j, --json      Output in JSON format
  -v, --version   version for nali

Use "nali [command] --help" for more information about a command.
```

### Update database

Without arguments, `nali update` updates the Chunzhen, ZX IPv6, ip2region (IPv4) and CDN databases. The larger optional databases `ip2region-ipv6` (~37 MB) and `ipinfo` (~24 MB) are only updated when they are selected (config `selected` or the `NALI_DB_*` environment variables) or already downloaded.

```
$ nali update
2026/09/30 09:51:48 正在下载最新 qqwry 数据库...
2026/09/30 09:51:49 qqwry 数据库下载成功: qqwry.dat
2026/09/30 09:51:49 正在下载最新 cdn 数据库...
2026/09/30 09:51:50 cdn 数据库下载成功: cdn.yml
...
```

Update specific databases (comma separated, aliases allowed); these are always updated:

```
$ nali update --db qqwry,cdn,ipinfo
```

Downloads are validated before they replace the local file, so a failed or invalid download keeps the existing database. Add `-v` to also update nali itself to the latest version.

### Specify database

Users can specify which database to use by setting the environment variables `NALI_DB_IP4`, `NALI_DB_IP6` and `NALI_DB_CDN`, or `selected` in the config file. Any name or alias from [Databases](#databases) can be used.

#### Windows

##### Use geoip db

```
set NALI_DB_IP4=geoip

or use powershell

$env:NALI_DB_IP4="geoip"
```

##### Use ipip db

```
set NALI_DB_IP6=ipip

or use powershell

$env:NALI_DB_IP6="ipip"
```

#### Linux

##### Use geoip db

```
export NALI_DB_IP4=geoip
```

##### Use ipip db

```
export NALI_DB_IP6=ipip
```

### Multilingual support

Set the environment variable `NALI_LANG` to choose the output language of the GeoIP2 / GeoLite2 / DB-IP databases; English is used when the database has no names in that language. The supported values are listed by the GeoIP2 database. Other databases provide a single language (IPinfo Lite is English).

```
# NALI_LANG=en NALI_DB_IP4=geoip nali 1.1.1.1
1.1.1.1 [Australia]
```

### Change directory

Set the environment variable `NALI_HOME` to specify the working directory where the configuration file and database are stored. You can also use absolute paths in the configuration file to specify other database paths.

Set the environment variable `NALI_CONFIG_HOME` to specify the configuration file directory and `NALI_DB_HOME` to specify the database file directory.

If no environment variable is specified, the XDG specification will be used, with the configuration file directory in `$XDG_CONFIG_HOME/nali` and the database file directory in `$XDG_DATA_HOME/nali`.

```
set NALI_HOME=D:\nalidb

or

export NALI_HOME=/home/nali
```

## Development

Requires Go >= 1.26.

```sh
make test      # go test -race ./...
make lint      # go vet + golangci-lint
make all-arch  # cross-compile every release platform
```

Tests need neither network access nor real databases. The MMDB fixtures in `testdata/mmdb` are produced by the separate module `testdata/mmdb/generate`; regenerate them with:

```sh
cd testdata/mmdb/generate && go run .
```

## Thanks

- [纯真QQIP离线数据库](http://www.cz88.net)
- [qqwry纯真数据库解析](https://github.com/yinheli/qqwry)
- [ZX公网ipv6数据库](https://ip.zxinc.org/ipquery/)
- [Geoip2 city数据库](https://www.maxmind.com/en/geoip2-precision-city-service)
- [geoip2-golang解析器](https://github.com/oschwald/geoip2-golang)
- [maxminddb-golang](https://github.com/oschwald/maxminddb-golang)
- [IPinfo Lite](https://ipinfo.io/lite)
- [CDN provider数据库](https://github.com/SukkaLab/cdn)
- [IPIP数据库](https://www.ipip.net/product/ip.html)
- [IPIP数据库解析](https://github.com/ipipdotnet/ipdb-go)
- [ip2region数据库](https://github.com/lionsoul2014/ip2region)
- [IP2Location DB3 LITE](https://lite.ip2location.com/database/db3-ip-country-region-city) *use the IPv6 BIN as it contains both IPv4 & IPv6*
- [Cobra CLI库](https://github.com/spf13/cobra)

Thanks to JetBrains for the Open Source License 

<a href="https://www.jetbrains.com/?from=nali">
  <img src="assets/GoLand.svg">
</a>

## Author

**Nali** © [zu1k](https://github.com/zu1k), Released under the [MIT](./LICENSE) License.<br>
