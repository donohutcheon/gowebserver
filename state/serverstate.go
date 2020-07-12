package state

import (
	"context"
	"github.com/donohutcheon/gowebserver/datalayer"
	"github.com/donohutcheon/gowebserver/provider/mail"
	"github.com/donohutcheon/gowebserver/provider/mail/mockmail"
	"github.com/gorilla/mux"
	"log"
	"sync"
)

type Channels struct {
	ConfirmUsers chan  datalayer.User
}

type Providers struct {
	Email      mail.Client
}

type ServerState struct {
	URL string
	Channels   Channels
	Context    context.Context
	DataLayer  datalayer.DataLayer
	Logger     *log.Logger
	ShutdownWG *sync.WaitGroup
	Router     *mux.Router
	Providers  Providers
	Cancel     context.CancelFunc
}

type MockCallbacks struct {
	MockMail mockmail.CallbackFunc
	MockMailWG *sync.WaitGroup
}

func NewMockCallbacks(callback mockmail.CallbackFunc) *MockCallbacks{
	m := new(MockCallbacks)
	m.MockMail = callback
	m.MockMailWG = new(sync.WaitGroup)
	return m
}