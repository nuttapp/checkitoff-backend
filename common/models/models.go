package models

type ServerFields struct {
	Hostname  string `json:"hostname"`
	IPAddress string `json:"ipAddress"`
}

type ClientFields struct {
	DeviceType string `json:"deviceType"`
	ID         string `json:"id"`
	OsVersion  string `json:"osVersion"`
}

type ListFields struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type UserFields struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type EventFields struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}
