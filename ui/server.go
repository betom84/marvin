package ui

import (
	"fmt"
	"marvin/alexa"
	"marvin/config"
	"marvin/ui/api"
	"net/http"
	"os"
	"path/filepath"
)

type Server struct {
	router      *router
	alexaServer *alexa.Server
}

func NewServer(a *alexa.Server) *Server {
	s := &Server{
		router:      newRouter(),
		alexaServer: a,
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.router.get("/api/alexa/state", api.HandleAlexaStateGet(s.alexaServer))
	s.router.post("/api/alexa/state", api.HandleAlexaStateSet(s.alexaServer))
	s.router.get("/api/log", api.HandleLogGet())
	s.router.get("/api/config", api.HandleConfigGet())
	s.router.get("/api/endpoint", api.HandleEndpointsGet())

	s.router.fallback = handleIndex()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := s.router.serveHTTP(w, r)

	if err != nil {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{ "error": "%s" }`, err.Error())))
	}
}

func handleIndex() handlerFuncWithError {
	var index = filepath.Join(fmt.Sprintf("%s/index.html", config.Get().UIRoot))

	return func(w http.ResponseWriter, r *http.Request) error {
		if config.Get().UIRoot == "" {
			return fmt.Errorf("UI not configured")
		}

		path, err := filepath.Abs(r.URL.Path)
		if err != nil {
			return err
		}

		path = filepath.Join(config.Get().UIRoot, path)

		_, err = os.Stat(path)
		if os.IsNotExist(err) {
			http.ServeFile(w, r, index)
			return nil
		}
		if err != nil {
			return err
		}

		http.FileServer(http.Dir(config.Get().UIRoot)).ServeHTTP(w, r)
		return nil
	}
}
