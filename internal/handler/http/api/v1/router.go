package v1

import (
	"github.com/gorilla/mux"
	"labra/pkg/middleware"
	"net/http"
)

func (h *Handler) GetVersion() string {
	return "v1"
}

func (h *Handler) GetContentType() string {
	return ""
}

func (h *Handler) AddRoutes(r *mux.Router) {
	r.Use(middleware.CORSMiddleware())
	r.HandleFunc("/account/profiles/{profile_id}/checkups/ocr", h.ScanResult).Methods(http.MethodOptions)

	r.HandleFunc("/checkups", h.UpdateCheckup).Methods(http.MethodOptions)
	r.HandleFunc("/account/profiles/{profile_id}/charts", h.Charts).Methods(http.MethodOptions)
	r.HandleFunc("/account/profiles/{profile_id}/checkups", h.CheckupList).Methods(http.MethodOptions)
	r.HandleFunc("/account/profiles", h.UserProfiles).Methods(http.MethodOptions)
	r.HandleFunc("/account/profiles", h.AddProfile).Methods(http.MethodOptions)
	r.HandleFunc("/patients/{id}", h.DeletePatient).Methods(http.MethodOptions)
	r.HandleFunc("/checkups/{id}", h.CheckupDetails).Methods(http.MethodOptions)
	r.HandleFunc("/dictionaries", h.GetDictionaries).Methods(http.MethodOptions)

	r.HandleFunc("/sign_in", h.SignIn).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/signup", h.SignUp).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/user/verify", h.VerifyUser).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/otp", h.SendOTP).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/session/refresh", h.RefreshToken).Methods(http.MethodPost, http.MethodOptions)

	authAPI := r.PathPrefix("").Subrouter()
	authAPI.Use(middleware.AuthorizationMiddleware(h.secret))

	authAPI.HandleFunc("/account/profiles/{profile_id}/checkups", h.CheckupList).Methods(http.MethodGet)
	authAPI.HandleFunc("/checkups", h.UpdateCheckup).Methods(http.MethodPatch)

	authAPI.HandleFunc("/account/profiles/{profile_id}/charts", h.Charts).Methods(http.MethodGet)
	authAPI.HandleFunc("/marker", h.Charts).Methods(http.MethodGet)
	authAPI.HandleFunc("/dictionaries", h.GetDictionaries).Methods(http.MethodGet)

	authAPI.HandleFunc("/account/profiles/{profile_id}/checkups/ocr", h.ScanResult).Methods(http.MethodPost)

	authAPI.HandleFunc("/checkups/{id}", h.CheckupDetails).Methods(http.MethodGet)

	authAPI.HandleFunc("/account/profiles", h.UserProfiles).Methods(http.MethodGet)
	authAPI.HandleFunc("/account/profiles", h.AddProfile).Methods(http.MethodPost)
	authAPI.HandleFunc("/patients/{id}", h.DeletePatient).Methods(http.MethodDelete)
}
