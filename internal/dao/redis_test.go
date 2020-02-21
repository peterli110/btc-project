package dao


import (
	"btc-project/internal/model"
	"context"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)


func TestSetQueryCache(t *testing.T) {
	Convey("TestSetQueryCache add query and data to redis", t, func() {
		for i := 0; i < 10; i++ {
			query := &model.BtcQuery{
				StartTime: time.Date(2020, 02, 10, i, 0, 0, 0, time.UTC),
				EndTime: time.Date(2020, 02, 10, i + 2, 0, 0, 0, time.UTC),
			}

			res := []*model.BtcResponse{
				{
					DateTime: "2020-02-10T08:00:00+00:00",
					Amount: float64(i) + 100,
				},
			}

			err := d.SetQueryCache(context.TODO(), query, res)
			So(err, ShouldBeNil)
		}
	})
}


func TestGetQueryCache(t *testing.T) {
	Convey("TestGetQueryCache get data from cache", t, func() {
		for i := 0; i < 10; i++ {
			query := &model.BtcQuery{
				StartTime: time.Date(2020, 02, 10, i, 0, 0, 0, time.UTC),
				EndTime: time.Date(2020, 02, 10, i + 2, 0, 0, 0, time.UTC),
			}
			_, ok := d.GetQueryCache(context.TODO(), query)
			So(ok, ShouldBeTrue)
		}
	})
}



func TestGetKeyFromQuery(t *testing.T) {
	Convey("TestGetKeyFromQuery verify the key is correct", t, func() {
		for i := 0; i < 10; i++ {
			query := &model.BtcQuery{
				StartTime: time.Date(2020, 02, 10, i, 0, 0, 0, time.UTC),
				EndTime: time.Date(2020, 02, 10, i + 2, 0, 0, 0, time.UTC),
			}

			correct := fmt.Sprintf("2020-02-10T%02d:00:00Z_2020-02-10T%02d:00:00Z", i, i + 2)
			res := getKeyFromQuery(query)
			So(res, ShouldEqual, correct)
		}
	})
}