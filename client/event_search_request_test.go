package client

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewTime(t *testing.T) {

	a := assert.New(t)

	millis := int64(1438167001716)

	toFormat := time.Unix(0, millis*int64(time.Millisecond))
	formatTime := toFormat.Format(timeFormatSearch)

	a.Equal("2015-07-29T13:50:01.716+00:00", formatTime)
}

func TestNewEventSearchEncoded(t *testing.T) {

	a := assert.New(t)

	millis := int64(1438167001716)
	millis2 := int64(1438167001716 + 200000)

	start := time.Unix(0, millis*int64(time.Millisecond))
	end := time.Unix(0, millis2*int64(time.Millisecond))

	b, err := NewEventSearchEncoded(start, end, 10, 0, ALL)
	a.NoError(err)
	a.Equal([]byte(`{"limit":10,"feature_name":"ALL","start_date":"2015-07-29T10:50:01.716+00:00","end_date":"2015-07-29T10:53:21.716+00:00","product":"SAEP","next":0}`), b)
}
