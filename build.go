package main

import (
	"fmt"
	"os"
	"time"

	"github.com/javapack-go/pack"
	"github.com/javapack-go/sourcefile"
)

func main() {
	args := os.Args
	pwd, _ := os.Getwd()
	//pwdSlice := strings.Split(pwd, ":")
	//pwdRoot := pwdSlice[0] + ":"
	vb := &sourcefile.VisitorBuilder{}
	visitor, _ := vb.Build(args)
	// ======================================================================
	before := time.Now().UnixNano()
	visitor.Process()
	visitor.ProcessClassFile()
	after := time.Now().UnixNano()
	fmt.Println(after - before)

	execCmd := &pack.ExecCmd{CmdProcess: "cmd.exe", CmdOps: "/k", ExecDir: pwd}
	execCmd.ExecutCmd("")

	c := &pack.Compile{
		FileMap:        visitor.FileMap,
		SourceFileStr:  visitor.SourceFileStr,
		LibSuffix:      visitor.LibSuffix,
		CompileCmdName: "javac",
		ProjPrefix:     visitor.ProjPrefix,
		//cmd命令 /c为执行完命令后自动关闭窗口，/k不关闭
		ExecCmd:        execCmd,
		ClassFileSlice: visitor.ClassFileSlice,
		ClassFileStr:   visitor.ClassFileStr,
	}
	var cco pack.CompileCmdOps
	var cerr error
	if cco, cerr = c.PreCompile(); cerr != nil {
		fmt.Errorf("build common compile command error: %s", cerr.Error())
		os.Exit(0)
	}
	c.Compile(cco)
}
