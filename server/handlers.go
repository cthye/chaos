package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"nessaj/ops"
)

var (
	version string = "1.0.0"
)

func mkResp(data interface{}) map[string]interface{} {
	return gin.H{
		"Code": 0,
		"Data": data,
		"Msg":  "",
	}
}
func mkErrResp(code int, msg string) map[string]interface{} {
	return gin.H{
		"Code": code,
		"Data": "",
		"Msg":  msg,
	}
}

func versionHandler(c *gin.Context) {
	c.JSON(200, mkResp(version))
}

type ChaosCreateBody struct {
	Op     string                 `json:"op" binding:"required"`
	Params map[string]interface{} `json:"params"`
}

type ChaosDestroyBody struct {
	Op string `json:"op" binding:"required"`
	Id string `json:"id" binding:"required"`
}

func chaosRunHandler(c *gin.Context) {
	var body ChaosCreateBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(200, mkErrResp(1, err.Error()))
		return
	}
	if op, ok := ops.AllOps[body.Op]; !ok {
		c.JSON(200, mkErrResp(1, "invalid operation"))
	} else {
		err = op.Validate(body.Params)
		if err != nil {
			c.JSON(200, mkErrResp(1, err.Error()))
			return
		}

		id, ret, err := op.Run(body.Params)
		if err != nil {
			c.JSON(200, mkErrResp(1, err.Error()))
			return
		}
		c.JSON(200, mkResp(map[string]interface{}{
			"Id":   id,
			"Info": ret,
		}))
	}
}

func chaosDestroyHandler(c *gin.Context) {
	var body ChaosDestroyBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(200, mkErrResp(1, err.Error()))
		return
	}
	if op, ok := ops.AllOps[body.Op]; !ok {
		c.JSON(200, mkErrResp(1, "invalid operation"))
	} else {
		ret, err := op.Destroy(body.Id)
		if err != nil {
			c.JSON(200, mkErrResp(1, err.Error()))
			return
		}
		c.JSON(200, mkResp(map[string]interface{}{
			"info": ret,
		}))
	}
}

type OperationInfo struct {
	Name string `json:"name"`
	Desc string `json:"description"`
}

type OperationDetail struct {
	OperationInfo
	Params ops.ParamSpec `json:"parameters"`
}

func chaosListHandler(c *gin.Context) {
	var result []OperationInfo
	for _, op := range ops.AllOps {
		result = append(result, OperationInfo{op.Name(), op.Desc()})
	}
	c.JSON(200, result)
}

func chaosDetailHandler(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(200, mkErrResp(1, "unexpected empty operation name"))
		return
	}
	if op, ok := ops.AllOps[name]; !ok {
		c.JSON(404, mkErrResp(1, fmt.Sprintf("operation %s not found", name)))
	} else {
		c.JSON(200, mkResp(OperationDetail{
			OperationInfo{
				Name: op.Name(),
				Desc: op.Desc(),
			},
			op.Params(),
		}))
	}
}
