# TJS Serial Display 命令说明

TJS Serial Display 是一个用于控制和管理 TJC 串口显示屏的命令行工具。

## 基本用法

```bash
tjs-serial-display [命令] [选项]
```

## 命令列表

### 1. list-ports

列出系统中所有可用的串口设备。

**语法：**
```bash
tjs-serial-display list-ports
```

**示例：**
```bash
$ tjs-serial-display list-ports
Found port: /dev/ttyUSB0
Found port: /dev/ttyUSB1
```

---

### 2. info

获取已连接的 TJC 串口显示屏设备信息。

**语法：**
```bash
tjs-serial-display info [-p|--port <port_name>] [-b|--baud <baud_rate>] [-a|--auto]
```

**参数：**
- `-p, --port <port_name>`: 指定串口设备路径
- `-b, --baud <baud_rate>`: 可选，波特率（默认：115200）
- `-a, --auto`: 自动遍历所有可用串口设备并尝试连接

**注意：** `--port` 和 `--auto` 二选一，如果都不指定则默认使用 `--auto`

**示例：**
```bash
# 指定串口设备
tjs-serial-display info --port /dev/ttyUSB0

# 自动检测串口设备
tjs-serial-display info --auto
tjs-serial-display info -a

# 不指定参数，默认自动检测
tjs-serial-display info
```

**输出信息包括：**
- 屏幕类型（非触摸屏/电阻屏/电容屏）
- 设备地址
- 设备型号
- 固件版本号
- 主控芯片编号
- 设备唯一编号
- Flash 存储大小

---

### 3. exec

执行 TJC 串口屏的原始指令。

**语法：**
```bash
tjs-serial-display exec <command> [-p|--port <port_name>] [-b|--baud <baud_rate>] [-a|--auto]
```

**参数：**
- `<command>`: 必需，要执行的 TJC 指令字符串
- `-p, --port <port_name>`: 指定串口设备路径
- `-b, --baud <baud_rate>`: 可选，波特率（默认：115200）
- `-a, --auto`: 自动遍历所有可用串口设备并尝试连接

**示例：**
```bash
# 获取当前页面
tjs-serial-display exec "sendme" --port /dev/ttyUSB0

# 跳转到页面 2
tjs-serial-display exec "page 2" --auto

# 打印文本框内容
tjs-serial-display exec "print t0.txt" -p /dev/ttyUSB0

# 设置文本内容
tjs-serial-display exec "t0.txt=\"Hello\"" --auto

# 隐藏控件
tjs-serial-display exec "vis t0,0" --auto

# 显示控件
tjs-serial-display exec "vis t0,1" --auto

# 模拟按下按钮
tjs-serial-display exec "click b0,1" --auto

# 模拟弹起按钮
tjs-serial-display exec "click b0,0" --auto
```

**常用 TJC 指令：**
- `page <id>`: 跳转到指定页面
- `sendme`: 获取当前页面 ID
- `print <target>`: 打印目标值
- `<component>.txt="value"`: 设置文本值
- `<component>.val=<number>`: 设置数值
- `vis <component>,0`: 隐藏控件
- `vis <component>,1`: 显示控件
- `click <component>,0`: 模拟弹起
- `click <component>,1`: 模拟按下
- `ref <component>`: 刷新控件
- `sleep=1`: 进入睡眠
- `sleep=0`: 退出睡眠

---

### 4. upgrade

升级连接设备的程序。

**语法：**
```bash
tjs-serial-display upgrade <tft_file> [-p|--port <port_name>] [-b|--baud <baud_rate>] [-a|--auto]
```

**参数：**
- `<tft_file>`: 必需，TFT 固件文件路径
- `-p, --port <port_name>`: 指定串口设备路径
- `-b, --baud <baud_rate>`: 可选，波特率（默认：115200）
- `-a, --auto`: 自动遍历所有可用串口设备并尝试连接

**示例：**
```bash
# 指定串口
tjs-serial-display upgrade program.tft --port /dev/ttyUSB0

# 自动检测
tjs-serial-display upgrade program.tft --auto
```

**注意事项：**
- 升级过程中请勿断开设备电源
- 确保程序文件与设备型号匹配
- 升级完成后设备会自动重启

---

### 5. help

显示帮助信息和命令用法。

