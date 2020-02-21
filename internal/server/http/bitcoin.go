package http

import (
	"btc-project/internal/model"
	"btc-project/library/ecode"
	"btc-project/library/log"
	bm "btc-project/library/net/http/blademaster"
	"strconv"
	"strings"
)

// addCoin add btc by timestamp
/*
	request body:
	{
		datetime: a time string can be parsed to time.Time
		amount: a float64 number larger than 0 and have no more than 8 decimal places
	}
 */
func addCoin(c *bm.Context) {
	var m model.BtcJSON
	if err := c.Bind(&m); err != nil {
		log.Error("Parse JSON Error: (%v)", err)
		return
	}

	// maximum decimal place is 8 since the smallest
	// unit of btc is 0.00000001
	if decimal := numDecPlaces(m.Amount); decimal > 8 {
		log.Error("too many decimal places: (%v)", m)
		c.JSON(nil, ecode.InvalidParams)
		c.Abort()
		return
	}

	// BTC is invented in 2009, so the year before 2008 is invalid
	if m.DateTime.Year() <= 2008 {
		log.Error("invalid date: (%v)", m)
		c.JSON(nil, ecode.InvalidParams)
		c.Abort()
		return
	}

	if err := svc.AddBitCoin(c, &m); err != nil {
		c.JSON(nil, err)
		c.Abort()
		return
	}

	c.JSON(nil, nil)
}


// queryCoin get btc by time range
/*
	request body:
	{
		startDatetime: a time string can be parsed to time.Time
		endDatetime: a time string can be parsed to time.Time
	}
	endDateTime must equal to or after startDatetime
*/
func queryCoin(c *bm.Context) {
	var m model.BtcQuery
	if err := c.Bind(&m); err != nil {
		log.Error("Parse JSON Error: (%v)", err)
		return
	}

	if m.EndTime.Before(m.StartTime) {
		log.Error("startTime after endTime: (%v)", m)
		c.JSON(nil, ecode.InvalidParams)
		c.Abort()
		return
	}

	// BTC is invented in 2009, so the year before 2008 is invalid
	if m.StartTime.Year() <= 2008 || m.EndTime.Year() <= 2008 {
		log.Error("invalid date: (%v)", m)
		c.JSON(nil, ecode.InvalidParams)
		c.Abort()
		return
	}

	d, err := svc.QueryBitCoin(c, &m)
	if err != nil {
		c.JSON(nil, err)
		c.Abort()
		return
	}
	c.JSON(d, nil)
}


// numDecPlaces will return the count of decimal places
func numDecPlaces(v float64) int {
	s := strconv.FormatFloat(v, 'f', -1, 64)
	i := strings.IndexByte(s, '.')
	if i > -1 {
		return len(s) - i - 1
	}
	return 0
}
