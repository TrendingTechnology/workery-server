package controllers

import (
	"net/http"
)

func (h *Controller) getVersion(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Workery Server v1.0"))
}

func (h *Controller) getAuthenticatedVersion(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Workery Server v1.0 with valid API Key"))
}
