// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package client

import (
	"encoding/json"
	"time"
)

const timeFormatSearch = "2006-01-02T15:04:05.000+00:00"

type eventSearchRequest struct {
	Limit      int    `json:"limit"`
	EventsType string `json:"feature_name"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Product    string `json:"product"`
	Next       int    `json:"next"`
}

// NewEventSearchEncoded create a []byte json encoded for doing a search request
func NewEventSearchEncoded(s, end time.Time, size, next int, t EventType) ([]byte, error) {
	//needed because Symantec endpoints gives random errors when using now.
	s = s.Add(-1 * 10 * time.Second)
	end = end.Add(-1 * 10 * time.Second)
	event := eventSearchRequest{
		StartDate:  s.UTC().Format(timeFormatSearch),
		EndDate:    end.UTC().Format(timeFormatSearch),
		Limit:      size,
		Next:       next,
		EventsType: t.String(),
		Product:    "SAEP",
	}

	jsonValue, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	return jsonValue, nil
}
