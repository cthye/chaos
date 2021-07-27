package ops

import (
	"testing"
	"github.com/mitchellh/mapstructure"
)

type TestOp struct {
	param TestOpParam
}

func (op TestOp) Name() string {
	return "test op"
}

func (op TestOp) Desc() string {
	return "test desc"
}

func (op TestOp) Validate(params map[string]interface{}) error {
	var param TestOpParam
	err := mapstructure.Decode(params, &param)
	if err != nil {
		return err
	}
	return nil
}

type TestOpParam struct {
	param1 int
	param2 string
	param3 float32
	param4 bool
}


func (op TestOp) Params() ParamSpec {
	return make(ParamSpec)
}

func (op TestOp) Init() error {
	return nil
}

func (op TestOp) Run(params map[string]interface{}) (string, string, error) {
	return "", "", nil
}

func (op TestOp) Destroy(id string) (string, error) {
	return "", nil
}

func Test_CheckParam(t *testing.T) {
	var params = make(map[string]interface{})
	params["param1"] = 1
	params["param2"] = "asdf"
	params["param3"] = 3.3
	params["param4"] = true

	op := TestOp{}
	IsOp(op)

	err := op.Validate(params)
	if err != nil {
		t.Error(err)
	}

	params = make(map[string]interface{})
	params["param2"] = "asdf"
	params["param3"] = 3.3
	params["param4"] = true

	err = op.Validate(params)
	if err != nil {
		t.Error(err)
	}
}
