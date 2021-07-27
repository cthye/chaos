package ops

import (
	"errors"
	"strings"
)

type ParamType string

const (
	PInt    ParamType = "int"
	PBool             = "bool"
	PString           = "string"
	PFloat            = "float64"
)

var AllOps = make(map[string]Op)

type Param struct {
	Iden     string    `json:"iden"`
	Desc     string    `json:"desc"`
	Type     ParamType `json:"type"`
	Required bool      `json:"required"`
}

type ParamSpec = map[string]Param

type Op interface {
	Name() string
	Desc() string
	Params() ParamSpec
	Init() error
	Validate(map[string]interface{}) error
	Run(map[string]interface{}) (string, string, error) // Validate must be called before Run
	Destroy(id string) (string, error)
}

func IsOp(op Op) bool {
	return true
}

func registerOp(op Op) error {
	err := op.Init()
	if err != nil {
		return err
	}
	AllOps[op.Name()] = op
	return nil
}

func Init() error {
	var errs []string
	for _, op := range []Op{EchoOp{}, CpuOp{}, MemOp{}, DiskBurnOp{},
		DiskFillOp{}, NetDelayOp{}, NetDNSOp{}, NetLossOp{}, ProcKillOp{}} {
		err := registerOp(op)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) != 0 {
		errStr := strings.Join(errs, ", ")
		return errors.New(errStr)
	} else {
		return nil
	}
}
