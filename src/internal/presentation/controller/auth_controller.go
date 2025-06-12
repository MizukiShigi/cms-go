package controller

import (
	"encoding/json"
	"net/http"

	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
	"github.com/MizukiShigi/cms-go/internal/presentation/helper"
	"github.com/MizukiShigi/cms-go/internal/usecase"

	"github.com/go-playground/validator/v10"
)

type AuthController struct {
	registerUserUsecase *usecase.RegisterUserUsecase
	loginUserUsecase    *usecase.LoginUserUsecase
	validator           *validator.Validate
}

func NewAuthController(registerUserUsecase *usecase.RegisterUserUsecase, loginUserUsecase *usecase.LoginUserUsecase) *AuthController {
	return &AuthController{
		registerUserUsecase: registerUserUsecase,
		loginUserUsecase:    loginUserUsecase,
		validator:           validator.New(),
	}
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (ac *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		myError := valueobject.NewMyError(valueobject.InvalidCode, "Invalid request payload")
		helper.RespondWithError(w, myError)
		return
	}

	err := ac.validator.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			myError := valueobject.NewMyError(valueobject.InvalidCode, err.Error())
			helper.RespondWithError(w, myError)
			return
		}
	}

	input := &usecase.RegisterUserInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := ac.registerUserUsecase.Execute(r.Context(), input)
	if err != nil {
		helper.RespondWithError(w, err)
		return
	}

	userResponse := UserResponse{
		ID:    output.ID,
		Name:  output.Name,
		Email: output.Email,
	}

	helper.RespondWithJSON(w, http.StatusCreated, userResponse)
}

func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		MyError := valueobject.NewMyError(valueobject.InvalidCode, "Invalid request payload")
		helper.RespondWithError(w, MyError)
		return
	}

	err := ac.validator.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			myError := valueobject.NewMyError(valueobject.InvalidCode, err.Error())
			helper.RespondWithError(w, myError)
			return
		}
	}

	input := &usecase.LoginUserInput{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := ac.loginUserUsecase.Execute(r.Context(), input)
	if err != nil {
		helper.RespondWithError(w, err)
		return
	}

	userResponse := UserResponse{
		ID:    output.User.ID,
		Name:  output.User.Name,
		Email: output.User.Email,
	}
	loginResponse := LoginResponse{
		Token: output.Token,
		User:  userResponse,
	}

	helper.RespondWithJSON(w, http.StatusOK, loginResponse)
}
