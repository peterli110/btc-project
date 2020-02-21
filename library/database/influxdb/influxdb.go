package influxdb

import (
	"btc-project/library/log"
	xtime "btc-project/library/time"
	influx "github.com/influxdata/influxdb1-client/v2"
	"time"
)

type Config struct {
	Addr    		string 				`toml:"addr"`
	Timeout 		xtime.Duration    	`toml:"timeout"`
	DBName			string				`toml:"databaseName"`
	Username		string				`toml:"username"`
	Password		string				`toml:"password"`
}

func NewInfluxDB(c *Config) *influx.Client {
	if c.Addr == "" {
		panic("influxdb address not defined")
	}

	if c.Timeout <= 0 {
		panic("timeout should be larger than 0")
	}

	db, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr: c.Addr,
		Timeout: time.Duration(c.Timeout),
		Username: c.Username,
		Password: c.Password,
	})

	if err != nil {
		log.Error("connect influxdb error(%v)", err)
		panic(err)
	}

	log.Info("influxdb connected")
	return &db
}