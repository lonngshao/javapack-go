package sourcefile

import "strings"
import "os"
import "fmt"
import "bytes"

type VisitorBuilder struct {
}

func (vb *VisitorBuilder) Build(args []string) (*Visitor, error) {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Errorf("get current runtime path error : %s", err.Error())
		return nil, err
	}
	//初始化一个遍历器的默认值
	visitor := &Visitor{RootPath: pwd,
		ProjPrefix:     "",
		LibSuffix:      "lib",
		FileMap:        make(map[string]string),
		ClassFileSlice: make([][]byte, 0, 50),
		ClassFileStr:   new(bytes.Buffer),
		SourceFileStr:  new(bytes.Buffer)}
	// =======================================================================
	var rootPath = ""
	var projPrefix = ""
	var libSuffix = ""
	for _, value := range args {
		if strings.Contains(value, "=") {
			argSlice := strings.Split(value, "=")
			if argSlice[0] == ROOTPATH {
				rootPath = argSlice[1]
			} else if argSlice[0] == PROJPREFIX {
				projPrefix = argSlice[1]
			} else if argSlice[0] == LIBSUFFIX {
				libSuffix = argSlice[1]
			}
		}
	}
	if len(rootPath) > 0 {
		visitor.RootPath = rootPath
	}
	if len(projPrefix) > 0 {
		visitor.ProjPrefix = projPrefix
	}
	if len(libSuffix) > 0 {
		visitor.LibSuffix = libSuffix
	}

	return visitor, nil
}
