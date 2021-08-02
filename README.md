<h1 align="center">
  <br>Nali<br>
</h1>

<h4 align="center">An offline tool for querying IP geographic information and CDN provider.</h4>

<p align="center">
  <a href="https://github.com/zu1k/nali/actions">
    <img src="https://img.shields.io/github/workflow/status/zu1k/nali/Go?style=flat-square" alt="Github Actions">
  </a>
  <a href="https://goreportcard.com/report/github.com/zu1k/nali">
    <img src="https://goreportcard.com/badge/github.com/zu1k/nali?style=flat-square">
  </a>
  <a href="https://github.com/zu1k/nali/releases">
    <img src="https://img.shields.io/github/release/zu1k/nali/all.svg?style=flat-square">
  </a>
</p>

#### [中文文档](https://github.com/zu1k/nali/blob/master/README_zh-CN.md)

## Origin

Inspired by Nali C version and nali-cli js version.

I want to query the IP geographic information and CDN service provider on the terminal, then found the Nali tool. Nali mean 'where' in Chinese, good name for this kind of tools.

However the C version has too few functions, and the js version is too big and the supported platforms are not complete, so I rewrite it in golang, add IPv6 support and Geoip2 database.

## Feature

- Chunzhen qqip database
- ZX ipv6 database
- Geoip2 city database
- IPIP free database
- Pipeline support
- Interactive query
- Offline query
- Both ipv4 and ipv6 supported
- CDN provider query
- Full platform support
- Color print

## Install

### Install from source

Nali Requires Go >= 1.14. You can build it from source:

```sh
$ go get -u -v github.com/zu1k/nali
```

### Install pre-build binariy

Pre-built binaries are available here: [release](https://github.com/zu1k/nali/releases)

Download the binary compatible with your platform, unpack and copy to the directory in path

### Install from docker

```
docker pull docker.pkg.github.com//zu1k/nali/nali:latest
```

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

### Use with dig

```
$ dig nali.lgf.im +short | nali
104.28.2.115 [美国 CloudFlare公司CDN节点]
104.28.3.115 [美国 CloudFlare公司CDN节点]
172.67.135.48 [美国 CloudFlare节点]
```

### Use with nslookup

```
$ nslookup nali.lgf.im 8.8.8.8 | nali
Server:         8.8.8.8 [美国 加利福尼亚州圣克拉拉县山景市谷歌公司DNS服务器]
Address:        8.8.8.8 [美国 加利福尼亚州圣克拉拉县山景市谷歌公司DNS服务器]#53

Non-authoritative answer:
Name:   nali.lgf.im
Address: 104.28.3.115 [美国 CloudFlare公司CDN节点]
Name:   nali.lgf.im
Address: 104.28.2.115 [美国 CloudFlare公司CDN节点]
Name:   nali.lgf.im
Address: 172.67.135.48 [美国 CloudFlare节点]
```

### Use with any other program

Because nali can read the contents of the `stdin` pipeline, it can be used with any program

```
bash abc.sh | nali
```

Nali will insert ip information after ip

### IPV6 support

Use like ipv4

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

### Help

```
$ nali --help
Usage:
  nali [flags]
  nali [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
  update      update chunzhen ip database

Flags:
      --gbk    Use GBK decoder
  -h, --help   help for nali

Use "nali [command] --help" for more information about a command.
```

### Update chunzhen IP database

```
$ nali update
2020/07/17 12:53:46 正在下载最新纯真 IP 库...
2020/07/17 12:54:05 已将最新的纯真 IP 库保存到本地 /root/.nali/qqwry.dat
```

### Use other database

Set environment variables `NALI_DB`

supported database:

- Geoip2 `['geoip', 'geoip2', 'geo']`
- Chunzhen `['chunzhen', 'qqip', 'qqwry']`
- IPIP `['ipip', 'ipipfree', 'ipip.net']`

#### Windows

##### Use geoip db

```
set NALI_DB=geoip

or use powershell

$env:NALI_DB="geoip"
```

##### Use ipip db

```
set NALI_DB=ipip

or use powershell

$env:NALI_DB="ipip"
```

#### Linux

##### Use geoip db

```
export NALI_DB=geoip
```

##### Use ipip db

```
export NALI_DB=ipip
```

### Change database directory

If the database directory is not specified, the database will be placed in `~/.nali`

Set environment variables `NALI_DB_HOME` to specify the database directory

```
set NALI_DB_HOME=D:\nalidb

or

export NALI_DB_HOME=/home/nali
```

## Thanks

- [纯真QQIP离线数据库](http://www.cz88.net/fox/ipdat.shtml)
- [qqwry mirror](https://qqwry.mirror.noc.one/)
- [qqwry纯真数据库解析](https://github.com/yinheli/qqwry)
- [ZX公网ipv6数据库](https://ip.zxinc.org/ipquery/)
- [Geoip2 city数据库](https://www.maxmind.com/en/geoip2-precision-city-service)
- [geoip2-golang解析器](https://github.com/oschwald/geoip2-golang)
- [CDN provider数据库](https://github.com/SukkaLab/cdn)
- [IPIP数据库](https://www.ipip.net/product/ip.html)
- [IPIP数据库解析](https://github.com/ipipdotnet/ipdb-go)
- [Cobra CLI库](https://github.com/spf13/cobra)
- [Nali-cli](https://github.com/SukkaW/nali-cli)

## License

MIT
