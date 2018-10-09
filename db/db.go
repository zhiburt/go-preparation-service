package db

//TODO ADD FIND PICTURE still in anather gorutine after insert in DB new preparation

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

const (
	DB_USER     = "user"
	DB_PASSWORD = "postgres"
	DB_NAME     = "test1"
)

var db *sql.DB

type (
	Preparation struct {
		Id               int
		Name             string
		Description      string
		Type             string
		ActiveIngredient string
		ImageURL         string
	}

	Supplier struct {
		Id          int
		Name        string
		Company     string
		Address     string
		Geolocation string
		Description string
	}

	PreparationOfSupplier struct {
		// Id int
		PreparationId int
		SupplierId    int
		Price         float64
	}
)

// ** BLOCK FOR PREPARATIONS

//AllPreparations will select all rows in DB
func GetAllPreparations() []*Preparation {
	// fmt.Println("# Querying")
	rows, err := db.Query("SELECT * FROM preparations")
	defer rows.Close()
	if err != nil {
		log.Errorf("PREPARATIONS::AllPreparations: %v", err)
		return nil
	}

	var preparations []*Preparation
	for rows.Next() {
		preparat := &Preparation{}
		err = rows.Scan(&preparat.Id, &preparat.Name, &preparat.Description, &preparat.Type, &preparat.ActiveIngredient, &preparat.ImageURL)
		checkErr(err)
		preparations = append(preparations, preparat)
	}

	log.Info("PREPARATIONS::AllPreparations: send all preparations")
	return preparations
}

//InsertIntoPreparations insert data into
func InsertIntoPreparations(p *Preparation) error {
	if p == nil {
		return fmt.Errorf("preparation must be no <nil>")
	}
	b, err := checkExistNamePreparation(p.Name)
	if err != nil {
		log.Errorf("InsertIntoPreparations with parametrs{name=%v} error: %v", p, err)
		return err
	}
	if b {
		return fmt.Errorf("Name %s exists in DB", p.Name)
	}

	var lastInsertId int
	err = db.QueryRow("INSERT INTO preparations(name, description,type, activeIngredient, imageURL) VALUES($1,$2,$3,$4,$5) returning id;",
		p.Name,
		p.Description,
		p.Type,
		p.ActiveIngredient,
		p.ImageURL).Scan(&lastInsertId) //scan for check and it's Copy & Past*)))
	if err != nil {
		return fmt.Errorf("scan error: %v", err)
	}

	log.Infof("preparations: insert %v into", p.Name)
	return nil
}

//FindPreparationByName find only one object
func FindPreparationByName(name string) *Preparation {
	rows, err := db.Query("SELECT * FROM preparations WHERE name = $1 LIMIT 1", name)
	if err != nil {
		__panic_text("FindPreparationByName SELECT err", err)
		return nil
	}
	defer rows.Close()

	if rows.Next() == false {
		log.Info("FindPreparationByName ", fmt.Sprintf("Preparation with NAME %q doesn't exist in DB", name))
		return nil
	}
	preparat := &Preparation{}
	err = rows.Scan(&preparat.Id, &preparat.Name, &preparat.Description, &preparat.Type, &preparat.ActiveIngredient, &preparat.ImageURL)
	if err != nil {
		__panic_text("FindPreparationByName", fmt.Errorf("scan error: %v", err))
		return nil
	}

	log.Infof("preparations: find %v in", name)
	return preparat
}

func FindPreparationById(id int) *Preparation {
	rows, err := db.Query("SELECT * FROM preparations WHERE id = $1 LIMIT 1", id)
	if err != nil {
		__panic_text("FindPreparationById SELECT err", err)
		return nil
	}
	defer rows.Close()

	if rows.Next() == false {
		log.Info("FindPreparationById ", fmt.Sprintf("Preparation with ID %q doesn't exist in DB", id))
		return nil
	}
	preparat := &Preparation{}
	err = rows.Scan(&preparat.Id, &preparat.Name, &preparat.Description, &preparat.Type, &preparat.ActiveIngredient, &preparat.ImageURL)
	if err != nil {
		__panic_text("FindPreparationById", fmt.Errorf("scan error: %v", err))
		return nil
	}

	log.Infof("preparations: find %v in", id)
	return preparat
}

