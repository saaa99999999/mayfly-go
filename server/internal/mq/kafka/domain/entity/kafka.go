package entity

import (
	"mayfly-go/internal/mq/kafka/kfm"
	"mayfly-go/pkg/model"
	"mayfly-go/pkg/utils/structx"
)

type Kafka struct {
	model.Model

	Code               string  `json:"code" gorm:"size:32;comment:code"`
	Name               string  `json:"name" gorm:"not null;size:50;comment:名称"`
	Hosts              string  `json:"hosts" gorm:"not null;size:500;comment:Kafka 连接地址，格式: host1:port1,host2:port2 或单个 broker"`
	Username           *string `json:"username" gorm:"size:100;comment:用户名"`
	Password           *string `json:"password" gorm:"size:100;comment:密码"`
	SshTunnelMachineId int     `json:"sshTunnelMachineId" gorm:"comment:ssh隧道的机器id"`
	SaslMechanism      *string `json:"saslMechanism" gorm:"comment:sasl机制"`
}

// 转换为kafkaInfo进行连接
func (k *Kafka) ToKafkaInfo() *kfm.KafkaInfo {
	mongoInfo := new(kfm.KafkaInfo)
	_ = structx.Copy(mongoInfo, k)
	return mongoInfo
}
