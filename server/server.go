package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"nessaj/config"
	"nessaj/constant"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupServer(conf *config.Config) *gin.Engine {
	r := gin.Default()
	r.Use(AuthenticationMiddleware(conf))

	r.GET("/version", versionHandler)
	r.GET("/chaos/list", chaosListHandler)
	r.GET("/chaos/detail/:name", chaosDetailHandler)

	r.POST("/chaos/run", chaosRunHandler)
	r.POST("/chaos/destroy", chaosDestroyHandler)

	return r
}

func Register(conf *config.Config) error {
	ip := conf.Host
	port := conf.Port
	// fmt.Printf("ip: %s", string(ip))
	reqBody, err := json.Marshal(map[string]string{
		"ip":     fmt.Sprintf("%s:%v", ip, port),
		"status": "running",
	})
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/register", "http://"+constant.ProxyURL), bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	respBytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respBytes))
	return nil
}

func RunServer(conf *config.Config) error {
	r := SetupServer(conf)
	if conf.ProxyAddr != "" {
		if err := Register(conf); err != nil {
			return fmt.Errorf("register failed: %s", err)
		}
	}
	return r.Run(fmt.Sprintf("%s:%d", conf.Host, conf.Port))
}
