package classpath

import "archive/zip"
import "errors"
import "io/ioutil"
import "path/filepath"

// 压缩包类路径
type ZipEntry struct {
	// 压缩包的绝对路径
	absPath string
	// 压缩包的读取器
	zipRC *zip.ReadCloser
}

// 创建一个ZipEntry对象
func newZipEntry(path string) *ZipEntry {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	// 返回一个ZipEntry对象，其中absPath是压缩包的绝对路径，zipRC是压缩包的读取器
	return &ZipEntry{absPath, nil}
}

// 读取类文件
func (zipE *ZipEntry) readClass(className string) ([]byte, Entry, error) {
	if zipE.zipRC == nil {
		err := zipE.openJar()
		if err != nil {
			return nil, nil, err
		}
	}

	classFile := zipE.findClass(className)
	if classFile == nil {
		return nil, nil, errors.New("class not found: " + className)
	}

	data, err := readClass(classFile)
	return data, zipE, err
}

// todo: close zip
func (zipE *ZipEntry) openJar() error {
	r, err := zip.OpenReader(zipE.absPath)
	if err == nil {
		zipE.zipRC = r
	}
	return err
}

func (zipE *ZipEntry) findClass(className string) *zip.File {
	for _, f := range zipE.zipRC.File {
		if f.Name == className {
			return f
		}
	}
	return nil
}

func readClass(classFile *zip.File) ([]byte, error) {
	rc, err := classFile.Open()
	if err != nil {
		return nil, err
	}
	// read class data
	data, err := ioutil.ReadAll(rc)
	rc.Close()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (zipE *ZipEntry) String() string {
	return zipE.absPath
}
