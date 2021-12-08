package consul

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/Mark-lupp/registry"
	"github.com/go-kit/kit/sd/consul"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

type ConsulAdapter struct {
	options *registry.Options
	client  consul.Client
	config  *consulapi.Config
	mutex   sync.Mutex
	// 服务实例缓存字段
	instancesMap sync.Map
}

func init() {
	// 注册插件
	registry.RegisterPlugin(&ConsulAdapter{})
}

func (c *ConsulAdapter) Name() string {

	return registry.GetRegistryType(registry.CONSUL)
}
func (c *ConsulAdapter) Init(ctx context.Context, opts ...registry.Option) (err error) {

	c.options = &registry.Options{}
	for _, opt := range opts {
		opt(c.options)
	}
	for _, address := range c.options.Addrs {
		client, err := consulapi.NewClient(&consulapi.Config{
			Address: address,
		})
		if err != nil {
			return err
		} else {
			if _, err = client.Status().Leader(); err != nil {
				return err
			}
			c.client = consul.NewClient(client)
			c.config = &consulapi.Config{
				Address: address,
			}
			continue
		}
	}

	return nil

}
func (c *ConsulAdapter) Register(ctx context.Context, service *registry.Service) (err error) {
	serviceRegistration := &consulapi.AgentServiceRegistration{}
	serviceCheck := &consulapi.AgentServiceCheck{}
	if c.options.RegistryPath != "" {
		serviceRegistration.Namespace = c.options.RegistryPath
	}
	// 检查失败多久删除服务
	if c.options.Timeout > 0 {
		seconds := fmt.Sprint(c.options.Timeout.Seconds()) + "s"
		serviceCheck.DeregisterCriticalServiceAfter = seconds
	}
	// 检查服务心跳间隔
	if c.options.HeartBeat > 0 {
		seconds := strconv.Itoa(int(c.options.HeartBeat)) + "s"
		serviceCheck.Interval = seconds
	}
	serviceRegistration.Name = service.Name
	for _, instance := range service.Instances {
		serviceRegistration.Address = instance.IP
		serviceRegistration.Port = instance.Port
		serviceRegistration.Meta = service.Meta
		serviceRegistration.Tags = service.Tags
		if instance.Weight > 0 {
			serviceRegistration.Weights = &consulapi.AgentWeights{Passing: instance.Weight}
		}
		serviceCheck.HTTP = "http://" + instance.IP + ":" + strconv.Itoa(instance.Port) + service.HealthCheckUrl
		serviceRegistration.Check = serviceCheck
		if err := c.client.Register(serviceRegistration); err != nil {
			return err
		}
	}
	return nil
}
func (c *ConsulAdapter) Unregister(ctx context.Context, service *registry.Service) (err error) {
	for _, instance := range service.Instances {
		serviceRegistration := &consulapi.AgentServiceRegistration{
			ID: instance.Id,
		}
		if err = c.client.Deregister(serviceRegistration); err != nil {
			return err
		}
	}

	return
}
func (c *ConsulAdapter) GetService(ctx context.Context, serviceName string) (service *registry.Service, err error) {
	//  该服务已监控并缓存
	instanceList, ok := c.instancesMap.Load(serviceName)
	if ok {
		return instanceList.(*registry.Service), nil
	}
	// 申请锁
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// 再次检查是否监控
	instanceList, ok = c.instancesMap.Load(serviceName)
	if ok {
		return instanceList.(*registry.Service), nil
	} else {
		// 注册监控
		go func() {
			params := make(map[string]interface{})
			params["type"] = "service"
			params["service"] = serviceName
			plan, _ := watch.Parse(params)
			plan.Handler = func(u uint64, i interface{}) {
				if i == nil {
					return
				}
				v, ok := i.([]*consulapi.ServiceEntry)
				if !ok {
					return // 数据异常，忽略
				}

				// 没有服务实例在线
				if len(v) == 0 {
					c.instancesMap.Store(serviceName, &registry.Service{})
				}

				var healthServices *registry.Service

				for _, service := range v {
					if service.Checks.AggregatedStatus() == consulapi.HealthPassing {
						healthServices.Instances = append(healthServices.Instances, newServiceInstance(service.Service))
					}
				}
				c.instancesMap.Store(serviceName, healthServices)
			}
			defer plan.Stop()
			plan.Run(c.config.Address)
		}()
	}
	return nil, nil
}
func newServiceInstance(service *consulapi.AgentService) *registry.Instance {

	rpcPort := service.Port - 1
	if service.Meta != nil {
		if rpcPortString, ok := service.Meta["rpcPort"]; ok {
			rpcPort, _ = strconv.Atoi(rpcPortString)
		}
	}
	return &registry.Instance{
		IP:       service.Address,
		Port:     service.Port,
		GrpcPort: rpcPort,
		Weight:   service.Weights.Passing,
	}

}
