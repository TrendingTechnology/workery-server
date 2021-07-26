package controllers

import (
	// "log"
	"encoding/json"
	"net/http"
	"strconv"

	// null "gopkg.in/guregu/null.v4"

	// "github.com/over55/workery-server/internal/models"
	"github.com/over55/workery-server/internal/idos"
	"github.com/over55/workery-server/internal/validators"
)

func (h *Controller) tenantGetEndpoint(w http.ResponseWriter, r *http.Request, idStr string) {
	// Extract the session details from our "Session" middleware.
	ctx := r.Context()
	role := uint64(ctx.Value("user_role").(int8))

	// Permission handling - If use is not administrator then error.
	if role != 1 {
		http.Error(w, "Forbidden - You are not an administrator", http.StatusForbidden)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m, err := h.TenantRepo.GetById(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ido := idos.NewTenantIDO(m)
	if err := json.NewEncoder(w).Encode(&ido); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Controller) tenantUpdateEndpoint(w http.ResponseWriter, r *http.Request, idStr string) {
	// Extract the session details from our "Session" middleware.
	ctx := r.Context()
	role := uint64(ctx.Value("user_role").(int8))

	// Permission handling - If use is not administrator then error.
	if role != 1 {
		http.Error(w, "Forbidden - You are not an administrator", http.StatusForbidden)
		return
	}

	// Lookup the tenant based on the `ID` or error.
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	m, err := h.TenantRepo.GetById(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if m == nil {
		http.Error(w, "Tenant does not exist", http.StatusNotFound)
		return
	}

	// Get the user `PUT` data from the HTTP request.
	var putData *idos.TenantIDO
	if err := json.NewDecoder(r.Body).Decode(&putData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isValid, errStr := validators.ValidateTenantSaveFromRequest(putData)
	if isValid == false {
		http.Error(w, errStr, http.StatusBadRequest)
		return
	}

	// The ID to lookup by.
	m.Id = putData.Id
	m.SchemaName = putData.SchemaName
	m.Name = putData.Name
	m.Timezone = putData.Timezone
	m.State = putData.State

	// Update our record.
	err = h.TenantRepo.UpdateById(ctx, m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

    // Return our result
	ido := idos.NewTenantIDO(m)
	if err := json.NewEncoder(w).Encode(&ido); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// func (h *Controller) deleteTenantById(w http.ResponseWriter, r *http.Request, idStr string) {
// 	// Extract the session details from our "Session" middleware.
// 	ctx := r.Context()
// 	role := uint64(ctx.Value("user_role_id").(uint8))
//
// 	// Permission handling - If use is not administrator then error.
// 	if role != 1 {
// 		http.Error(w, "Forbidden - You are not an administrator", http.StatusForbidden)
// 		return
// 	}
// 	userId, err := strconv.ParseUint(idStr, 10, 64)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
//
// 	// Check to see if the tenant already exists.
// 	dto, err := h.FoundationService.GetUserById(userId)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	if dto == nil {
// 		http.Error(w, "Does not exist", http.StatusNotFound)
// 		return
// 	}
//
// 	// Perform the delete operation.
// 	w.WriteHeader(http.StatusNoContent)
// 	response, err := h.FoundationService.DeleteTenantById(userId)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
//
// 	if err := json.NewEncoder(w).Encode(&response); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }
