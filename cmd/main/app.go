package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"math/rand"
	"net/http"
	"os"
	"sharedstate/src/adapters"
	"strconv"
	"time"
)

type App struct {
	router *mux.Router
	server *http.Server
}

func NewApp() *App {
	return &App{}
}

func (a *App) Start(port string) error {
	a.router = mux.NewRouter()

	rand.Seed(time.Now().UnixNano())
	str := strconv.Itoa(rand.Intn(1000))
	adapter := adapters.NewHttpHandler(str)
	a.router.Handle("/hello"+str, adapter).Methods(http.MethodGet)
	x := make(chan bool)
	a.router.HandleFunc("/kill"+str, killServer(str, x, a)).Methods(http.MethodGet)
	a.server = &http.Server{
		Addr:    "localhost:" + port,
		Handler: a.router,
	}
	fmt.Println(str, " is listening on ", a.server.Addr)
	go func() {
		for true {
			z := make(chan bool)
			select {
			case up := <-x:
				if !up {
					os.Exit(0)

				}
			case _ = <-z:
				continue
			}
			time.Sleep(5 * time.Second)
			fmt.Println("still listening in ", str)
		}
	}()
	err := a.server.ListenAndServe()
	if err != nil {
		fmt.Println("failed to start server", err.Error())
		return err
	}

	defer a.server.Close()
	return nil

}

func killServer(str string, x chan bool, a *App) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("filling replica ", str)

		writer.WriteHeader(200)
		writer.Write([]byte("killing " + str))
		x <- false
		a.server.Close()
	}
}
