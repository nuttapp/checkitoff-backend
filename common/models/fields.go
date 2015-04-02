package models

import "time"

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
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Title     string    `json:"title"`
	Users     []string  `json:"users"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Msg struct {
	Method   string `json:"method"`
	Resource string `json:"resource"`
	ID       string `json:"id"`
}
