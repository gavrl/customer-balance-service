package handler

import (
	"github.com/gavrl/app/pkg/formatter"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

type validateErrorResponse struct {
	Errors []formatter.ValidationError `json:"errors"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}

func newValidateErrorResponse(c *gin.Context, f *formatter.JSONFormatter, verr validator.ValidationErrors) {
	errors := f.Descriptive(verr)
	logrus.Errorf("validation errors: %+v", errors)
	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, validateErrorResponse{errors})
}
