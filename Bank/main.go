package main

import (
	"fmt"

	"log"
	"net/http"
	"time"

	"github.com/abdulmajid18/dbhelper"
	"github.com/abdulmajid18/modelfunctions"
	"github.com/gorilla/mux"
)

func main() {

	fmt.Println("Heyyy  Initializng the Database ! ")
	db, err := dbhelper.InitDB()
	if err != nil {
		fmt.Println(err)
	}

	dbclient := &modelfunctions.DBClient{Db: db}

	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/createAccount/", dbclient.CreateAccount).Methods("POST")
	r.HandleFunc("/getAccount/", dbclient.GetAccount).Methods("GET")
	r.HandleFunc("/listAccounts/", dbclient.ListAccounts).Methods("GET")
	r.HandleFunc("/updateAccount/", dbclient.UpdateAccount).Methods("PATCH")
	r.HandleFunc("/deleteAccount/", dbclient.DeleteAccount).Methods("DELETE")
	r.HandleFunc("/deposit/", dbclient.DepositMoney).Methods("POST")
	r.HandleFunc("/withdraw/", dbclient.WithdrawMoney).Methods("POST")
	r.HandleFunc("/sendMoney/", dbclient.SendMoney).Methods("POST")

	server := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())

}