func DeletePreparation(p *Preparation) int {
	if p == nil {
		log.Warning("DeletePreparation: ", "parametr is <nil>")
		return 0
	}
	var stmt *sql.Stmt
	var err error
	var res sql.Result
	if p.Name != "" {
		stmt, err = db.Prepare("delete from preparations where name=$1")
		checkErr(err)

		res, err = stmt.Exec(p.Name)
	} else {
		stmt, err = db.Prepare("delete from preparations where id=$1")
		checkErr(err)

		res, err = stmt.Exec(p.Id)
	}
	if err != nil {
		__panic_text("DeletePreparation", fmt.Errorf("DeletePreparation-Exec error: %v", err))
		return 0
	}

	affect, err := res.RowsAffected()
	checkErr(err)

	log.Infof("preparations: delete name=%s with Id=%d in", p.Name, p.Id)
	return int(affect)
}

func UpdatePreparation(p *Preparation) bool {
	// fmt.Println("# Updating")
	if p == nil {
		log.Warning("UpdatePreparation: ", "parametr is <nil>")
		return false
	}

	var stmt *sql.Stmt
	var err error
	var res sql.Result
	if p.Name != "" {
		stmt, err = db.Prepare("update preparations set description=$2, type=$3, activeIngredient=$4, imageURL=$5 where name=$1")
		checkErr(err)

		res, err = stmt.Exec(p.Name, p.Description, p.Type, p.ActiveIngredient, p.ImageURL)
	} else {
		stmt, err = db.Prepare("update preparations set name=$2, description=$3, type=$4, activeIngredient=$5, imageURL=$6 where id=$1")
		checkErr(err)

		res, err = stmt.Exec(p.Id, p.Name, p.Description, p.Type, p.ActiveIngredient, p.ImageURL)
	}
	if err != nil {
		__panic_text("UpdatePreparation: ", fmt.Errorf("UpdatePreparation-Exec error: %v", err))
		return false
	}

	//IT can comment
	affect, err := res.RowsAffected()
	checkErr(err)
	if affect != 1 {
		log.Warning("UpdatePreparation: We haven't changed anything preparation: %v")
		return false
	}

	// fmt.Println(affect, " rows changed")
	log.Infof("PREPARATIONS: update %v in", p)
	return true
}

// ** BLOCK FOR SUPPLIERS

func InsertIntoSuppliers(s *Supplier) bool {
	if s == nil {
		log.Warning("InsertIntoSuppliers: ", "parametr is <nil>")
		return false
	}
	b, err := checkExistsSupplier(s)
	if err != nil {
		__panic_text("InsertIntoSuppliers", err)
		return false
	}
	if b {
		log.Infof("Name %q exists in DB", s.Name)
		return false
	}

	var lastInsertId int
	err = db.QueryRow("INSERT INTO suppliers(name, company, address, geolocation, description) VALUES($1,$2,$3, $4, $5) returning id;",
		s.Name,
		s.Company,
		s.Address,
		s.Geolocation,
		s.Description).Scan(&lastInsertId) //scan for check and it's Copy & Past*)))
	if err != nil {
		__panic_text("InsertIntoSuppliers", fmt.Errorf("scan error: %v", err))
		return false
	}

	return true
}

func GetAllSuppliers() []*Supplier {
	// fmt.Println("# Querying")
	rows, err := db.Query("SELECT * FROM suppliers")
	defer rows.Close()
	if err != nil {
		__panic_text("SUPPLIERS:: ", fmt.Errorf("GetAllSuppliers: %v", err))
		return nil
	}

	var spl []*Supplier
	for rows.Next() {
		s := &Supplier{}
		err = rows.Scan(&s.Id, &s.Name, &s.Company, &s.Address, &s.Geolocation, &s.Description)
		checkErr(err)
		spl = append(spl, s)
	}

	return spl
}

