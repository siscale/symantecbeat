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
