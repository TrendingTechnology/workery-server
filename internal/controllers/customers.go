package controllers

import (
	// "encoding/json"
	"encoding/json"
	"log"
	"net/http"
	// "time"
	"strconv"

	// "github.com/google/uuid"
	null "gopkg.in/guregu/null.v4"

	"github.com/over55/workery-server/internal/idos"
	"github.com/over55/workery-server/internal/models"
)

func (h *Controller) customersListEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantId := uint64(ctx.Value("user_tenant_id").(uint64))
	// userId := uint64(ctx.Value("user_id").(uint64))

	// Extract our parameters from the URL.
	offsetParamString := r.FormValue("offset")
	offsetParam, _ := strconv.ParseUint(offsetParamString, 10, 64)
	limitParamString := r.FormValue("limit")
	limitParam, _ := strconv.ParseUint(limitParamString, 10, 64)
	if limitParam == 0 || limitParam > 500 {
		limitParam = 100
	}
	searchString := r.FormValue("search")
	sortOrderString := r.FormValue("sort_order")
	if sortOrderString == "" {
		sortOrderString = "ASC"
	}
	sortFieldString := r.FormValue("sort_field")
	if sortFieldString == "" {
		sortFieldString = "last_name"
	}

	// Start by defining our base listing filter and then append depending on
	// different cases.
	f := models.LiteCustomerFilter{
		TenantId:  tenantId,
		SortField: sortFieldString,
		SortOrder: sortOrderString,
		Search:    null.NewString(searchString, searchString != ""),
		Offset:    offsetParam,
		Limit:     limitParam,
	}

	// // For debugging purposes only.
	// log.Println("TenantId", f.TenantId)
	// log.Println("Search", f.Search)
	// log.Println("Offset", f.Offset)
	// log.Println("Limit", f.Limit)
	// log.Println("SortOrder", f.SortOrder)
	// log.Println("SortField", f.SortField)

	arrCh := make(chan []*models.LiteCustomer)
	countCh := make(chan uint64)

	go func() {
		arr, err := h.LiteCustomerRepo.ListByFilter(ctx, &f)
		if err != nil {
			log.Println("WARNING: customersListEndpoint|ListByFilter|err:", err.Error())
			arrCh <- nil
			return
		}
		arrCh <- arr[:]
	}()

	go func() {
		count, err := h.LiteCustomerRepo.CountByFilter(ctx, &f)
		if err != nil {
			log.Println("WARNING: customersListEndpoint|CountByFilter|err:", err.Error())
			countCh <- 0
			return
		}
		countCh <- count
	}()

	arr, count := <-arrCh, <-countCh

	res := idos.NewLiteCustomerListResponseIDO(arr, count)

	if err := json.NewEncoder(w).Encode(&res); err != nil { // [2]
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Controller) customerGetEndpoint(w http.ResponseWriter, r *http.Request, idStr string) {
	defer r.Body.Close()

	//
	// Get the customer based on the primary key.
	//

	// Extract the session details from our "Session" middleware.
	ctx := r.Context()
	tenantId := uint64(ctx.Value("user_tenant_id").(uint64))
	role_id := uint64(ctx.Value("user_role_id").(int8))

	// Permission handling - If use is not administrator then error.
	if role_id != 1 {
		http.Error(w, "Forbidden - You are not an administrator", http.StatusForbidden)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mCh := make(chan *models.Customer)

	go func() {
		m, err := h.CustomerRepo.GetById(ctx, id)
		if err != nil {
			mCh <- nil
		} else {
			mCh <- m
		}
	}()

	//
	// Get all the tags with this customer.
	//

	tCh := make(chan []*models.CustomerTag)

	go func() {
		// Lookup the tags that belong to the customer.
		f := &models.CustomerTagFilter{
			TenantId:   tenantId,
			SortField:  "tag_id",
			SortOrder:  "ASC",
			CustomerId: null.NewInt(int64(id), id != 0),
			Offset:     0,
			Limit:      1000,
		}
		tags, err := h.CustomerTagRepo.ListByFilter(ctx, f)
		if err != nil {
			tCh <- nil
		} else {
			tCh <- tags
		}
	}()

	//
	// Wait for all our database requests to finish and then attach the
	// compiled data.
	//

	m, tags := <-mCh, <-tCh
	m.Tags = tags

	//
	// Serialize the data.
	//

	ido := idos.NewCustomerIDO(m)
	if err := json.NewEncoder(w).Encode(&ido); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
