package vo

import (
	"mayfly-go/pkg/model"
)

type ContainerConf struct {
	model.Model
	model.ExtraData

	Addr   string `json:"addr"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Remark string `json:"remark"`
}

func (c *ContainerConf) GetCode() string {
	return c.Code
}
