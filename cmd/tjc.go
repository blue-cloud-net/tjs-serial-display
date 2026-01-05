package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/blue-cloud-net/tjc-serial-display/internal/client"
	"github.com/blue-cloud-net/tjc-serial-display/internal/serial"
	"github.com/blue-cloud-net/tjc-serial-display/pkg/models"
)

const (
	defaultBaudRate = 115200
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "list-ports":
		handleListPorts()
	case "info":
		handleInfo(os.Args[2:])
	case "exec":
		handleExec(os.Args[2:])
	case "upgrade":
		handleUpgrade(os.Args[2:])
	case "help":
		if len(os.Args) > 2 {
			printCommandHelp(os.Args[2])
		} else {
			printUsage()
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func handleListPorts() {
	ports, err := serial.ListPorts()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing ports: %v\n", err)
		os.Exit(1)
	}

	if len(ports) == 0 {
		fmt.Println("No serial ports found!")
		return
	}

	fmt.Println("Available serial ports:")
	for _, port := range ports {
		fmt.Printf("  %s\n", port)
	}
}

func handleInfo(args []string) {
	fs := flag.NewFlagSet("info", flag.ExitOnError)
	port := fs.String("port", "", "Serial port path")
	portShort := fs.String("p", "", "Serial port path (short)")
	baud := fs.Int("baud", defaultBaudRate, "Baud rate")
	baudShort := fs.Int("b", defaultBaudRate, "Baud rate (short)")
	auto := fs.Bool("auto", false, "Auto detect serial port")
	autoShort := fs.Bool("a", false, "Auto detect serial port (short)")

	fs.Parse(args)

	portName := getStringFlag(*port, *portShort)
	baudRate := getIntFlag(*baud, *baudShort, defaultBaudRate)
	autoDetect := *auto || *autoShort

	if portName == "" && !autoDetect {
		autoDetect = true
	}

	if portName != "" && autoDetect {
		fmt.Fprintf(os.Stderr, "Error: --port and --auto cannot be used together\n")
		os.Exit(1)
	}

	var c *client.TjcDisplayClient

	if autoDetect {
		var err error
		c, err = autoDetectDevice(baudRate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Device found on port: %s\n\n", c.PortName)
	} else {
		c = &client.TjcDisplayClient{
			PortName: portName,
			BaudRate: baudRate,
		}
	}
	defer c.Close()

	info, err := c.GetDeviceInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting device info: %v\n", err)
		os.Exit(1)
	}

	printDeviceInfo(info)
}

func handleExec(args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: exec command requires a command string\n")
		fmt.Fprintf(os.Stderr, "Usage: tjs-serial-display exec <command> [-p|--port <port>] [-b|--baud <rate>] [-a|--auto]\n")
		os.Exit(1)
	}

	cmdString := args[0]
	fs := flag.NewFlagSet("exec", flag.ExitOnError)
	port := fs.String("port", "", "Serial port path")
	portShort := fs.String("p", "", "Serial port path (short)")
	baud := fs.Int("baud", defaultBaudRate, "Baud rate")
	baudShort := fs.Int("b", defaultBaudRate, "Baud rate (short)")
	auto := fs.Bool("auto", false, "Auto detect serial port")
	autoShort := fs.Bool("a", false, "Auto detect serial port (short)")

	fs.Parse(args[1:])

	portName := getStringFlag(*port, *portShort)
	baudRate := getIntFlag(*baud, *baudShort, defaultBaudRate)
	autoDetect := *auto || *autoShort

	if portName == "" && !autoDetect {
		autoDetect = true
	}

	if portName != "" && autoDetect {
		fmt.Fprintf(os.Stderr, "Error: --port and --auto cannot be used together\n")
		os.Exit(1)
	}

	var c *client.TjcDisplayClient

	if autoDetect {
		var err error
		c, err = autoDetectDevice(baudRate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		c = &client.TjcDisplayClient{
			PortName: portName,
			BaudRate: baudRate,
		}
	}
	defer c.Close()

	// 判断命令类型
	if strings.HasPrefix(cmdString, "print ") {
		// print 命令需要返回结果
		result, err := c.Prints(strings.TrimPrefix(cmdString, "print "))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(result)
	} else if cmdString == "sendme" {
		// sendme 返回当前页面
		page, err := c.GetPage()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Current page: %d\n", page)
	} else {
		// 其他命令直接执行
		err := c.ExecuteCommand(cmdString)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Command executed successfully")
	}
}

func handleUpgrade(args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: upgrade command requires a TFT file path\n")
		fmt.Fprintf(os.Stderr, "Usage: tjs-serial-display upgrade <tft_file> [-p|--port <port>] [-b|--baud <rate>] [-a|--auto]\n")
		os.Exit(1)
	}

	tftFile := args[0]
	fs := flag.NewFlagSet("upgrade", flag.ExitOnError)
	port := fs.String("port", "", "Serial port path")
	portShort := fs.String("p", "", "Serial port path (short)")
	baud := fs.Int("baud", defaultBaudRate, "Baud rate")
	baudShort := fs.Int("b", defaultBaudRate, "Baud rate (short)")
	auto := fs.Bool("auto", false, "Auto detect serial port")
	autoShort := fs.Bool("a", false, "Auto detect serial port (short)")

	fs.Parse(args[1:])

	portName := getStringFlag(*port, *portShort)
	baudRate := getIntFlag(*baud, *baudShort, defaultBaudRate)
	autoDetect := *auto || *autoShort

	if portName == "" && !autoDetect {
		autoDetect = true
	}

	if portName != "" && autoDetect {
		fmt.Fprintf(os.Stderr, "Error: --port and --auto cannot be used together\n")
		os.Exit(1)
	}

	var c *client.TjcDisplayClient

	if autoDetect {
		var err error
		c, err = autoDetectDevice(baudRate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Device found on port: %s\n", c.PortName)
	} else {
		c = &client.TjcDisplayClient{
			PortName: portName,
			BaudRate: baudRate,
		}
	}
	defer c.Close()

	fmt.Printf("Upgrading device with file: %s\n", tftFile)

	err := c.Upgrade(tftFile, func(progress *models.UpgradeProgress) {
		bar := progressBar(progress.Percentage, 50)
		speed := formatBytes(progress.Speed)
		fmt.Printf("\r[%s] %.1f%% (%s/%s) %s/s",
			bar,
			progress.Percentage,
			formatBytes(progress.Current),
			formatBytes(progress.Total),
			speed,
		)
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "\nUpgrade failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nUpgrade completed successfully!")
}

// autoDetectDevice 自动检测并连接设备
func autoDetectDevice(baudRate int) (*client.TjcDisplayClient, error) {
	ports, err := serial.ListPorts()
	if err != nil {
		return nil, fmt.Errorf("failed to list ports: %w", err)
	}

	if len(ports) == 0 {
		return nil, fmt.Errorf("no serial ports found")
	}

	for _, port := range ports {
		c := &client.TjcDisplayClient{
			PortName: port,
			BaudRate: baudRate,
		}

		// 尝试连接并获取设备信息
		_, err := c.GetDeviceInfo()
		if err == nil {
			return c, nil
		}
		c.Close()
	}

	return nil, fmt.Errorf("no TJC device found on any port")
}

// 辅助函数
func getStringFlag(flag1, flag2 string) string {
	if flag1 != "" {
		return flag1
	}
	return flag2
}

func getIntFlag(flag1, flag2, defaultValue int) int {
	if flag1 != defaultValue {
		return flag1
	}
	if flag2 != defaultValue {
		return flag2
	}
	return defaultValue
}

func printDeviceInfo(info *models.DeviceInfo) {
	screenTypes := map[int]string{
		0: "Non-touch screen",
		1: "Resistive screen",
		2: "Capacitive screen",
	}

	screenType := screenTypes[info.Type]
	if screenType == "" {
		screenType = fmt.Sprintf("Unknown (%d)", info.Type)
	}

	fmt.Println("Device Information:")
	fmt.Printf("  Screen Type:      %s\n", screenType)
	fmt.Printf("  Device Address:   %s\n", info.Address)
	fmt.Printf("  Model:            %s\n", info.Model)
	fmt.Printf("  Firmware Version: %d\n", info.FirmwareVersion)
	fmt.Printf("  MCU Number:       %d\n", info.MainControlChipNumber)
	fmt.Printf("  Serial Number:    %s\n", info.Number)
	fmt.Printf("  Flash Size:       %s\n", formatBytes(int64(info.FlashSize)))
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func progressBar(percentage float64, width int) string {
	filled := int(percentage / 100 * float64(width))
	if filled > width {
		filled = width
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}

func printUsage() {
	fmt.Println("TJS Serial Display - TJC Serial Display Controller")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  tjs-serial-display <command> [options]")
	fmt.Println()
	fmt.Println("Available Commands:")
	fmt.Println("  list-ports          List all available serial ports")
	fmt.Println("  info                Get device information")
	fmt.Println("  exec <command>      Execute TJC command")
	fmt.Println("  upgrade <file>      Upgrade device firmware")
	fmt.Println("  help [command]      Show help for a command")
	fmt.Println()
	fmt.Println("Global Options:")
	fmt.Println("  -p, --port <name>   Serial port path")
	fmt.Println("  -b, --baud <rate>   Baud rate (default: 115200)")
	fmt.Println("  -a, --auto          Auto detect serial port")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  tjs-serial-display list-ports")
	fmt.Println("  tjs-serial-display info --auto")
	fmt.Println("  tjs-serial-display exec \"page 2\" -p /dev/ttyUSB0")
	fmt.Println("  tjs-serial-display upgrade program.tft --auto")
	fmt.Println()
	fmt.Println("For more information, use: tjs-serial-display help <command>")
}

func printCommandHelp(command string) {
	switch command {
	case "list-ports":
		fmt.Println("Usage: tjs-serial-display list-ports")
		fmt.Println()
		fmt.Println("List all available serial ports on the system.")
	case "info":
		fmt.Println("Usage: tjs-serial-display info [-p|--port <port>] [-b|--baud <rate>] [-a|--auto]")
		fmt.Println()
		fmt.Println("Get information from the connected TJC display device.")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  -p, --port <name>   Serial port path (e.g., /dev/ttyUSB0)")
		fmt.Println("  -b, --baud <rate>   Baud rate (default: 115200)")
		fmt.Println("  -a, --auto          Auto detect device on available ports")
	case "exec":
		fmt.Println("Usage: tjs-serial-display exec <command> [-p|--port <port>] [-b|--baud <rate>] [-a|--auto]")
		fmt.Println()
		fmt.Println("Execute a TJC command on the display device.")
		fmt.Println()
		fmt.Println("Common Commands:")
		fmt.Println("  page <id>                 Jump to page")
		fmt.Println("  sendme                    Get current page ID")
		fmt.Println("  print <target>            Print target value")
		fmt.Println("  <comp>.txt=\"value\"        Set text value")
		fmt.Println("  vis <comp>,0              Hide component")
		fmt.Println("  vis <comp>,1              Show component")
		fmt.Println("  click <comp>,0            Click up")
		fmt.Println("  click <comp>,1            Click down")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  -p, --port <name>   Serial port path")
		fmt.Println("  -b, --baud <rate>   Baud rate (default: 115200)")
		fmt.Println("  -a, --auto          Auto detect device")
	case "upgrade":
		fmt.Println("Usage: tjs-serial-display upgrade <tft_file> [-p|--port <port>] [-b|--baud <rate>] [-a|--auto]")
		fmt.Println()
		fmt.Println("Upgrade the device firmware with a TFT file.")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  -p, --port <name>   Serial port path")
		fmt.Println("  -b, --baud <rate>   Baud rate (default: 115200)")
		fmt.Println("  -a, --auto          Auto detect device")
		fmt.Println()
		fmt.Println("Warning: Do not disconnect power during upgrade!")
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
	}
}