func UpdateSupplier(s *Supplier) bool {
	if s == nil {
		log.Warning("InsertIntoSuppliers: ", "parametr is <nil>")
		return false
	}

	stmt, err := db.Prepare("update suppliers set name=$2, company=$3, address=$4, geolocation=$5, description=$6 where id=$1")
	checkErr(err)

	res, err := stmt.Exec(s.Id, s.Name, s.Company, s.Address, s.Geolocation, s.Description)
	if err != nil {
		log.Error("UpdateSupplier: exec error")
		__panic_text("SUPPLIERS:: ", fmt.Errorf("UpdateSupplier-Exec error: %v", err))
		return false
	}

	affect, err := res.RowsAffected()
	checkErr(err)
	if affect != 1 {
		log.Error("UpdateSupplier: affect != 1")
		__panic_text("SUPPLIERS:: ", fmt.Errorf("UpdatePreparation error: we change %d more then 1 or less value ", affect))
		return false
	}

	log.Info("SUPPLIERS: Deleted one with id", s.Id)
	return true
}

func FindSupplierByNameAndAdress(name string, address string) *Supplier {
	rows, err := db.Query("SELECT * FROM suppliers WHERE name = $1 AND address=$2 LIMIT 1", name, address)
	if err != nil {
		__panic_text("SELECT err", err)
	}
	defer rows.Close()

	if rows.Next() == false {
		log.Infof("SUPPLIERS: Supplier with NAME %q and ADDRESS %q doesn't exist in DB", name, address)
		return nil
	}
	s := &Supplier{}
	err = rows.Scan(&s.Id, &s.Name, &s.Company, &s.Address, &s.Geolocation, &s.Description)
	if err != nil {
		__panic_text("FindSupplierByNameAndAdress:: ", fmt.Errorf("scan error: %v", err))
		return nil
	}

	log.Infof("SUPPLIERS: find by %s and %s in", name, address)
	return s
}

func FindSupplierById(id int) *Supplier {
	rows, err := db.Query("SELECT * FROM suppliers WHERE id = $1 LIMIT 1", id)
	if err != nil {
		__panic_text("SELECT err", err)
	}
	defer rows.Close()

	if rows.Next() == false {
		log.Infof("FindSupplierById: Supplier with id %d doesn't exist in DB", id)
		return nil
	}
	s := &Supplier{}
	err = rows.Scan(&s.Id, &s.Name, &s.Company, &s.Address, &s.Geolocation, &s.Description)
	if err != nil {
		__panic_text("FindSupplierById:: ", fmt.Errorf("scan error: %v", err))
		return nil
	}

	log.Infof("SUPPLIERS::FindSupplierById find by id %d", id)
	return s
}

func FindSuppliersByCompany(company string) []*Supplier {
	rows, err := db.Query("SELECT * FROM suppliers WHERE company=$1", company)
	if err != nil {
		__panic_text("SELECT err", err)
	}
	defer rows.Close()

	var resp []*Supplier
	for rows.Next() {
		s := &Supplier{}
		err = rows.Scan(&s.Id, &s.Name, &s.Company, &s.Address, &s.Geolocation, &s.Description)
		if err != nil {
			__panic_text("FindSuppliersByCompany:: ", fmt.Errorf("scan error: %v", err))
			return nil
		}

		resp = append(resp, s)
	}

	log.Infof("SUPPLIERS: find by company %v in", company)
	return resp
}

func DeleteSupplier(p *Supplier) int {
	if p == nil {
		log.Warning("DeleteSupplier: ", "parametr is <nil>")
		return 0
	}
	// fmt.Println("# Deleting")
	stmt, err := db.Prepare("delete from suppliers where id=$1")
	checkErr(err)

	res, err := stmt.Exec(p.Id)
	if err != nil {
		__panic_text("DeleteSupplier", fmt.Errorf("DeleteSupplier-Exec error: %v", err))
		return 0
	}

	affect, err := res.RowsAffected()
	checkErr(err)

	log.Infof("SUPLIERS: delete %v with Id=%d in", p.Name, p.Id)
	return int(affect)
}

// * * BLOCK FOR SUPPLIERS_PREPARATIONS
// FOR MANY TO MANY

