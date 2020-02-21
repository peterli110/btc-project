package dao

import (
	"btc-project/conf"
	"btc-project/internal/model"
	"btc-project/library/ecode"
	"btc-project/library/log"
	"context"
	"fmt"
	influx "github.com/influxdata/influxdb1-client/v2"
	"time"
)

// queryString will query on [startTime, endTime)
const queryString = "select value from (select CUMULATIVE_SUM(SUM(value)) as value," +
	"last(timestamp) as ts from amount where time < '%v' " +
	"group by time(1h)) where ts >= %v"

// InsertCoin save a new record in db
// amount and timestamp will be saved in int64,
// timestamp is additional variable to help query the target data.
func (dao *Dao) InsertCoin(c context.Context, m *model.BtcData) error {
	pt, err := influx.NewPoint(
		"amount",
		nil,
		map[string]interface{}{"value": m.Amount, "timestamp": m.DateTime.Unix()},
		m.DateTime)

	if err != nil {
		log.Error("create new point: (%v)", err)
		return ecode.ServerErr
	}

	bps, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database: conf.Conf.InfluxDBConfig.DBName,
	})

	if err != nil {
		log.Error("create new batchpoint: (%v)", err)
		return ecode.ServerErr
	}

	bps.AddPoint(pt)
	err = dao.InfluxDB.Write(bps)
	if err != nil {
		log.Error("influxdb write: (%v)", err)
		return ecode.ServerErr
	}

	return nil
}


// QueryCoin will return the result of query
func (dao *Dao) QueryCoin(c context.Context, m *model.BtcQuery) ([][]interface{}, error) {
	qs := fmt.Sprintf(queryString, m.EndTime.UTC().Format(time.RFC3339), m.StartTime.UTC().Unix())
	q := influx.NewQuery(qs, conf.Conf.InfluxDBConfig.DBName, "RFC3339")
	resp, err := dao.InfluxDB.Query(q)
	if err != nil {
		log.Error("query: (%v)", err)
		return nil, ecode.ServerErr
	}
	if resp.Error() != nil {
		log.Error("query response: (%v)", resp.Error())
		return nil, ecode.ServerErr
	}

	// response length should be 1 even if there are no result
	if len(resp.Results) == 0 {
		log.Error("response length is 0: (%v)", qs)
		return nil, ecode.ServerErr
	}

	// no result
	if len(resp.Results[0].Series) == 0 {
		return nil, nil
	}

	return resp.Results[0].Series[0].Values, nil
}