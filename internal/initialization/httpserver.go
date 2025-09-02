package initialization

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ducklawrence05/go-test-backend-api/pkg/logger"
	"go.uber.org/zap"
)

const (
	defaultReadTimeout     = 10 * time.Second
	defaultWriteTimeout    = 20 * time.Second
	defaultIdleTimeout     = 60 * time.Second
	defaultShutdownTimeout = 3 * time.Second
)

func NewServer(port int, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		IdleTimeout:  defaultIdleTimeout,
	}
}

func RunServer(server *http.Server, logger logger.Interface) {
	go func() {
		fmt.Println("Listening and serving HTTP on", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Error ListenAndServe():", zap.Error(err))
		}
	}()

	//  kill, Ctrl + C
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down server...")

	// context to grateful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown:", zap.Error(err))
	}

	fmt.Println("Server exited properly")
}
