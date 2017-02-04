package pack

import (
	"os/exec"
)

type ExecCmd struct {
	CmdName string
	CmdOps  string
}

func (ec *ExecCmd) executCmd(cmdParams string) {
	cmdName := ec.CmdName
	cmd := exec.Command(cmdName, ec.CmdOps, cmdParams)
	cmd.Run()
}
