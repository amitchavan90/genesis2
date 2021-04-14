package genesis

import (
	"bytes"
	"context"
	"fmt"
	"genesis/config"
	"net/http"
	"text/template"

	"github.com/caddyserver/caddy"
	// http driver for caddy
	_ "github.com/caddyserver/caddy/caddyhttp"
	"go.uber.org/zap"
)

const caddyfileTemplate = `
{{ .adminCaddyAddr}} {
	tls off
    proxy /api/ {{ .apiAddr }} {
		transparent
		websocket
		timeout 10m
	}
	index {{ .adminIndexFile }} index.html
    root {{ .adminRootPath }}
    rewrite { 
        if {path} not_match ^/api
        to {path} /
    }
}
{{ .consumerCaddyAddr}} {
	tls off
    proxy /api/ {{ .apiAddr }} {
		transparent
		websocket
		timeout 10m
    }
	index {{ .consumerIndexFile }} index.html
    root {{ .consumerRootPath }}
    rewrite { 
        if {path} not_match ^/api
        to {path} /
    }
}
`

// LoadbalancerService for caddy style load balancing
type LoadbalancerService struct {
	Log *zap.SugaredLogger
}

// Run the API service
func (s *LoadbalancerService) Run(ctx context.Context, apiAddr string, options *config.LoadBalancer, version string) error {
	s.Log.Infow("Starting caddy")
	caddy.AppName = "Genesis"
	caddy.AppVersion = version
	caddy.Quiet = true
	t := template.Must(template.New("CaddyFile").Parse(caddyfileTemplate))
	data := map[string]string{
		"apiAddr":           apiAddr,
		"adminCaddyAddr":    options.AdminAddr,
		"adminIndexFile":    options.AdminIndexFile,
		"adminRootPath":     options.AdminRootPath,
		"consumerCaddyAddr": options.ConsumerAddr,
		"consumerIndexFile": options.ConsumerIndexFile,
		"consumerRootPath":  options.ConsumerRootPath,
	}

	result := &bytes.Buffer{}
	err := t.Execute(result, data)
	if err != nil {
		return fmt.Errorf("execute template: %w", err)
	}
	caddyfile := &caddy.CaddyfileInput{
		Contents:       result.Bytes(),
		Filepath:       "Caddyfile",
		ServerTypeName: "http",
	}

	instance, err := caddy.Start(caddyfile)
	if err != nil {
		return fmt.Errorf("start caddy: %w", err)
	}

	go func() {
		select {
		case <-ctx.Done():
			s.Log.Info("Stopping caddy")
			err := instance.Stop()
			if err != nil {
				fmt.Println(err)
			}
		}
	}()
	instance.Wait()
	return nil
}

// APIService for long running
type APIService struct {
	Addr string
	Log  *zap.SugaredLogger
}

// Run the API service
func (s *APIService) Run(ctx context.Context, controller http.Handler) error {
	s.Log.Infow("Starting API")

	server := &http.Server{
		Addr:    s.Addr,
		Handler: controller,
	}

	go func() {
		select {
		case <-ctx.Done():
			s.Log.Info("Stopping API")
			err := server.Shutdown(ctx)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	return server.ListenAndServe()
}
