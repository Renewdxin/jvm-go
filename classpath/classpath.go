package classpath

import "os"
import "path/filepath"

type Classpath struct {
	// 启动类路径
	bootClasspath Entry
	// 扩展类路径
	extClasspath Entry
	// 用户类路径
	userClasspath Entry
}

// Parse 解析启动类路径和扩展类路径
func Parse(jreOption, cpOption string) *Classpath {
	cp := &Classpath{}
	// 优先使用用户输入的-Xjre选项作为jre目录。如果没有输入该选项，则在当前目录下寻找jre目录
	cp.parseBootAndExtClasspath(jreOption)
	// 使用用户输入的-cp选项或-classpath选项作为用户类路径。如果没有输入该选项，则使用当前目录作为用户类路径
	cp.parseUserClasspath(cpOption)
	return cp
}

// 解析启动类路径和扩展类路径
func (cp *Classpath) parseBootAndExtClasspath(jreOption string) {
	// 获取jre目录
	jreDir := getJreDir(jreOption)

	// jre/lib/*
	jreLibPath := filepath.Join(jreDir, "lib", "*")
	cp.bootClasspath = newWildcardEntry(jreLibPath)

	// jre/lib/ext/*
	jreExtPath := filepath.Join(jreDir, "lib", "ext", "*")
	cp.extClasspath = newWildcardEntry(jreExtPath)
}

// 获取jre目录
func getJreDir(jreOption string) string {
	if jreOption != "" && exists(jreOption) {
		return jreOption
	}
	if exists("./jre") {
		return "./jre"
	}
	if jh := os.Getenv("JAVA_HOME"); jh != "" {
		return filepath.Join(jh, "jre")
	}
	panic("Can not find jre folder!")
}

// 判断路径是否存在
func exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (cp *Classpath) parseUserClasspath(cpOption string) {
	if cpOption == "" {
		cpOption = "."
	}
	cp.userClasspath = newEntry(cpOption)
}

// ReadClass
// className: fully/qualified/ClassName
func (cp *Classpath) ReadClass(className string) ([]byte, Entry, error) {
	className = className + ".class"
	if data, entry, err := cp.bootClasspath.readClass(className); err == nil {
		return data, entry, err
	}
	if data, entry, err := cp.extClasspath.readClass(className); err == nil {
		return data, entry, err
	}
	return cp.userClasspath.readClass(className)
}

func (cp *Classpath) String() string {
	return cp.userClasspath.String()
}
