package ops

import (
	"nessaj/utils"
)

type CpuOp struct {
	param CpuParam
}

type CpuParam struct {
	CpuCount   string `json:"cpuCount"`
	CpuPercent string `json:"cpuPercent"`
	CpuList    string `json:"cpuList"`
	Timeout    string `json:"timeout"`
}

func (op CpuOp) Name() string {
	return "cpu"
}

func (op CpuOp) Desc() string {
	return "params for burncpu: --cpu-count --cpu-percent"
}

func (op CpuOp) Validate(params map[string]interface{}) error {
	var param CpuParam
	err := utils.DecodeParams(params, &param)
	if err != nil {
		return err
	}
	return nil
}

func (op CpuOp) Params() ParamSpec {
	param := make(ParamSpec)

	param["CpuCount"] = Param{
		Iden:     "CpuCount",
		Desc:     "the number of cpus",
		Type:     PString,
		Required: false,
	}

	param["CpuPercent"] = Param{
		Iden:     "CpuPercent",
		Desc:     "percent of cpu burned",
		Type:     PString,
		Required: false,
	}

	param["CpuList"] = Param{
		Iden:     "CpuList",
		Desc:     "cpus allowed burning",
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

func (op CpuOp) Init() error {
	return nil
}

func (op CpuOp) Run(params map[string]interface{}) (string, string, error) {
	// get parameters
	var param CpuParam
	err := utils.DecodeParams(params, &param)
	if err != nil {
		return "", "", err
	}

	args := []string{"create", "cpu", "load"}
	if param.CpuCount != "" {
		args = append(args, "--cpu-count", param.CpuCount)
	}
	if param.CpuPercent != "" {
		args = append(args, "--cpu-percent", param.CpuPercent)
	}
	if param.CpuList != "" {
		args = append(args, "--cpu-list", param.CpuList)
	}
	if param.Timeout != "" {
		args = append(args, "--timeout", param.Timeout)
	}

	// run executable
	out, _, err := utils.RunExec(args, param.Timeout)

	//get uid
	uid, err := utils.GetUID(out.Bytes())
	if err != nil {
		return "", "", err
	}
	return uid, out.String(), err
}

func (op CpuOp) Destroy(id string) (string, error) {
	// using blade destroy cmd
	out, err := utils.DestroyExec(id)
	return out.String(), err
}
