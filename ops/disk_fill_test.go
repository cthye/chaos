package ops

import (
	"testing"
)

func TestRunDiskFill(t *testing.T) {
	var params = map[string]interface{}{
		"Path":         "~/Desktop",
		"Size":         "",
		"Reserve":      "1024",
		"Percent":      "90",
		"RetainHandle": false,
		"Timeout":      "10",
	}

	op := DiskFillOp{}
	IsOp(op)

	err := op.Validate(params)
	if err != nil {
		t.Error(err)
	}
	_, _, err = op.Run(params)
	if err != nil {
		t.Errorf("execute disk fill failed: %s", err)
	}
}

func TestDestroyDiskFill(t *testing.T) {
	var params = map[string]interface{}{
		"Path":         "~/Desktop",
		"Size":         "40000",
		"Reserve":      "1024",
		"Percent":      "",
		"RetainHandle": false,
		"Timeout":      "10",
	}

	op := DiskFillOp{}
	IsOp(op)

	err := op.Validate(params)
	if err != nil {
		t.Error(err)
	}
	uid, _, err := op.Run(params)
	if err != nil {
		t.Errorf("execute disk fill failed: %s", err)
	}

	_, err = op.Destroy(uid)
	if err != nil {
		t.Errorf("execute disk fill failed: %s", err)
	}
}
