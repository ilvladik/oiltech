package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"oiltech/internal/dtos"
	"oiltech/internal/usecases"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Handler struct {
	datasetUsecase        *usecases.DatasetUsecase
	algorithmUsecase      *usecases.AlgorithmUsecase
	transformationUsecase *usecases.TransformationUsecase
	runUsecase            *usecases.RunUsecase
}

func NewHandler(
	datasetUsecase *usecases.DatasetUsecase,
	algorithmUsecase *usecases.AlgorithmUsecase,
	transformationUsecase *usecases.TransformationUsecase,
	runUsecase *usecases.RunUsecase,
) *Handler {
	return &Handler{
		datasetUsecase:        datasetUsecase,
		algorithmUsecase:      algorithmUsecase,
		transformationUsecase: transformationUsecase,
		runUsecase:            runUsecase,
	}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8090", "http://127.0.0.1:8090"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/healthz", h.health)

	r.Route("/datasets", func(r chi.Router) {
		r.Get("/", h.listDatasets)
		r.Post("/", h.createDataset)
		r.Get("/{code}", h.getDataset)
		r.Get("/{code}/rows", h.datasetRows)
	})

	r.Route("/algorithms", func(r chi.Router) {
		r.Get("/", h.listAlgorithms)
		r.Post("/", h.createAlgorithm)
		r.Get("/{code}", h.getAlgorithm)
	})

	r.Route("/transformations", func(r chi.Router) {
		r.Get("/", h.listTransformations)
		r.Post("/", h.createTransformation)
		r.Get("/tree", h.transformationTree)
		r.Get("/{code}", h.getTransformation)
		r.Get("/{code}/runs", h.transformationRuns)
		r.Post("/{code}/run", h.runTransformation)
	})
	return r
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) listDatasets(w http.ResponseWriter, r *http.Request) {
	response, err := h.datasetUsecase.ListDatasets(r.Context())
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) createDataset(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateDatasetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteAPIError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}
	response, err := h.datasetUsecase.CreateDataset(r.Context(), req)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusCreated, response)
}

func (h *Handler) getDataset(w http.ResponseWriter, r *http.Request) {
	response, err := h.datasetUsecase.GetDataset(r.Context(), chi.URLParam(r, "code"))
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) datasetRows(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	response, err := h.datasetUsecase.GetRows(
		r.Context(),
		chi.URLParam(r, "code"),
		r.URL.Query().Get("field"),
		r.URL.Query().Get("value"),
		limit,
	)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) listAlgorithms(w http.ResponseWriter, r *http.Request) {
	response, err := h.algorithmUsecase.ListAlgorithms(r.Context())
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) createAlgorithm(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateAlgorithmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteAPIError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}
	response, err := h.algorithmUsecase.CreateAlgorithm(r.Context(), req)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusCreated, response)
}

func (h *Handler) getAlgorithm(w http.ResponseWriter, r *http.Request) {
	response, err := h.algorithmUsecase.GetAlgorithm(r.Context(), chi.URLParam(r, "code"))
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) listTransformations(w http.ResponseWriter, r *http.Request) {
	response, err := h.transformationUsecase.ListTransformations(r.Context())
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) createTransformation(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateTransformationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteAPIError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}
	response, err := h.transformationUsecase.CreateTransformation(r.Context(), req)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusCreated, response)
}

func (h *Handler) getTransformation(w http.ResponseWriter, r *http.Request) {
	response, err := h.transformationUsecase.GetTransformation(r.Context(), chi.URLParam(r, "code"))
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) transformationTree(w http.ResponseWriter, r *http.Request) {
	response, err := h.transformationUsecase.GetTree(r.Context())
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) runTransformation(w http.ResponseWriter, r *http.Request) {
	response, err := h.runUsecase.RunTransformation(r.Context(), chi.URLParam(r, "code"))
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) transformationRuns(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	response, err := h.runUsecase.ListRuns(r.Context(), chi.URLParam(r, "code"), limit)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, response)
}
