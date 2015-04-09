package config

type Config struct {
	Hostname           string // host name of this server
	IPAddress          string // ip of this server
	NSQProducerTCPAddr string
	NSQTopic           string
}
