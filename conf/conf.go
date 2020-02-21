package conf

import (
	"btc-project/library/conf/paladin"
	"btc-project/library/database/influxdb"
	redisclient "btc-project/library/database/redis"
	bm "btc-project/library/net/http/blademaster"
)

// config variables
var (
	Conf = &Config{}
	FileName = "config.toml"
)


// Config server config
type Config struct {
	Server    			*bm.ServerConfig    	`toml:"server"`
	InfluxDBConfig		*influxdb.Config		`toml:"influxdb"`
	RedisConfig     	*redisclient.Config     `toml:"redis"`
	Constants 			*Constants 				`toml:"constants"`
}

// Constants server constants
type Constants struct {
	IsDev        		bool   					`toml:"isDev"`
}



// Init initialize the conf file
func Init() (err error) {
	if err = paladin.Init(); err != nil {
		return
	}

	if err = paladin.Get(FileName).UnmarshalTOML(&Conf); err != nil {
		return
	}

	return
}