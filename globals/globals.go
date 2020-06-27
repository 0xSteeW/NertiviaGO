package globals

import (
	globalconstants "nertivia/api/constants"
)

type Constants struct {
	EndpointURL string
	WebsocketURL string
	EndpointUser string
	EndpointChannel string
	EndpointServer string
}

func ReadConstants() *Constants {
	cst := new(Constants)
	cst.EndpointURL = globalconstants.EndpointURL
	cst.WebsocketURL = globalconstants.WebsocketURL
	cst.EndpointUser = globalconstants.EndpointUser
	cst.EndpointChannel = globalconstants.EndpointChannel
	cst.EndpointServer = globalconstants.EndpointServer
	return cst
}
