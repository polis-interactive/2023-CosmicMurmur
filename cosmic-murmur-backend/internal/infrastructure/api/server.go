package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http"
	"sync"
)

type Server struct {
	router       *gin.Engine
	srv          *http.Server
	shutdown     bool
	shutdownLock sync.Mutex
	port         int
}

func NewServer(cfg Config) (*Server, error) {

	htmlPath, err := checkIfIsHtmlRoot(cfg.GetWebServerRootDirectory(), cfg.GetProgramName())
	if err != nil {
		return nil, err
	}

	if cfg.GetWebServerIsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile(htmlPath, true)))

	return &Server{
		router:   router,
		port:     cfg.GetWebServerPort(),
		shutdown: true,
	}, nil
}

func (s *Server) Startup() error {

	s.shutdownLock.Lock()
	defer s.shutdownLock.Unlock()

	if s.shutdown == false {
		return errors.New("FrontendServer, Startup: Tried to startup server twice")
	}

	addr := fmt.Sprintf("0.0.0.0:%d", s.port)
	log.Println(fmt.Sprintf("FrontendServer, Startup: listening at %s", addr))

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("FrontendServer, Startup: Failed to listen: %v", err)
		return err
	}
	log.Println(listener)
	return nil
}
