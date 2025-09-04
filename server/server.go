package server

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/lm-gautam-bhagat/minio-server/api/packages/minio"
	"github.com/lm-gautam-bhagat/minio-server/api/packages/router"
	"github.com/lm-gautam-bhagat/minio-server/config"
	"github.com/lm-gautam-bhagat/minio-server/log"
)

type Server struct {
	cfg *config.ConfigObject
}

func NewServer(cfg *config.ConfigObject) *Server {
	return &Server{
		cfg: cfg,
	}
}

func (s *Server) MapRoutes(r **gin.Engine) {
	handlers := []router.HTTPHandlerProvider{
		minio.NewMinioModule(),
	}
	engine := *r // dereference *gin.Engine

	for _, handler := range handlers {
		for _, hldr := range handler.GetHTTPHandler() {
			wrapped := APIWrap(hldr.Handler)
			switch hldr.Method {
			case http.MethodGet:
				engine.GET(hldr.Path, wrapped)
			case http.MethodPost:
				engine.POST(hldr.Path, wrapped)
			case http.MethodPut:
				engine.PUT(hldr.Path, wrapped)
			case http.MethodDelete:
				engine.DELETE(hldr.Path, wrapped)
			default:
				engine.Handle(hldr.Method, hldr.Path, wrapped)
			}
		}
	}
}

func (s *Server) Setup() error {
	err := log.Init(s.cfg.Logger.FilePath, s.cfg.Logger.Level)
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

func (s *Server) StartServer() {
	router := gin.Default()

	s.MapRoutes(&router)

	addr := ":8080"
	if s.cfg != nil && s.cfg.Network.Host != "" && s.cfg.Network.Port != "" {
		addr = s.cfg.Network.Host + ":" + s.cfg.Network.Port
	}

	srv := http.Server{
		Addr:    addr,
		Handler: router,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	select {
	case <-errChan:
		fmt.Println("server error")
		log.Close()
	case <-sigChan:
		fmt.Println("intrupt is called")
		log.Close()
	}

	srv.Close()
}

func APIWrap(h func(*router.SessionContext)) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := router.NewSessionContext(c)
		h(ctx)
	}
}
