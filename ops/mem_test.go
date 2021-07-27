package ops

import (
	"testing"
)

func TestRunMem(t *testing.T) {
	var params = map[string]interface{}{
		"MemPercent": "50",
		"Mode":       "ram",
		"Reserve":    "200",
		"Rate":       "100",
		"Timeout":    "30",
	}

	op := MemOp{}
	IsOp(op)

	err := op.Validate(params)
	if err != nil {
		t.Error(err)
	}
	_, _, err = op.Run(params)
	if err != nil {
		t.Errorf("execute memory burn failed: %s", err)
	}

}

func TestDestroyMem(t *testing.T) {
	var params = map[string]interface{}{
		"MemPercent": "50",
		"Mode":       "ram",
		"Reserve":    "200",
		"Rate":       "100",
		"Timeout":    "30",
	}

	op := MemOp{}
	IsOp(op)

	err := op.Validate(params)
	if err != nil {
		t.Error(err)
	}
	uid, _, err := op.Run(params)
	if err != nil {
		t.Errorf("execute memory burn failed: %s", err)
	}

	_, err = op.Destroy(uid)
	if err != nil {
		t.Errorf("stop memory burn failed: %s", err)
	}
}
