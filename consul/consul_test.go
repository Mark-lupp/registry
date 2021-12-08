package consul

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Mark-lupp/registry"
)

func TestRegister(t *testing.T) {
	Iregistry, err := registry.InitRegistry(context.TODO(), "consul", registry.WithAddrs([]string{"127.0.0.1:8501"}),
		registry.WithTimeout(30*time.Second),
		registry.WithHeartBeat(5),
	)
	if err != nil {
		t.Error(err)
		return
	}
	service := &registry.Service{
		Name: "test",
	}
	service.Instances = append(service.Instances, &registry.Instance{
		Id:   "test-1",
		IP:   "127.0.0.1",
		Port: 8801,
	},
		&registry.Instance{
			Id:   "test-2",
			IP:   "127.0.0.1",
			Port: 8802,
		})
	if err = Iregistry.Register(context.TODO(), service); err != nil {
		t.Error(err)
		return
	}
}
func TestDeRegister(t *testing.T) {
	Iregistry, err := registry.InitRegistry(context.TODO(), "consul", registry.WithAddrs([]string{"127.0.0.1:8501"}),
		registry.WithTimeout(time.Second),
		registry.WithHeartBeat(5),
	)
	if err != nil {
		t.Error(err)
		return
	}
	service := &registry.Service{
		Name: "test",
	}
	service.Instances = append(service.Instances, &registry.Instance{
		Id:   "test-1",
		IP:   "127.0.0.1",
		Port: 8801,
	},
		&registry.Instance{
			Id:   "test-2",
			IP:   "127.0.0.1",
			Port: 8802,
		})
	if err := Iregistry.Unregister(context.TODO(), service); err != nil {
		t.Errorf("DeRegister Fail:%s", err.Error())
		return
	}
}
func TestGetService(t *testing.T) {
	Iregistry, err := registry.InitRegistry(context.TODO(), "consul", registry.WithAddrs([]string{"127.0.0.1:8501"}),
		registry.WithTimeout(time.Second),
		registry.WithHeartBeat(5),
	)
	if err != nil {
		t.Error(err)
		return
	}
	service, err := Iregistry.GetService(context.TODO(), "test")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(service)
}
