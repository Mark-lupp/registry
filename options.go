package registry

import (
	"time"
)

type Options struct {
	Addrs   []string      // 注册中心地址
	Timeout time.Duration // 服务过期时间，检查失败后多久删除服务
	// example:  /xxx_company/app/kuaishou/service_A/10.192.1.1:8801
	// example:  /xxx_company/app/kuaishou/service_A/10.192.1.2:8801
	RegistryPath string // 注册路径或命令空间
	HeartBeat    int64  // 健康检查间隔
}

type Option func(opts *Options)

func WithTimeout(timeout time.Duration) Option {
	return func(opts *Options) {
		opts.Timeout = timeout
	}
}

func WithAddrs(addrs []string) Option {
	return func(opts *Options) {
		opts.Addrs = addrs
	}
}

func WithRegistryPath(path string) Option {
	return func(opts *Options) {
		opts.RegistryPath = path
	}
}

func WithHeartBeat(heartHeat int64) Option {
	return func(opts *Options) {
		opts.HeartBeat = heartHeat
	}
}
