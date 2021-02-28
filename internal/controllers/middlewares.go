package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/over55/workery-server/internal/utils"
)

// Middleware will split the full URL path into slash-sperated parts and save to
// the context to flow downstream in the app for this particular request.
func (h *BaseHandler) URLProcessorMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Split path into slash-separated parts, for example, path "/foo/bar"
		// gives p==["foo", "bar"] and path "/" gives p==[""]. Our API starts with
		// "/api/v1", as a result we will start the array slice at "3".
		p := strings.Split(r.URL.Path, "/")[3:]

		// log.Println(p) // For debugging purposes only.

		// Open our program's context based on the request and save the
		// slash-seperated array from our URL path.
		ctx := r.Context()
		ctx = context.WithValue(ctx, "url_split", p)

		// Flow to the next middleware.
		fn(w, r.WithContext(ctx))
	}
}

func (h *BaseHandler) JWTProcessorMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		reqToken := r.Header.Get("Authorization")

		if reqToken != "" {

			// Special thanks to "poise" via https://stackoverflow.com/a/44700761
			splitToken := strings.Split(reqToken, "Bearer ")

			if len(splitToken) < 2 {
				http.Error(w, "not properly formatted authorization header", http.StatusBadRequest)
				return
			}

			reqToken = splitToken[1]

			// log.Println(reqToken) // For debugging purposes only.

			sessionUuid, err := utils.ProcessJWTToken(h.SecretSigningKeyBin, reqToken)
			if err == nil {
				// Update our context to save our JWT token content information.
				ctx = context.WithValue(ctx, "is_authorized", true)
				ctx = context.WithValue(ctx, "session_uuid", sessionUuid)

				// Flow to the next middleware with our JWT token saved.
				fn(w, r.WithContext(ctx))
				return
			}
			log.Println("JWTProcessorMiddleware | ProcessJWT | err", err)
		}

		// Flow to the next middleware without anything done.
		ctx = context.WithValue(ctx, "is_authorized", false)
		fn(w, r.WithContext(ctx))
	}
}

func (h *BaseHandler) AuthorizationMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get our authorization information.
		isAuthorized, ok := ctx.Value("is_authorized").(bool)
		if ok && isAuthorized {
			sessionUuid := ctx.Value("session_uuid").(string)

			//TODO: HANDLE CASE WHEN REDIS DB IS CLEARED BUT TOKEN AND USER RECORD ARE VALID.

			// Lookup our user profile in the session or return 500 error.
			user, err := h.SessionManager.GetUser(ctx, sessionUuid)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// If system administrator disabled the user account then we need
			// to generate a 403 error letting the user know their account has
			// been disabled and you cannot access the protected API endpoint.
			if user.State == 0 {
				http.Error(w, "Account disabled - please contact admin", http.StatusForbidden)
				return
			}

			// Save our user information to the context.
			ctx = context.WithValue(ctx, "user", user)
		}

		fn(w, r.WithContext(ctx))
	}
}

func (h *BaseHandler) PaginationMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Open our program's context based on the request and save the
		// slash-seperated array from our URL path.
		ctx := r.Context()

		// Setup our variables for the paginator.
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pageTokenString := r.FormValue("page_token")
		pageSizeString := r.FormValue("page_size")

		// Convert to unsigned 64-bit integer.
		pageToken, err := strconv.ParseUint(pageTokenString, 10, 64)
		if err != nil {
			// DEVELOPERS NOTE: ALWAYS DEFINE 100 IF NOT SPECIFIED OR ERROR.
			pageToken = 0
		}
		pageSize, err := strconv.ParseUint(pageSizeString, 10, 64)
		if err != nil {
			// DEVELOPERS NOTE: ALWAYS DEFINE 100 IF NOT SPECIFIED OR ERROR.
			pageSize = 100
		}

		// Attach the 'page' parameter value to our context to be used.
		ctx = context.WithValue(ctx, "pageTokenParm", pageToken)
		ctx = context.WithValue(ctx, "pageSizeParam", pageSize)

		// Flow to the next middleware.
		fn(w, r.WithContext(ctx))
	}
}

func (h *BaseHandler) AttachMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	// Attach our middleware handlers here. Please note that all our middleware
	// will start from the bottom and proceed upwards.
	// Ex: `URLProcessorMiddleware` will be executed first and
	//     `AuthorizationMiddleware` will be executed last.
	fn = h.AuthorizationMiddleware(fn) // Note: Must be above `JWTProcessorMiddleware`.
	fn = h.JWTProcessorMiddleware(fn)
	fn = h.PaginationMiddleware(fn)
	fn = h.URLProcessorMiddleware(fn)

	return func(w http.ResponseWriter, r *http.Request) {
		// Flow to the next middleware.
		fn(w, r)
	}
}
