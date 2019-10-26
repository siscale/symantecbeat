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

// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period       time.Duration `config:"period"`
	ApiURL       string        `config:"url"`
	CustomerID   string        `config:"customer_id"`
	DomainID     string        `config:"domain_id"`
	ClientID     string        `config:"client_id"`
	ClientSecret string        `config:"client_secret"`
	BatchSize    int           `config:"batch_size"`
	StartDate    time.Duration `config:"start_date"`
}

var DefaultConfig = Config{
	Period:    5 * time.Minute,
	StartDate: 60 * time.Minute,
	BatchSize: 1000,
	ApiURL:    "https://usea1.r3.securitycloud.symantec.com/r3_epmp_i",
}
