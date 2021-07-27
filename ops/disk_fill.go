package ops

import (
	"errors"
	"nessaj/utils"
)

type DiskFillOp struct {
	param DiskFillParam
}

type DiskFillParam struct {
	Path         string
	Size         string
	Reserve      string
	Percent      string
	RetainHandle bool
	Timeout      string
}

func (op DiskFillOp) Name() string {
	return "diskfill"
}

func (op DiskFillOp) Desc() string {
	return "params for fill disk"
}

func (op DiskFillOp) Validate(params map[string]interface{}) error {
	var param DiskFillParam
	err := utils.DecodeParams(params, &param)
	if err != nil {
		return err
	}
	return nil
}

func (op DiskFillOp) Params() ParamSpec {
	param := make(ParamSpec)

	param["Path"] = Param{
		Iden:     "Path",
		Desc:     "the path of directory to run test (default /)",
		Type:     PString,
		Required: false, // not necessary but strongly recommended
	}

	param["Size"] = Param{
		Iden:     "Size",
		Desc:     "the size of disk space used to run test(default 10M)", // size or percent
		Type:     PString,
		Required: false,
	}

	param["Reserve"] = Param{
		Iden:     "Reserve",
		Desc:     "the size of disk space reserved",
		Type:     PString,
		Required: false,
	}

	param["Percent"] = Param{
		Iden:     "Percent",
		Desc:     "the percent of disk usage", // size or percent
		Type:     PString,
		Required: false,
	}

	param["RetainHandle"] = Param{
		Iden:     "RetainHandle",
		Desc:     "whether to save the content filled in disk",
		Type:     PBool,
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

func (op DiskFillOp) Init() error {
	return nil
}

func (op DiskFillOp) Run(params map[string]interface{}) (string, string, error) {
	// get parameters
	var param DiskFillParam
	err := utils.DecodeParams(params, &param)
	if err != nil {
		return "", "", errors.New("mapstructure")
	}

	args := []string{"create", "disk", "fill"}
	if param.Path != "" {
		args = append(args, "--path", param.Path)
	}
	if param.Size != "" {
		args = append(args, "--size", param.Size)
	}
	if param.Reserve != "" {
		args = append(args, "--reserve", param.Reserve)
	}
	if param.Percent != "" {
		args = append(args, "--percent", param.Percent)
	}
	if param.RetainHandle {
		args = append(args, "--retain-handle")
	}
	if param.Timeout != "" {
		args = append(args, "--timeout", param.Timeout)
	}

	// run executable
	out, _, err := utils.RunExec(args, param.Timeout)

	// get uid
	uid, err := utils.GetUID(out.Bytes())
	if err != nil {
		return "", "", err
	}

	return uid, out.String(), err
}

func (op DiskFillOp) Destroy(id string) (string, error) {
	// using blade destroy cmd
	out, err := utils.DestroyExec(id)
	return out.String(), err
}
