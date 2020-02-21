package dao

import (
	"btc-project/internal/model"
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestInsertCoin(t *testing.T) {
	Convey("TestInsertCoin insert data into influxdb", t, func() {
		for i := 0; i < 10; i++ {
			data := &model.BtcData{
				DateTime: time.Date(2020, 02, 10, i, 0, 0, 0, time.UTC),
				Amount: int64((i + 1) * 100000000),
			}

			err := d.InsertCoin(context.TODO(), data)
			So(err, ShouldBeNil)
		}
	})
}



func TestQueryCoin(t *testing.T) {
	Convey("TestGetQueryCache get data from influxdb", t, func() {
		for i := 0; i < 10; i++ {
			query := &model.BtcQuery{
				StartTime: time.Date(2020, 02, 10, i, 0, 0, 0, time.UTC),
				EndTime: time.Date(2020, 02, 10, i + 2, 0, 0, 0, time.UTC),
			}
			res, err := d.QueryCoin(context.TODO(), query)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeEmpty)
		}
	})
}