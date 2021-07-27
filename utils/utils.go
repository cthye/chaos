package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"log"
	"nessaj/config"
	"os"
	"os/exec"
	"strings"
	"time"
)

var BladeBinPath string

func Init(conf *config.Config) error {
	BladeBinPath = conf.ChaosbladeBin
	return nil
}

func ChkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

func GetBladePath() string {
	return BladeBinPath
}

func DebugMsg(format string, args ...interface{}) {
	msg := fmt.Sprintf("%v ", time.Now().Format("15:04:05.00000"))
	log.Printf(msg+strings.Trim(format, "\r\n ")+"\n", args...)
}

func DecodeParams(param map[string]interface{}, params interface{}) error {
	config := &mapstructure.DecoderConfig{
		ErrorUnused: true,
	}
	config.Result = params
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(param)
}

func GetUID(out []byte) (string, error) {
	var res map[string]interface{}

	err := json.Unmarshal(out, &res)
	if err != nil {
		return "", err
	}

	uid := res["result"].(string)
	return uid, nil
}

func RunExec(args []string, timeout string) (bytes.Buffer, bytes.Buffer, error) {
	binPath := GetBladePath()
	cmd := exec.Command(binPath, args...)
	if timeout != "" {
		duration, _ := time.ParseDuration(timeout + "s")
		ctx, cancel := context.WithTimeout(context.Background(), duration)
		defer cancel()
		cmd = exec.CommandContext(ctx, binPath, args...)
	}

	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return out, stderr, fmt.Errorf("excecute command failed, %v, %v, %v", err, out.String(), stderr.String())
	}
	return out, stderr, nil
}

func DestroyExec(id string) (bytes.Buffer, error) {
	binPath := GetBladePath()
	args := []string{"destroy", id}
	cmd := exec.Command(binPath, args...)

	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return out, fmt.Errorf("kill command failed, %v, %v, %v", err, out.String(), stderr.String())
	}
	return out, nil
}
