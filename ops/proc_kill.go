package ops

import (
	"nessaj/utils"
)

type ProcKillOp struct {
	param ProcKillParam
}

type ProcKillParam struct {
	Process    string
	ProcessCmd string
	Count      string
	Signal     string
	Timeout    string
}

func (op ProcKillOp) Name() string {
	return "processkill"
}

func (op ProcKillOp) Desc() string {
	return "test system resilience when processes non-exsisted"
}

func (op ProcKillOp) Validate(params map[string]interface{}) error {
	var param ProcKillParam
	err := utils.DecodeParams(params, &param)
	if err != nil {
		return err
	}
	return nil
}

func (op ProcKillOp) Params() ParamSpec {
	param := make(ParamSpec)

	param["Process"] = Param{
		Iden:     "Process",
		Desc:     "proces keyword",
		Type:     PString,
		Required: false,
	}

	param["ProcessCmd"] = Param{
		Iden:     "ProcessCmd",
		Desc:     "process command",
		Type:     PString,
		Required: false,
	}

	param["Count"] = Param{
		Iden:     "Count",
		Desc:     "the number of process to be killed (0 for unlimited)",
		Type:     PString,
		Required: false,
	}

	param["Signal"] = Param{
		Iden:     "Signal",
		Desc:     "default 9",
		Type:     PString,
		Required: false,
	}

	param["Timeout"] = Param{
		Iden:     "Timeout",
		Desc:     "runtime limitation",
		Type:     PString,
		Required: false,
	}

	return param
}

func (op ProcKillOp) Init() error {
	return nil
}

func (op ProcKillOp) Run(params map[string]interface{}) (string, string, error) {
	// get parameters
	var param ProcKillParam
	err := utils.DecodeParams(params, &param)
	if err != nil {
		return "", "", err
	}

	args := []string{"create", "process", "kill"}
	if param.Process != "" {
		args = append(args, "--process", param.Process)
	}
	if param.ProcessCmd != "" {
		args = append(args, "--process-cmd", param.ProcessCmd)
	}
	if param.Count != "" {
		args = append(args, "--count", param.Count)
	}
	if param.Signal != "" {
		args = append(args, "--signal", param.Signal)
	}
	if param.Timeout != "" {
		args = append(args, "--timeout", param.Timeout)
	}

	// run executable
	out, _, err := utils.RunExec(args, param.Timeout)
	if err != nil {
		return "", "", err
	}

	//get uid
	uid, err := utils.GetUID(out.Bytes())
	if err != nil {
		return "", "", err
	}
	return uid, out.String(), err
}

func (op ProcKillOp) Destroy(id string) (string, error) {
	// using blade destroy cmd
	out, err := utils.DestroyExec(id)
	return out.String(), err
}
