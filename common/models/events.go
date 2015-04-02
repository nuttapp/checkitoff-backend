package models

const (
	// Catch fat finger mistakes where an event type is not equal to it's intended event type
	// ex: a CreateListMsg should never have an EventType "UpdateListEvent"
	InvalidMsgMethodError   = "Invalid msg method"
	InvalidMsgResourceError = "Invalid msg resource"

	// base event fields
	MissingClientIDError         = "Msg client id cannot be empty"
	MissingClientDeviceTypeError = "Msg client deviceType cannot be empty"
	MissingUserIDError           = "Msg user id cannot be empty"
	MissingEventIDError          = "Msg id cannot be empty"
	MissingMsgMethodError        = "Msg method cannot be empty"
	MissingMsgResourceError      = "Msg resource cannot be empty"
	MissingServerHostnameError   = "Msg server hostname cannot be empty"
	MissingServerIPAddressError  = "Msg server ip address cannot be empty"

	MissingListIDError    = "List id cannot be emtpy"
	MissingListTitleError = "List title cannot be empty"
)

type Event interface {
	IsReadyToBeSaved() error
}
