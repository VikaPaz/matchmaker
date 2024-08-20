package matchmaker

import (
	"encoding/json"
	"net/http"

	"github.com/VikaPaz/matchmaker/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type Service interface {
	AddPlayer(req models.AddRequest) (models.Player, error)
}

type Handler struct {
	service Service
	log     *logrus.Logger
}

func NewHandler(service Service, logger *logrus.Logger) *Handler {
	return &Handler{
		service: service,
		log:     logger,
	}
}

func (rs *Handler) Router() chi.Router {
	r := chi.NewRouter()

	r.Post("/users", rs.new)

	return r
}

// @Summary Adding a user.
// @Description Handles request to add a new user
// @Tags matching
// @Accept json
// @Produce json
// @Param request body models.AddRequest true "Player"
// @Success 200 {object} models.Player "Created Player"
// @Failure 400
// @Failure 500
// @Router /matchmaker/users [post]
func (rs *Handler) new(w http.ResponseWriter, r *http.Request) {
	req := models.AddRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		rs.log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		rs.log.Error("Name is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := rs.service.AddPlayer(req)
	if err != nil {
		rs.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rs.log.Infof("Added player: %v", resp)

	data, err := json.Marshal(resp)
	if err != nil {
		rs.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		rs.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
