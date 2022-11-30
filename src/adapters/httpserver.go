package adapters

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"net/http"
	"strconv"
	"time"
)

type HttpHandler struct {
	ID    string
	redis *redis.Client
}

func NewHttpHandler(ID string) *HttpHandler {
	return &HttpHandler{
		ID: ID,
		redis: redis.NewClient(&redis.Options{
			Addr:     "redis:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
}

func (h HttpHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	fmt.Println("handling req in ", h.ID)

	name := req.URL.Query().Get("name")
	conn := h.redis.Conn(req.Context())
	defer conn.Close()

	timeLeft, err := conn.Get(req.Context(), name).Time()
	if err == redis.Nil {
		conn.Set(req.Context(), name, time.Now().Add(60*time.Second), time.Hour)
		_, err = writer.Write([]byte("welcome " + name + ". Please wait 60 seconds..."))
		return
	} else if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte(err.Error()))
		return
	}
	if timeLeft.After(time.Now()) {
		writer.WriteHeader(418)
		writer.Write([]byte("please wait, still processing " + name + ". " + strconv.Itoa(int(time.Now().Sub(timeLeft).Seconds())) + " seconds left."))
	} else {
		conn.Del(req.Context(), name)
	}

	if err != nil {
		writer.WriteHeader(418)
	}
}
