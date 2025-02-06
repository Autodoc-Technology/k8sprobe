package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Autodoc-Technology/k8sprobe"
)

func main() {
	// create a probe that is initially valid and going to invalid in 10 seconds
	p := k8sprobe.NewProbe(true)
	go func() {
		<-time.After(10 * time.Second)
		p.SetValid(false, "probe failed due to timeout")
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
				_ = resp.Body.Close()
				if resp.StatusCode != http.StatusOK {
					slog.Error("probe failed", "status", resp.Status, "code", resp.StatusCode)
					stop()
				}
			}
		}
	}()

	<-ctx.Done()
}
