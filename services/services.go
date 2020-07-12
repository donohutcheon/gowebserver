package services

import (
	"github.com/donohutcheon/gowebserver/services/users"
	"github.com/donohutcheon/gowebserver/state"
)

func StartServices(state *state.ServerState) {
	state.ShutdownWG.Add(1)
	go users.ConfirmUsersForever(state)
}