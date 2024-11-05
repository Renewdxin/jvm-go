package main

import "flag"
import "fmt"
import "os"

// java [-options] class [args...]

// Cmd 结构体存储命令行参数。
type Cmd struct {
	helpFlag         bool     // -help 或 -? 选项，打印帮助信息
	versionFlag      bool     // -version 选项，打印版本信息并退出
	verboseClassFlag bool     // -verbose 或 -verbose:class 选项，启用类加载的详细输出
	verboseInstFlag  bool     // -verbose:inst 选项，启用指令执行的详细输出
	cpOption         string   // -classpath 或 -cp 选项，指定类路径
	XjreOption       string   // -Xjre 选项，指定JRE路径
	class            string   // 要执行的类名
	args             []string // 传递给main方法的参数
}

// parseCmd 解析命令行参数并返回 Cmd 结构体。
func parseCmd() *Cmd {
	cmd := &Cmd{}

	// 设置 usage 函数，在遇到无效参数时打印使用方法。
	flag.Usage = printUsage

	// 定义命令行选项。
	flag.BoolVar(&cmd.helpFlag, "help", false, "打印帮助信息")
	flag.BoolVar(&cmd.helpFlag, "?", false, "打印帮助信息")
	flag.BoolVar(&cmd.versionFlag, "version", false, "打印版本信息并退出")
	flag.BoolVar(&cmd.verboseClassFlag, "verbose", false, "启用详细输出（类加载信息）")       // 更明确的描述
	flag.BoolVar(&cmd.verboseClassFlag, "verbose:class", false, "启用详细输出（类加载信息）") // 同上
	flag.BoolVar(&cmd.verboseInstFlag, "verbose:inst", false, "启用详细输出（指令执行信息）")  // 指明是指令执行信息
	flag.StringVar(&cmd.cpOption, "classpath", "", "指定类路径")
	flag.StringVar(&cmd.cpOption, "cp", "", "指定类路径")
	flag.StringVar(&cmd.XjreOption, "Xjre", "", "指定JRE路径")

	// 解析命令行选项。
	flag.Parse()

	// 获取非选项参数（类名和程序参数）。
	args := flag.Args()
	if len(args) > 0 {
		cmd.class = args[0] // 第一个非选项参数是类名
		cmd.args = args[1:] // 后续的非选项参数是程序参数
	}

	return cmd
}

// printUsage 打印使用方法。
func printUsage() {
	fmt.Printf("Usage: %s [-options] class [args...]\n", os.Args[0])
	// flag.PrintDefaults()  // 可以选择取消注释，打印详细的选项说明
}
