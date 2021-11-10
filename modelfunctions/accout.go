package modelfunctions

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type DBClient struct {
	Db *sql.DB
}

func (driver *DBClient) CreateAccount(w http.ResponseWriter, r *http.Request) {
	postBody, _ := ioutil.ReadAll(r.Body)
	var account Account
	err := json.Unmarshal(postBody, &account)
	if err != nil {
		fmt.Println(err)
	}
	account, err = driver.InsertAccount(account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(account)
		w.Write(response)
	}
}

func (driver *DBClient) InsertAccount(account Account) (Account, error) {
	statement := "INSERT INTO accounts (owner, balance, currency)" +
		"VALUES ($1, $2, $3) RETURNING id, owner, balance, currency, created_at"
	stmt, err := driver.Db.Prepare(statement)
	if err != nil {
		return account, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(account.Owner, account.Balance, account.Currency)
	err = row.Scan(&account.ID, &account.Owner, &account.Balance, &account.Currency, &account.CreatedAt)
	if err != nil {
		fmt.Println(err)
	}
	return account, err
}

func (driver *DBClient) GetAccount(w http.ResponseWriter, r *http.Request) {
	postBody, _ := ioutil.ReadAll(r.Body)
	var account Account
	err := json.Unmarshal(postBody, &account)
	if err != nil {
		fmt.Println(err)
	}
	account, err = driver.ReadAccount(account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(account)
		w.Write(response)
	}

}

func (driver *DBClient) ReadAccount(account Account) (Account, error) {
	statement := "SELECT * FROM accounts WHERE id = $1 LIMIT 1"
	stmt, err := driver.Db.Prepare(statement)
	if err != nil {
		return account, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(account.ID)

	err = row.Scan(&account.ID, &account.Owner, &account.Balance, &account.Currency, &account.CreatedAt)
	if err != nil {
		fmt.Println(err)
	}
	return account, err
}

func (driver *DBClient) ListAccounts(w http.ResponseWriter, r *http.Request) {
	postBody, _ := ioutil.ReadAll(r.Body)
	var listparams ListAccountsParams
	err := json.Unmarshal(postBody, &listparams)
	if err != nil {
		fmt.Println(err)
	}
	items, err := driver.ReadAccounts(listparams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(items)
		w.Write(response)
	}

}

func (driver *DBClient) ReadAccounts(listparams ListAccountsParams) ([]Account, error) {
	statement := "SELECT  * FROM accounts ORDER BY id LIMIT $1 OFFSET $2"
	stmt, err := driver.Db.Prepare(statement)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(listparams.Limit, listparams.Offset)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	items := []Account{}
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Balance,
			&i.Currency,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil

}

func (driver *DBClient) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	postBody, _ := ioutil.ReadAll(r.Body)
	var account Account
	err := json.Unmarshal(postBody, &account)
	if err != nil {
		fmt.Println(err)
	}
	account, err = driver.QueryUpdateAccount(account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(account)
		w.Write(response)
	}
}

func (driver *DBClient) QueryUpdateAccount(account Account) (Account, error) {
	statement := "UPDATE accounts  SET balance = $1 WHERE id = $2 RETURNING id, owner, balance, currency, created_at"
	stmt, err := driver.Db.Prepare(statement)
	if err != nil {
		return account, err
	}
	defer stmt.Close()
	fmt.Println(account)
	row := stmt.QueryRow(account.Balance, account.ID)
	err = row.Scan(&account.ID, &account.Owner, &account.Balance, &account.Currency, &account.CreatedAt)
	if err != nil {
		fmt.Println(err)
	}
	return account, err
}

func (driver *DBClient) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	postBody, _ := ioutil.ReadAll(r.Body)
	var account Account
	err := json.Unmarshal(postBody, &account)
	if err != nil {
		fmt.Println(err)
	}
	account, err = driver.QueryDeleteAccount(account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(account)
		w.Write(response)
	}

}

func (driver *DBClient) QueryDeleteAccount(account Account) (Account, error) {
	statement := "DELETE FROM accounts WHERE id = $1 RETURNING id, owner"
	stmt, err := driver.Db.Prepare(statement)
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()
	row := stmt.QueryRow(account.ID)
	err = row.Scan(&account.ID, &account.Owner)
	if err != nil {
		fmt.Println(err)
	}
	return account, err

}
