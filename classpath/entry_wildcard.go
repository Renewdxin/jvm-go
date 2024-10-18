package classpath

import "os"
import "path/filepath"
import "strings"

// 通配符类路径
func newWildcardEntry(path string) CompositeEntry {
	baseDir := path[:len(path)-1] // remove *
	compositeEntry := []Entry{}

	// 遍历baseDir目录下的所有文件和目录
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 跳过baseDir目录
		if info.IsDir() && path != baseDir {
			return filepath.SkipDir
		}
		// 如果文件名以.jar或.JAR结尾，则创建一个ZipEntry对象并将其添加到compositeEntry中
		if strings.HasSuffix(path, ".jar") || strings.HasSuffix(path, ".JAR") {
			jarEntry := newZipEntry(path)
			compositeEntry = append(compositeEntry, jarEntry)
		}
		return nil
	}
	// 遍历baseDir目录下的所有文件和目录
	filepath.Walk(baseDir, walkFn)

	return compositeEntry
}