**语法：**
```bash
tjs-serial-display help [command]
```

**参数：**
- `[command]`: 可选，查看特定命令的详细帮助

**示例：**
```bash
# 显示所有命令列表
tjs-serial-display help

# 显示特定命令的详细帮助
tjs-serial-display help info
tjs-serial-display help exec
```

---

## 全局选项

以下选项可用于大多数需要串口连接的命令：

| 选项 | 说明 | 默认值 | 示例 |
|------|------|--------|------|
| `-p, --port <port_name>` | 串口设备路径 | 无 | `-p /dev/ttyUSB0` 或 `--port /dev/ttyUSB0` |
| `-b, --baud <baud_rate>` | 串口波特率 | 115200 | `-b 9600` 或 `--baud 9600` |
| `-a, --auto` | 自动遍历检测串口设备 | - | `-a` 或 `--auto` |

**端口选择规则：**
- 如果指定了 `--port`，则使用指定的串口设备
- 如果指定了 `--auto` 或 `-a`，则自动遍历所有可用串口并尝试连接
- 如果两者都未指定，默认使用 `--auto` 模式
- `--port` 和 `--auto` 不能同时使用

### 支持的波特率

常用波特率包括：2400, 4800, 9600, 19200, 38400, 57600, 115200, 230400

---

## 使用示例

### 示例 1：检测和连接设备

```bash
# 1. 列出可用串口
tjs-serial-display list-ports

# 2. 自动检测并获取设备信息（推荐）
tjs-serial-display info --auto

# 3. 指定串口获取设备信息
tjs-serial-display info --port /dev/ttyUSB0
```

### 示例 2：使用 exec 命令控制设备

```bash
# 跳转到页面 2
tjs-serial-display exec "page 2" --auto

# 获取当前页面
tjs-serial-display exec "sendme" --auto

# 设置文本内容
tjs-serial-display exec "t0.txt=\"Hello World\"" -p /dev/ttyUSB0

# 读取文本内容
tjs-serial-display exec "print t0.txt" --auto
```

### 示例 3：控制控件显示和操作

```bash
# 隐藏文本框 t0
tjs-serial-display exec "vis t0,0" --auto

# 显示按钮 b1
tjs-serial-display exec "vis b1,1" --auto

# 模拟按下按钮 b0
tjs-serial-display exec "click b0,1" -p /dev/ttyUSB0

# 模拟弹起按钮 b0
tjs-serial-display exec "click b0,0" --auto
```

### 示例 4：升级固件

```bash
# 指定串口升级
tjs-serial-display upgrade program.tft --port /dev/ttyUSB0

# 自动检测串口升级
tjs-serial-display upgrade program.tft --auto
```

### 示例 5：快速开始（推荐新手）

```bash
# 不需要指定任何端口，工具会自动检测并连接

# 获取设备信息
tjs-serial-display info

# 跳转到页面 1
tjs-serial-display exec "page 1"

# 设置文本
tjs-serial-display exec "t0.txt=\"测试\""

# 读取数据
tjs-serial-display exec "print t0.txt"
```

---

## 设备响应码说明

工具会自动处理设备的各种响应码，包括：

### 成功响应
- `0x01`: 指令成功执行

### 错误响应
- `0x00`: 无效指令
- `0x02`: 控件 ID 无效
- `0x03`: 页面 ID 无效
- `0x04`: 图片 ID 无效
- `0x05`: 字库 ID 无效
- `0x06`: 文件操作失败
- `0x09`: CRC 校验失败
- `0x11`: 波特率设置无效
- `0x1A`: 变量名称无效
- `0x1E`: 参数数量无效
- `0x24`: 串口缓冲区溢出

### 事件响应
- `0x65`: 控件点击事件
- `0x66`: 页面 ID 返回
- `0x67`: 触摸坐标数据
- `0x70`: 字符串数据返回
- `0x71`: 数值数据返回
- `0x86`: 自动进入睡眠模式
- `0x87`: 自动唤醒
- `0x88`: 系统启动成功

---

## 故障排除

### 串口权限问题

在 Linux 系统上，可能需要将用户添加到 `dialout` 组：

```bash
sudo usermod -a -G dialout $USER
# 注销后重新登录生效
```

或临时授予权限：

```bash
sudo chmod 666 /dev/ttyUSB0
```
