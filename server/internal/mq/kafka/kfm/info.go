package kfm

import (
	"context"
	"fmt"
	machineapp "mayfly-go/internal/machine/application"
	"mayfly-go/pkg/logx"
	"mayfly-go/pkg/utils/netx"
	"net"
	"strings"
	"time"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/plain"
	"github.com/twmb/franz-go/pkg/sasl/scram"
)

type KafkaInfo struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`

	// Kafka 连接地址，格式: host1:port1,host2:port2 或单个 broker
	Hosts string `json:"-"`

	Username string `json:"-"`
	Password string `json:"-"`

	SshTunnelMachineId int `json:"-"` // ssh隧道机器id

	// SASL 配置
	SaslEnabled   bool   `json:"-"`
	SaslMechanism string `json:"-"` // PLAIN, SCRAM-SHA-256, SCRAM-SHA-512
}

func (mi *KafkaInfo) Conn() (*KafkaConn, error) {
	opts := mi.BuildBaseOpts()

	// 创建客户端
	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka client: %w", err)
	}

	kac := kadm.NewClient(client)

	// 测试连接 - 获取 broker 列表来验证连接
	brokers, err := kac.ListBrokers(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to list brokers: %w", err)
	}
	if len(brokers) == 0 {
		client.Close()
		return nil, fmt.Errorf("no available brokers")
	}

	logx.Infof("连接 kafka: %s", mi.Hosts)

	return &KafkaConn{
		Id:      getConnId(mi.Id),
		Info:    mi,
		Client:  client,
		Configs: opts,
		Ac:      kac,
	}, nil
}

func (mi *KafkaInfo) BuildBaseOpts() []kgo.Opt {

	// 构建基础配置（复用 KafkaInfo.Conn 的逻辑）
	opts := []kgo.Opt{
		kgo.SeedBrokers(strings.Split(mi.Hosts, ",")...),
		kgo.DialTimeout(8 * time.Second),
		kgo.RecordPartitioner(kgo.ManualPartitioner()),
	}

	// 配置认证
	if mi.Username != "" {
		scramAuth := &scram.Auth{
			User: mi.Username,
			Pass: mi.Password,
		}
		plainAuth := &plain.Auth{
			User: mi.Username,
			Pass: mi.Password,
		}
		if mi.SaslMechanism != "" {
			switch strings.ToUpper(mi.SaslMechanism) {
			case "SCRAM-SHA-256":
				opts = append(opts, kgo.SASL(scramAuth.AsSha256Mechanism()))
			case "SCRAM-SHA-512":
				opts = append(opts, kgo.SASL(scramAuth.AsSha512Mechanism()))
			default:
				opts = append(opts, kgo.SASL(plainAuth.AsMechanism()))
			}
		} else {
			opts = append(opts, kgo.SASL(plainAuth.AsMechanism()))
		}
	}

	// SSH 隧道
	if mi.SshTunnelMachineId > 0 {
		stm, err := machineapp.GetMachineApp().GetSshTunnelMachine(context.Background(), mi.SshTunnelMachineId)
		if err != nil {
			logx.Errorf("获取 ssh隧道失败：%v", err)
		} else {
			dialFn := func(ctx context.Context, network, address string) (net.Conn, error) {
				sshConn, err := stm.GetDialConn(network, address)
				if err != nil {
					return nil, err
				}
				return &netx.WrapSshConn{Conn: sshConn}, nil
			}
			opts = append(opts, kgo.Dialer(dialFn))
		}
	}

	return opts
}

type KafkaSshProxyDialer struct {
	machineId int
}

func (sd *KafkaSshProxyDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	stm, err := machineapp.GetMachineApp().GetSshTunnelMachine(ctx, sd.machineId)
	if err != nil {
		return nil, err
	}
	if sshConn, err := stm.GetDialConn(network, address); err == nil {
		// 将ssh conn包装，否则内部设置超时会报错,ssh conn不支持设置超时会返回错误: ssh: tcpChan: deadline not supported
		return &netx.WrapSshConn{Conn: sshConn}, nil
	} else {
		return nil, err
	}
}

// Dial 实现 sarama.Dialer 接口
func (sd *KafkaSshProxyDialer) Dial(network, address string) (net.Conn, error) {
	return sd.DialContext(context.Background(), network, address)
}

// 生成kafka连接id
func getConnId(id uint64) string {
	if id == 0 {
		return ""
	}
	return fmt.Sprintf("kafka:%d", id)
}

// CreateTopicParam 创建Topic参数
type CreateTopicParam struct {
	TopicName         string             `json:"topic" binding:"required" `             // Topic 名称
	NumPartitions     int32              `json:"numPartitions" binding:"required" `     // 分区数
	ReplicationFactor int16              `json:"replicationFactor" binding:"required" ` // 副本数
	ConfigEntries     map[string]*string `json:"configEntries"`                         // 配置项
}

// CreatePartitionsParam 修改分区参数
type CreatePartitionsParam struct {
	TopicName     string `json:"topic"`         // Topic 名称
	NumPartitions int    `json:"numPartitions"` // 分区数
}

// ConsumeMessageParam 消费消息参数
type ConsumeMessageParam struct {
	Topic          string `json:"topic"`
	Group          string `json:"group"`                    // 消费组 ID，为空则不使用消费组              // Topic 名称
	Partition      int32  `json:"partition" default:"-1"`   // 分区号，-1 表示所有分区
	Number         int    `json:"number" default:"10"`      // 消费消息数量
	PullTimeout    int    `json:"pullTimeout" default:"10"` // 拉取超时时间（秒）
	Decompression  string `json:"decompression"`            // 解压方式：gzip, lz4, zstd, snappy
	Decode         string `json:"decode"`                   // 解码方式：Base64
	Earliest       bool   `json:"earliest"`                 // 是否从最早开始消费，false 则从最新
	StartTime      string `json:"startTime"`                // 消费起始时间戳（毫秒），0 表示不使用，优先级高于 Earliest
	CommitOffset   bool   `json:"commitOffset"`             // 是否提交消费位点
	IsolationLevel string `json:"isolationLevel"`           // 读取消息的隔离级别，默认为 read_committed
}

// ConsumeMessageResult 消费的消息结果
type ConsumeMessageResult struct {
	Id            int               `json:"id"`
	Offset        int64             `json:"offset"`
	Partition     int32             `json:"partition"`
	Key           string            `json:"key"`
	Value         string            `json:"value"`
	Timestamp     string            `json:"timestamp"`
	Topic         string            `json:"topic"`
	Headers       map[string]string `json:"headers"`
	LeaderEpoch   int32             `json:"leaderEpoch"`
	ProducerEpoch int16             `json:"producerEpoch"`
	ProducerID    int64             `json:"producerID"`
}

// ProduceMessageParam 生产消息参数
type ProduceMessageParam struct {
	Topic       string              `json:"topic"`       // Topic 名称
	Key         string              `json:"key"`         // 消息 Key（可选）
	Partition   int32               `json:"partition"`   // 指定分区号（可选，-1 表示自动选择）
	Value       string              `json:"value"`       // 消息内容
	Headers     []map[string]string `json:"headers"`     // 消息 Headers（可选）
	Times       int                 `json:"times"`       // 发送次数
	Compression string              `json:"compression"` // 压缩方式：gzip, lz4, zstd, snappy
}

type RecordHeader struct {
	Key   string
	Value string
}

// ProduceMessageResult 生产消息结果
type ProduceMessageResult struct {
	Partition int32 `json:"partition"`
	Offset    int64 `json:"offset"`
}

type BrokerInfo struct {
	Id   int32   `json:"id"`
	Addr string  `json:"addr"`
	Rack *string `json:"rac"`
}
