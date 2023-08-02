package wiring

import (
	"github.com/mjedari/vgang-project/src/app/configs"
	"net"
)

func (w *Wire) GetRedisUrl() string {
	return net.JoinHostPort(w.Configs.Server.Host, w.Configs.Server.Port)
}

func (w *Wire) GetServerConfig() string {
	return net.JoinHostPort(w.Configs.Server.Host, w.Configs.Server.Port)
}

func (w *Wire) GetRedisConfig() configs.RedisConfig {
	return w.Configs.Redis
}
