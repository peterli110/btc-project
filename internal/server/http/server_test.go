package http

import (
	"btc-project/conf"
	"btc-project/internal/service"
	"flag"
	influx "github.com/influxdata/influxdb1-client/v2"
	"os"
	"testing"

	bm "btc-project/library/net/http/blademaster"
	"github.com/gavv/httpexpect"
)


var instance *bm.Engine

func TestMain(m *testing.M) {
	if err := os.Setenv("TZ", "Greenwich"); err != nil {
		panic(err)
	}

	if err := flag.Set("conf", "../../../configs/"); err != nil {
		panic(err)
	}

	flag.Parse()
	conf.FileName = "config_test.toml"
	if err := conf.Init(); err != nil {
		panic(err)
	}


	svc := service.New(conf.Conf)
	instance = New(conf.Conf, svc)

	// reset test db
	q := influx.NewQuery("DROP DATABASE " + conf.Conf.InfluxDBConfig.DBName, "", "")
	resp, err := svc.Dao.InfluxDB.Query(q)
	if err != nil {
		panic(err)
	}
	if resp.Error() != nil {
		panic(resp.Error())
	}

	q = influx.NewQuery("CREATE DATABASE " + conf.Conf.InfluxDBConfig.DBName, "", "")
	resp, err = svc.Dao.InfluxDB.Query(q)
	if err != nil {
		panic(err)
	}
	if resp.Error() != nil {
		panic(resp.Error())
	}



	m.Run()
	os.Exit(0)
}


func httpTester(t *testing.T) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://" + conf.Conf.Server.Addr,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
}

