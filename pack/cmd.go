package pack

import "os/exec"
import "fmt"

type ExecCmd struct {
	CmdProcess string
	CmdOps     string
	ExecDir    string
}

func (ec *ExecCmd) ExecutCmd(cmdParams string) {
	cmd := exec.Command(ec.CmdProcess, ec.CmdOps, ec.ExecDir, cmdParams)
	if err := cmd.Run(); err != nil {
		fmt.Println("call cmd error:", err)
	}
}
