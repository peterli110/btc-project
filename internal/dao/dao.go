package dao

import (
	"btc-project/conf"
	"btc-project/library/database/influxdb"
	redisclient "btc-project/library/database/redis"
	"btc-project/library/log"
	"context"
	"github.com/go-redis/redis"
	influx "github.com/influxdata/influxdb1-client/v2"
	"time"
)

// Dao object
type Dao struct {
	InfluxDB			influx.Client
	Redis 				*redis.Client
}


// New init DAO
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		InfluxDB: *influxdb.NewInfluxDB(c.InfluxDBConfig),
		Redis: redisclient.NewRedis(c.RedisConfig),
	}


	// create database if not exist
	q := influx.NewQuery("CREATE DATABASE " + c.InfluxDBConfig.DBName, "", "")
	resp, err := dao.InfluxDB.Query(q)
	if err != nil {
		panic(err)
	}

	if resp.Error() != nil {
		panic(resp.Error())
	}
	return
}



// Close will disconnect from DB
func (dao *Dao) Close() {
	err := dao.InfluxDB.Close()
	if err != nil {
		log.Error("disconnect influxdb error(%v)", err)
	}
	err = dao.Redis.Close()
	if err != nil {
		log.Error("disconnect redis error(%v)", err)
	}
	log.Info("disconnected")
}



// Ping check the connection of DB.
func (dao *Dao) Ping(ctx context.Context) (err error) {
	if _, _, err = dao.InfluxDB.Ping(1 * time.Second); err != nil {
		return
	}

	if _, err = dao.Redis.Ping().Result(); err != nil {
		return
	}

	return
}
