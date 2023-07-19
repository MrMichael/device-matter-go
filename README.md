# Device Matter Go
## 简介
基于[Matter](https://csa-iot.org/all-solutions/matter/) SDK， 根据x86平台编译所得chip-tool命令行工具，对其进行EdgeX设备微服务的封装，使应用程序其可以通过Restful API控制Matter Device。



## 使用方法

#### 1）工程目录说明

```shell
├── Attribution.txt
├── bin
├── CHANGELOG.md
├── cmd
│   ├── credentials		# Matter PAA证书
│   ├── chip-app1	  # 已编译的matter simulated device可执行文件（x86）
│   ├── chip-tool		# 已编译的matter chip-tool可执行文件（x86）
│   ├── main.go
│   └── res		# 微服务配置文件
│       ├── configuration.toml
│       ├── devices
│       │   ├── device.matter.chiptool.toml
│       │   └── device.matter.simulator.toml
│       └── profiles
│           ├── device.matter.chiptool.yaml
│           └── device.matter.simulator.yaml
├── Dockerfile	# Docker镜像构建文件
├── go.mod
├── go.sum
├── GOVERNANCE.md
├── internal
│   └── driver	# 设备微服务驱动文件
│       ├── chiptool.go
│       └── driver.go
├── Jenkinsfile
├── LICENSE
├── Makefile
├── README.md
├── vendor
└── version.go
```

#### 2）Device Profile

- device.matter.chiptool.toml

  ```yaml
  deviceResources:
  -
    name: "cmdParams"
    isHidden: false
    description: "the parameters for chip-tool(matter command-line tool) to send"
    properties:
      valueType: "StringArray"
      readWrite: "RW"
  ```

  - 定义一个命令行工具设备，提供StringArray类型的参数cmdParams，能将接收到的字符串数组转化为chip-tool命令行参数，然后执行，最后将返回执行结果。

- device.matter.simulator.yaml

  ```yaml
  ...
  deviceCommands:
    -
      name: "ParseSetupPayload"
      isHidden: false
      readWrite: "RW"
      resourceOperations:
        - { deviceResource: "qrcode_payload"}
    -
      name: "CommissionIntoWiFiOverBT"
      isHidden: false
      readWrite: "RW"
      resourceOperations:
        - { deviceResource: "node_id"}
        - { deviceResource: "ssid" }
        - { deviceResource: "password" }
        - { deviceResource: "pin_code" }
        - { deviceResource: "discriminator" }
        - { deviceResource: "timeout" , defaultValue: "10"}
    -
      name: "CommissionWithQRCode"
      isHidden: false
      readWrite: "RW"
      resourceOperations:
        - { deviceResource: "node_id"}
        - { deviceResource: "qrcode_payload" }
    -
      name: "RemoveFromFabric"
      isHidden: false
      readWrite: "RW"
      resourceOperations:
        - { deviceResource: "node_id"}
    -
      name: "Cluster_On_Off_Toggle"
      isHidden: false
      readWrite: "RW"
      resourceOperations:
        - { deviceResource: "node_id"}
        - { deviceResource: "endpoint_id", defaultValue: "1"}
        - { deviceResource: "on_off_toggle", defaultValue: "1"}
    -
      name: "Cluster_Read_On_Off"
      isHidden: false
      readWrite: "RW"
      resourceOperations:
        - { deviceResource: "node_id"}
        - { deviceResource: "endpoint_id", defaultValue: "1"}
  ```

- 定义一个Matter灯泡类型设备，提供qrcode_payload解析、设备匹配入网及退网、设备开关控制，设备状态读取等命令功能。

- 每一个命令均支持读写操作，通过读操作可以获取上一次写操作的返回结果。

#### 3）driver.chiptool

通过os/exec执行chiptool，根据不同功能输入相应参数，并解析返回结果。

```go
...
func (s *Driver) chipToolParseSetupPayload(payload string) error {

	executable := "./chip-tool"
	args := []string{"payload","parse-setup-payload", payload}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, executable, args...)
	out, err := cmd.CombinedOutput()

	if ctx.Err() == nil {
		if err == nil {
			cmd2 := exec.Command("grep", "-oE", "VendorID.*|ProductID.*|Discovery.*|Long.*|Passcode.*")
			cmd2.Stdin = bytes.NewReader(out)
			out2, _ := cmd2.Output()
			s.parseSetupPayloadResp = string(out2)
		} else {
			s.parseSetupPayloadResp = "Result failure"
			return err
		}
	} else {
		s.parseSetupPayloadResp = "Exec failure"
		return ctx.Err()
	}
	return nil
}
```

#### 4）程序编译及运行

```shell
# 关闭安全模式
$ export EDGEX_SECURITY_SECRET_STORE=false
# 由于edgex服务组件使用容器以bridge网络启动，因此需要修改SERVICE_HOST为docker0 ip，否则会绑定127.0.0.1
$ export SERVICE_HOST="172.17.0.1"

# 编译工程
cd device-matter-go
make build

# 运行程序
cd cmd
./device-matter

# 构建egdex-matter镜像
cd ..
make docker
```



## License
[Apache-2.0](LICENSE)
