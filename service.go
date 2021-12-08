package registry

const (
	CONSUL = iota
	ETCD
	ZOOKEEPER
	NACOS
)

var RegistryType = map[int]string{
	CONSUL:    "consul",
	ETCD:      "etcd",
	ZOOKEEPER: "zookeeper",
	NACOS:     "nacos",
}

// 服务抽象
type Service struct {
	Name           string            `json:"name"`      // 服务名称
	Instances      []*Instance       `json:"instances"` // 服务实例
	HealthCheckUrl string            // 健康检查地址
	Meta           map[string]string // 服务元数据
	Tags           []string          // 服务标签
}

// 服务实例的抽象
type Instance struct {
	Id       string `json:"id"`     // 实例ID
	IP       string `json:"ip"`     // 实例地址
	Port     int    `json:"port"`   // 实例端口
	Weight   int    `json:"weight"` // 实例权重
	GrpcPort int
}

func GetRegistryType(t int) string {
	str, ok := RegistryType[t]
	if ok {
		return str
	}
	return ""
}
