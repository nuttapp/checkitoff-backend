package config

type Config struct {
	Hostname           string // host name of this server
	IPAddress          string // ip of this server
	Role               string // api
	NSQProducerTCPAddr string
	NSQLookupdHTTPAddr string
	NSQTopic           string
}

func NewConfig() *Config {
	config := &Config{
		Hostname:           "127.0.0.1",
		IPAddress:          "127.0.0.1",
		Role:               "api",
		NSQProducerTCPAddr: "127.0.0.1:4150",
		NSQLookupdHTTPAddr: "127.0.0.1:4161",
		NSQTopic:           "api_messages",
	}
	return config
}
