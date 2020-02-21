[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0)
[![Build Status](https://travis-ci.org/peterli110/btc-project.svg?branch=master)](https://travis-ci.org/peterli110/btc-project)

# btc-project

## About this repo

This repository is the interview project for Anymind Group.

## Directories

| directory            | description                                                  |
| -------------------- | ------------------------------------------------------------ |
| cmd                  | The main entry of server.                                    |
| conf                 | Config initialization.                                       |
| config               | This directory includes config files.                        |
| internal/server/http | Routers and http handlers, which can handle incoming request and verify parameters. |
| internal/service     | Implementation of service logics.                            |
| internal/dao         | Data access object will provide the access of database.      |
| library              | This directory includes frequently used libraries such as the http server and database driver. |
| vendor               | Go vendor, includes the third party libraries from go packages. |
|                      |                                                              |

## Frameworks and core dependencies

HTTP Framework: [blademaster](https://godoc.org/github.com/bilibili/kratos/pkg/net/http/blademaster) is an open source web framework based on [Gin](https://github.com/gin-gonic/gin), which is developed by [bilibili](https://github.com/bilibili) and is an excellent framework for production environment.

Cache: [redis](https://redis.io/) is the most widespread memory cache in the world.

Database: [influxDB](https://www.influxdata.com/developers/) is an open source time series database.

## Installation

1. Make sure [redis](https://redis.io/) and [influxDB](https://www.influxdata.com/developers/) are installed and configured.
2. Clone this project `git clone https://github.com/peterli110/btc-project.git $GOPATH/src/btc-project`
3. Modify the `[redis]` and `[influxdb]` part in `./configs/config.toml` if necessary.
4. Build the server by `cd cmd` and `go build`.
5. Run the server by `./cmd`

## Test

1. Make sure [redis](https://redis.io/) and [influxDB](https://www.influxdata.com/developers/) are installed and configured.
2. Modify the `[redis]` and `[influxdb]` part in `./configs/config_test.toml` if necessary.
3. Test with `go test ./internal/... -v`

## Instruction

##### API

This server provides 2 APIs which accept requests with `Content-Type: application/json`.

**Please note that since BTC is invented at the beginning of 2009, so all the time string must after 2008 and before or equal to current time.**

`POST /coin/add` 

```json
{
	"datetime": "2020-02-19T12:34:00+00:00",
	"amount": 1.23
}
```
while datetime is a valid RFC 3339 time string,

amount is a positive float number which is the amount of btc.

Sample request:

```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"datetime":"2020-02-19T10:28:00+00:00","amount":751.12}' \
  http://localhost:8000/coin/add
```

Response will be:

```json
{
    "code": 0
}
```



`POST /coin/search` 

```json
{
	"startDatetime": "2020-02-19T12:34:00+00:00",
	"endDatetime": "2020-02-19T12:34:00+00:00"
}
```
while startDatetime and endDatetime are valid RFC 3339 time strings,

and endDatetime must after startDatetime.


Sample request:

```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"startDatetime":"2020-02-19T08:12:00+00:00","endDatetime":"2020-02-20T08:32:00+00:00"}' \
  http://localhost:8000/coin/search
```

Response will be:

```json
{
    "code": 0,
    "data": [
        {
            "datetime": "2020-02-19T10:00:00+00:00",
            "amount": 751.12
        },
        {
            "datetime": "2020-02-19T12:00:00+00:00",
            "amount": 1051.12
        }
    ]
}
```



##### Response Format

Response will be always with a status http.StatusOK(200), and response body will be formatted in JSON, a sample response will be

```json
{
		"code": 0,
		"data": [
		            "data will include search result as the response of /coin/search",
		            "and will not exist as the response of /coin/add"
                ]
}
```

Response Code

To simplify the project, I only use 3 response codes for convenience:

| Code | Explanation           |
| ---- | --------------------- |
| 0    | Request success       |
| -1   | Invalid parameters    |
| -500 | Server internal error |

## Constraints

Since the smallest unit of BTC is `10^-8`, and the float calculation will lose precision, I converted the BTC to the smallest unit and use int64 to store and calculate. The InfluxDB  is also use int64 to store integer values. 

The maximum value of BTC that this project could store is `MaxInt64 * 10^-8 = 9 * 10^10`.

