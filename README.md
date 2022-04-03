<h1 align="center">
  <br>Nali<br>
</h1>

<h4 align="center">一个查询IP地理信息和CDN提供商的离线终端工具.</h4>

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

#### [English](https://github.com/zu1k/nali/blob/master/README_en.md)

## 来源

该工具受 Nali C版本 和 nali-cli js版本的启发.

我想要在终端对IP地理信息和CDN服务提供商进行查询，发现了Nali这个工具，Nali与哪里谐音，非常适合这类工具

经过简单的使用，我发现最初的C语言版本功能缺失，而js版本包实在大的恐怖、而支持的平台非常有限，所以我用golang重写了这个工具，在原有功能的基础上增加了对IPv6的支持，并且增加了Geoip2数据库


## 功能

- 支持多种数据库
  - 纯真 IPv4 离线数据库
  - ZX IPv6 离线数据库
  - Geoip2 城市数据库 (可选)
  - IPIP 数据库 (可选)
  - ip2region 数据库 (可选)
- CDN 服务提供商查询
- 支持管道处理
- 支持交互式查询
- 同时支持IPv4和IPv6
- 支持多语言
- 查询完全离线
- 全平台支持
- 支持彩色输出

## 安装

### 从源码安装

Nali 需要预先安装 Go >= 1.18. 安装后可以从源码安装软件:

```sh
$ go install github.com/zu1k/nali
```

### 下载预编译的可执行程序

可以从Release页面下载预编译好的可执行程序: [Release](https://github.com/zu1k/nali/releases)

你需要选择适合你系统和硬件架构的版本下载，解压后可直接运行

## 使用说明

### 查询一个IP的地理信息

```
$ nali 1.2.3.4
1.2.3.4 [澳大利亚 APNIC Debogon-prefix网络]
```

#### 或者 使用 `管道`

```
$ echo IP 6.6.6.6 | nali
IP 6.6.6.6 [美国 亚利桑那州华楚卡堡市美国国防部网络中心]
```

### 同时查询多个IP的地理信息

```
$ nali 1.2.3.4 4.3.2.1 123.23.3.0
1.2.3.4 [澳大利亚 APNIC Debogon-prefix网络]
4.3.2.1 [美国 新泽西州纽瓦克市Level3Communications]
123.23.3.0 [越南 越南邮电集团公司]
```

### 交互式查询

使用 `exit` 或  `quit` 退出查询

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

### 与 `dig` 命令配合使用

需要你系统中已经安装好 dig 程序

```
$ dig nali.zu1k.com +short | nali
104.28.2.115 [美国 CloudFlare公司CDN节点]
104.28.3.115 [美国 CloudFlare公司CDN节点]
172.67.135.48 [美国 CloudFlare节点]
```

### 与 `nslookup` 命令配合使用

需要你系统中已经安装好 nslookup 程序

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

### 与任意程序配合使用

因为 nali 支持管道处理，所以可以和任意程序配合使用

```
bash abc.sh | nali
```

Nali 将在 IP后面插入IP地理信息，CDN域名后面插入CDN服务提供商信息

### 支持IPv6

和 IPv4 用法完全相同

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

### 查询 CDN 服务提供商

因为 CDN 服务通常使用 CNAME 的域名解析方式，所以推荐与 `nslookup` 或者 `dig` 配合使用，在已经知道 CNAME 后可单独使用

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

## 用户交互

### 查看帮助

```
$ nali --help
Usage:
  nali [flags]
  nali [command]

Available Commands:
  help        Help about any command
  update      update chunzhen ip database

Flags:
  -h, --help     help for nali
  -t, --toggle   Help message for toggle

Use "nali [command] --help" for more information about a command.
```

### 更新数据库

```
$ nali update
2020/07/17 12:53:46 正在下载最新纯真 IP 库...
2020/07/17 12:54:05 已将最新的纯真 IP 库保存到本地 /root/.nali/qqwry.dat
```

### 自选数据库

用户可以指定使用哪个数据库，需要设置环境变量： `NALI_DB_IP4`、`NALI_DB_IP6` 或者两个同时设置

支持的变量内容:

- Geoip2 `['geoip', 'geoip2', 'geo']`
- Chunzhen `['chunzhen', 'qqip', 'qqwry']`
- IPIP `['ipip', 'ipipfree', 'ipip.net']`
- Ip2Resion `['ip2region', 'region', 'i2r']`

#### Windows平台

##### 使用geoip数据库

```
set NALI_DB_IP4=geoip

或者使用 powershell

$env:NALI_DB_IP4="geoip"
```

##### 使用ipip数据库

```
set NALI_DB_IP6=ipip

或者使用 powershell

$env:NALI_DB_IP6="ipip"
```

#### Linux平台

##### 使用geoip数据库

```
export NALI_DB_IP4=geoip
```

##### 使用ipip数据库

```
export NALI_DB_IP4=ipip
```

### 多语言支持

通过修改环境变量 `NALI_LANG` 来指定使用的语言，当使用非中文语言时仅支持GeoIP2这个数据库

该参数可设置的值见 GeoIP2 这个数据库的支持列表

```
# NALI_LANG=en nali 1.1.1.1
1.1.1.1 [Australia]
```

### 更换数据库目录

如果未指定数据库存放目录，数据库默认将存放在 `~/.nali`

设置环境变量 `NALI_DB_HOME` 来指定数据库目录

```
set NALI_DB_HOME=D:\nalidb

or

export NALI_DB_HOME=/home/nali
```

## 感谢列表

- [纯真QQIP离线数据库](http://www.cz88.net/fox/ipdat.shtml)
- [qqwry mirror](https://qqwry.mirror.noc.one/)
- [qqwry纯真数据库解析](https://github.com/yinheli/qqwry)
- [ZX公网ipv6数据库](https://ip.zxinc.org/ipquery/)
- [Geoip2 city数据库](https://www.maxmind.com/en/geoip2-precision-city-service)
- [geoip2-golang解析器](https://github.com/oschwald/geoip2-golang)
- [CDN provider数据库](https://github.com/SukkaLab/cdn)
- [IPIP数据库](https://www.ipip.net/product/ip.html)
- [IPIP数据库解析](https://github.com/ipipdotnet/ipdb-go)
- [ip2region数据库](https://github.com/lionsoul2014/ip2region)
- [Cobra CLI库](https://github.com/spf13/cobra)
- [Nali-cli](https://github.com/SukkaW/nali-cli)

感谢 JetBrains 提供开源项目免费License 

<a href="https://www.jetbrains.com/?from=nali">
  <img src="assets/GoLand.svg">
</a>

## 作者

**Nali** © [zu1k](https://github.com/zu1k), 遵循 [MIT](./LICENSE) 证书.<br>

> Blog [@zu1k](https://zu1k.com) · GitHub [@zu1k](https://github.com/zu1k) · Twitter [@zu1k_lv](https://twitter.com/zu1k_lv) · Telegram Channel [@peekfun](https://t.me/peekfun)

## Star统计

[![Stargazers over time](https://starchart.cc/zu1k/nali.svg)](https://starchart.cc/zu1k/nali)