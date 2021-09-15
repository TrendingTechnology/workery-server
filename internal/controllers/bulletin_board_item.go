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

func (h *Controller) bulletinBoardItemsListEndpoint(w http.ResponseWriter, r *http.Request) {
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
	stateParamString := r.FormValue("state")
	stateParam, _ := strconv.ParseUint(stateParamString, 10, 64)

	// Start by defining our base listing filter and then append depending on
	// different cases.
	f := models.LiteBulletinBoardItemFilter{
		TenantId:  tenantId,
		SortField: sortFieldString,
		SortOrder: sortOrderString,
		Search:    null.NewString(searchString, searchString != ""),
		Offset:    offsetParam,
		Limit:     limitParam,
		States:    []int8{int8(stateParam)},
	}

	// // For debugging purposes only.
	// log.Println("TenantId", f.TenantId)
	// log.Println("Search", f.Search)
	// log.Println("Offset", f.Offset)
	// log.Println("Limit", f.Limit)
	// log.Println("SortOrder", f.SortOrder)
	// log.Println("SortField", f.SortField)

	arrCh := make(chan []*models.LiteBulletinBoardItem)
	countCh := make(chan uint64)

	go func() {
		arr, err := h.LiteBulletinBoardItemRepo.ListByFilter(ctx, &f)
		if err != nil {
			log.Println("WARNING: bulletinBoardItemsListEndpoint|ListByFilter|err:", err.Error())
			arrCh <- nil
			return
		}
		arrCh <- arr[:]
	}()

	go func() {
		count, err := h.LiteBulletinBoardItemRepo.CountByFilter(ctx, &f)
		if err != nil {
			log.Println("WARNING: bulletinBoardItemsListEndpoint|CountByFilter|err:", err.Error())
			countCh <- 0
			return
		}
		countCh <- count
	}()

	arr, count := <-arrCh, <-countCh

	res := idos.NewLiteBulletinBoardItemListResponseIDO(arr, count)

	if err := json.NewEncoder(w).Encode(&res); err != nil { // [2]
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
