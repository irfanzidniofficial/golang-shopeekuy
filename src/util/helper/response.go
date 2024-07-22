package helper

import (
	"golang-shopeekuy/src/util/repository/model"
	"net/http"

	"github.com/thedevsaddam/renderer"
)

const (
	SUCCESS_MESSSAGE string = "Success"
)

func HandleResponse(w http.ResponseWriter, render *renderer.Render, statusCode int, message interface{}, data interface{}) {
	response := model.BaseResponse{
		Message: message,
		Data:    data,
	}

	render.JSON(w, statusCode, response)
}
