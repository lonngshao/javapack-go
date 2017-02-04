package pack

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

//编译命令常量
const (
	CmdName       string = "javac"
	BaseSourceDir string = "sourceLists"
)

//编译命令常用选项参数
type CompileCmdOps struct {
	Sourcepath string
	//类路径，依赖jar包路径
	Classpath string
	//输出编译器正在执行的操作
	Verbose string
	//指定存放生成的类文件的位置
	D string
	//编码
	Encoding string

	SourceFiles string
}

func (cco CompileCmdOps) toString() (string, error) {
	// return json.Marshal(cco)
	return cco.D + " " + cco.Verbose + " " + cco.Encoding + " " + cco.Classpath + " " + cco.Sourcepath + " " + cco.SourceFiles, nil
}

//编译执行时的公用数据
type Compile struct {
	FileMap        map[string]string
	SourceFileStr  *bytes.Buffer
	LibSuffix      string
	ProjPrefix     string
	CompileCmdName string
	ExecCmd        *ExecCmd
	ClassFileSlice [][]byte
	ClassFileStr   *bytes.Buffer
}

//编译预处理，构建编译命令的通用部分，构建编译源文件列表（蛋疼的javac只能编译单个文件或者源文件列表，草）
func (c *Compile) PreCompile() (CompileCmdOps, error) {
	var classpath = ""
	fm := c.FileMap
	if fm == nil {
		cco := CompileCmdOps{}
		return cco, errors.New("the file map to compile is nil")
	}
	// projLib := strings.Replace(c.ProjPrefix+c.LibSuffix, "--", "-", -1)
	// if _, ok := fm[projLib]; ok {
	// 	classpath = fm[projLib]
	// }
	//classpath只能精确到jar,构建并生成类文件命令列表
	// cfs := c.ClassFileSlice
	cfstr := c.ClassFileStr
	// classFile := bytes.Join(cfs, []byte(";"))
	classpath = "-classpath " + cfstr.String()
	verbose := "-verbose"
	encoding := "-encoding utf8"
	cco := CompileCmdOps{Sourcepath: "", Classpath: classpath, Verbose: verbose, Encoding: encoding}
	pwd, perr := os.Getwd()
	if perr != nil {
		return cco, perr
	}
	var fileSep = "/"
	var sourceDir = ""
	if strings.Contains(pwd, "\\") {
		fileSep = "\\"
	}
	if strings.LastIndex(pwd, "\\") == len(pwd)-1 || strings.LastIndex(pwd, "/") == len(pwd)-1 {
		sourceDir = pwd + BaseSourceDir
	} else {
		sourceDir = pwd + fileSep + BaseSourceDir
	}
	if _, oerr := os.Open(sourceDir); os.IsNotExist(oerr) {
		os.Mkdir(sourceDir, os.ModeDir)
	}
	//构建并生成源文件
	for v, k := range fm {
		var sourceFile = sourceDir + fileSep + v + ".txt"
		var out *os.File
		var oerr error
		os.Remove(sourceFile)
		if out, oerr = os.OpenFile(sourceFile, syscall.O_RDWR, os.ModeAppend); oerr != nil && os.IsNotExist(oerr) {
			if out, oerr = os.Create(sourceFile); oerr != nil {
				return cco, oerr
			}
			filepath.Walk(k, func(path string, info os.FileInfo, err error) error {
				return buildSources(out, path, err)
			})
		}

	}
	return cco, nil
}

func buildSources(out *os.File, path string, err error) error {
	if err != nil {
		return err
	}
	if strings.HasSuffix(path, ".java") {
		out.WriteString(path + "\r\n")
	}
	return err
}

//执行编译方法，需要根据sourcepath的不同路径修改CompileCmdOps值
func (c *Compile) Compile(cco CompileCmdOps) error {
	fMap := c.FileMap
	sfs := c.SourceFileStr
	for k, v := range fMap {
		if v == cco.Classpath {
			return nil
		}
		cco.Sourcepath = "-sourcepath " + sfs.String()
		var fileSep = "/"
		if strings.Contains(v, "\\") {
			fileSep = "\\"
		}
		distPath := v + fileSep + "classes"
		os.Remove(distPath)
		if _, perr := os.Stat(distPath); perr != nil && os.IsNotExist(perr) {
			os.Mkdir(distPath, os.ModeDir)
		}
		cco.D = "-d " + v + fileSep + "classes"
		var pwd string
		var perr error
		if pwd, perr = os.Getwd(); perr != nil {
			return perr
		}
		cco.SourceFiles = "@" + pwd + fileSep + BaseSourceDir + fileSep + k + ".txt"

		var ccoStr string
		var err error
		if ccoStr, err = cco.toString(); err != nil {
			return errors.New("translate cmd args error")
		}
		jc := c.ExecCmd
		fmt.Println("build exec command str :" + CmdName + " " + ccoStr)
		fmt.Println("==================================================")
		jc.executCmd(CmdName + " " + ccoStr)
	}
	return nil
}