func InsertIntoPreparationOfSuppliers(s *Supplier, p *Preparation, price float64) bool {
	if s == nil || p == nil {
		log.Warning("InsertIntoPreparationOfSuppliers: ", "parametrs is <nil> must be no <nil>")
		return false
	}
	b, err := checkExistsPreparationOfSupplier(s, p)
	if err != nil {
		log.Error("InsertIntoPreparationOfSuppliers: ", err)
		return false
	}
	if b {
		log.Warningf("supplier %s has preparation %q (exists in DB)", s.Name, p.Name)
		return false
	}

	var lastInsertId int
	err = db.QueryRow("INSERT INTO suppliers_preparations(preparation_id, supplier_id, price) VALUES($1,$2,$3) returning id;",
		p.Id,
		s.Id,
		price).Scan(&lastInsertId) //scan for check and it's Copy & Past*)))
	if err != nil {
		log.Errorf("InsertIntoPreparationOfSuppliers: %v", fmt.Errorf("scan error: %v", err))
		//ERORR IF ID > N
		//__panic_text("InsertIntoPreparationOfSuppliers", fmt.Errorf("scan error: %v", err))
		return false
	}

	return true
}

func GetAllPreparationOfSupplier(s *Supplier) []*PreparationOfSupplier {
	if s == nil {
		log.Warning("GetAllPreparationOfSupplier: ", "parametr is <nil>")
		return nil
	}
	rows, err := db.Query("SELECT preparation_id, supplier_id, price  FROM suppliers_preparations WHERE supplier_id = $1", s.Id)
	defer rows.Close()
	if err != nil {
		__panic_text("SUPPLIERS:: ", fmt.Errorf("GetAllPreparationOfSupplier: %v", err))
		return nil
	}

	var spl []*PreparationOfSupplier
	for rows.Next() {
		s := &PreparationOfSupplier{}
		err = rows.Scan(&s.PreparationId, &s.SupplierId, &s.Price)
		checkErr(err)
		spl = append(spl, s)
	}

	return spl
}

func GetAllPreparationOfSupplierByPrep(p *Preparation) []*PreparationOfSupplier {
	if p == nil {
		log.Warning("GetAllPreparationOfSupplierByPrep: ", "parametr is <nil>")
		return nil
	}
	rows, err := db.Query("SELECT preparation_id, supplier_id, price  FROM suppliers_preparations WHERE preparation_id = $1", p.Id)
	defer rows.Close()
	if err != nil {
		__panic_text("SUPPLIERS:: ", fmt.Errorf("GetAllPreparationOfSupplierByPrep: %v", err))
		return nil
	}

	var spl []*PreparationOfSupplier
	for rows.Next() {
		s := &PreparationOfSupplier{}
		err = rows.Scan(&s.PreparationId, &s.SupplierId, &s.Price)
		checkErr(err)
		spl = append(spl, s)
	}

	return spl
}

func DeletePreparationOfSupplier(s *Supplier, p *Preparation) int {
	if p == nil || s == nil {
		log.Warning("DeletePreparationOfSupplier: ", "parametr is <nil>")
		return 0
	}
	// fmt.Println("# Deleting")
	stmt, err := db.Prepare("delete from suppliers_preparations where preparation_id=$1 and supplier_id = $2")
	checkErr(err)

	res, err := stmt.Exec(p.Id, s.Id)
	if err != nil {
		__panic_text("PreparationOfSupplier", fmt.Errorf("DeletePreparationOfSupplier-Exec error: %v", err))
		return 0
	}

	affect, err := res.RowsAffected()
	checkErr(err)

	log.Infof("PreparationOfSupplier: delete %v with Id=%d in", p.Name, p.Id)
	return int(affect)
}

//GetAllGoods withaut ID
func GetAllGoods() []*PreparationOfSupplier {
	// fmt.Println("# Querying")
	rows, err := db.Query("SELECT supplier_id, preparation_id, price FROM suppliers_preparations")
	defer rows.Close()
	if err != nil {
		__panic_text("suppliers_preparations:: ", fmt.Errorf("GetAllGoods: %v", err))
		return nil
	}

	var spl []*PreparationOfSupplier
	for rows.Next() {
		s := &PreparationOfSupplier{}
		err = rows.Scan(&s.SupplierId, &s.PreparationId, &s.Price)
		checkErr(err)
		spl = append(spl, s)
	}

	return spl
}

