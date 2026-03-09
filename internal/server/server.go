package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/xydac/xbridge/internal/build"
	"github.com/xydac/xbridge/internal/config"
	"github.com/xydac/xbridge/internal/engine"
	"github.com/xydac/xbridge/internal/server/handlers"
)

// Server is the xbridge HTTP server.
type Server struct {
	Router       *chi.Mux
	Config       *config.Config
	BuildManager *build.Manager
	WorkDir      string
}

// New creates and configures a new xbridge server.
func New(workDir string, cfg *config.Config) *Server {
	s := &Server{
		Router:       chi.NewRouter(),
		Config:       cfg,
		BuildManager: build.NewManager(workDir, cfg),
		WorkDir:      workDir,
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	r := s.Router

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)

	if s.Config.Server.Key != "" {
		r.Use(AuthMiddleware(s.Config.Server.Key))
	}

	// Health
	r.Get("/health", handlers.HandleHealth())

	// Build
	bh := &handlers.BuildHandlers{Manager: s.BuildManager}
	r.Post("/build", bh.HandleBuild())
	r.Post("/build/clean", bh.HandleCleanBuild())
	r.Get("/build/{id}", bh.HandleBuildStatus())
	r.Get("/build/{id}/logs", bh.HandleBuildLogs())
	r.Get("/build/{id}/artifact", bh.HandleBuildArtifact())

	// Simulator
	r.Get("/simulators", handlers.HandleListSimulators())
	r.Post("/simulators/boot", handlers.HandleBootSimulator())
	r.Post("/simulators/shutdown", handlers.HandleShutdownSimulator())
	r.Get("/simulators/{udid}/screenshot", handlers.HandleScreenshot())
	r.Post("/simulators/{udid}/install", handlers.HandleInstallApp())
	r.Post("/simulators/{udid}/launch", handlers.HandleLaunchApp())
	r.Post("/simulators/{udid}/openurl", handlers.HandleOpenURL())
	r.Get("/simulators/{udid}/logs", handlers.HandleSimulatorLogs())

	// Git (only if project is a git repo)
	if engine.IsGitRepo(s.WorkDir) {
		gh := &handlers.GitHandlers{WorkDir: s.WorkDir}
		r.Get("/git/status", gh.HandleGitStatus())
		r.Post("/git/pull", gh.HandleGitPull())
		r.Post("/git/checkout", gh.HandleGitCheckout())
	}
}

// ListenAndServe starts the server on the configured port.
func (s *Server) ListenAndServe() error {
	addr := fmt.Sprintf(":%d", s.Config.Server.Port)
	log.Printf("Server running on %s", addr)
	return http.ListenAndServe(addr, s.Router)
}
