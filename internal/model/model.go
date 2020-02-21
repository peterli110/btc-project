package model

import "time"

// BtcJSON is the struct to handle the format of request
type BtcJSON struct {
	DateTime		time.Time		`json:"datetime" validate:"required,lte"`
	Amount			float64			`json:"amount" validate:"gt=0"`
}

// BtcResponse is the struct to handle the format of response
type BtcResponse struct {
	DateTime		string			`json:"datetime"`
	Amount			float64			`json:"amount"`
}

// BtcData convert the amount to int64 which can be stored in database
type BtcData struct {
	DateTime		time.Time		`bson:"datetime"`
	Amount			int64			`bson:"amount"`
}

// BtcQuery is the struct to handle search request
type BtcQuery struct {
	StartTime		time.Time		`json:"startDatetime" validate:"required,lte"`
	EndTime			time.Time		`json:"endDatetime" validate:"required,lte"`
}