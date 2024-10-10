package https

import (
	"encoding/json"
	"go-axiata/config"
	"go-axiata/model"
	"go-axiata/pkg/helper"
	"go-axiata/project"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type handler struct {
	cfg     *config.Config
	logger  zerolog.Logger
	usecase project.Usecase
}

func NewHandlers(cfg *config.Config, usecase project.Usecase, logger zerolog.Logger) Handlers {
	return &handler{
		cfg:     cfg,
		logger:  logger,
		usecase: usecase,
	}
}

func (h *handler) GetPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:

			var req model.Request
			query := r.URL.Query()
			req.Limit = query.Get("limit")
			req.Page = query.Get("page")
			req.Tag = query.Get("tag")

			filter := helper.ParsePaginationSearch(req)

			res, err := h.usecase.GetPost(r.Context(), filter)
			if err != nil {
				helper.RespondJSON(w, http.StatusBadRequest, false, err.Error(), nil, nil)
				return
			} else {
				helper.RespondJSON(w, http.StatusOK, true, "", filter, res)
				return
			}

		case http.MethodPost:
			claims, ok := r.Context().Value("claims").(*model.Claims)
			if !ok {
				helper.RespondJSON(w, http.StatusInternalServerError, false, "Cannot get user", nil, nil)
				return
			}

			if claims.Role != "user" {
				helper.RespondJSON(w, http.StatusUnauthorized, false, "Access denied. You do not have authorization", nil, nil)
				return
			}

			var reqPost model.ReqPost

			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&reqPost)
			if err != nil {
				helper.RespondJSON(w, http.StatusBadRequest, false, err.Error(), nil, nil)
			}

			res, err := h.usecase.CreatePost(r.Context(), reqPost)
			if err != nil {
				helper.RespondJSON(w, http.StatusBadRequest, false, err.Error(), reqPost, nil)
			} else {
				helper.RespondJSON(w, http.StatusOK, true, "Data successfully created", reqPost, res)
			}

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func (h *handler) DetailPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodGet:

			parts := strings.Split(r.URL.Path, "/")
			idStr := parts[3]
			id, _ := uuid.Parse(idStr)

			res, err := h.usecase.DetailPost(r.Context(), id)
			if err != nil {
				if err.Error() == "not_found" {
					helper.RespondJSON(w, http.StatusNotFound, false, "Data Not Found", nil, nil)
				} else {
					helper.RespondJSON(w, http.StatusBadRequest, false, err.Error(), nil, nil)
				}
			} else {
				helper.RespondJSON(w, http.StatusOK, true, "", nil, res)
			}

		case http.MethodDelete:

			parts := strings.Split(r.URL.Path, "/")
			idStr := parts[3]
			id, _ := uuid.Parse(idStr)

			err := h.usecase.DeletePost(r.Context(), id)
			if err != nil {
				if err.Error() == "not_found" {
					helper.RespondJSON(w, http.StatusNotFound, false, "Data Not Found", nil, nil)
				} else {
					helper.RespondJSON(w, http.StatusBadRequest, false, err.Error(), nil, nil)
				}
			} else {
				helper.RespondJSON(w, http.StatusOK, true, "Data successfully deleted", id, nil)
			}

		case http.MethodPut:

			claims, ok := r.Context().Value("claims").(*model.Claims)
			if !ok {
				helper.RespondJSON(w, http.StatusInternalServerError, false, "Cannot get user", nil, nil)
				return
			}

			if claims.Role != "admin" {
				helper.RespondJSON(w, http.StatusUnauthorized, false, "Access denied. You do not have authorization", nil, nil)
				return
			}

			parts := strings.Split(r.URL.Path, "/")
			idStr := parts[3]
			id, _ := uuid.Parse(idStr)

			var reqPost model.ReqPost

			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&reqPost)
			if err != nil {
				helper.RespondJSON(w, http.StatusBadRequest, false, err.Error(), nil, nil)
			}

			err = h.usecase.UpdatePost(r.Context(), id, reqPost)
			if err != nil {
				if err.Error() == "not_found" {
					helper.RespondJSON(w, http.StatusNotFound, false, "Data Not Found", nil, nil)
				} else {
					helper.RespondJSON(w, http.StatusBadRequest, false, err.Error(), nil, nil)
				}
			} else {
				helper.RespondJSON(w, http.StatusOK, true, "Data successfully updated", reqPost, id)
			}

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
