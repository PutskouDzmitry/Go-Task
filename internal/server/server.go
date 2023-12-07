package server

import (
	"context"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strconv"
	"task/internal/config"
	"task/pkg/utils/validate"
	"time"
)

type serverConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read-timeout"`
	WriteTimeout time.Duration `yaml:"write-timeout"`
}

type Server struct {
	config *serverConfig
	log    *zap.Logger
	Engine *gin.Engine
	http   *http.Server
}

// newConfigServer initializes a new serverConfig struct based on the provided config.Provider.
//
// It takes a config.Provider as a parameter and returns a pointer to a serverConfig struct and an error.
func newConfigServer(cfg *config.Config) (*serverConfig, error) {

	return &serverConfig{
		Host:         cfg.HTTP.Host,
		Port:         cfg.HTTP.Port,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}, nil

}

// New initializes a new Server instance.
//
// It takes a config.Provider and a *zap.Logger as parameters.
// It returns a *Server and an error.
func New(cfg *config.Config, logger *zap.Logger) (*Server, error) {
	var app *gin.Engine

	// Set gin mode, which:
	//   - debug mode means show stack info.
	if os.Getenv("GIN_MODE") == "debug" {
		app = gin.Default()
		app.Use(gin.Recovery())
		gin.SetMode(gin.DebugMode)

	} else {
		app = gin.New()
		gin.SetMode(gin.ReleaseMode)
		// Add a ginzap middleware, which:
		//   - Logs all requests, like a combined access and error log.
		//   - Logs to stdout.
		//   - RFC3339 with UTC time format.
		app.Use(ginzap.Ginzap(logger, time.RFC3339, true))
		// Logs all panic to error log
		//   - stack means whether output the stack info.
		app.Use(ginzap.RecoveryWithZap(logger, true))
	}

	// Get config server
	// - provider is a uber config
	sConfig, err := newConfigServer(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get serverConfig config")
	}

	// Set max request body size to 40MB
	app.MaxMultipartMemory = 40 * 1024 * 1024

	binding.Validator = new(validate.DefaultValidator)

	app.StaticFS("/uploads", http.Dir("./uploads"))

	// Set http server params
	server := &http.Server{
		Addr:              sConfig.Host + ":" + strconv.Itoa(sConfig.Port),
		Handler:           app,
		ReadTimeout:       sConfig.ReadTimeout,
		ReadHeaderTimeout: sConfig.ReadTimeout,
		WriteTimeout:      sConfig.WriteTimeout,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
	}

	return &Server{
		config: sConfig,
		log:    logger,
		Engine: app,
		http:   server,
	}, nil
}

// Start starts the server.
//
// It listens and serves HTTP requests.
// If an error occurs, it checks if the error is http.ErrServerClosed.
// If it is, it logs a warning message "server is closed".
// Otherwise, it logs a fatal message "failed to run server" and includes the error details.
func (s *Server) Start() error {
	s.log.Info("starting server on " + s.config.Host + ":" + strconv.Itoa(s.config.Port))
	return errors.Wrap(s.http.ListenAndServe(), "missing start server")
}

// Stop stops the server.
//
// This function does not take any parameters.
// It does not have a return type.
func (s *Server) Stop(ctx context.Context) {
	if err := s.http.Shutdown(ctx); err != nil {
		s.log.Fatal("failed to shutdown server", zap.Error(err))
	}
}
