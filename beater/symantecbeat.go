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

package beater

import (
	"fmt"
	"time"

	"github.com/marian-craciunescu/symantecbeat/client"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/marian-craciunescu/symantecbeat/config"
)

// Symantecbeat configuration.
type Symantecbeat struct {
	done     chan struct{}
	config   config.Config
	client   beat.Client
	smClient client.SymantecClient
	lastRun  time.Time
}

// New creates an instance of symantecbeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}
	logp.Info("using config %v", c)
	sm := client.NewSymantecClient(c.ApiURL, c.CustomerID, c.DomainID, c.ClientID, c.ClientSecret)

	bt := &Symantecbeat{
		done:     make(chan struct{}),
		config:   c,
		smClient: sm,
		lastRun:  time.Now().UTC().Add(-1 * c.StartDate),
	}
	return bt, nil
}

// Run starts symantecbeat.
func (bt *Symantecbeat) Run(b *beat.Beat) error {
	logp.Info("symantecbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
			{
				logp.Info("Starting ticker cycle at time=%s", time.Now().Format(time.RFC3339))
				err := bt.smClient.GetOauthToken()
				if err != nil {
					logp.Err("Error getting the accesToken.Check credentials error=%s", err.Error())
					continue
				}

				end := time.Now().UTC()
				for eventType := range client.AllTypes {
					t := client.EventType(eventType)
					mapStrArr, err := bt.smClient.DoRequest(bt.lastRun, end, t, bt.config.BatchSize)
					if err != nil {
						logp.Err("Error while doing request.Err=%s", err.Error())
					} else {
						for _, mapStr := range mapStrArr {
							ts := time.Now()
							if err == nil {
								event := beat.Event{
									Timestamp: ts,
									Fields:    mapStr,
								}
								bt.client.Publish(event)
							}

						}

					}
				}
				bt.lastRun = end
			}

		}

	}
}

// Stop stops symantecbeat.
func (bt *Symantecbeat) Stop() {
	_ = bt.client.Close()
	close(bt.done)
}
