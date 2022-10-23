package main

import (
	"net/http"
	"os"

	"676f.dev/zinc"
	"github.com/gorilla/mux"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/michaelpeterswa/ratelimit-demo/internal/handlers"
	"github.com/michaelpeterswa/ratelimit-demo/internal/kv"
	"go.uber.org/zap"
)

func main() {
	k := koanf.New(".")

	logger, err := zinc.InitLogger("dev")
	if err != nil {
		panic(err)
	}

	if err := k.Load(file.Provider(os.Getenv("CONFIG_FILE")), yaml.Parser()); err != nil {
		logger.Fatal("error loading config file", zap.Error(err))
	}

	redisClient := kv.NewRedisClient(k.String("redis.url"), k.String("redis.port"))
	idh := handlers.NewIDHandler(redisClient.Client, k)

	r := mux.NewRouter()
	r.HandleFunc("/", idh.Handle).Methods(http.MethodPost)
	http.Handle("/", r)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		logger.Fatal("error starting server", zap.Error(err))
	}

}
