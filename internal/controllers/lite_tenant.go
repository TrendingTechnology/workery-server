package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	null "gopkg.in/guregu/null.v4"

	"github.com/over55/workery-server/internal/idos"
	"github.com/over55/workery-server/internal/models"
)

func (h *Controller) liteTenantsListEndpoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx := r.Context()

	// Permission handling.
	role_id := uint64(ctx.Value("user_role_id").(int8))
	if role_id != 1 {
		http.Error(w, "Forbidden - You are not an administrator", http.StatusForbidden)
		return
	}

	// Fetch our URL parameters saved by our "Pagination" middleware.
	pageToken := ctx.Value("pageTokenParm").(uint64)
	pageSize := ctx.Value("pageSizeParam").(uint64)
	if pageSize == 0 || pageSize > 500 {
		pageSize = 100
	}

	//
	// Extract our parameters from the URL and create our filter.
	//

	stateParamString := r.FormValue("state")
	stateParam, stateParamErr := strconv.ParseInt(stateParamString, 10, 64)

	// Start by defining our base listing filter and then append depending on
	// different cases.
	f := &models.LiteTenantFilter{
		State:      null.NewInt(stateParam, stateParamErr == nil),
		LastSeenId: pageToken,
		Limit:      pageSize,
	}

	//
	// Submit our filter query to our database.
	//

	resultCh := make(chan []*models.LiteTenant)
	countCh := make(chan uint64)
	go func() {
		results, err := h.LiteTenantRepo.ListByFilter(ctx, f)
		if err != nil {
			log.Println("WARNING | liteTenantsListEndpoint | ListByFilter | err:", err)
		}
		resultCh <- results[:]
	}()
	go func() {
		count, err := h.LiteTenantRepo.CountByFilter(ctx, f)
		if err != nil { // For debugging purposes only.
			log.Println("WARNING | liteTenantsListEndpoint | CountByFilter | err:", err)
		}
		countCh <- count
	}()

	// Block the main function until we have results from our concurrently
	// running `goroutines`.
	results, count := <-resultCh, <-countCh

	// Take our data-layer results, serialize, and send to the user.
	responseData := idos.NewLiteTenantListResponseIDO(results, count)
	b, err := json.Marshal(&responseData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
