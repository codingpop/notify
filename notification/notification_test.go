package notification_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/codingpop/refurbed/notification"
)

func TestNotification(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	appErrs := make(chan error, 1)

	n := notification.New(context.Background(), server.URL, time.Duration(0), 1, appErrs)

	n.Enqueue("hi")
	n.Enqueue("how")
	n.Enqueue("sdfkjsdf")

	err := <-appErrs

	if !errors.Is(err, notification.ErrRequest) {
		t.Errorf("want: %s, got: %s", notification.ErrRequest, err)
	}
}
