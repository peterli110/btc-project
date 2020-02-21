package http

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func TestAddCoin(t *testing.T) {
	e := httpTester(t)

	invalidReq := []map[string]interface{}{
		{
			"datetime": "2020-02-10",
			"amount": 1.23,
		},
		{
			"datetime": "2019-10-05T14:48:01+01:00",
			"amount": 100.123456789,
		},
		{
			"datetime": "3019-10-05T14:48:01+01:00",
			"amount": 1000000000000000000,
		},
		{
			"datetime": "1019-10-05T14:48:01+01:00",
			"amount": 1.23,
		},
		{
			"date": "2019-10-05T14:48:01+01:00",
			"amount": 1.34,
		},
	}

	for _, v := range invalidReq {
		e.POST("/coin/add").WithHeader("Content-Type", "application/json").
			WithJSON(v).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsKey("code").
			ValueEqual("code", -1)
	}
}



func TestQueryCoin(t *testing.T) {
	e := httpTester(t)

	invalidReq := []map[string]interface{}{
		{
			"startDatetime": "2020-02-10",
			"endDatetime": "2019-10-05T14:48:01+01:00",
		},
		{
			"startDatetime": "2019-10-05T14:49:01+01:00",
			"endDatetime": "2019-10-05T14:48:01+01:00",
		},
		{
			"startDatetime": "2007-10-04T14:48:01+01:00",
			"endDatetime": "2019-10-05T14:48:01+01:00",
		},
		{
			"startDatetime": "2019-10-05T14:48:01+01:00",
			"endDatetime": "2050-10-05T14:48:01+01:00",
		},
	}

	for _, v := range invalidReq {
		e.POST("/coin/search").WithHeader("Content-Type", "application/json").
			WithJSON(v).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsKey("code").
			ValueEqual("code", -1)
	}
}


func TestNumDecPlaces(t *testing.T) {
	Convey("TestNumDecPlaces to test number of decimal places", t, func() {
		que := []float64{123, 1.234, 1.23456, 1.2345678901, 12345.6}
		ans := []int{0, 3, 5, 10, 1}

		for i, v := range que {
			r := numDecPlaces(v)
			So(r, ShouldEqual, ans[i])
		}
	})
}