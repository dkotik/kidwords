package http

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/dkotik/htadaptor"
	"github.com/dkotik/kidwords/service/secret/mock"
)

var livePort = flag.Int("livePort", 0, "start an HTTP service on port")

func TestMain(m *testing.M) {
	exitCode := m.Run()
	if exitCode == 0 && livePort != nil && *livePort != 0 {
		mux, err := newTestService(
			mock.New(),
			WithAdaptor(
				htadaptor.New(
					htadaptor.WithErrorHandler(htadaptor.ErrorHandlerFunc(
						func(w http.ResponseWriter, r *http.Request, err error) error {
							w.WriteHeader(http.StatusInternalServerError)
							_, _ = fmt.Fprintln(w, err.Error())
							return nil
						})),
				),
			),
		)
		if err != nil {
			panic(err)
		}
		_ = http.ListenAndServe(
			fmt.Sprintf("localhost:%d", *livePort),
			mux,
		)
	}
	os.Exit(exitCode)
}
