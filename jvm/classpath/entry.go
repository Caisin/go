package classpath

import "os"
import "strings"

//路径分隔符
const pathListSeparator = string(os.PathListSeparator)

/**
可以把类路径想象成一个大的整体，
它由启
动类路径、
扩展类路径和
用户类路径三个小路径构成。
三个小路径又分别由更小的路径构成。
是不是很像组合模式（composite pattern）
用组合模式来设计和实现类路径。
*/
type Entry interface {
	readClass(className string) ([]byte, Entry, error)
	String() string
}

//newEntry（）函数根据参数创建不同类型的Entry实例
func newEntry(path string) Entry {
	if strings.Contains(path, pathListSeparator) {
		//综合
		return newCompositeEntry(path)
	}
	if strings.HasSuffix(path, "*") {
		//通配符
		return newWildcardEntry(path)
	}
	if strings.HasSuffix(path, ".jar") || strings.HasSuffix(path, ".JAR") || strings.HasSuffix(path, ".zip") || strings.HasSuffix(path, ".ZIP") {
		//jar
		return newZipEntry(path)
	}
	//路径
	return newDirEntry(path)
}
