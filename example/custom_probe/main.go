package main

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Autodoc-Technology/k8sprobe"
)

// CustomProbe represents a structure to validate external connectivity status for probing operations.
type CustomProbe struct {
	isExternalConnectionValid bool
}

// IsValid returns true if the external connection is valid, otherwise returns false.
func (c CustomProbe) IsValid() (bool, k8sprobe.Cause) {
	cause := "OK"
	if !c.isExternalConnectionValid {
		cause = "external connection is not valid"
	}
	return c.isExternalConnectionValid, cause
}

func main() {
	// create a probe that is initially valid and going to invalid in 10 seconds
	p := &CustomProbe{isExternalConnectionValid: true}
	go func() {
		<-time.After(10 * time.Second)
		p.isExternalConnectionValid = false
	}()

	m := k8sprobe.NewManager()
	m.RegisterProbe(k8sprobe.LivenessProbe, p)

	// create http server with health probe handler
	http.Handle("/healthz/"+k8sprobe.UrlPathValue, k8sprobe.NewHttpHandler(m))
	go func() {
		if err := http.ListenAndServe(":8089", http.DefaultServeMux); err != nil {
			return
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// check probe status every 2 seconds using http client
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(2 * time.Second):
				resp, err := http.Get("http://localhost:8089/healthz/" + k8sprobe.LivenessProbe.String())
				if err != nil {
					_ = resp.Body.Close()
					return
				}
				slog.Info("probe http check result", "status", resp.Status, "code", resp.StatusCode)
				if resp.StatusCode != http.StatusOK {
					bytes, _ := io.ReadAll(resp.Body)
					slog.Error(
						"probe failed",
						"status", resp.Status,
						"code", resp.StatusCode,
						"body", string(bytes),
					)
					stop()
				}
				_ = resp.Body.Close()
			}
		}
	}()

	<-ctx.Done()
}
