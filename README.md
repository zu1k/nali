<h1 align="center">
  <br>Nali<br>
</h1>

<h4 align="center">一个查询IP地理信息和CDN提供商的离线终端工具.</h4>

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

#### [English](https://github.com/zu1k/nali/blob/master/README_en.md)

## 功能

- 支持多种数据库
  - 纯真 IPv4 离线数据库
  - ZX IPv6 离线数据库
  - ip2region IPv4 / IPv6 数据库 (可选)
  - IPinfo Lite 数据库 (可选)
  - Geoip2 城市数据库 (可选)
  - IPIP 数据库 (可选)
  - DB-IP 数据库 (可选)
  - IP2Location DB3 LITE 数据库 (可选)
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

Nali 需要预先安装 Go >= 1.26. 安装后可以从源码安装软件:

```sh
$ go install github.com/zu1k/nali@latest
```

### 下载预编译的可执行程序

可以从Release页面下载预编译好的可执行程序: [Release](https://github.com/zu1k/nali/releases)

你需要选择适合你系统和硬件架构的版本下载，解压后可直接运行

### Arch 系 Linux

我们在 Aur 中发布了 3 个相关的包:

- `nali-go`: Release 版本，安装时编译
- `nali-go-bin`: Release 版本，预编译的二进制文件
- `nali-go-git`: 最新的 master 分支版本，安装时编译

### Debian / Ubuntu

从 [Release](https://github.com/zu1k/nali/releases) 页面下载对应架构（amd64、i386、armhf、arm64）的 `.deb` 包后安装：

```sh
$ sudo apt install ./nali_*_amd64.deb
```

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

程序第一次运行后，会在 config 目录生成配置文件 `config.yaml` (使用 `nali info` 来查看具体信息)，配置文件定义了数据库信息，默认用户无需进行修改

数据库格式默认如下：

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

其中，`languages` 和 `types` 表示该数据库支持的语言和查询类型。 如果你需要增加数据库，需小心修改配置文件，如果有任何问题，欢迎提 issue 询问。

旧版本生成的配置文件无需修改：配置文件中找不到的数据库会使用内置的默认定义。

### 数据库一览

| 数据库 | 名称 / 别名 | 查询类型 | 默认文件 | 自动下载 |
| --- | --- | --- | --- | --- |
| 纯真 IPv4 | `qqwry`, `chunzhen` | IPv4（默认） | `qqwry.dat` | ✓ |
| ZX IPv6 | `zxipv6wry`, `zxipv6`, `zx` | IPv6（默认） | `zxipv6wry.db` | ✓ |
| ip2region IPv4 | `ip2region`, `i2r` | IPv4 | `ip2region.xdb` | ✓ |
| ip2region IPv6 | `ip2region-ipv6`, `i2r-ipv6` | IPv6 | `ip2region_v6.xdb` | ✓ |
| IPinfo Lite | `ipinfo` | IPv4、IPv6 | `ipinfo_lite.mmdb` | ✓ |
| GeoIP2 / GeoLite2 City | `geoip`, `geoip2`, `geolite`, `geolite2` | IPv4、IPv6 | `GeoLite2-City.mmdb` | 需手动下载 |
| DB-IP | `dbip`, `db-ip` | IPv4、IPv6 | `dbip.mmdb` | 需手动下载 |
| IPIP | `ipip` | IPv4、IPv6 | `ipipfree.ipdb` | 需手动下载 |
| IP2Location DB3 LITE | `ip2location` | IPv4、IPv6 | `IP2LOCATION-LITE-DB3.IPV6.BIN` | 需手动下载 |
| CDN | `cdn` | 域名（默认） | `cdn.yml` | ✓ |

支持自动下载的数据库在首次使用、本地文件不存在时会自动下载；需手动下载的数据库请将文件放入 `nali info` 显示的数据目录，或在配置文件中用绝对路径指定 `file`。

### 使用 IPinfo Lite

[IPinfo Lite](https://ipinfo.io/lite) 提供免费的国家和 ASN 数据，同时支持 IPv4 和 IPv6，首次使用时自动下载：

```
$ NALI_DB_IP4=ipinfo NALI_DB_IP6=ipinfo nali 8.8.8.8 2001:4860:4860::8888
8.8.8.8 [United States AS15169 Google LLC]
2001:4860:4860::8888 [United States AS15169 Google LLC]

$ NALI_DB_IP4=ipinfo nali --json 8.8.8.8
{"type":0,"ip":"8.8.8.8","text":"United States AS15169 Google LLC","source":"ipinfo","info":{"network":"8.8.8.0/24","country":"United States","country_code":"US","continent":"North America","continent_code":"NA","asn":"AS15169","as_name":"Google LLC","as_domain":"google.com"}}
```

长期使用可在配置文件中将 `selected.ipv4`、`selected.ipv6` 设为 `ipinfo`。

文本输出包含国家、ASN 和 AS 名称；JSON 输出包含 `network`、`country`、`country_code`、`continent`、`continent_code`、`asn`、`as_name`、`as_domain`。数据库记录未提供 `network` 字段时，使用命中的 CIDR 网段。名称为数据库中的英文原文。

`format` 为 `mmdb` 的数据库会根据元数据自动识别是否为 IPinfo 格式，因此也可以使用登录 [IPinfo](https://ipinfo.io/dashboard/downloads) 后下载的 `ipinfo_lite.mmdb`。默认下载地址是社区维护的镜像 [NetworkCats/IPinfoLite-Download](https://github.com/NetworkCats/IPinfoLite-Download)，可以在配置文件中通过 `download-urls` 改为其他地址：

```yaml
- name: ipinfo
  format: mmdb
  file: ipinfo_lite.mmdb
  download-urls:
  - https://github.com/NetworkCats/IPinfoLite-Download/releases/latest/download/ipinfo_lite.mmdb
  languages: [en]
  types: [IPv4, IPv6]
```

IPinfo Lite 数据以 [CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/) 许可发布，公开使用其数据时请注明来源 IPinfo。

### 使用 ip2region IPv6

ip2region 分为 IPv4（`ip2region`）和 IPv6（`ip2region-ipv6`）两个数据库文件，首次使用时自动下载：

```
$ NALI_DB_IP4=ip2region NALI_DB_IP6=ip2region-ipv6 nali 223.5.5.5 240e::1
223.5.5.5 [中国|浙江省|杭州市|阿里|CN]
240e::1 [中国|北京市|北京市|电信|CN]
```

旧版本的 `ip2region.xdb` 文件仍可直接使用；配置文件中失效的旧下载地址会在启动时自动迁移到 `ip2region_v4.xdb`。

### 查看帮助

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
      --color string   Colorize output: auto, always or never (NO_COLOR is honored in auto mode) (default "auto")
      --gbk            Decode input as GBK (detected automatically on Chinese Windows)
  -h, --help           help for nali
  -j, --json           Output in JSON format
  -v, --version        version for nali

Use "nali [command] --help" for more information about a command.
```

### 更新数据库

不带参数时更新纯真、ZX IPv6、ip2region (IPv4) 和 CDN 数据库。体积较大的可选数据库 `ip2region-ipv6`（约 37 MB）和 `ipinfo`（约 24 MB）只在被选用（配置文件的 `selected` 或 `NALI_DB_*` 环境变量）或本地已存在时才会更新。

```
$ nali update
2026/09/30 09:51:48 正在下载最新 qqwry 数据库...
2026/09/30 09:51:49 qqwry 数据库下载成功: qqwry.dat
2026/09/30 09:51:49 正在下载最新 cdn 数据库...
2026/09/30 09:51:50 cdn 数据库下载成功: cdn.yml
...
```

或者指定数据库（逗号分隔，支持别名），指定的数据库总会更新：

```
$ nali update --db qqwry,cdn,ipinfo
```

本地数据库已是最新时不会重复下载（显示 `数据库已是最新版本`）。下载的数据会先校验再替换本地文件，下载失败或内容无效时保留原有的数据库。加上 `-v` 参数会同时将 nali 更新到最新版本。

### 自选数据库

用户可以指定使用哪个数据库，需要设置环境变量 `NALI_DB_IP4`、`NALI_DB_IP6`、`NALI_DB_CDN`，或者修改配置文件中的 `selected`。可以使用[数据库一览](#数据库一览)中的任意名称或别名。

#### 同时查询多个数据库

用逗号分隔多个数据库名称，每个数据库的结果显示在各自的方括号中：

```
$ NALI_DB_IP4=qqwry,ipinfo nali 8.8.8.8
8.8.8.8 [美国–加利福尼亚州–圣克拉拉–山景城 谷歌公司DNS服务器] [United States AS15169 Google LLC]
```

JSON 输出中 `source`、`text`、`info` 仍为第一个数据库的结果，另外 `results` 数组包含所有数据库的结果（仅在选择了多个数据库时出现）。

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

通过环境变量 `NALI_LANG` 指定 GeoIP2 / GeoLite2 / DB-IP 数据库的输出语言，数据库中没有对应语言时使用英文。该参数可设置的值见 GeoIP2 数据库的支持列表。其他数据库只提供单一语言（例如 IPinfo Lite 为英文）。

```
# NALI_LANG=en NALI_DB_IP4=geoip nali 1.1.1.1
1.1.1.1 [Australia]
```

### 工作目录

设置环境变量 `NALI_HOME` 来指定工作目录，配置文件和数据库存放在工作目录下。也可在配置文件中使用绝对路径指定其他数据库路径。

设置环境变量 `NALI_CONFIG_HOME` 来指定配置文件目录，`NALI_DB_HOME` 来执行数据库文件目录

如果未指定相关环境变量，将使用 XDG 规范，配置文件目录在 `$XDG_CONFIG_HOME/nali`，数据库文件目录在 `$XDG_DATA_HOME/nali`

```
set NALI_HOME=D:\nali

or

export NALI_HOME=/var/nali
```

## 常见问题

### 在管道中间使用时没有颜色

默认只在输出到终端时显示颜色。使用 `--color always` 强制输出颜色，例如 `cmd | nali --color always | less -R`；使用 `--color never` 或设置环境变量 `NO_COLOR=1` 关闭颜色。

### Windows 下中文乱码

在简体中文 Windows（代码页 936）中，nali 会自动把非 UTF-8 的输入按 GBK 解码，`tracert`、`nslookup` 等命令的输出可以直接通过管道交给 nali。如果仍然乱码，可以加 `--gbk` 强制按 GBK 解码输入。

### `tail -f log | jq . | nali` 最后几行不输出

nali 读到一行就立即输出。`jq` 的输出不是终端时会缓冲，所以最后几行停留在 jq 中，使用 `jq --unbuffered` 即可：

```sh
tail -f access.log | jq --unbuffered . | nali
```

### 筛选 CDN 结果

域名的 CDN 识别结果在 JSON 输出中 `source` 为 `cdn`，可以用 jq 筛选：

```sh
nali --json < domains.txt | jq -c 'select(.source == "cdn")'
```

## 开发

需要 Go >= 1.26。

```sh
make test      # go test -race ./...
make lint      # go vet + golangci-lint
make all-arch  # 交叉编译所有发布平台
```

测试不依赖网络和真实数据库。MMDB 测试数据位于 `testdata/mmdb`，由独立模块 `testdata/mmdb/generate` 生成，修改后可重新生成：

```sh
cd testdata/mmdb/generate && go run .
```

## 感谢列表

- [纯真QQIP离线数据库](http://www.cz88.net)
- [qqwry纯真数据库解析](https://github.com/yinheli/qqwry)
- [ZX公网ipv6数据库](https://ip.zxinc.org/ipquery/)
- [Geoip2 city数据库](https://www.maxmind.com/en/geoip2-precision-city-service)
- [geoip2-golang解析器](https://github.com/oschwald/geoip2-golang)
- [maxminddb-golang解析器](https://github.com/oschwald/maxminddb-golang)
- [IPinfo Lite数据库](https://ipinfo.io/lite)
- [CDN provider数据库](https://github.com/SukkaLab/cdn)
- [IPIP数据库](https://www.ipip.net/product/ip.html)
- [IPIP数据库解析](https://github.com/ipipdotnet/ipdb-go)
- [ip2region数据库](https://github.com/lionsoul2014/ip2region)
- [IP2Location DB3 LITE](https://lite.ip2location.com/database/db3-ip-country-region-city)
- [Cobra CLI库](https://github.com/spf13/cobra)

感谢 JetBrains 提供开源项目免费License 

<a href="https://www.jetbrains.com/?from=nali">
  <img src="assets/GoLand.svg">
</a>

## 作者

**Nali** © [zu1k](https://github.com/zu1k), 遵循 [MIT](./LICENSE) 证书.<br>
