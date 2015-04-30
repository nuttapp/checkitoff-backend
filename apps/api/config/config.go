package config

type NSQ struct {
	ProducerTCPAddr string
	LookupdHTTPAddr string
	PubTopic        string
	SubTopic        string
	SubChannel      string
}

type Server struct {
	Hostname  string // host name of this server
	IPAddress string // ip of this server
	Role      string // api
}

type Config struct {
	Server *Server
	NSQ    *NSQ
	Redis  *NSQ
}

func NewConfig() *Config {
	config := &Config{
		Server: &Server{
			Hostname:  "api01",
			IPAddress: "127.0.0.1",
			Role:      "api",
		},
		NSQ: &NSQ{
			ProducerTCPAddr: "127.0.0.1:4150",
			LookupdHTTPAddr: "127.0.0.1:4161",
			PubTopic:        "api_messages",
			SubTopic:        "api_replies",
			SubChannel:      "client_handler",
		},
	}
	return config
}
