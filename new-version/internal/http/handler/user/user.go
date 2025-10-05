package user

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"new-version/internal/config"
	"new-version/internal/contract/user"
	userDto "new-version/internal/contract/user"
	mw "new-version/internal/http-server/middleware"
	"new-version/pkg/httphelpers"
	userSvc "new-version/internal/service/user"
	"new-version/pkg/json"
	"time"
)

type Handler interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
	LoginUser(w http.ResponseWriter, r *http.Request)
	LogoutUser(w http.ResponseWriter, r *http.Request)
}

type DefaultHandler struct {
	log *slog.Logger
	svc userSvc.Service
	cfg *config.Security
}

func (u *DefaultHandler) RegisterRoutes(mux *http.ServeMux, log *slog.Logger) {
	logMw := mw.LoggerMiddleware(log)
	authMw := mw.AuthMiddleware(u.cfg)

	mux.Handle("POST /user/register", logMw(http.HandlerFunc(u.RegisterUser)))
	mux.Handle("POST /user/login", logMw(http.HandlerFunc(u.LoginUser)))
	mux.Handle("POST /user/logout",
		authMw(logMw(http.HandlerFunc(u.LogoutUser)), httphelpers.USER_LVL))
}

func New(
	log *slog.Logger,
	srv userSvc.Service,
	cfg *config.Security,
) *DefaultHandler {
	return &DefaultHandler{
		log: log,
		svc: srv,
		cfg: cfg,
	}
}

// Register adds a new user to library.
// @ID registerUser
// @Summary Register
// @Tags user
// @Description register a new user
// @Accept json
// @Produce json
// @Param req body user.Request true "UserCreate"
// @Success 201 {object} httphelpers.Response
// @Failure 400 {object} httphelpers.Response
// @Failure 500 {object} httphelpers.Response
// @Failure default {object} httphelpers.Response
// @Router /user/register [post]
func (u *DefaultHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	const op = "modules.user.handler.Register"

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()
	defer r.Body.Close()

	var req userDto.Request

	if err := json.ReadRequestBody(r, &req); err != nil {
		json.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := u.svc.Register(ctx, req)
	if err != nil {
		json.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.WriteSuccess(w, "user has been registered", nil, http.StatusCreated)
}

// Login allows a user to sign in and get permissions.
// @ID loginUser
// @Summary Login
// @Tags user
// @Description login user
// @Accept json
// @Produce json
// @Param req body user.Request true "UserLogin"
// @Success 200 {object} httphelpers.Response
// @Failure 400 {object} httphelpers.Response
// @Failure 500 {object} httphelpers.Response
// @Failure default {object} httphelpers.Response
// @Router /user/login [post]
func (u *DefaultHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	const op = "modules.user.handler.Login"

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()
	defer r.Body.Close()

	var req user.Request

	if err := json.ReadRequestBody(r, &req); err != nil {
		json.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := u.svc.Login(ctx, req)
	if err != nil {
		json.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println(token)

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(u.cfg.AccessTokenExpire),
	})

	json.WriteSuccess(w, "successful login", nil, http.StatusOK)
}

// Logout allows a user to sign out from system and to be protected.
// @ID logoutUser
// @Summary Logout
// @Tags user
// @Description logout user
// @Produce json
// @Success 200 {object} httphelpers.Response
// @Failure 400 {object} httphelpers.Response
// @Failure 500 {object} httphelpers.Response
// @Failure default {object} httphelpers.Response
// @Router /user/logout [post]
func (u *DefaultHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	const op = "modules.user.handler.Logout"

	defer r.Body.Close()

	cookie, err := r.Cookie("access_token")
	if err != nil {
		json.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookie.Name,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})

	json.WriteSuccess(w, "successful logout", nil, http.StatusOK)
}
