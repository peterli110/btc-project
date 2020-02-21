package blademaster

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"btc-project/library/ecode"
	"btc-project/library/log"
	"btc-project/library/net/metadata"
)

// Logger is logger  middleware
func Logger() HandlerFunc {
	const noUser = "no_user"
	return func(c *Context) {
		var bodyBuffer []byte
		if c.Request.Body != nil {
			bodyBuffer, _ = ioutil.ReadAll(c.Request.Body) // after this operation body will equal 0
			// Restore the io.ReadCloser to request
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBuffer))
		}

		now := time.Now()
		ip := metadata.String(c, metadata.RemoteIP)
		req := c.Request
		path := req.URL.Path
		params := string(bodyBuffer)
		var quota float64
		if deadline, ok := c.Context.Deadline(); ok {
			quota = time.Until(deadline).Seconds()
		}

		c.Next()

		// ignore OPTIONS
		if c.Request.Method == "OPTIONS" {
			return
		}

		err := c.Error
		cerr := ecode.Cause(err)

		if cerr.Code() == 1 {
			err = nil
		}

		dt := time.Since(now)
		//caller := metadata.String(c, metadata.Caller)
		caller := c.Request.Header.Get(_httpHeaderUser)
		if caller == "" {
			caller = noUser
		}

		stats.Incr(caller, path[1:], strconv.FormatInt(int64(cerr.Code()), 10))
		stats.Timing(caller, int64(dt/time.Millisecond), path[1:])

		lf := log.Infov
		errmsg := ""
		isSlow := dt >= (time.Millisecond * 500)
		if err != nil {
			errmsg = err.Error()
			lf = log.Warnv
			if cerr.Code() < 0 {
				lf = log.Errorv
			}
		} else {
			if isSlow {
				lf = log.Warnv
			}
		}
		lf(c,
			log.KVString("method", req.Method),
			log.KVString("ip", ip),
			log.KVString("user", caller),
			log.KVString("path", path),
			log.KVString("params", params),
			log.KVInt("ret", cerr.Code()),
			log.KVString("msg", cerr.Message()),
			log.KVString("stack", fmt.Sprintf("%+v", err)),
			log.KVString("err", errmsg),
			log.KVFloat64("timeout_quota", quota),
			log.KVFloat64("ts", dt.Seconds()),
			log.KVString("source", "http-access-log"),
		)
	}
}
