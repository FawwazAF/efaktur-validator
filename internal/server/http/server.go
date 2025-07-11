package http

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Index   indexHandlerInterface
	Efaktur efakturHandlerInterface
}

type indexHandlerInterface interface {
	HandlerIndex(gCtx *gin.Context)
}

type efakturHandlerInterface interface {
	HandlerValidateEfaktur(c *gin.Context)
}

type Server struct {
	handler Handler
	router  *gin.Engine
}

func NewServer(handler Handler) *Server {
	return &Server{
		handler: handler,
	}
}

func (s *Server) Start(addr string) {
	s.router = gin.Default()
	s.registerHandler()
	srv := &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	go func() {
		log.Println("Running http server at port", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Failed to run http server: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
