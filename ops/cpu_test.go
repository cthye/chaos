package ops

import (
	"testing"
)

func TestRunCpu(t *testing.T) {
	// TODO: verify cpu usage and timeout (top?)

	var params = map[string]interface{}{
		"cpuCount":   "0",
		"cpuPercent": "0",
		"cpuList":    "",
		"timeout":    "",
	}

	op := CpuOp{}
	IsOp(op)

	err := op.Validate(params)
	if err != nil {
		t.Error(err)
	}
	_, _, err = op.Run(params)
	if err != nil {
		t.Errorf("execute cpu burn failed: %s", err)
	}

	params = map[string]interface{}{
		"cpuCount":   "2",
		"cpuPercent": "50",
		"cpuList":    "",
		"timeout":    "",
	}

	err = op.Validate(params)
	if err != nil {
		t.Error(err)
	}
	_, _, err = op.Run(params)
	if err != nil {
		t.Errorf("execute cpu burn failed: %s", err)
	}

	err = op.Validate(params)
	if err != nil {
		t.Error(err)
	}
	params = map[string]interface{}{
		"cpuCount":   "1",
		"cpuPercent": "50",
		"cpuList":    "1,3",
		"timeout":    "100",
	}

	_, _, err = op.Run(params)
	if err != nil {
		t.Errorf("execute cpu burn failed: %s", err)
	}
}

func TestDestroyCpu(t *testing.T) {
	var params = map[string]interface{}{
		"cpuCount":   "2",
		"cpuPercent": "50",
		"cpuList":    "1,2",
		"timeout":    "100",
	}

	op := CpuOp{}
	IsOp(op)

	err := op.Validate(params)
	if err != nil {
		t.Error(err)
	}
	uid, _, err := op.Run(params)
	if err != nil {
		t.Errorf("execute cpu burn failed: %s", err)
	}

	t.Log(uid)

	_, err = op.Destroy(uid)
	if err != nil {
		t.Errorf("stop cpu burn failed: %s", err)
	}
}
