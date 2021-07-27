混沌工程agent在目标机器上启动http服务器，接收来自管理端或proxy的指令执行故障注入。

## development

### 证书生成

程序启动时需要指定鉴权用公钥，本地开发测试时可以通过如下方式生成密钥和公钥：

```sh
openssl ecparam -out ec_key.pem -name prime256v1 -genkey
openssl ec -in ec_key.pem -pubout -out ec_pub.pem
```

### 相关package

- 配置读取使用[viper](https://github.com/spf13/viper)库
- 命令行参数解析暂定使用[pflag](https://github.com/spf13/pflag)
- http服务使用[gin](https://github.com/gin-gonic/gin)


## 鉴权

程序启动时解析公钥，并对所有请求进行签名验证(JWT)。作为请求方需要按如下方式添加签名信息：

1. 准备private key
2. 构造签名的内容：`{"iat": 1594092544, "exp": 1594092559}`，其中`iat`代表`issue at`，`exp`代表过期时间，建议设置两者相差不超过30秒。详情参考JWT的[RFC](https://tools.ietf.org/html/rfc7519#page-9)
3. 使用private key对上面内容按照JWT规范签名，算法使用`ES256` (ECDSA P-256 curve, SHA-256 hash algo)，得到`token`
4. Header头带上签名`Authorization: Bearer <token>`


### snippet

python

```python
import time
import requests
import jwt

key = open('./ec_key.pem').read()

url = ''
now = int(time.time())
auth = jwt.encode({'iat': now, 'exp': now + 10}, key, 'ES256').decode()
print(auth)
r = requests.get(url, headers={'Authorization': f'Bearer {auth}'})
print(r.status_code)
```

go

```go
*Omit the error report*

package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	prikey_file := "/path/to/key/pem"
	prikey_bytes, _ := ioutil.ReadFile(prikey_file)
	_, rest := pem.Decode(prikey_bytes)
	keyContent, _ := pem.Decode(rest)
	key, _ := x509.ParseECPrivateKey(keyContent.Bytes)
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"exp": time.Now().Add(time.Second * 30),
		"iat": time.Now(),
	})

	tokenString, _ := token.SignedString(key)
	bearer := "Bearer " + tokenString
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/version", "http://127.0.0.1:1337"), nil)
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, _ := client.Do(req)
	fmt.Println(resp.StatusCode)
}
```

## Proxy安装
**编译**
```bash
go build
```
**配置&运行**
```bash 
sudo ./nessaj_proxy --sender_key_file /path/to/sender/key
```

## Agent安装
**编译**
```bash
go build
```

**配置&运行**
```bash 
sudo ./nessaj --pubkey_file ec_pub.pem --chaosblade_bin /path/to/chaosblade
```
- 公钥可以通过pubkey或者pubkey_file配置
- 另外还需要配置chaosblade_bin, 指向chaosblade-[version]/blade


## APIs

指令的参数或者实现方式请参考[Chaosblade](https://chaosblade-io.gitbook.io/chaosblade-help-zh-cn/)

### Create

创建一个混沌实验

**Url**: /chaos/run

**Method**: POST

**Content-Type**: json

**Params**:

| name   | description                                                  | type   | required |
| ------ | ------------------------------------------------------------ | ------ | -------- |
| op     | 实验操作类型（cpu/memory/diskburn/diskfill/netdelay/netdns/netloss/processkill） | string | Y        |
| params | 操作的参数                                                   | json   | Y        |

#### (op) cpu

CPU 负载实验

**params**

| name       | description                                  | type   | required |
| ---------- | -------------------------------------------- | ------ | -------- |
| timeout    | 设定运行时长，单位是秒，通用参数             | string | N        |
| cpuCount   | 指定 CPU 满载的个数                          | string | N        |
| cpuList    | 指定 CPU 满载的具体核，核索引从 0 开始 (0-3) | string | N        |
| cpuPercent | 指定 CPU 负载百分比，取值在 0-100            | string | N        |

P.S.

cpu-list不为空时，cpu-count为1

如以上参数均为空，默认全部核满载，无timeout

**Response**

所有operation的response都是一致的，后面不再累叙。

正常返回: 

```json
{
    "Code":0,
    "Data":{
        "Id":"8734b7c3b720d751",
        "Info":"{
        	'code': 200, 
        	'result': '8734b7c3b720d751', 
        	'success': True
    	}",
	},
    "Msg":""
}
```

错误返回:

```json
{
	"Code": 1,
	"Data": "",
	"Msg": "...."
}
```

**snippet**

python

```python
payload = {
    "op": "cpu",
    "params": {
        "cpuCount": "1",
        "cpuPercent": "50",
        "cpuList": "1-3",
        "timeout": "10"
    }
}

# headers 为上面鉴权的header
r_run = requests.post(host+'/chaos/run', json=payload, headers=headers)
ret_id = r_run.json()['Data']['Id']
print(ret_id)
ret_info = json.loads(r_run.json()['Data']['Info'])
print(ret_info)
```

#### (op) memory

内存占用实验

**params**

| name       | description                                                  | type   | required |
| ---------- | ------------------------------------------------------------ | ------ | -------- |
| timeout    | 设定运行时长，单位是秒，通用参数                             | string | N        |
| memPercent | 内存使用率，取值是 0 到 100 的整数                           | string | N        |
| mode       | 内存占用模式，有 ram 和 cache 两种。ram 采用代码实现，可控制占用速率，优先推荐此模式；cache 是通过挂载tmpfs实现；默认值是cache | string | N        |
| reserve    | 保留内存的大小，单位是MB，如果 mem-percent 参数存在，则优先使用 mem-percent 参数 | string | N        |
| rate       | 内存占用速率，单位是 MB/S，仅在 mode=ram时生效               | string | N        |

P.S.

此场景触发内存占用满，即使指定了 --timeout 参数，也可能出现通过 blade 工具无法恢复的情况，可通过**重启机器**解决！！！推荐**指定内存百分比**！

**snippet**

```python
payload = {
    "op": "memory",
    "params": {
        "memPercent": "10",
		"mode":       "ram",
		"reserve":    "200",
		"rate":       "100",
		"timeout":    "30",
    }
}
```

#### (op) diskburn

磁盘读写 io 负载实验

**params**

| name    | description                                                  | type   | required |
| ------- | ------------------------------------------------------------ | ------ | -------- |
| timeout | 设定运行时长，单位是秒，通用参数                             | string | N        |
| path    | 指定提升磁盘 io 的目录，会作用于其所在的磁盘上，默认值是 /   | string | N        |
| read    | 触发提升磁盘读 IO 负载，会创建 600M 的文件用于读，销毁实验会自动删除 | bool   | Y        |
| write   | 触发提升磁盘写 IO 负载，会根据块大小的值来写入一个文件，比如块大小是 10，则固定的块的数量是 100，则会创建 1000M 的文件，销毁实验会自动删除 | bool   | Y        |
| size    | 块大小, 单位是 M, 默认值是 10，一般不需要修改，除非想更大的提高 io 负载 | string | N        |

P.S.

read/write模式二选一（必须至少触发一种方式）

**snippet**

```python
payload = {
    "op": "diskburn",
    "params": {
        "path":    "/home/cthye/Desktop",
		"size":    "500",
		"write":   False,
		"read":    True,
		"timeout": "30"
    }
}
```

#### (op) diskfill

磁盘填充负载实验

**params**

| name         | description                                                  | type   | required |
| ------------ | ------------------------------------------------------------ | ------ | -------- |
| timeout      | 设定运行时长，单位是秒，通用参数                             | string | N        |
| path         | 指定提升磁盘 io 的目录，会作用于其所在的磁盘上，默认值是 /   | string | N        |
| reserve      | 保留磁盘大小，单位是MB。取值是不包含单位的正整数。如果 size、percent、reserve 参数都存在，优先级是 percent > reserve > size | string | N        |
| percent      | 指定磁盘使用率，取值是不带%号的正整数                        | string | N        |
| retainHandle | 是否保留填充                                                 | bool   | N        |
| size         | 块大小, 单位是 M, 默认值是 10，一般不需要修改，除非想更大的提高 io 负载 | string | N        |

**snippet**

```python
payload = {
    "op": "diskfill",
    "params": {
        "path":         "/home/cthye/Desktop",
		"size":         "500",
		"reserve":      "1024",
		"retainHandle": False,
		"timeout":      "10",
    }
}
```

#### (op) netdelay

磁盘填充负载实验

**params**

| name           | description                                                  | type   | required |
| -------------- | ------------------------------------------------------------ | ------ | -------- |
| timeout        | 设定运行时长，单位是秒，通用参数                             | string | N        |
| desIP          | 目标 IP. 支持通过子网掩码来指定一个网段的IP地址, 例如 192.168.1.0/24. 则 192.168.1.0~192.168.1.255 都生效。你也可以指定固定的 IP，如 192.168.1.1 或者 192.168.1.1/32，也可以通过都号分隔多个参数，例如 192.168.1.1,192.168.2.1 | string | N        |
| excludePort    | 排除掉的端口，默认会忽略掉通信的对端端口，目的是保留通信可用。可以指定多个，使用逗号分隔或者连接符表示范围，例如 22,8000 或者 8000-8010。 这个参数不能与 localPort 或者 remotePort 参数一起使用 | string | N        |
| excludeIP      | 排除受影响的 IP，支持通过子网掩码来指定一个网段的IP地址, 例如 192.168.1.0/24. 则 192.168.1.0~192.168.1.255 都生效。你也可以指定固定的 IP，如 192.168.1.1 或者 192.168.1.1/32，也可以通过都号分隔多个参数，例如 192.168.1.1,192.168.2.1 | string | N        |
| interface      | 网卡设备，例如 eth0                                          | string | Y        |
| localPort      | 本地端口，一般是本机暴露服务的端口。可以指定多个，使用逗号分隔或者连接符表示范围，例如 80,8000-8080 | string | N        |
| offset         | 延迟时间上下浮动的值, 单位是毫秒                             | string | N        |
| remoteport     | 远程端口，一般是要访问的外部暴露服务的端口。可以指定多个，使用逗号分隔或者连接符表示范围，例如 80,8000-8080 | string | N        |
| time           | 延迟时间，单位是毫秒 (必要参数)                              | string | Y        |
| force          | 强制覆盖已有的 tc 规则，请务必在明确之前的规则可覆盖的情况下使用 | bool   | N        |
| ignorePeerPort | 针对添加 excludePort 参数，报 ss 命令找不到的情况下使用，忽略排除端口 | bool   | N        |

P.S.

- 如果不指定端口、ip 参数，而是整个网卡延迟，切记要添加 timeout 参数或者excludePort 参数 (防止因延迟时间设置太长，造成机器无法连接的情况，如果真实发生此问题，重启机器即可恢复)

- Q: {"code":604,"success":false,"error":"RTNETLINK answers: File exists\n exit status 2 exit status 1"} A： 网络相关的场景实验已存在，销毁原有的后再执行。

- 无法销毁实验，则可以通过下面指令恢复实验：

  ```bash
  $ tc filter del dev [interface] parent 1: prio 4
  $ tc qdisc del dev [interface] root
  ```

  interface: 实验的网卡

**snippet**

```python
payload = {
    "op": "netdelay",
    "params": {
        "desIP": "220.181.38.148", 
	    "interface": "enp0s31f6", 
	    "time"           :"6000",
	    "timeout"        :"30"
    }
}
```

----

TO BE CONTINUED



