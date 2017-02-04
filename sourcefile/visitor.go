package sourcefile

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

// 打包项目路径常量
const (
	//打包项目根路径
	ROOTPATH   string = "rootpath"
	PROJPREFIX string = "projprefix"
	LIBSUFFIX  string = "libsuffix"
)

// 目录文件遍历访问器
type Visitor struct {
	RootPath       string
	ProjPrefix     string
	LibSuffix      string
	FileMap        map[string]string
	SourceFileStr  *bytes.Buffer
	ClassFileSlice [][]byte
	ClassFileStr   *bytes.Buffer
}

// 具体处理遍历源文件目录文件方法
func (v *Visitor) Visit(path string, info os.FileInfo, err error) error {
	fm := v.FileMap
	sfs := v.SourceFileStr
	rootPath := v.RootPath
	projPrefix := v.ProjPrefix
	// libSuffix := v.LibSuffix
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return nil
	}
	//最后一个分隔符之前的部分,应该==rootpath
	pathDir := filepath.Dir(path)
	//最后一个分隔符之后的部分,应该以projPrefix为前缀
	pathBase := filepath.Base(path)
	noSepPathDir := strings.Replace(strings.Replace(pathDir, "\\", ".", -1), "/", ".", -1)
	noSepPathBase := strings.Replace(strings.Replace(pathBase, "\\", ".", -1), "/", ".", -1)
	noSepRootPath := strings.Replace(strings.Replace(rootPath, "\\", ".", -1), "/", ".", -1)
	if noSepPathDir != noSepRootPath || !strings.HasPrefix(noSepPathBase, projPrefix) {
		return nil
	}
	if _, ok := fm[info.Name()]; ok {
		return nil
	}
	fm[info.Name()] = path
	sfs.WriteString(path)
	sfs.WriteString(";")
	return nil
}

func (v *Visitor) VisitClassFile(path string, info os.FileInfo, err error) error {
	cfs := v.ClassFileSlice
	cfstr := v.ClassFileStr
	if err != nil {
		return err
	}
	if !info.IsDir() {
		jarName := info.Name()
		if strings.HasSuffix(jarName, ".jar") {
			cfs = append(cfs, []byte(path))
			cfstr.WriteString(path)
			cfstr.WriteString(";")
		}
	}
	return nil
}

func (v *Visitor) Process() {
	filepath.Walk(v.RootPath, func(path string, info os.FileInfo, err error) error {
		return v.Visit(path, info, err)
	})
}

func (v *Visitor) ProcessClassFile() {
	rootPath := v.RootPath
	projPrefix := v.ProjPrefix
	libSuffix := v.LibSuffix
	classPath := rootPath + "\\" + projPrefix + libSuffix
	filepath.Walk(classPath, func(path string, info os.FileInfo, err error) error {
		return v.VisitClassFile(path, info, err)
	})
}
