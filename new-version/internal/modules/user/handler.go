package user

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"new-version/internal/config"
	mw "new-version/internal/http-server/middleware"
	"new-version/internal/modules/common"
	"new-version/pkg/json"
	"time"
)

type UserHandler interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
	LoginUser(w http.ResponseWriter, r *http.Request)
	LogoutUser(w http.ResponseWriter, r *http.Request)
}

type UserHandlerImpl struct {
	log *slog.Logger
	srv UserService
	cfg *config.Security
}

func (u *UserHandlerImpl) RegisterRoutes(mux *http.ServeMux, log *slog.Logger) {
	logMw := mw.LoggerMiddleware(log)
	authMw := mw.AuthMiddleware(u.cfg)

	mux.Handle("POST /user/register", logMw(http.HandlerFunc(u.RegisterUser)))
	mux.Handle("POST /user/login", logMw(http.HandlerFunc(u.LoginUser)))
	mux.Handle("POST /user/logout", authMw(logMw(http.HandlerFunc(u.LogoutUser)), common.USER_ACCESS_LEVEL))
}

func NewUserHandler(log *slog.Logger, srv UserService, cfg *config.Security) *UserHandlerImpl {
	return &UserHandlerImpl{
		log: log,
		srv: srv,
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
// @Param req body UserCreate true "UserCreate"
// @Success 201 {object} httphelpers.Response
// @Failure 400 {object} httphelpers.Response
// @Failure 500 {object} httphelpers.Response
// @Failure default {object} httphelpers.Response
// @Router /user/register [post]
func (u *UserHandlerImpl) RegisterUser(w http.ResponseWriter, r *http.Request) {
	const op = "modules.user.handler.Register"

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()
	defer r.Body.Close()

	var req UserCreate

	if err := json.ReadRequestBody(r, &req); err != nil {
		json.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := u.srv.Register(ctx, req)
	if err != nil {
		json.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.WriteSuccess(w, "registered user", map[string]string{"id": id.String()}, http.StatusCreated)
}

// Login allows a user to sign in and get permissions.
// @ID loginUser
// @Summary Login
// @Tags user
// @Description login user
// @Accept json
// @Produce json
// @Param req body UserLogin true "UserLogin"
// @Success 200 {object} httphelpers.Response
// @Failure 400 {object} httphelpers.Response
// @Failure 500 {object} httphelpers.Response
// @Failure default {object} httphelpers.Response
// @Router /user/login [post]
func (u *UserHandlerImpl) LoginUser(w http.ResponseWriter, r *http.Request) {
	const op = "modules.user.handler.Login"

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()
	defer r.Body.Close()

	var req UserLogin

	if err := json.ReadRequestBody(r, &req); err != nil {
		json.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := u.srv.Login(ctx, req)
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
func (u *UserHandlerImpl) LogoutUser(w http.ResponseWriter, r *http.Request) {
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
