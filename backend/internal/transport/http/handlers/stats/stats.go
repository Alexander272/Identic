package stats

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Alexander272/Identic/backend/internal/access"
	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/Alexander272/Identic/backend/internal/models/response"
	"github.com/Alexander272/Identic/backend/internal/services"
	"github.com/Alexander272/Identic/backend/internal/transport/middleware"
	"github.com/Alexander272/Identic/backend/pkg/error_bot"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type commonParamsSetter interface {
	SetActorID(*uuid.UUID)
	SetStartDate(*time.Time)
	SetEndDate(*time.Time)
	SetLimit(int)
	SetOffset(int)
}

type searchLogsAdapter struct {
	*models.GetSearchLogsDTO
}

func (a *searchLogsAdapter) SetActorID(v *uuid.UUID)   { a.GetSearchLogsDTO.ActorID = v }
func (a *searchLogsAdapter) SetStartDate(v *time.Time) { a.GetSearchLogsDTO.StartDate = v }
func (a *searchLogsAdapter) SetEndDate(v *time.Time)   { a.GetSearchLogsDTO.EndDate = v }
func (a *searchLogsAdapter) SetLimit(v int)            { a.GetSearchLogsDTO.Limit = v }
func (a *searchLogsAdapter) SetOffset(v int)           { a.GetSearchLogsDTO.Offset = v }

type activityLogsAdapter struct {
	*models.GetAllActivityLogsDTO
}

func (a *activityLogsAdapter) SetActorID(v *uuid.UUID)   { a.GetAllActivityLogsDTO.ActorID = v }
func (a *activityLogsAdapter) SetStartDate(v *time.Time) { a.GetAllActivityLogsDTO.StartDate = v }
func (a *activityLogsAdapter) SetEndDate(v *time.Time)   { a.GetAllActivityLogsDTO.EndDate = v }
func (a *activityLogsAdapter) SetLimit(v int)            { a.GetAllActivityLogsDTO.Limit = v }
func (a *activityLogsAdapter) SetOffset(v int)           { a.GetAllActivityLogsDTO.Offset = v }

type userLoginsAdapter struct {
	*models.GetUserLoginsDTO
}

func (a *userLoginsAdapter) SetActorID(v *uuid.UUID)   {}
func (a *userLoginsAdapter) SetStartDate(v *time.Time) { a.GetUserLoginsDTO.StartDate = v }
func (a *userLoginsAdapter) SetEndDate(v *time.Time)   { a.GetUserLoginsDTO.EndDate = v }
func (a *userLoginsAdapter) SetLimit(v int)            { a.GetUserLoginsDTO.Limit = v }
func (a *userLoginsAdapter) SetOffset(v int)           { a.GetUserLoginsDTO.Offset = v }

type Handler struct {
	service services.Statistic
}

func NewHandler(service services.Statistic) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.Statistic, middleware *middleware.Middleware) {
	handler := NewHandler(service)

	stats := api.Group("/statistics")
	{
		search := stats.Group("", middleware.CheckPermissions(access.Reg.R(access.ResourceSearch).Read()))
		{
			search.GET("/search", handler.getSearch)
		}

		activity := stats.Group("", middleware.CheckPermissions(access.Reg.R(access.ResourceActivity).Read()))
		{
			activity.GET("/activity", handler.getActivity)
		}

		userLogins := stats.Group("", middleware.CheckPermissions(access.Reg.R(access.ResourceUser).Read()))
		{
			userLogins.GET("/logins", handler.getUserLogins)
		}
	}
}

func parseCommonParams(c *gin.Context, setter commonParamsSetter) error {
	if actorId := c.Query("actorId"); actorId != "" {
		id, err := uuid.Parse(actorId)
		if err != nil {
			return err
		}
		setter.SetActorID(&id)
	}
	if startDate := c.Query("startDate"); startDate != "" {
		t, err := time.Parse(time.RFC3339, startDate)
		if err != nil {
			return err
		}
		setter.SetStartDate(&t)
	}
	if endDate := c.Query("endDate"); endDate != "" {
		t, err := time.Parse(time.RFC3339, endDate)
		if err != nil {
			return err
		}
		setter.SetEndDate(&t)
	}
	if limit := c.Query("limit"); limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return err
		}
		setter.SetLimit(l)
	}
	if offset := c.Query("offset"); offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return err
		}
		setter.SetOffset(o)
	}
	return nil
}

func (h *Handler) getSearch(c *gin.Context) {
	req := &searchLogsAdapter{&models.GetSearchLogsDTO{}}

	if err := parseCommonParams(c, req); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Некорректные данные")
		return
	}

	data, err := h.service.GetSearch(c, req.GetSearchLogsDTO)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), req)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data, Total: len(data)})
}

func (h *Handler) getActivity(c *gin.Context) {
	req := &activityLogsAdapter{&models.GetAllActivityLogsDTO{}}

	if err := parseCommonParams(c, req); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Некорректные данные")
		return
	}

	data, err := h.service.GetActivity(c, req.GetAllActivityLogsDTO)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), req)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data, Total: len(data)})
}

func (h *Handler) getUserLogins(c *gin.Context) {
	req := &userLoginsAdapter{&models.GetUserLoginsDTO{}}

	if err := parseCommonParams(c, req); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Некорректные данные")
		return
	}

	data, err := h.service.GetLastUserLogin(c, req.GetUserLoginsDTO)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), req)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data, Total: len(data)})
}
