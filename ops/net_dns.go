package ops

import (
	"nessaj/utils"
)

type NetDNSOp struct {
	param NetDNSParam
}

type NetDNSParam struct {
	Domain  string
	IP      string
	Timeout string
}

func (op NetDNSOp) Name() string {
	return "netdns"
}

func (op NetDNSOp) Desc() string {
	return "modify local hosts or dns resolv"
}

func (op NetDNSOp) Validate(params map[string]interface{}) error {
	var param NetDNSParam
	err := utils.DecodeParams(params, &param)
	if err != nil {
		return err
	}
	return nil
}

func (op NetDNSOp) Params() ParamSpec {
	param := make(ParamSpec)

	param["Domain"] = Param{
		Iden:     "Domain",
		Desc:     "domain",
		Type:     PString,
		Required: true,
	}

	param["IP"] = Param{
		Iden:     "IP",
		Desc:     "ip",
		Type:     PString,
		Required: true,
	}

	param["Timeout"] = Param{
		Iden:     "Timeout",
		Desc:     "runtime limitation",
		Type:     PString,
		Required: false,
	}

	return param
}

func (op NetDNSOp) Init() error {
	return nil
}

func (op NetDNSOp) Run(params map[string]interface{}) (string, string, error) {
	// get parameters
	var param NetDNSParam
	err := utils.DecodeParams(params, &param)
	if err != nil {
		return "", "", err
	}

	args := []string{"create", "network", "dns"}
	if param.Domain != "" {
		args = append(args, "--domain", param.Domain)
	}
	if param.IP != "" {
		args = append(args, "--ip", param.IP)
	}
	if param.Timeout != "" {
		args = append(args, "--timeout", param.Timeout)
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

func (op NetDNSOp) Destroy(id string) (string, error) {
	// using blade destroy cmd
	out, err := utils.DestroyExec(id)
	return out.String(), err
}
