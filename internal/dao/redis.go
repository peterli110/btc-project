package dao

import (
	"btc-project/internal/model"
	"btc-project/library/log"
	"context"
	"encoding/json"
	"time"
)


const cacheExpire = 5 * time.Minute


// SetQueryCache will store the result of query in redis if redis is available
// Since the result is not real-time sensitive, we can set the expiration time
// to 5 minutes to reduce the pressure of database.
func (dao *Dao) SetQueryCache(c context.Context, m *model.BtcQuery, v []*model.BtcResponse) error {
	key := getKeyFromQuery(m)

	value, err := json.Marshal(v)
	if err != nil {
		log.Error("JSON Marshal: (%v)", err)
		return err
	}

	if _, err := dao.Redis.WithContext(c).Set(key, value, cacheExpire).Result(); err != nil {
		log.Error("redis set: (%v)", err)
		return err
	}

	return nil
}


// GetQueryCache will try to get result of query in redis
// if redis is disabled or no result found or an error occurred,
// the second return value will be false
func (dao *Dao) GetQueryCache(c context.Context, m *model.BtcQuery) ([]*model.BtcResponse, bool) {
	var v []*model.BtcResponse
	key := getKeyFromQuery(m)

	value, err := dao.Redis.WithContext(c).Get(key).Result()
	if err != nil {
		if err.Error() != "redis: nil" {
			log.Error("redis get: (%v)", err)
			return nil, false
		}
		return nil, false
	}

	if err = json.Unmarshal([]byte(value), &v); err != nil {
		log.Error("JSON Unmarshal: (%v)", err)
		return nil, false
	}

	return v, true
}


// getKeyFromQuery will generate redis key from query
// with format "startTime_endTime" in RFC3339 format
// e.g. "2006-01-02T15:04:05Z07:00_2006-01-02T15:04:05Z07:00"
// timezone will be converted to UTC
func getKeyFromQuery(m *model.BtcQuery) string {
	start := m.StartTime.UTC().Format(time.RFC3339)
	end := m.EndTime.UTC().Format(time.RFC3339)
	return start + "_" + end
}