func main() {

	// fmt.Println("# Inserting values")

	// var lastInsertId int
	// err = db.QueryRow("INSERT INTO company(name,age,address, salary, test) VALUES($1,$2,$3,$4,$5) returning id;", "astaxie", 31, "Minsk obl 22a", 102.23, "test string").Scan(&lastInsertId)
	// checkErr(err)
	// fmt.Println("last inserted id =", lastInsertId)

	// fmt.Println("# Updating")
	// stmt, err := db.Prepare("update company set name=$1 where id=$2")
	// checkErr(err)

	// res, err := stmt.Exec("astaxieupdateVVV22", lastInsertId)
	// checkErr(err)

	// affect, err := res.RowsAffected()
	// checkErr(err)

	// fmt.Println(affect, "rows changed")

	// fmt.Println("# Querying")
	// rows, err := db.Query("SELECT * FROM company")
	// defer rows.Close()
	// checkErr(err)

	// for rows.Next() {
	// 	var id int
	// 	var name string
	// 	var age int
	// 	var address string
	// 	var salary float64
	// 	var test string
	// 	err = rows.Scan(&id, &name, &age, &address, &salary, &test)
	// 	checkErr(err)
	// 	fmt.Println("id | name | age | address | salary | test ")
	// 	fmt.Printf("%3v | %8v | %6v | %6v | %#8v | %8v\n", id, name, age, address, salary, test)
	// }

	// fmt.Println("# Deleting")
	// stmt, err = db.Prepare("delete from company where id=$1")
	// checkErr(err)

	// res, err = stmt.Exec(lastInsertId)
	// checkErr(err)

	// affect, err = res.RowsAffected()
	// checkErr(err)

	// fmt.Println(affect, "rows changed")
}

func checkExistNamePreparation(n string) (bool, error) {
	test, err := db.Query("SELECT id FROM preparations WHERE name = $1", n)
	if err != nil {
		return false, fmt.Errorf("checkExistNamePreparation with parametr %q : %v", n, err)
		__panic_text("SELECT err", err)

	}
	defer test.Close()
	if test.Next() != false {
		return true, nil
	}

	return false, nil
}

func checkExistsSupplier(s *Supplier) (bool, error) {
	test, err := db.Query("SELECT id FROM suppliers WHERE name = $1 AND address = $2 LIMIT 1", s.Name, s.Address)
	if err != nil {
		return false, fmt.Errorf("checkExistSupplier error with parametr <name: %s > <address: %s > : %v", s.Name, s.Address, err)
	}
	defer test.Close()
	if test.Next() != false {
		return true, nil
	}

	return false, nil
}

func checkExistsPreparationOfSupplier(s *Supplier, p *Preparation) (bool, error) {
	test, err := db.Query("SELECT id FROM suppliers_preparations WHERE preparation_id = $1 AND supplier_id=$2", p.Id, s.Id)
	if err != nil {
		return false, fmt.Errorf("checkExistNamePreparation with parametrs(supplier_id=%q, preparation_id=%q)  : %v", s.Id, p.Id, err)
		__panic_text("SELECT err", err)

	}
	defer test.Close()
	if test.Next() != false {
		return true, nil
	}

	return false, nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func __panic_text(text string, err error) {
	panic(text + ": " + err.Error())
}

//log init
func init() {
	log = logrus.WithFields(logrus.Fields{
		"DB": "postgre",
	})
	log.Logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	log.Logger.Level = logrus.InfoLevel

	// log.Logger.SetFormatter(&logrus.JSONFormatter{})
	// lf, err := os.OpenFile("./logs/_db/log.json.log", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	// if err != nil {
	// 	panic(err)
	// }
	// //defer lf.Close()
	// log.Logger.SetOutput(lf)
}

//I DON'T KNOW HAW TO DO IT NOW
//init db WITHAUT CLOSE!!!!
func init() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	var err error
	db, err = sql.Open("postgres", dbinfo)
	checkErr(err)
	//	// defer db.Close()
}
