package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/xsadia/secred/internal"
	"github.com/xsadia/secred/repository"
)

func (s *Server) GetWareHouseItemsHandler(w http.ResponseWriter, r *http.Request) {
	ah := r.Header.Get("Authorization")

	token, err := internal.ValidateAuthHeader(ah)

	if err != nil {
		internal.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	_, err = internal.VerifyToken(token)

	if err != nil {
		internal.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var count, start int
	startQuery := r.URL.Query().Get("start")
	countQuery := r.URL.Query().Get("count")

	if startQuery == "" {
		start = 0
	} else {
		start, _ = strconv.Atoi(startQuery)
	}

	if countQuery == "" {
		count = 10
	} else {
		count, _ = strconv.Atoi(countQuery)
	}

	items, err := repository.GetWarehouseItems(s.DB, start, count)

	if err != nil {
		internal.RespondWithError(w, http.StatusInternalServerError, internalServerError)
		return
	}

	internal.RespondWithJSON(w, http.StatusOK, items)
}

func (s *Server) GetWareHouseItemHandler(w http.ResponseWriter, r *http.Request) {
	ah := r.Header.Get("Authorization")

	token, err := internal.ValidateAuthHeader(ah)

	if err != nil {
		internal.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	_, err = internal.VerifyToken(token)

	if err != nil {
		internal.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var wi repository.WarehouseItem

	vars := mux.Vars(r)

	wi.Id = vars["id"]

	if err = wi.GetWarehouseItemById(s.DB); err != nil {
		internal.RespondWithError(w, http.StatusNotFound, "Item not found")
		return
	}

	internal.RespondWithJSON(w, http.StatusOK, wi)
}

func (s *Server) CreateWarehouseItemHandler(w http.ResponseWriter, r *http.Request) {
	ah := r.Header.Get("Authorization")

	token, err := internal.ValidateAuthHeader(ah)

	if err != nil {
		internal.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	_, err = internal.VerifyToken(token)

	if err != nil {
		internal.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var wi repository.WarehouseItem

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&wi); err != nil {
		internal.RespondWithError(w, http.StatusBadRequest, invalidRequestPayloadError)
		return
	}

	defer r.Body.Close()

	if err := wi.CreateWarehouseItem(s.DB); err != nil {
		internal.RespondWithError(w, http.StatusConflict, "Item already registered")
		return
	}

	internal.RespondWithJSON(w, http.StatusCreated, wi)
}

func (s *Server) UpdateWarehouseItemHandler(w http.ResponseWriter, r *http.Request) {
	ah := r.Header.Get("Authorization")

	token, err := internal.ValidateAuthHeader(ah)

	if err != nil {
		internal.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	_, err = internal.VerifyToken(token)

	if err != nil {
		internal.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	vars := mux.Vars(r)

	var wi repository.WarehouseItem

	decoder := json.NewDecoder(r.Body)

	if err = decoder.Decode(&wi); err != nil {
		internal.RespondWithError(w, http.StatusBadRequest, invalidRequestPayloadError)
		return
	}

	defer r.Body.Close()

	wi.Id = vars["id"]

	if err = wi.GetWarehouseItemById(s.DB); err != nil {
		internal.RespondWithError(w, http.StatusNotFound, "Item not found")
		return
	}

	if err = wi.UpdateWarehouseItem(s.DB); err != nil {
		internal.RespondWithError(w, http.StatusInternalServerError, internalServerError)
		return
	}

	internal.RespondWithJSON(w, http.StatusNoContent, nil)
}

func (s *Server) DeleteWarehouseItemHandler(w http.ResponseWriter, r *http.Request) {
	ah := r.Header.Get("Authorization")

	token, err := internal.ValidateAuthHeader(ah)

	if err != nil {
		internal.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	_, err = internal.VerifyToken(token)

	if err != nil {
		internal.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	vars := mux.Vars(r)

	var wi repository.WarehouseItem

	wi.Id = vars["id"]

	if err = wi.GetWarehouseItemById(s.DB); err != nil {
		internal.RespondWithError(w, http.StatusNotFound, "Item not found")
		return
	}

	if err = wi.DeleteWarehouseItem(s.DB); err != nil {
		internal.RespondWithError(w, http.StatusInternalServerError, internalServerError)
		return
	}

	internal.RespondWithJSON(w, http.StatusNoContent, nil)
}

func (s *Server) UploadCSVWarehouse(w http.ResponseWriter, r *http.Request) {
	ah := r.Header.Get("Authorization")

	token, err := internal.ValidateAuthHeader(ah)

	if err != nil {
		internal.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	_, err = internal.VerifyToken(token)

	if err != nil {
		internal.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	r.ParseMultipartForm(10 << 20)

	file, header, err := r.FormFile("csvFile")

	if err != nil {
		internal.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = internal.CheckMIMEContentType(header.Header); err != nil {
		internal.RespondWithError(w, http.StatusUnsupportedMediaType, err.Error())
		return
	}

	defer file.Close()

	tempFile, err := ioutil.TempFile("./temp", "upload-*.csv")

	if err != nil {
		internal.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		internal.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	tempFile.Write(fileBytes)

	wil, err := internal.ParseCSV(tempFile.Name())

	if err != nil {
		internal.RespondWithError(w, http.StatusNotAcceptable, err.Error())
		return
	}

	for _, wi := range wil {
		go func(curr repository.WarehouseItem) {
			curr.UpSertWarehouseItem(s.DB)
		}(wi)
	}

	internal.RespondWithJSON(w, http.StatusNoContent, nil)
}
