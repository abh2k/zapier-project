package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"zapier-project/internal/deployments"

	"github.com/gin-gonic/gin"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type successResponse struct {
	Data any `json:"data"`
}

type errorResponse struct {
	Error APIError `json:"error"`
}

type DeploymentsHandler struct {
	store  *deployments.Store
	logger *slog.Logger
}

func NewDeploymentsHandler(store *deployments.Store, logger *slog.Logger) *DeploymentsHandler {
	return &DeploymentsHandler{store: store, logger: logger}
}

func (h *DeploymentsHandler) List(c *gin.Context) {
	service := c.Query("service")
	status := c.Query("status")
	if status != "" && !deployments.IsValidStatus(status) {
		respondError(
			c,
			http.StatusBadRequest,
			"invalid_status",
			fmt.Sprintf("invalid status %q; valid values: %s", status, strings.Join(deployments.ValidStatuses(), ", ")),
		)
		return
	}

	results := h.store.List(service, status)
	h.logger.Info("listed deployments",
		"service", service,
		"status", status,
		"count", len(results),
	)
	c.JSON(http.StatusOK, successResponse{Data: results})
}

func (h *DeploymentsHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	deployment, found := h.store.GetByID(id)
	if !found {
		h.logger.Warn("deployment not found", "id", id)
		respondError(c, http.StatusNotFound, "not_found", "deployment not found")
		return
	}

	h.logger.Info("fetched deployment", "id", id, "service", deployment.Service, "status", deployment.Status)
	c.JSON(http.StatusOK, successResponse{Data: deployment})
}

func respondError(c *gin.Context, status int, code, message string) {
	c.JSON(status, errorResponse{
		Error: APIError{
			Code:    code,
			Message: message,
		},
	})
}
