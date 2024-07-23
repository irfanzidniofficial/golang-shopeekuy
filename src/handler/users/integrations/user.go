package integrations

import (
	"golang-shopeekuy/src/util/helper"
	"golang-shopeekuy/src/util/helper/integrations"
	"golang-shopeekuy/src/util/repository/model/users"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/thedevsaddam/renderer"
)

type userDto interface {
	Register(bReq users.User) (*uuid.UUID, error)
}

type userDtoIntegration interface {
	UserDataSignUp(state, code string) (*users.OauthUserData, error)
	GetUsers(bReq users.User) (*users.User, error)
}

type Handler struct {
	userDto            userDto
	userDtoIntegration userDtoIntegration
	render             *renderer.Render
}

func NewHandler(userDto userDto, userDtoIntegration userDtoIntegration, render *renderer.Render) *Handler {
	return &Handler{
		userDto:            userDto,
		userDtoIntegration: userDtoIntegration,
		render:             render,
	}
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, integrations.SSOSignUp.AuthCodeURL(integrations.RandomString), http.StatusTemporaryRedirect) //
}

func handleOauthCallback(
	w http.ResponseWriter,
	r *http.Request,
	render *renderer.Render,
	dto userDto,
	integration userDtoIntegration,
	userDataFunc func(state, code string) (*users.OauthUserData, error),
	register bool) {
	state, code := r.FormValue("state"), r.FormValue("code")
	if state == "" || code == "" {
		helper.HandleResponse(w, render, http.StatusConflict, "state  or code is nul", nil)
		return
	}
	userData, err := userDataFunc(state, code)
	if err != nil {
		helper.HandleResponse(w, render, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if register {
		checkUser, err := integration.GetUsers(users.User{
			Email: userData.Email,
		})
		if err != nil {
			helper.HandleResponse(w, render, http.StatusInternalServerError, err, nil)
		}

		if checkUser!= nil {
            helper.HandleResponse(w, render, http.StatusFound, "user already registered", nil)
            return
        }

		userName := strings.ReplaceAll(strings.ToLower(userData.GivenName), " ", " ")

		response, err := dto.Register(users.User{
			Email:    userData.Email,
			Username: userName,
			Role:     "Admin",
			CategoryPreferences: []string{
				"Baju",
				"Celana",
			},
			Address: "Jakarta",
		})
		if err != nil {
			helper.HandleResponse(w, render, http.StatusInternalServerError, err, nil)
			return 

		}
		helper.HandleResponse(w, render, http.StatusOK, helper.SUCCESS_MESSSAGE, response)
	}

}
