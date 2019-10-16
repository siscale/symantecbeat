// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period time.Duration `config:"period"`
	ApiURL string `config:"url"`
	CustomerID  string `config:"customer_id"`
	DomainID  string `config:"domain_id"`
	ClientID  string `config:"client_id"`
	ClientSecret  string `config:"client_sercret"`
}

var DefaultConfig = Config{
	Period: 1 * time.Second,
}
