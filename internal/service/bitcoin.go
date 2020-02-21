package service

import (
	"btc-project/internal/model"
	"btc-project/library/ecode"
	"btc-project/library/log"
	"context"
	"encoding/json"
	"github.com/shopspring/decimal"
	"time"
)


var (
	convertRatio = decimal.NewFromInt(100000000)
)

// AddBitCoin will insert data into database.
// Since float calculation might lose precision,
// we'll use int64 to store the amount data.
func (s *Service) AddBitCoin(c context.Context, m *model.BtcJSON) error {
	data := &model.BtcData{
		DateTime: m.DateTime.UTC(),
		Amount: btcFloatToInt(m.Amount),
	}

	return s.Dao.InsertCoin(c, data)
}


// QueryBitCoin will get the query result
func (s *Service) QueryBitCoin(c context.Context, m *model.BtcQuery) (r []*model.BtcResponse, err error) {
	alignCoinQuery(m)

	// read from cache
	var ok bool
	if r, ok = s.Dao.GetQueryCache(c, m); ok {
		return r, nil
	}

	// cache not found, read from db
	result, err := s.Dao.QueryCoin(c, m)
	if err != nil {
		return nil, err
	}

	r = make([]*model.BtcResponse, len(result))
	for i, v := range result {
		// v should be the slice with 2 elements,
		// v[0] is time and v[1] is value,
		// type is json.Number

		// use ok to avoid panic
		timeVal, ok := v[0].(json.Number)
		if !ok {
			log.Error("type assertion: (%v)", v)
			return nil, ecode.ServerErr
		}
		amountVal, ok := v[1].(json.Number)
		if !ok {
			log.Error("type assertion: (%v)", v)
			return nil, ecode.ServerErr
		}

		// get values
		tm, err := timeVal.Int64()
		if err != nil {
			log.Error("parse int64: (%v)", err)
			return nil, ecode.ServerErr
		}
		amount, err := amountVal.Int64()
		if err != nil {
			log.Error("parse int64: (%v)", err)
			return nil, ecode.ServerErr
		}

		// parse time and set amount,
		// hard code time zone to match output requirement
		// current timezone is set to the Greenwich
		// standard time. If use Z07:00, timezone in +0
		// will be omitted automatically.
		r[i] = &model.BtcResponse{
			DateTime: time.Unix(0, tm).Format("2006-01-02T15:04:05+00:00"),
			Amount: btcIntToFloat(amount),
		}
	}

	// set cache, err can be omitted
	_ = s.Dao.SetQueryCache(c, m, r)
	return r, nil
}



// btcFloatToInt will convert float64 to int64
// based on smallest unit of bitcoin.
// also use decimal library to avoid loss of precision
// for example 1051.12 * 100000000 = 105111999999
func btcFloatToInt(f float64) int64 {
	s := decimal.NewFromFloat(f)
	s = s.Mul(convertRatio)
	return s.IntPart()
}

// btcIntToFloat will convert int64 to float64
// based on smallest unit of bitcoin.
func btcIntToFloat(i int64) float64 {
	s := decimal.NewFromInt(i)
	s = s.Div(convertRatio)
	res, _ := s.Float64()
	return res
}


// alignCoinQuery will rotate the startTime and endTime
// to the beginning of its hour. Also it helps to caching
// the results.
func alignCoinQuery(m *model.BtcQuery) {
	m.StartTime = roundTimeHelper(m.StartTime)
	m.EndTime = roundTimeHelper(m.EndTime)
}

// roundTimeHelper helper function to align query
func roundTimeHelper(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
}