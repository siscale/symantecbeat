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

const timeFormat = "2006-01-02T15:04:05.999Z"

type eventExportRequest struct {
	BatchSize  int    `json:"batchSize"`
	EventsType string `json:"type"`
	StartDate  string `json:"startDate"`
	EndDate    string `json:"endDate"`
}

// NewEventExportEncoded create a []byte json encoded for doing a deprecated export request
func NewEventExportEncoded(s, end time.Time, size int, t EventType) ([]byte, error) {
	event := eventExportRequest{
		StartDate:  s.Format(timeFormat),
		EndDate:    end.Format(timeFormat),
		BatchSize:  size,
		EventsType: t.String(),
	}

	jsonValue, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	return jsonValue, nil
}
