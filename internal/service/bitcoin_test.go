package service

import (
	"btc-project/internal/model"
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

/*
	2020/02/10 00:00 get 10 btc
	2020/02/10 01:00 get 20 btc
	2020/02/10 02:00 get 30 btc
	2020/02/10 03:00 get 40 btc
	2020/02/10 04:00 get 50 btc
 */
func TestAddBitCoin(t *testing.T) {
	Convey("TestAddBitCoin insert data into influxdb", t, func() {
		for i := 0; i < 5; i++ {
			data := &model.BtcJSON{
				DateTime: time.Date(2020, 02, 10, i, 0, 0, 0, time.UTC),
				Amount: float64((i + 1) * 10),
			}

			err := s.AddBitCoin(context.TODO(), data)
			So(err, ShouldBeNil)
		}
	})
}


/*
	Query from  [2020/02/10 00:00, 2020/02/10 03:00)
			    [2020/02/10 00:00, 2020/02/10 01:00) ===> 10 btc
				[2020/02/10 01:00, 2020/02/10 02:00) ===> 10 + 20 = 30 btc
				[2020/02/10 02:00, 2020/02/10 03:00) ===> 30 + 30 = 60 btc
 */
func TestQueryBitCoin(t *testing.T) {
	Convey("TestQueryBitCoin get data and verify correctness", t, func() {
		query := &model.BtcQuery{
			StartTime: time.Date(2020, 02, 10, 0, 20, 25, 0, time.UTC),
			EndTime: time.Date(2020, 02, 10, 3, 30, 50, 0, time.UTC),
		}
		res, err := s.QueryBitCoin(context.TODO(), query)
		So(err, ShouldBeNil)

		// verify result is correct
		So(res, ShouldHaveLength, 3)
		So(res[0].DateTime, ShouldEqual, "2020-02-10T00:00:00+00:00")
		So(res[0].Amount, ShouldEqual, 10)
		So(res[1].DateTime, ShouldEqual, "2020-02-10T01:00:00+00:00")
		So(res[1].Amount, ShouldEqual, 30)
		So(res[2].DateTime, ShouldEqual, "2020-02-10T02:00:00+00:00")
		So(res[2].Amount, ShouldEqual, 60)

		Convey("When there is no result", func() {
			query := &model.BtcQuery{
				StartTime: time.Date(2010, 02, 10, 0, 20, 25, 0, time.UTC),
				EndTime: time.Date(2010, 02, 10, 3, 30, 50, 0, time.UTC),
			}
			res, err := s.QueryBitCoin(context.TODO(), query)
			So(err, ShouldBeNil)
			So(res, ShouldHaveLength, 0)
		})
	})
}


func TestBtcFloatToInt(t *testing.T) {
	Convey("TestBtcFloatToInt convert by decimal", t, func() {
		f := float64(1051.12)
		errResult := f * 100000000
		correct := btcFloatToInt(f)
		So(errResult, ShouldNotEqual, 105112000000)
		So(correct, ShouldEqual, 105112000000)
	})
}


func TestBtcIntToFloat(t *testing.T) {
	Convey("TestBtcIntToFloat convert by decimal", t, func() {
		i := int64(105112000000)
		correct := btcIntToFloat(i)
		So(correct, ShouldEqual, 1051.12)
	})
}