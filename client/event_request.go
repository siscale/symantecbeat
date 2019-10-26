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

type EventType int

const (
	AGENT_FRAMEWORK EventType = iota
	APP_CONTROL
	APP_CONTROL_LITE
	APP_CONTROL_WHITELIST
	APP_ISOLATION
	BEHAVIORAL_ANALYSIS
	COMPLIANCE
	DATA_PROTECTION
	DECEPTION
	DETECTION_MONITORING
	DETECTION_RESPONSE
	DEVICE_CONTROL
	EXPLOIT_PROTECTION
	FIREWALL
	LOCATION_MANAGEMENT
	MALWARE_PROTECTION
	NETWORK_INTEGRITY
	NETWORK_IPS
	POLICY_MANAGER
	ROAMING_CLIENT
	TAMPER_PROTECTION
	TDAD_PROTECT
	TELEMETRY
	VR_ASSESSMENT
	VR_REMEDIATION
	WEB_SECURITY
)

var AllTypes = []EventType{
	AGENT_FRAMEWORK,
	APP_CONTROL,
	APP_CONTROL_LITE,
	APP_CONTROL_WHITELIST,
	APP_ISOLATION,
	BEHAVIORAL_ANALYSIS,
	COMPLIANCE,
	DATA_PROTECTION,
	DECEPTION,
	DETECTION_MONITORING,
	DETECTION_RESPONSE,
	DEVICE_CONTROL,
	EXPLOIT_PROTECTION,
	FIREWALL,
	LOCATION_MANAGEMENT,
	MALWARE_PROTECTION,
	NETWORK_INTEGRITY,
	NETWORK_IPS,
	POLICY_MANAGER,
	ROAMING_CLIENT,
	TAMPER_PROTECTION,
	TDAD_PROTECT,
	TELEMETRY,
	VR_ASSESSMENT,
	VR_REMEDIATION,
	WEB_SECURITY,
}

func (t EventType) String() string {
	names := [...]string{
		"AGENT FRAMEWORK ",
		"APP CONTROL",
		"APP CONTROL LITE",
		"APP CONTROL WHITELIST",
		"APP ISOLATION",
		"BEHAVIORAL ANALYSIS",
		"COMPLIANCE",
		"DATA PROTECTION",
		"DECEPTION",
		"DETECTION MONITORING",
		"DETECTION RESPONSE",
		"DEVICE CONTROL",
		"EXPLOIT PROTECTION",
		"FIREWALL",
		"LOCATION MANAGEMENT",
		"MALWARE PROTECTION",
		"NETWORK INTEGRITY",
		"NETWORK IPS",
		"POLICY MANAGER",
		"ROAMING CLIENT",
		"TAMPER PROTECTION",
		"TDAD PROTECT",
		"TELEMETRY",
		"VR ASSESSMENT",
		"VR REMEDIATION",
		"WEB SECURITY",
	}

	if t < AGENT_FRAMEWORK || t > WEB_SECURITY {
		return "Unknown"
	}
	return names[t]
}

const timeFormat = "2006-01-02T15:04:05.999Z"

type eventRequest struct {
	BatchSize  int    `json:"batchSize"`
	EventsType string `json:"type"`
	StartDate  string `json:"startDate"`
	EndDate    string `json:"endDate"`
}

func NewEventEncoded(s, end time.Time, size int, t EventType) ([]byte, error) {
	event := eventRequest{
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
