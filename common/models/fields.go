package models

type Server struct {
	Hostname  string `json:"hostname"`
	IPAddress string `json:"ipAddress"`
}

type Client struct {
	DeviceType string `json:"deviceType"`
	ID         string `json:"id"`
	OsVersion  string `json:"osVersion"`
}

type List struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type EventFields struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}
