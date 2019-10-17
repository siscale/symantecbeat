package beater

import (
	"fmt"
	"github.com/marian-craciunescu/symantecbeat/client"
	"time"

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
	logp.Info("using config %s", c)
	sm := client.NewSymantecClient(c.ApiURL, c.CustomerID, c.DomainID, c.ClientID, c.ClientSecret)

	bt := &Symantecbeat{
		done:     make(chan struct{}),
		config:   c,
		smClient: sm,
		lastRun:  c.StartDate,
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
				err := bt.smClient.GetOauthToken()
				if err != nil {
					logp.Err("Error getting the accesToken.Check credentials", err.Error())
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
							tsa, err := mapStr.GetValue("timestamp")
							if err == nil {
								ts, err = time.Parse("2006-01-02T15:04:05.999Z", tsa.(string))
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
				}
				bt.lastRun = end
			}

		}

	}
}

// Stop stops symantecbeat.
func (bt *Symantecbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
