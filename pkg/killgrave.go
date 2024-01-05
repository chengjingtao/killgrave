package pkg

import (
	"fmt"
	killgrave "github.com/friendsofgo/killgrave/internal"
	server "github.com/friendsofgo/killgrave/internal/server/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/afero"
	"log"
	"net/http"
)

type Server = server.Server
type Config = killgrave.Config

// RunServer will run killgrave mock server
func RunServer(cfg killgrave.Config) server.Server {
	_defaultStrictSlash := true
	router := mux.NewRouter().StrictSlash(_defaultStrictSlash)
	httpAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	httpServer := http.Server{
		Addr:    httpAddr,
		Handler: handlers.CORS(server.PrepareAccessControl(cfg.CORS)...)(router),
	}

	proxyServer, err := server.NewProxy(cfg.Proxy.Url, cfg.Proxy.Mode)
	if err != nil {
		log.Fatal(err)
	}

	imposterFs := server.NewImposterFS(afero.NewOsFs())
	s := server.NewServer(
		cfg.ImpostersPath,
		router,
		&httpServer,
		proxyServer,
		cfg.Secure,
		imposterFs,
	)
	if err := s.Build(); err != nil {
		log.Fatal(err)
	}

	s.Run()
	return s
}
