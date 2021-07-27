package ops

import (
	"errors"
	"nessaj/utils"
)

type DiskBurnOp struct {
	param DiskBurnParam
}

type DiskBurnParam struct {
	Path    string
	Size    string
	Write   bool
	Read    bool
	Timeout string
}

func (op DiskBurnOp) Name() string {
	return "diskburn"
}

func (op DiskBurnOp) Desc() string {
	return "params for burn disk"
}

func (op DiskBurnOp) Validate(params map[string]interface{}) error {
	var param DiskBurnParam
	err := utils.DecodeParams(params, &param)

	if err != nil {
		return err
	}
	return nil
}

func (op DiskBurnOp) Params() ParamSpec {
	param := make(ParamSpec)

	param["Path"] = Param{
		Iden:     "Path",
		Desc:     "the path of directory to run test (default /)",
		Type:     PString,
		Required: false, // not necessary but strongly recommended
	}

	param["Size"] = Param{
		Iden:     "Size",
		Desc:     "the size of disk space used to run test(default 10M)",
		Type:     PString,
		Required: false,
	}

	param["Write"] = Param{
		Iden:     "Write",
		Desc:     "write disk",
		Type:     PBool,
		Required: false,
	}

	param["Read"] = Param{
		Iden:     "Read",
		Desc:     "read disk",
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

func (op DiskBurnOp) Init() error {
	return nil
}

func (op DiskBurnOp) Run(params map[string]interface{}) (string, string, error) {
	// get parameters
	var param DiskBurnParam
	err := utils.DecodeParams(params, &param)

	if err != nil {
		return "", "", errors.New("mapstructure")
	}

	args := []string{"create", "disk", "burn"}
	if param.Read {
		args = append(args, "--read")
	} else {
		args = append(args, "--write")
	}
	if param.Path != "" {
		args = append(args, "--path", param.Path)
	}
	if param.Size != "" {
		args = append(args, "--size", param.Size)
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

func (op DiskBurnOp) Destroy(id string) (string, error) {
	// using blade destroy cmd
	out, err := utils.DestroyExec(id)
	return out.String(), err
}
