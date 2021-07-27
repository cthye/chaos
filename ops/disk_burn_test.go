package ops

import (
	"testing"
)

func TestRunDiskBurn(t *testing.T) {
	var params = map[string]interface{}{
		// ? why can't relative path to ~
		"Path":    "/home/cthye/Desktop",
		"Size":    "500",
		"Write":   false,
		"Read":    true,
		"Timeout": "",
	}

	op := DiskBurnOp{}
	IsOp(op)

	err := op.Validate(params)
	if err != nil {
		t.Error(err)
	}
	uid, _, err := op.Run(params)
	if err != nil {
		t.Errorf("execute disk burn failed: %s", err)
	}

	params = map[string]interface{}{
		"Path":    "",
		"Size":    "500",
		"Write":   true,
		"Read":    false,
		"Timeout": "10",
	}

	err = op.Validate(params)
	if err != nil {
		t.Error(err)
	}
	uid, _, err = op.Run(params)
	if err != nil {
		t.Errorf("execute disk burn failed: %s", err)
	}

	_, err = op.Destroy(uid)
	if err != nil {
		t.Errorf("execute disk burn failed: %s", err)
	}
}
