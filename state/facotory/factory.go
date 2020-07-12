package facotory

import (
	"context"
	"github.com/donohutcheon/gowebserver/datalayer"
	"github.com/donohutcheon/gowebserver/datalayer/mockdatalayer"
	"github.com/donohutcheon/gowebserver/provider/mail"
	"github.com/donohutcheon/gowebserver/provider/mail/mailtrap"
	"github.com/donohutcheon/gowebserver/provider/mail/mockmail"
	"github.com/donohutcheon/gowebserver/router"
	"github.com/donohutcheon/gowebserver/server"
	"github.com/donohutcheon/gowebserver/services"
	"github.com/donohutcheon/gowebserver/state"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"
)

type environment int

const (
	staging environment = 0
	prod environment = 1
)



var (
	//CertFile environment variable for CertFile
	CertFile = os.Getenv("CERT_FILE")
	//KeyFile environment variable for KeyFile
	KeyFile = os.Getenv("KEY_FILE")
)

func newState(env environment, logger *log.Logger, mainThreadWG *sync.WaitGroup) (*state.ServerState, error) {
	ctx, cancel := context.WithCancel(context.Background())
	dataLayer, err := datalayer.New()
	if err != nil {
		return nil, err
	}

	s := &state.ServerState{
		URL: os.Getenv("URL"),
		Channels: state.Channels{
			ConfirmUsers: make(chan datalayer.User, 1),
		},
		Context: ctx,
		Logger:    logger,
		DataLayer: dataLayer,
		ShutdownWG: new(sync.WaitGroup),
		Router: mux.NewRouter(),
		Cancel: cancel,
	}

	if env == prod {
		s.Providers.Email = mail.Client(mailtrap.New(s))
	} else {
		s.Providers.Email = mail.Client(mailtrap.New(s))
	}

	services.StartServices(s)

	mainThreadWG.Add(2)
	go handleSignals(s, mainThreadWG)
	go runServer(s, mainThreadWG)

	return s, nil
}

func NewForProduction(logger *log.Logger, mainThreadWG *sync.WaitGroup) (*state.ServerState, error) {
	s, err := newState(prod, logger, mainThreadWG)
	if err != nil {
		return s, err
	}

	return s, nil
}

func NewForStaging(logger *log.Logger, mainThreadWG *sync.WaitGroup) (*state.ServerState, error) {
	s, err := newState(staging, logger, mainThreadWG)
	if err != nil {
		return s, err
	}

	return s, nil
}

func NewForTesting(t *testing.T, callbacks *state.MockCallbacks) *state.ServerState {
	logger := log.New(os.Stdout, "microservice", log.LstdFlags|log.Lshortfile)
	ctx := context.Background()
	mockDataLayer := mockdatalayer.New(t)

	mail := &mockmail.MockClient{
		T:            t,
		Context:      ctx,
		CallbackFunc: callbacks.MockMail,
		Group:        callbacks.MockMailWG,
	}

	r := mux.NewRouter()
	state := &state.ServerState{
		Channels: state.Channels{
			ConfirmUsers: make(chan datalayer.User, 1),
		},
		Context:    ctx,
		Logger:     logger,
		ShutdownWG: new(sync.WaitGroup),
		DataLayer:  mockDataLayer,
		Router:     r,
		Providers: state.Providers{
			Email: mockmail.New(mail),
		},
	}

	h := router.NewHandlers(state)
	err := h.SetupRoutes(r)
	require.NoError(t, err)

	srv := server.New(r, "", "0")
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	services.StartServices(state)
	go func() {
		err := srv.Serve(l)
		require.NoError(t, err)
	} ()
	state.URL = "http://" + strings.Replace(l.Addr().String(), "[::]", "localhost", 1)

	return state
}

func runServer(state *state.ServerState, mainThreadWG *sync.WaitGroup) {
	defer mainThreadWG.Done()

	logger := state.Logger
	h := router.NewHandlers(state)

	router := state.Router
	err := h.SetupRoutes(router)
	if err != nil {
		logger.Fatalf("Could not start router %s", err.Error())
	}

	//ServiceAddress address to listen on
	bindAddress := os.Getenv("BIND_ADDRESS")
	port        := os.Getenv("PORT")
	logger.Printf("Server Binding to %s:%s", bindAddress, port)
	srv := server.New(router, bindAddress, port)

	go func() {
		// TODO: Put back in for TLS
		/*err := srv.ListenAndServeTLS(CertFile, KeyFile)*/
		err := srv.ListenAndServe() //Launch the app, visit localhost:8000/api
		if err != nil && err != http.ErrServerClosed{
			logger.Fatalf("Server failed to start %s", err.Error())
		}
	}()

	logger.Printf("server started")

	<-state.Context.Done()

	logger.Printf("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		logger.Fatalf("server Shutdown Failed: %s", err.Error())
	}

	logger.Printf("server exited properly")
}

func handleSignals(state *state.ServerState, mainThreadWG *sync.WaitGroup) {
	defer mainThreadWG.Done()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("waiting for system call...")
	signalChan := <-c
	log.Printf("system call: %+v", signalChan)
	// Close all channels here and then wait for the wait group to unlock.
	close(state.Channels.ConfirmUsers)
	state.ShutdownWG.Wait() //Wait for consumers to finish processing messages and exit
	state.Cancel()
}