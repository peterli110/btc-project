package service

import (
	"btc-project/conf"
	"flag"
	influx "github.com/influxdata/influxdb1-client/v2"
	"os"
	"testing"
)

var (
	s *Service
)

func TestMain(m *testing.M) {
	if err := os.Setenv("TZ", "Greenwich"); err != nil {
		panic(err)
	}

	if err := flag.Set("conf", "../../configs/"); err != nil {
		panic(err)
	}

	flag.Parse()
	conf.FileName = "config_test.toml"
	if err := conf.Init(); err != nil {
		panic(err)
	}


	s = New(conf.Conf)

	// reset test db
	q := influx.NewQuery("DROP DATABASE " + conf.Conf.InfluxDBConfig.DBName, "", "")
	resp, err := s.Dao.InfluxDB.Query(q)
	if err != nil {
		panic(err)
	}
	if resp.Error() != nil {
		panic(resp.Error())
	}

	q = influx.NewQuery("CREATE DATABASE " + conf.Conf.InfluxDBConfig.DBName, "", "")
	resp, err = s.Dao.InfluxDB.Query(q)
	if err != nil {
		panic(err)
	}
	if resp.Error() != nil {
		panic(resp.Error())
	}

	m.Run()
	os.Exit(0)
}
