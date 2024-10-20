package main

import (
	"fmt"
	"jvm-go/classfile"
  "jvm-go/classpath"
)


func main() {
  cmd := parseCmd()
  if cmd.versionFlag {
    fmt.Println("version 0.0.1")
  } else if cmd.helpFlag || cmd.class == "" {
    printUsage()
  } else {
    startJVM(cmd)
  }
}


func loadClass(className string, cp *classpath.ClassPath) *classfile.ClassFile {
  classData, _, err := cp.ReadClass(className)
  if err != nil {
    panic(err)
  }
  cf, err := classfile.Parse(classData)
  if err != nil {
    panic(err)
  }
  return cf
}


func printClassInfo(cf *classfile.ClassFile) {
  fmt.Printf("version: %v.%v\n", cf.MajorVersion(), cf.MinorVersion())
  fmt.Printf("constants count: %v\n", len(cf.ConstantPool()))
  fmt.Printf("access flags: 0x%x\n", cf.AccessFlags())
  fmt.Printf("this class: %v\n", cf.ClassName())
  fmt.Printf("super class: %v\n", cf.SuperClassName())
  fmt.Printf("interfaces: %v\n", cf.InterfaceNames())
  fmt.Printf("fields count: %v\n", len(cf.Fields()))
  for _, f := range cf.Fields() {
    fmt.Printf("   %s\n", f.Name())
  }
  fmt.Printf("methods count: %v\n", len(cf.Methods()))
  for _, m := range cf.Methods() {
    fmt.Printf("   %s\n", m.Name())
  }
}

  