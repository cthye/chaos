package config

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io/ioutil"
	"nessaj/constant"
	"net/http"
	"os"
	"path/filepath"
)

type Config struct {
	Host          string
	Port          uint16
	Pubkey        *ecdsa.PublicKey
	ChaosbladeBin string
	Verbose       bool
	ProxyAddr     string
}

func MkConfig(host string, port uint16, pubkey *ecdsa.PublicKey, chaosbladeBin string, verbose bool, proxyAddr string) Config {
	return Config{
		Host:          host,
		Port:          port,
		Pubkey:        pubkey,
		ChaosbladeBin: chaosbladeBin,
		Verbose:       verbose,
		ProxyAddr:     proxyAddr,
	}
}

var config_path string

func init() {
	exeDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	chaosbladeDefaultBin := filepath.Join(exeDir, constant.ChaosbladeDefaultFolder, "blade")
	flag.String("host", "127.0.0.1", "host ip to bind")
	flag.Uint16P("port", "p", 1337, "port used to bind")
	flag.String("pubkey", "", "public key in PEM format, mutual exclusive to pubkey_file")
	flag.String("pubkey_file", "", "public key file with content in PEM format, mutual exclusive to pubkey")
	flag.String("chaosblade_bin", chaosbladeDefaultBin, "bin path of chaosblade (blade)")
	flag.BoolP("verbose", "v", false, "verbose mode")
	flag.StringVarP(&config_path, "config", "c", "nessaj.yaml", "config file path")
	flag.String("proxy_addr", "", "proxy binded to agent")
}

type RespJson struct {
	Code int
	Data map[string]string
	Msg  string
}

func Parse() (*Config, error) {
	// env
	viper.SetEnvPrefix("nessaj")
	viper.BindEnv("host")
	viper.BindEnv("port")
	viper.BindEnv("chaosblade_bin")

	// flag
	flag.Parse()
	err := viper.BindPFlags(flag.CommandLine)
	if err != nil {
		return nil, err
	}

	// config file
	viper.SetConfigFile(config_path)
	viper.ReadInConfig() // ignore error

	var pubkey_content []byte
	proxyAddr := viper.GetString("proxy_addr")
	if proxyAddr != "" {
		// get key from proxy
		constant.ProxyURL = proxyAddr
		resp, err := http.Get("http://" + constant.ProxyURL + "/pub")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var respJson RespJson
		if err = json.NewDecoder(resp.Body).Decode(&respJson); err != nil {
			return nil, err
		}
		pubkey_content = []byte(respJson.Data["pub"])
	} else {
		pubkey := viper.GetString("pubkey")
		pubkey_file := viper.GetString("pubkey_file")
		if pubkey == "" && pubkey_file == "" {
			return nil, errors.New("pubkey and pubkey_file can not both be empty")
		}
		if pubkey != "" && pubkey_file != "" {
			return nil, errors.New("only one of pubkey and pubkey_file can be specified")
		}
		if pubkey != "" {
			pubkey_content = []byte(pubkey)
		} else {
			pubkey_bytes, err := ioutil.ReadFile(pubkey_file)
			if err != nil {
				return nil, err
			}
			pubkey_content = pubkey_bytes
		}
	}

	block, _ := pem.Decode(pubkey_content)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	switch pub := pub.(type) {
	case *ecdsa.PublicKey:
		fmt.Println("public key parsed successful")
	default:
		return nil, errors.New(fmt.Sprintf("unexpected public key type %T", pub))
	}
	config := MkConfig(viper.GetString("host"), uint16(viper.GetUint("port")),
		pub.(*ecdsa.PublicKey), viper.GetString("chaosblade_bin"), viper.GetBool("verbose"), proxyAddr)
	return &config, nil
}
