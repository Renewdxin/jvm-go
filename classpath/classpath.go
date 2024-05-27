package classpath

import (
    "os"
    "path/filepath"
)

type ClassPath struct {
	bootClassPath Entry
	extClassPath Entry
	userClassPath Entry
}

func Parse(jreOption, cpOption string) *ClassPath {
	cp := &ClassPath{}
	cp.parseBootAndExtClassPath(jreOption)
	cp.parseUserClassPath(cpOption)
	return cp
}

func (self *ClassPath) parseBootAndExtClassPath(jreOption string) {
	jreDir := getJreDir(jreOption)
	
	// jre/lib/*
	jreLibPath := filepath.Join(jreDir, "lib", "*")

	self.bootClassPath = newWildcardEntry(jreLibPath)
	
	// jre/lib/ext/*
	jreExtPath := filepath.Join(jreDir, "lib", "ext", "*")
	self.extClassPath = newWildcardEntry(jreExtPath)
}


func (self *ClassPath) parseUserClassPath(cpOption string) {
	// 如果没有，则使用当前目录
	if cpOption == "" {
	  cpOption = "."
	}
	self.userClassPath = newEntry(cpOption)
}

func getJreDir(jreOption string) string {
	// user first
	if jreOption != "" && exists(jreOption) {
	  return jreOption
	}
	// 当前目录下寻找
	if exists("./jre") {
		return "./jre"
	}
	// system JAVA_HOME
	if jh := os.Getenv("JAVA_HOME"); jh != "" {
		return filepath.Join(jh, "jre")
	}
	panic("Can not find jre folder! ")
}

func exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
	  if os.IsNotExist(err) {
		  return false
	  }
	}
	return true
}

// ReadClass()方法依次从启动类路径、扩展类路径和用户类路径中搜索class文件
func (self *ClassPath) ReadClass(className string) ([]byte, Entry, error) {
	className = className + ".class"
	if data, entry, err := self.bootClassPath.readClass(className); err == nil {
	  return data, entry, err
	}
	if data, entry, err := self.extClassPath.readClass(className); err == nil {
	  return data, entry, err
	}
	return self.userClassPath.readClass(className)
}

func (self *ClassPath) String() string {
	return self.userClassPath.String()
}