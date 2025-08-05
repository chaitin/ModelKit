package v1

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/chaitin/ModelKit/domain"
	"github.com/chaitin/ModelKit/modelkit"
	"github.com/chaitin/ModelKit/pkg/log"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ModelKitPandaAdapter struct {
	usecase      domain.ModelKit
	isApmEnabled bool
	baseLogger   *log.Logger
}

func NewModelKitPandaAdapter(
	echo *echo.Echo,
	logger *log.Logger,
	isApmEnabled bool,
) *ModelKitPandaAdapter {
	m := &ModelKitPandaAdapter{
		usecase:      modelkit.NewModelKit(),
		isApmEnabled: isApmEnabled,
		baseLogger:   logger.WithModule("http_base_handler"),
	}

	return m
}

// get provider supported model list
//
//	@Summary		get provider supported model list
//	@Description	get provider supported model list
//	@Tags			model
//	@Accept			json
//	@Produce		json
//	@Param			params	query		domain.GetProviderModelListReq	true	"get supported model list request"
//	@Success		200		{object}	domain.Response{data=domain.GetProviderModelListResp}
//	@Router			/api/v1/model/provider/supported [get]
func (h *ModelKitPandaAdapter) GetProviderSupportedModelList(c echo.Context) error {
	var req domain.GetProviderModelListReq
	if err := c.Bind(&req); err != nil {
		return h.NewResponseWithError(c, "invalid request", err)
	}
	if err := c.Validate(&req); err != nil {
		return h.NewResponseWithError(c, "invalid request", err)
	}
	ctx := c.Request().Context()

	models, err := h.usecase.PandaModelList(ctx, &req)
	if err != nil {
		return h.NewResponseWithError(c, "get user model list failed", err)
	}
	return h.NewResponseWithData(c, models)
}

func (h *ModelKitPandaAdapter) NewResponseWithError(c echo.Context, msg string, err error) error {
	traceID := ""
	if h.isApmEnabled {
		span := trace.SpanFromContext(c.Request().Context())
		traceID = span.SpanContext().TraceID().String()
		span.SetAttributes(attribute.String("error", fmt.Sprintf("%+v", err)), attribute.String("msg", msg))
	} else {
		traceID = uuid.New().String()
	}
	h.baseLogger.LogAttrs(c.Request().Context(), slog.LevelError, msg, slog.String("trace_id", traceID), slog.Any("error", err))
	return c.JSON(http.StatusOK, domain.Response{
		Success: false,
		Message: fmt.Sprintf("%s [trace_id: %s]", msg, traceID),
	})
}

func (h *ModelKitPandaAdapter) NewResponseWithData(c echo.Context, data any) error {
	return c.JSON(http.StatusOK, domain.Response{
		Success: true,
		Data:    data,
	})
}
