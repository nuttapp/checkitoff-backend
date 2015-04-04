package models

type Event interface {
	IsReadyToBeSaved() error
}
