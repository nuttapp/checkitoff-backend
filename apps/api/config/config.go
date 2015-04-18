package config

type Config struct {
	Hostname           string // host name of this server
	IPAddress          string // ip of this server
	Role               string // api
	NSQProducerTCPAddr string
	NSQLookupdHTTPAddr string
	NSQPubTopic        string
	NSQSubTopic        string
}

func NewConfig() *Config {
	config := &Config{
		Hostname:           "api01",
		IPAddress:          "127.0.0.1",
		Role:               "api",
		NSQProducerTCPAddr: "127.0.0.1:4150",
		NSQLookupdHTTPAddr: "127.0.0.1:4161",
		NSQPubTopic:        "api_messages",
		NSQSubTopic:        "api01",
	}
	return config
}
