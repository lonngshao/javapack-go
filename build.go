package main

import (
	"bytes"
	"fmt"
	"lsn/packjavawithgo/pack"
	"lsn/packjavawithgo/sourcefile"
	"os"
	"time"
)

func Slicetest() {
	var a []string = make([]string, 0, 50)
	var b [][]byte = make([][]byte, 0, 50)
	b = append(b, []byte("4"))
	b = append(b, []byte("5"))
	// a := []string{"1", "2", "3"}
	a = append(a, "4")
	fmt.Println(a)
	fmt.Println(string(bytes.Join(b, []byte(";"))))

}

func main() {
	// Slicetest()
	args := os.Args
	vb := &sourcefile.VisitorBuilder{}
	visitor, _ := vb.Build(args)
	// ======================================================================
	before := time.Now().UnixNano()
	visitor.Process()
	visitor.ProcessClassFile()
	after := time.Now().UnixNano()
	fmt.Println(after - before)

	c := &pack.Compile{
		FileMap:        visitor.FileMap,
		SourceFileStr:  visitor.SourceFileStr,
		LibSuffix:      visitor.LibSuffix,
		CompileCmdName: "javac",
		ProjPrefix:     visitor.ProjPrefix,
		ExecCmd:        &pack.ExecCmd{CmdName: "cmd.exe", CmdOps: "/c"},
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
