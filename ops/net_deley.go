package ops

import (
	"errors"
	"nessaj/utils"
)

type NetDelayOp struct {
	param NetDelayParam
}

type NetDelayParam struct {
	DesIP          string
	ExcludePort    string
	ExcludeIP      string
	Interface      string
	LocalPort      string
	Offset         string
	RemotePort     string
	Time           string
	Force          bool
	IgnorePeerPort bool
	Timeout        string
}

func (op NetDelayOp) Name() string {
	return "netdelay"
}

func (op NetDelayOp) Desc() string {
	return "test resilience for network delay"
}

func (op NetDelayOp) Validate(params map[string]interface{}) error {
	var param NetDelayParam
	err := utils.DecodeParams(params, &param)
	if err != nil {
		return err
	}
	return nil
}

func (op NetDelayOp) Params() ParamSpec {
	param := make(ParamSpec)

	param["DesIP"] = Param{
		Iden:     "DesIP",
		Desc:     "target ip (Support for mask, like 192.168.1.0/24. Support for mutiply IPs (splitted by ','): 192.168.1.1,192.168.2.1)",
		Type:     PString,
		Required: false,
	}

	param["ExcludePort"] = Param{
		Iden:     "ExcludePort",
		Desc:     "exclude the port for testing and keep it available",
		Type:     PString,
		Required: false,
	}

	param["ExcludeIP"] = Param{
		Iden:     "ExcludeIP",
		Desc:     "exclude the IP for testing and keep it available",
		Type:     PString,
		Required: false,
	}

	param["Interface"] = Param{
		Iden:     "Interface",
		Desc:     "network device",
		Type:     PString,
		Required: true,
	}

	param["LocalPort"] = Param{
		Iden:     "LocalPort",
		Desc:     "local port",
		Type:     PString,
		Required: false,
	}

	param["Offset"] = Param{
		Iden:     "Offset",
		Desc:     "delay time (ms)",
		Type:     PString,
		Required: true,
	}

	param["RemotePort"] = Param{
		Iden:     "RemotePort",
		Desc:     "remote port",
		Type:     PString,
		Required: false,
	}

	param["Force"] = Param{
		Iden:     "Force",
		Desc:     "force to substitute the current traffic control rule",
		Type:     PBool,
		Required: false,
	}

	param["Time"] = Param{
		Iden:     "Time",
		Desc:     "delay time",
		Type:     PBool,
		Required: true,
	}

	param["IgnorePeerPort"] = Param{
		Iden:     "IgnorePeerPort",
		Desc:     "针对添加 --exclude-port 参数，报 ss 命令找不到的情况下使用，忽略排除端口",
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

func (op NetDelayOp) Init() error {
	return nil
}

func (op NetDelayOp) Run(params map[string]interface{}) (string, string, error) {
	// get parameters
	var param NetDelayParam
	err := utils.DecodeParams(params, &param)
	if err != nil {
		return "", "", err
	}

	args := []string{"create", "network", "delay"}
	if param.DesIP != "" {
		args = append(args, "--destination-ip", param.DesIP)
	}
	if param.ExcludePort != "" {
		args = append(args, "--exclude-port", param.ExcludePort)
	}
	if param.ExcludeIP != "" {
		args = append(args, "--exclude-ip", param.ExcludeIP)
	}
	if param.Interface != "" {
		args = append(args, "--interface", param.Interface)
	} else {
		// * can be ommitted. (blade should report error, too)
		return "", "", errors.New("missing interface")
	}
	if param.LocalPort != "" {
		args = append(args, "--local-port", param.LocalPort)
	}
	if param.Offset != "" {
		args = append(args, "--offset", param.Offset)
	}
	if param.RemotePort != "" {
		args = append(args, "--remote-port", param.RemotePort)
	}
	if param.Time != "" {
		args = append(args, "--time", param.Time)
	} else {
		// * can be ommitted. (blade should report error, too)
		return "", "", errors.New("missing delay time")
	}
	if param.Force {
		args = append(args, "--force")
	}
	if param.IgnorePeerPort {
		args = append(args, "--ignore-peer-port")
	}
	if param.Timeout != "" {
		args = append(args, "--timeout", param.Timeout)
	}
	if param.RemotePort == "" && param.LocalPort == "" && param.DesIP == "" && param.ExcludePort == "" && param.Timeout == "" {
		// not specify certain testing ports or ips will cause the delay of the whole network device
		return "", "", errors.New("specify port or IP for testing or setup --timeout or --exclude-port")
	}

	// run executable
	out, _, err := utils.RunExec(args, param.Timeout)
	if err != nil {
		return "", "", err
	}

	//get uid
	uid, err := utils.GetUID(out.Bytes())
	if err != nil {
		return "", "", err
	}
	return uid, out.String(), err
}

func (op NetDelayOp) Destroy(id string) (string, error) {
	// using blade destroy cmd
	out, err := utils.DestroyExec(id)
	return out.String(), err
}
