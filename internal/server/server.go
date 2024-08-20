package server

import (
	"github.com/VikaPaz/matchmaker/internal/server/matchmaker"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type ImplServer struct {
	service matchmaker.Service
	log     *logrus.Logger
}

func NewServer(s matchmaker.Service, logger *logrus.Logger) *ImplServer {
	return &ImplServer{
		service: s,
		log:     logger,
	}
}

func (i *ImplServer) Handlers() *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	m := matchmaker.NewHandler(i.service, i.log)

	r.Mount("/matchmaker", m.Router())

	return r
}
