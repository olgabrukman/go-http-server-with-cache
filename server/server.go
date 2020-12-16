package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"go-http-server-with-cache/cache"
)

type GracefulServer struct {
	srv   *http.Server
	wg    *sync.WaitGroup
	cache cache.Cache
}

func NewServer(wg *sync.WaitGroup, port int, cache cache.Cache) *GracefulServer {
	address := fmt.Sprintf(":%d", port)
	//nolint: exhaustivestruct
	srv := &http.Server{Addr: address}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		clientID, ok := r.URL.Query()["clientId"]

		id, err := strconv.Atoi(clientID[0])
		if !ok || len(clientID[0]) < 1 || err != nil {
			log.Println("url param 'clientId' is missing/not an integer")
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		success := cache.Increment(id)

		if success {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	})

	go func() {
		defer wg.Done() // let main know we are done cleaning up

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			panic(err)
		}
	}()

	log.Println("started http server")

	return &GracefulServer{srv, wg, cache}
}

func (srv *GracefulServer) StopServer() {
	if err := srv.srv.Shutdown(context.TODO()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}

	// wait for goroutine started in startHttpServer() to stop
	srv.wg.Wait()
	srv.cache.StopCleanUp()

	http.DefaultServeMux = new(http.ServeMux)

	log.Println("stopped http server")
}
