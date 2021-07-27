package ops

import (
	"errors"
	"nessaj/utils"
)

type MemOp struct {
	param MemParam
}

type MemParam struct {
	MemPercent string
	Mode       string
	Reserve    string
	Rate       string
	Timeout    string
}

func (op MemOp) Name() string {
	return "memory"
}

func (op MemOp) Desc() string {
	return "params for burn memory: --mem-percent --mode --reserve --rate --timeout"
}

func (op MemOp) Validate(params map[string]interface{}) error {
	var param MemParam
	err := utils.DecodeParams(params, &param)
	if err != nil {
		return err
	}
	return nil
}

func (op MemOp) Params() ParamSpec {
	param := make(ParamSpec)

	param["MemPercent"] = Param{
		Iden:     "MemPercent",
		Desc:     "the percent of ocupied memory",
		Type:     PString,
		Required: true, // not necessary but strongly recommended
	}

	param["Mode"] = Param{
		Iden:     "Mode",
		Desc:     "ram or cache (default cache)",
		Type:     PString,
		Required: false,
	}

	param["Reserve"] = Param{
		Iden:     "Reserve",
		Desc:     "memory size reserved for test (use mem-percent in higher priority if exists)",
		Type:     PString,
		Required: false,
	}

	param["Rate"] = Param{
		Iden:     "Rate",
		Desc:     "memory usage rate (only works in ram mode)",
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

func (op MemOp) Init() error {
	return nil
}

func (op MemOp) Run(params map[string]interface{}) (string, string, error) {
	// get parameters
	var param MemParam
	err := utils.DecodeParams(params, &param)
	if err != nil {
		return "", "", errors.New("mapstructure")
	}

	args := []string{"create", "mem", "load"}
	if param.MemPercent != "" {
		args = append(args, "--mem-percent", param.MemPercent)
	} else {
		// to avoid crashing, set default: 50%
		args = append(args, "--mem-percent", "50")
	}
	if param.Mode != "" {
		args = append(args, "--mode", param.Mode)
	}
	if param.Reserve != "" {
		args = append(args, "--reserve", param.Reserve)
	}
	if param.Rate != "" {
		args = append(args, "--rate", param.Rate)
	}
	if param.Timeout != "" {
		args = append(args, "--timeout", param.Timeout)
	}

	// run excutable
	out, _, err := utils.RunExec(args, param.Timeout)

	// get uid
	uid, err := utils.GetUID(out.Bytes())
	if err != nil {
		return "", "", err
	}

	return uid, out.String(), err
}

func (op MemOp) Destroy(id string) (string, error) {
	// using blade destroy cmd
	out, err := utils.DestroyExec(id)
	return out.String(), err
}
