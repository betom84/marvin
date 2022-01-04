package ui

import (
	"fmt"
	"marvin/alexa"
	"marvin/config"
	"marvin/logger"
	"marvin/metrics"
	"marvin/ui/api"
	"net/http"
	"os"
	"path/filepath"
)

type Server struct {
	router         *router
	alexaServer    *alexa.Server
	logMultiWriter *logger.LogMultiWriter
}

func NewServer(a *alexa.Server, mw *logger.LogMultiWriter) *Server {
	s := &Server{
		router:         newRouter(),
		alexaServer:    a,
		logMultiWriter: mw,
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.router.get("/api/alexa/state", api.HandleAlexaStateGet(s.alexaServer))
	s.router.post("/api/alexa/state", api.HandleAlexaStateSet(s.alexaServer))
	s.router.get("/api/log", api.HandleLogGet())
	s.router.put("/api/log", api.HandleLogPut())
	s.router.get("/api/log/socket", api.HandleLogSocket(s.logMultiWriter))
	s.router.get("/api/config", api.HandleConfigGet())
	s.router.get("/api/endpoint", api.HandleEndpointsGet())

	s.router.get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		metrics.Handler().ServeHTTP(w, r)
	})

	s.router.fallback = handleIndex()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered ", r)

			w.Header().Add("content-type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{ "error": "%s" }`, r)))
		}
	}()

	s.router.serveHTTP(w, r)
}

func handleIndex() http.HandlerFunc {
	var index = filepath.Join(fmt.Sprintf("%s/index.html", config.Get().UIRoot))

	return func(w http.ResponseWriter, r *http.Request) {
		if config.Get().UIRoot == "" {
			panic(fmt.Errorf("UI not configured"))
		}

		path, err := filepath.Abs(r.URL.Path)
		if err != nil {
			panic(err)
		}

		path = filepath.Join(config.Get().UIRoot, path)

		_, err = os.Stat(path)
		if os.IsNotExist(err) {
			http.ServeFile(w, r, index)
			return
		}

		if err != nil {
			panic(err)
		}

		http.FileServer(http.Dir(config.Get().UIRoot)).ServeHTTP(w, r)
	}
}
