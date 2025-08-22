package {{.ProjectPackageName}}

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"{{.ProjectPackage}}/pkg/auth"
	"{{.ProjectPackage}}/pkg/ui"
)

// Server represents the main application server.
type Server struct {
	config        *Config
	logger        *zap.Logger
	metrics       *Metrics
	metricsServer *MetricsServer
	mainServer    *http.Server
	auth          *auth.OIDCAuth
}

// NewServer creates a new server instance.
func NewServer(ctx context.Context, config *Config) (server *Server, err error) {
	// Initialize logger
	logger, err := NewLogger(config.LogLevel)
	if err != nil {
		return nil, err
	}

	// Initialize metrics
	metrics := NewMetrics()

	// Create metrics server
	metricsServer, err := NewMetricsServer(config.MetricsAddress, logger, metrics)
	if err != nil {
		return nil, err
	}

	server = &Server{
		config:        config,
		logger:        logger,
		metrics:       metrics,
		metricsServer: metricsServer,
	}

	// Initialize OIDC auth if configured
	if config.OIDCClientID != "" && config.OIDCClientSecret != "" {
		authConfig := &auth.Config{
			ClientID:     config.OIDCClientID,
			ClientSecret: config.OIDCClientSecret,
			RedirectURL:  config.OIDCRedirectURL,
			IssuerURL:    config.OIDCIssuerURL,
			CookieDomain: config.OIDCCookieDomain,
			CookieSecure: config.OIDCCookieSecure,
			StaticToken:  config.OIDCStaticToken,
		}

		server.auth, err = auth.NewOIDCAuth(ctx, authConfig, logger)
		if err != nil {
			logger.Error("Failed to initialize OIDC auth", zap.Error(err))
			// Don't fail startup, just log the error and continue without auth
			server.auth = nil
		} else {
			logger.Info("OIDC authentication enabled")
		}
	} else {
		logger.Info("OIDC authentication disabled - credentials not provided")
	}

	// Set up main HTTP server
	server.setupMainServer()

	return server, nil
}

// setupMainServer configures the main HTTP server.
func (s *Server) setupMainServer() {
	router := mux.NewRouter()

	// Register auth routes if auth is enabled
	if s.auth != nil {
		s.auth.RegisterRoutes(router)
	}

	// Register API routes
	s.registerAPIRoutes(router)

	// Register UI routes (with optional auth protection)
	s.registerUIRoutes(router)

	s.mainServer = &http.Server{
		Addr:           s.config.ServerAddress,
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}
}

// registerAPIRoutes registers API endpoints.
func (s *Server) registerAPIRoutes(router *mux.Router) {
	// API routes that might need authentication
	apiRouter := router.PathPrefix("/api").Subrouter()

	if s.auth != nil {
		// Protected API routes
		apiRouter.HandleFunc("/user", s.auth.RequireAuth(s.userAPIHandler)).Methods("GET")
	} else {
		// Unprotected fallback
		apiRouter.HandleFunc("/user", s.userAPIHandler).Methods("GET")
	}

	// Public API routes (no auth required)
	apiRouter.HandleFunc("/status", s.metrics.InstrumentHandler("/api/status", s.statusHandler)).Methods("GET")
}

// registerUIRoutes registers UI routes.
func (s *Server) registerUIRoutes(router *mux.Router) {
	uiHandler := s.metrics.InstrumentHandler("/", ui.Handler())

	if s.auth != nil {
		// Protect UI routes with authentication
		router.PathPrefix("/").Handler(s.auth.RequireAuth(uiHandler))
	} else {
		// Serve UI without authentication
		ui.RegisterRoutes(router)
	}
}

// Run starts the server and handles graceful shutdown.
func (s *Server) Run(ctx context.Context) (err error) {
	// Create cancellable context for graceful shutdown
	//nolint:ineffassign,staticcheck,wastedassign // ctx reassignment is intentional for cancellation
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Enhanced startup logging
	s.logger.Info("🚀 {{.ProjectName}} Service Starting",
		zap.String("version", "1.0.0"),
		zap.String("build", "development"),
	)
	s.logger.Info("📡 Component Services:",
		zap.String("spa_server", fmt.Sprintf("http://%s", s.config.ServerAddress)),
		zap.String("metrics_server", fmt.Sprintf("http://%s", s.config.MetricsAddress)),
	)
	s.logger.Info("🔗 Available Endpoints:",
		zap.String("spa_ui", fmt.Sprintf("http://%s/", s.config.ServerAddress)),
		zap.String("api_status", fmt.Sprintf("http://%s/api/status", s.config.ServerAddress)),
		zap.String("api_user", fmt.Sprintf("http://%s/api/user", s.config.ServerAddress)),
		zap.String("metrics", fmt.Sprintf("http://%s/metrics", s.config.MetricsAddress)),
		zap.String("health", fmt.Sprintf("http://%s/healthz", s.config.MetricsAddress)),
		zap.String("readiness", fmt.Sprintf("http://%s/readyz", s.config.MetricsAddress)),
	)

	// Channel to receive OS signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// WaitGroup to coordinate server shutdown
	var wg sync.WaitGroup

	// Start metrics server
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.logger.Info("Starting metrics server", zap.String("address", s.config.MetricsAddress))
		metricsErr := s.metricsServer.Start()
		if metricsErr != nil {
			s.logger.Error("Metrics server error", zap.Error(metricsErr))
		}
	}()

	// Start main server
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.logger.Info("Starting main server", zap.String("address", s.config.ServerAddress))
		mainErr := s.mainServer.ListenAndServe()
		if mainErr != nil && mainErr != http.ErrServerClosed {
			s.logger.Error("Main server error", zap.Error(mainErr))
		}
	}()

	// Wait for shutdown signal
	sig := <-sigCh
	s.logger.Info("Received shutdown signal", zap.String("signal", sig.String()))

	// Cancel context to signal shutdown
	cancel()

	// Gracefully shutdown servers
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	s.logger.Info("Shutting down servers...")

	// Shutdown main server
	shutdownErr := s.mainServer.Shutdown(shutdownCtx)
	if shutdownErr != nil {
		s.logger.Error("Error shutting down main server", zap.Error(shutdownErr))
	}

	// Shutdown metrics server
	shutdownErr = s.metricsServer.Shutdown(shutdownCtx)
	if shutdownErr != nil {
		s.logger.Error("Error shutting down metrics server", zap.Error(shutdownErr))
	}

	// Wait for all goroutines to finish
	wg.Wait()

	s.logger.Info("Server shutdown complete")
	return nil
}

// userAPIHandler handles user API requests.
func (s *Server) userAPIHandler(w http.ResponseWriter, r *http.Request) {
	var userEmail = "anonymous"

	// Try to get user from context (set by auth middleware)
	if userInfo := r.Context().Value("user"); userInfo != nil {
		if user, ok := userInfo.(*auth.UserInfo); ok {
			userEmail = user.Email
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"user":"` + userEmail + `"}`))
}

// statusHandler handles status API requests.
func (s *Server) statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok","service":"{{.ProjectName}}"}`))
}
