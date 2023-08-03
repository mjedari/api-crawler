package configs

var Config Configuration

type Server struct {
	Host string
	Port string
}

type RedisConfig struct {
	Host             string
	Port             string
	User             string
	Pass             string
	ServicesRedisKey string
}

type Collector struct {
	Concurrent  bool
	Interval    int
	SplitFactor int
}

type Credentials struct {
	UserName string
	Password string
}

type RateLimiter struct {
	Active bool
	Rate   int
	Period uint64
}

type OriginRemote struct {
	Products string
	Login    string
	BaseURL  string
}

type Configuration struct {
	Server       Server
	Redis        RedisConfig
	RateLimiter  RateLimiter
	Credentials  Credentials
	Collector    Collector
	OriginRemote OriginRemote
}
