package e2e

import (
	"fmt"
	"log"
	"net"
	"net/http/httptest"
	"testing"
	"time"

	httpexpect "github.com/gavv/httpexpect/v2"
	"github.com/jordanlanch/modak-test/api/route"
	"github.com/jordanlanch/modak-test/bootstrap"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

// Setup set up database and server for E2E test
func Setup(t *testing.T, cassetteName string) (expect *httpexpect.Expect, teardown func()) {

	t.Helper()

	app := bootstrap.App("../test/.env")

	env := app.Env

	timeout := time.Duration(env.ContextTimeout) * time.Second

	// Create new VCR cassette
	recorder, err := recorder.New(cassetteName)
	if err != nil {
		log.Fatal(err)
	}

	// Use the recorder for all requests
	httpClient := recorder.GetDefaultClient()

	router := route.Setup(env, timeout, app.Rdb)

	srv := httptest.NewUnstartedServer(router)
	listener, err := net.Listen("tcp", "127.0.0.1:42783")
	if err != nil {
		if listener, err = net.Listen("tcp6", "[::1]:0"); err != nil {
			panic(fmt.Sprintf("httptest: failed to listen on a port: %v", err))
		}
	}
	srv.Listener = listener
	srv.Start()
	expect = httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  srv.URL,
		Client:   httpClient,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	return expect, func() {
		app.Rdb.Close()
		recorder.Stop()
	}
}
