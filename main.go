package main

import (
	"fmt"
	"nessaj/config"
	"nessaj/ops"
	"nessaj/server"
	"nessaj/utils"
	"os"
)

func main() {
	conf, err := config.Parse()
	utils.ChkErr(err)
	err = ops.Init()
	if err != nil {
		fmt.Println("operation initialization has error: ", err)
	}
	err = utils.Init(conf)
	if err != nil {
		fmt.Println("utils init failed (very unlikely to happen)")
	}
	fmt.Printf("got %d operations\n", len(ops.AllOps))
	err = server.RunServer(conf)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	}
}
