package server

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	mapiadmin "github.com/lm-gautam-bhagat/minio-server/api/packages/admin"
	minioclient "github.com/lm-gautam-bhagat/minio-server/api/packages/minio"
	"github.com/lm-gautam-bhagat/minio-server/api/router"
	"github.com/lm-gautam-bhagat/minio-server/config"
	"github.com/lm-gautam-bhagat/minio-server/log"
	"github.com/lm-gautam-bhagat/minio-server/storage"
)

type Server struct {
	cfg       *config.ConfigObject
	strClient *storage.StorageClient
	strAdmin  *storage.StorageAdmin
}

func NewServer(cfg *config.ConfigObject, strClient *storage.StorageClient, strAdmin *storage.StorageAdmin) *Server {
	return &Server{
		cfg: cfg, strClient: strClient, strAdmin: strAdmin,
	}
}

func (s *Server) MapRoutes(r **gin.Engine) {
	handlers := []router.HTTPHandlerProvider{
		minioclient.NewMinioModule(s.strClient),
		mapiadmin.NewMinioAdminModule(s.strAdmin, s.strClient),
	}
	engine := *r // dereference *gin.Engine

	for _, handler := range handlers {
		for _, hldr := range handler.GetHTTPHandler() {
			wrapped := APIWrap(hldr.Handler)
			path := fmt.Sprintf("/api/v%d/%s", hldr.Version, hldr.Path)
			switch hldr.Method {
			case http.MethodGet:
				engine.GET(path, wrapped)
			case http.MethodPost:
				engine.POST(path, wrapped)
			case http.MethodPut:
				engine.PUT(path, wrapped)
			case http.MethodDelete:
				engine.DELETE(path, wrapped)
			default:
				engine.Handle(hldr.Method, path, wrapped)
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

	log.Info("running: ", fmt.Sprintf("http://%s/", addr))
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

	case <-sigChan:
		fmt.Println("intrupt is called")

	}

	srv.Close()
	log.Close()
}

func APIWrap(h func(*router.SessionContext)) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := router.NewSessionContext(c)
		h(ctx)
	}
}
