package classpath

import (
    "os"
    "path/filepath"
    "strings"
)

func newWildcardEntry(path string) CompositeEntry {
	baseDir := path[:len(path)-1]
	composite := []Entry{}
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != baseDir {
			return filepath.SkipDir
		}
		// 根据后缀名选出jar文件
		if strings.HasSuffix(path, ".jar") || strings.HasSuffix(path, ".JAR") {
			jarEn := newZipEntry(path)
			composite = append(composite, jarEn)
		}
		return nil
	}

	filepath.Walk(baseDir, walkFn)
	return composite
}