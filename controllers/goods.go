package controllers

import (
	"fmt"
	"net/http"

	"uni/coorse/db"
)

type (
	//ALL SUPPLIER's goods
	AllSuppliersGoods struct {
		SupplierID int `valid:"numeric, required" schema:"supplier_id"`
	}

	//ALL PREPARATIONS AS GOODS
	AllPreparationsGoods struct {
		PreparationID int `valid:"numeric, required" schema:"preparation_id"`
	}

	InsertGoodRequest struct {
		PreparationID int     `valid:"numeric, required" schema:"preparation_id"`
		SupplierID    int     `valid:"numeric, required" schema:"supplier_id"`
		Price         float64 `valid:"float, required" schema:"price"`
	}

	DeleteGoodRequest struct {
		PreparationID int `valid:"numeric, required" schema:"preparation_id"`
		SupplierID    int `valid:"numeric, required" schema:"supplier_id"`
	}
)

func AllGoodsController(w http.ResponseWriter, r *http.Request) {

	goods := db.GetAllGoods()

	b, err := GetOkJSON(goods)
	if err != nil {
		err := &ApiError{fmt.Errorf("server error : %v", err), http.StatusInternalServerError}
		w.WriteHeader(err.StatusCode)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func AllSuppliersGoodsController(w http.ResponseWriter, r *http.Request) {
	sr := &AllSuppliersGoods{}

	err := decode(sr, r.URL.Query())
	if __error_handle(&w, err) {
		return
	}
	_, err = valid(sr)
	if __error_handle(&w, err) {
		return
	}

	goods := db.GetAllPreparationOfSupplier(&db.Supplier{Id: sr.SupplierID})

	b, err := GetOkJSON(goods)
	if err != nil {
		err := &ApiError{fmt.Errorf("server error : %v", err), http.StatusInternalServerError}
		w.WriteHeader(err.StatusCode)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
func AllPreparationsGoodsController(w http.ResponseWriter, r *http.Request) {
	sr := &AllPreparationsGoods{}

	err := decode(sr, r.URL.Query())
	if __error_handle(&w, err) {
		return
	}
	_, err = valid(sr)
	if __error_handle(&w, err) {
		return
	}

	// //found by name
	// if sr.SupplierID < 0 {
	// 	e := &ApiError{fmt.Errorf("InvalidRequest -> id must be > 0", err), http.StatusInternalServerError}
	// 	w.WriteHeader(e.StatusCode)
	// 	w.Write([]byte(e.Err.Error()))
	// 	return
	// }

	goods := db.GetAllPreparationOfSupplierByPrep(&db.Preparation{Id: sr.PreparationID})

	b, err := GetOkJSON(goods)
	if err != nil {
		err := &ApiError{fmt.Errorf("server error : %v", err), http.StatusInternalServerError}
		w.WriteHeader(err.StatusCode)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func InsertGoodController(w http.ResponseWriter, r *http.Request) {
	sr := &InsertGoodRequest{}

	r.ParseForm()
	err := decode(sr, r.Form)
	if __error_handle(&w, err) {
		return
	}
	_, err = valid(sr)
	if __error_handle(&w, err) {
		return
	}

	b := db.InsertIntoPreparationOfSuppliers(
		&db.Supplier{Id: sr.SupplierID},
		&db.Preparation{Id: sr.PreparationID},
		sr.Price)

	if !b {
		err := &ApiError{fmt.Errorf("good exists you need to delete an old good"), http.StatusInternalServerError}
		w.WriteHeader(err.StatusCode)
		w.Write([]byte(err.Error()))
		return
	}

	WriteOkStatus(&w)
}

func DeleteGoodController(w http.ResponseWriter, r *http.Request) {
	sr := &DeleteGoodRequest{}

	r.ParseForm()
	err := decode(sr, r.Form)
	if __error_handle(&w, err) {
		return
	}
	_, err = valid(sr)
	if __error_handle(&w, err) {
		return
	}

	i := db.DeletePreparationOfSupplier(
		&db.Supplier{Id: sr.SupplierID},
		&db.Preparation{Id: sr.PreparationID})

	if i == 0 {
		err := &ApiError{fmt.Errorf("good doesn't exists"), http.StatusInternalServerError}
		w.WriteHeader(err.StatusCode)
		w.Write([]byte(err.Error()))
		return
	}

	WriteOkStatus(&w)
}
