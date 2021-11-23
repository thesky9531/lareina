package gomicro

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	grpcc "github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/server"
	grpcs "github.com/micro/go-micro/v2/server/grpc"
	"github.com/micro/go-micro/v2/web"
	"github.com/micro/go-plugins/registry/kubernetes/v2"
	"github.com/thesky9531/lareina/log"
	"strings"
	"time"
)

type Config struct {
	ListenAddr       string
	Registry         string
	RegisterTTL      int
	RegisterInterval int
	RegisterAddr     string
}

func NewRegistry(cf *Config) registry.Registry {
	switch strings.ToLower(strings.TrimSpace(cf.Registry)) {
	case "etcd":
		return etcd.NewRegistry(
			registry.Addrs(cf.RegisterAddr),
		)
	case "kubernetes":
		return kubernetes.NewRegistry()
	default:
		return etcd.NewRegistry(
			registry.Addrs(cf.RegisterAddr),
		)
	}
}

func NewServer(cf *Config, serviceName string, afterStop func() error) micro.Service {
	if cf.ListenAddr == "" {
		cf.ListenAddr = ":8081"
	}
	opt := []micro.Option{
		micro.Name(serviceName),
		micro.Server(grpcs.NewServer(server.Address(cf.ListenAddr), server.Name(serviceName))),
		micro.Client(grpcc.NewClient()),
		micro.RegisterTTL(time.Second * time.Duration(cf.RegisterTTL)),
		micro.RegisterInterval(time.Second * time.Duration(cf.RegisterInterval)),
		micro.Registry(
			NewRegistry(cf),
		),
		micro.AfterStop(afterStop),
	}
	return micro.NewService(opt...)
}

func NewClient(cf *Config) client.Client {
	service := micro.NewService(
		micro.Registry(
			NewRegistry(cf),
		),
	)
	client1 := service.Client()
	err := client1.Init(
		client.RequestTimeout(30*time.Second),
		client.Retries(0),
		client.DialTimeout(30*time.Second),
	)
	if err != nil {
		log.ErrLog("", err)
	}
	return client1
}

func NewClientOperationReportOnly(cf *Config) client.Client {
	service := micro.NewService(
		micro.Registry(
			NewRegistry(cf),
		),
	)
	client1 := service.Client()
	err := client1.Init(
		client.RequestTimeout(600*time.Second),
		client.Retries(0),
		client.DialTimeout(600*time.Second),
	)
	if err != nil {
		log.ErrLog("", err)
	}
	return client1
}

func NewWebServer(cf *Config, webName string, afterStop func() error) web.Service {
	if cf.ListenAddr == "" {
		cf.ListenAddr = ":8080"
	}
	opt := []web.Option{
		web.Name(webName),
		web.RegisterTTL(time.Second * time.Duration(cf.RegisterTTL)),
		web.RegisterInterval(time.Second * time.Duration(cf.RegisterInterval)),
		web.Address(cf.ListenAddr),
		web.Registry(
			NewRegistry(cf),
		),
		web.AfterStop(afterStop),
	}
	return web.NewService(opt...)
}

// NewWebSelector 选择器初始化
func NewWebSelector(cf *Config) selector.Selector {
	return selector.NewSelector(selector.Registry(NewRegistry(cf)))
}
