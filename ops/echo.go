package ops

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"os"
	"strconv"
)

type EchoOp struct {
	param EchoParam
}

func (op EchoOp) Name() string {
	return "echo"
}

func (op EchoOp) Desc() string {
	return "this is poc and just return whatever param passed in"
}

type EchoParam struct {
	Txt string `mapstructure:"txt"`
}

func (op EchoOp) Validate(params map[string]interface{}) error {
	var param EchoParam
	err := mapstructure.Decode(params, &param)
	if err != nil {
		return err
	}
	if param.Txt == "" {
		return errors.New("empty txt value")
	}
	return nil
}

func (op EchoOp) Params() ParamSpec {
	param := make(ParamSpec)
	param["txt"] = Param{
		Iden:     "txt",
		Desc:     "echo content",
		Type:     PString,
		Required: true,
	}
	return param
}

func (op EchoOp) Init() error {
	return nil
}

func (op EchoOp) Run(params map[string]interface{}) (string, string, error) {
	var param EchoParam
	_ = mapstructure.Decode(params, &param)

	pid := strconv.Itoa(os.Getpid())

	return pid, param.Txt, nil
}

func (op EchoOp) Destroy(pid string) (string, error) {
	return "not applicable", nil
}
