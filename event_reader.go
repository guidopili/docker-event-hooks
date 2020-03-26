package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
)

const SocketName = "/var/run/docker.sock"

type DockerEvent struct {
	ID     string `json:"id,omitempty"`
	From   string `json:"from,omitempty"`
	Type   string
	Action string
	Actor  struct {
		ID         string
		Attributes map[string]string
	}
	Scope string `json:"scope,omitempty"`
	TimeNano int64 `json:"timeNano,omitempty"`
}

func eventsReader(ctx context.Context, jsonConfig JsonConfig)  (<-chan DockerEvent, <-chan error)  {
	messages := make(chan DockerEvent)
	errs := make(chan error, 1)

	started := make(chan struct{})

	go func() {
		defer close(errs)

		handleErr := func(err error) bool {
			if err != nil {
				close(started)
				errs <- err
				return true
			}

			return false
		}

		info, err := os.Stat(SocketName)
		if os.IsNotExist(err) || info.Mode() & os.ModeSocket != os.ModeSocket {
			close(started)
			errs <- err
			return
		}

		client := http.Client{
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("unix", SocketName)
				},
			},
		}

		r, err := client.Head("http://docker/_ping")
		if handleErr(err) {
			return
		}

		apiVersion := r.Header.Get("Api-Version")

		if apiVersion == "" {
			apiVersion = "1.40"
		}

		response, err := client.Get(fmt.Sprintf("http://docker/v%s/events?type=container&type=network", apiVersion))

		if handleErr(err) {
			return
		}

		defer response.Body.Close()
		decoder := json.NewDecoder(response.Body)

		close(started)
		for {
			select {
			case <-ctx.Done():
				errs <- ctx.Err()
				return
			default:
				var event DockerEvent
				if err := decoder.Decode(&event); err != nil {
					errs <- err
					return
				}

				select {
				case messages <- event:
				case <-ctx.Done():
					errs <- ctx.Err()
					return
				}
			}
		}
	}()
	<-started

	return messages, errs
}