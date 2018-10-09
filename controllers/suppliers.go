package controllers

import (
	"fmt"
	"net/http"

	"uni/coorse/db"
	// "../db"
)

//AllPreparationsController Get all preparation
//json responce
func AllSuppliersController(w http.ResponseWriter, r *http.Request) {
	all := db.GetAllSuppliers()

	m := make(map[string][]*db.Supplier)
	m["data"] = all
	b, err := GetOkJSON(all)
	if err != nil {
		err := &ApiError{fmt.Errorf("server error : %v", err), http.StatusInternalServerError}
		w.WriteHeader(err.StatusCode)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

type (
	FindSupplierByAddressAndNameRequest struct {
		Address string `valid:", required" schema:"address"`
		Name    string `valid:", required" schema:"name"`
	}
	FindSupplierByIdRequest struct {
		Id int `valid:"numeric, required" schema:"id"`
	}

	FindSupplierByCompanyRequest struct {
		Company string `valid:", required" schema:"company"`
	}

	InsertSupplierRequest struct {
		Name        string `valid:", required" schema:"name"`
		Company     string `valid:", optional" schema:"company"`
		Address     string `valid:", required" schema:"address"`
		Geolocation string `valid:", optional" schema:"geolocation"`
		Description string `valid:", optional" schema:"description"`
	}

	UpdateSupplierRequest struct {
		Id          int    `valid:"numeric, required" schema:"id"`
		Name        string `valid:", required" schema:"name"`
		Company     string `valid:", required" schema:"company"`
		Address     string `valid:", required" schema:"address"`
		Geolocation string `valid:", optional" schema:"geolocation"`
		Description string `valid:", optional" schema:"description"`
	}

	DeleteSupplierRequest struct {
		Id          int    `valid:"numeric, optional" schema:"id"`
		Name        string `valid:", optional" schema:"name"`
		Company     string `valid:", optional" schema:"company"`
		Address     string `valid:", optional" schema:"address"`
		Geolocation string `valid:", optional" schema:"geolocation"`
		Description string `valid:", optional" schema:"description"`
	}
)

func FindSupplierByIdController(w http.ResponseWriter, r *http.Request) {
	sr := &FindSupplierByIdRequest{}

	err := decode(sr, r.URL.Query())
	if __error_handle(&w, err) {
		return
	}
	_, err = valid(sr)
	if __error_handle(&w, err) {
		return
	}

	//found by name
	// if sr.Id < 0 {
	// 	e := &ApiError{fmt.Errorf("InvalidRequest -> id must be > 0"), http.StatusInternalServerError}
	// 	w.WriteHeader(e.StatusCode)
	// 	w.Write([]byte(e.Error()))
	// 	return
	// }

	supplier := db.FindSupplierById(sr.Id)

	b, err := GetOkJSON(supplier)
	if err != nil {
		err := &ApiError{fmt.Errorf("server error : %v", err), http.StatusInternalServerError}
		w.WriteHeader(err.StatusCode)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func FindSupplierByNameController(w http.ResponseWriter, r *http.Request) {
	sr := &FindSupplierByAddressAndNameRequest{}

	err := decode(sr, r.URL.Query())
	if __error_handle(&w, err) {
		return
	}
	_, err = valid(sr)
	if __error_handle(&w, err) {
		return
	}

	//found by name
	// if sr.Name == "" || sr.Address == "" {
	// 	e := &ApiError{fmt.Errorf("InvalidRequest : name and address must be not empty"), http.StatusInternalServerError}
	// 	w.WriteHeader(e.StatusCode)
	// 	w.Write([]byte(e.Error()))
	// 	return
	// }

	supplier := db.FindSupplierByNameAndAdress(sr.Name, sr.Address)

	b, err := GetOkJSON(supplier)
	if err != nil {
		err := &ApiError{fmt.Errorf("server error : %v", err), http.StatusInternalServerError}
		w.WriteHeader(err.StatusCode)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func FindSuppliersByCompanyController(w http.ResponseWriter, r *http.Request) {
	sr := &FindSupplierByCompanyRequest{}

	err := decode(sr, r.URL.Query())
	if __error_handle(&w, err) {
		return
	}
	_, err = valid(sr)
	if __error_handle(&w, err) {
		return
	}

	//found by name
	// if sr.Company == "" {
	// 	e := &ApiError{fmt.Errorf("InvalidRequest : name and address must be not empty"), http.StatusInternalServerError}
	// 	w.WriteHeader(e.StatusCode)
	// 	w.Write([]byte(e.Err.Error()))
	// 	return
	// }

	suppliers := db.FindSuppliersByCompany(sr.Company)

	b, err := GetOkJSON(suppliers)
	if err != nil {
		err := &ApiError{fmt.Errorf("server error : %v", err), http.StatusInternalServerError}
		w.WriteHeader(err.StatusCode)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func InsertSupplierController(w http.ResponseWriter, r *http.Request) {
	sr := &InsertSupplierRequest{}

	r.ParseForm()
	err := decode(sr, r.Form)
	if __error_handle(&w, err) {
		return
	}
	_, err = valid(sr)
	if __error_handle(&w, err) {
		return
	}

	//found by name
	// if sr.Company == "" || sr.Name == "" || sr.Address == "" {
	// 	e := &ApiError{fmt.Errorf("InvalidRequest -> name, company and address must be not empty"), http.StatusInternalServerError}
	// 	w.WriteHeader(e.StatusCode)
	// 	w.Write([]byte(e.Error()))
	// 	return
	// }

	stat := db.InsertIntoSuppliers(&db.Supplier{
		Name:        sr.Name,
		Company:     sr.Company,
		Address:     sr.Address,
		Geolocation: sr.Geolocation,
		Description: sr.Description,
	})

	if !stat {
		err := &ApiError{fmt.Errorf("db: supplier exists in"), http.StatusInternalServerError}
		w.WriteHeader(err.StatusCode)
		w.Write([]byte(err.Error()))
		return
	}

	WriteOkStatus(&w)
}

func UpdateSupplierController(w http.ResponseWriter, r *http.Request) {
	sr := &UpdateSupplierRequest{}

	r.ParseForm()
	err := decode(sr, r.Form)
	if __error_handle(&w, err) {
		return
	}
	_, err = valid(sr)
	if __error_handle(&w, err) {
		return
	}

	//found by name
	// if sr.Company == "" || sr.Name == "" || sr.Address == "" || sr.Id == 0 {
	// 	e := &ApiError{fmt.Errorf("InvalidRequest : name and address, company, id must be not empty"), http.StatusInternalServerError}
	// 	w.WriteHeader(e.StatusCode)
	// 	w.Write([]byte(e.Error()))
	// 	return
	// }

	stat := db.UpdateSupplier(&db.Supplier{
		Id:          sr.Id,
		Name:        sr.Name,
		Company:     sr.Company,
		Address:     sr.Address,
		Geolocation: sr.Geolocation,
		Description: sr.Description,
	})

	if !stat {
		err := &ApiError{fmt.Errorf("db: supplier doesn't exists in"), http.StatusInternalServerError}
		w.WriteHeader(err.StatusCode)
		w.Write([]byte(err.Error()))
		return
	}

	WriteOkStatus(&w)
}

func DeleteSupplierController(w http.ResponseWriter, r *http.Request) {
	sr := &DeleteSupplierRequest{}

	r.ParseForm()
	err := decode(sr, r.Form)
	if __error_handle(&w, err) {
		return
	}
	// _, err = valid(sr)
	// if __error_handle(&w, err) {
	// 	return
	// }

	//found by name
	if sr.Id == 0 && (sr.Name == "" || sr.Address == "") {
		e := &ApiError{fmt.Errorf("InvalidRequest -> id or name + address must be not empty"), http.StatusInternalServerError}
		w.WriteHeader(e.StatusCode)
		w.Write([]byte(e.Error()))
		return
	}

	stat := db.DeleteSupplier(&db.Supplier{
		Id:          sr.Id,
		Name:        sr.Name,
		Company:     sr.Company,
		Address:     sr.Address,
		Geolocation: sr.Geolocation,
		Description: sr.Description,
	})

	if stat == 0 {
		err := &ApiError{fmt.Errorf("db: supplier doesn't exists in"), http.StatusInternalServerError}
		w.WriteHeader(err.StatusCode)
		w.Write([]byte(err.Error()))
		return
	}

	WriteOkStatus(&w)
}
