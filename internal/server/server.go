package server

import (
	"FinUslugi/internal/config"
	"FinUslugi/internal/handlers"
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type Server struct {
	DB     *pgxpool.Pool
	Logger *zap.SugaredLogger
}

func NewServer(logger *zap.SugaredLogger) (*Server, error) {
	db, err := connectDB()
	if err != nil {
		return nil, fmt.Errorf("could not connect to the database: %v", err)
	}

	srv := &Server{
		DB:     db,
		Logger: logger,
	}

	return srv, nil
}

func connectDB() (*pgxpool.Pool, error) {
	cfg, err := config.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed parse config: %w", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.DBName,
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %v", err)
	}

	dbpool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	return dbpool, nil
}

func (s *Server) HandleMaterials(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handlers.GetAllMaterials(s.DB, s.Logger, w, r)
	case http.MethodPost:
		handlers.CreateMaterial(s.DB, s.Logger, w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) HandleMaterial(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/materials/"):]

	switch r.Method {
	case http.MethodGet:
		handlers.GetMaterialByID(s.DB, s.Logger, w, r, id)
	case http.MethodPut:
		handlers.UpdateMaterial(s.DB, s.Logger, w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
