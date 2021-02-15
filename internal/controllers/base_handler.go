package controllers

import (
    "net/http"

    // "github.com/over55/workery-server/internal/repositories"
    "github.com/over55/workery-server/internal/models"
	"github.com/over55/workery-server/internal/session"
)

type BaseHandler struct {
    SecretSigningKeyBin []byte
    UserRepo models.UserRepository
    SessionManager *session.SessionManager
}

func NewBaseHandler(
    keyBin []byte,
    ur models.UserRepository,
    sm *session.SessionManager,
) (*BaseHandler) {
    return &BaseHandler{
        SecretSigningKeyBin: keyBin,
        UserRepo: ur,
        SessionManager: sm,
    }
}

func (h *BaseHandler) HandleRequests(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    // Get our URL paths which are slash-seperated.
    ctx := r.Context()
    p := ctx.Value("url_split").([]string)
    n := len(p)

    // Get our authorization information.
    isAuthsorized, ok := ctx.Value("is_authorized").(bool)

    switch {
    case n == 1 && p[0] == "version" && r.Method == http.MethodGet:
        if ok && isAuthsorized {
            h.getAuthenticatedVersion(w, r)
        } else {
            h.getVersion(w, r)
        }
    default:
        http.NotFound(w, r)
    }
}
