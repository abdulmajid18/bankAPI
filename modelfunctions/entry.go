package modelfunctions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (driver *DBClient) DepositMoney(w http.ResponseWriter, r *http.Request) {
	postBody, _ := ioutil.ReadAll(r.Body)
	var depositparam DepositMoneyParams

	err := json.Unmarshal(postBody, &depositparam)
	if err != nil {
		fmt.Println(err)
	}
	account, err := driver.DepositMoneyQuery(depositparam)
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

func (driver *DBClient) DepositMoneyQuery(depositparam DepositMoneyParams) (DepositInfo, error) {
	var account DepositInfo
	statement := "SELECT balance from accounts WHERE owner = $1"
	stmt, err := driver.Db.Prepare(statement)
	if err != nil {
		return account, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(depositparam.Owner)

	err = row.Scan(&account.Balance)
	if err != nil {
		fmt.Println(err)
	}
	statement = "UPDATE accounts SET balance = $1  WHERE owner = $2 RETURNING id, owner, balance, currency, created_at"
	stmt, err = driver.Db.Prepare(statement)

	if err != nil {
		return account, err
	}
	defer stmt.Close()

	account.Balance = account.Balance + depositparam.Amount
	row = stmt.QueryRow(account.Balance, depositparam.Owner)
	err = row.Scan(&account.ID, &account.Owner, &account.Balance, &account.Currency, &account.CreatedAt)
	if err != nil {
		fmt.Println(err)
	}
	statement = "INSERT INTO deposits (owner, amount, reference) VALUES ($1, $2, $3) RETURNING owner, amount, reference, created_at"
	stmt, err = driver.Db.Prepare(statement)
	if err != nil {
		return account, err
	}
	defer stmt.Close()
	row = stmt.QueryRow(depositparam.Owner, depositparam.Amount, depositparam.Reference)
	err = row.Scan(&account.Owner, &account.Amount, &account.Reference, &account.CreatedAt)
	if err != nil {
		fmt.Println(err)
	}

	return account, err
}

func (driver *DBClient) WithdrawMoney(w http.ResponseWriter, r *http.Request) {
	postBody, _ := ioutil.ReadAll(r.Body)
	var withdrawparam WithdrawMoneyParams

	err := json.Unmarshal(postBody, &withdrawparam)
	if err != nil {
		fmt.Println(err)
	}
	account, err := driver.WithdrawMoneyQuery(withdrawparam)
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

func (driver *DBClient) WithdrawMoneyQuery(withdrawparam WithdrawMoneyParams) (WithdrawInfo, error) {
	var account WithdrawInfo
	statement := "SELECT balance from accounts WHERE owner = $1"
	stmt, err := driver.Db.Prepare(statement)
	if err != nil {
		return account, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(withdrawparam.Owner)

	err = row.Scan(&account.Balance)
	if err != nil {
		fmt.Println(err)
	}
	statement = "UPDATE accounts SET balance = $1  WHERE owner = $2 RETURNING id, owner, balance, currency, created_at"
	stmt, err = driver.Db.Prepare(statement)

	if err != nil {
		return account, err
	}
	defer stmt.Close()

	account.Balance = account.Balance - withdrawparam.Amount
	row = stmt.QueryRow(account.Balance, withdrawparam.Owner)
	err = row.Scan(&account.ID, &account.Owner, &account.Balance, &account.Currency, &account.CreatedAt)
	if err != nil {
		fmt.Println(err)
	}
	statement = "INSERT INTO deposits (owner, amount, reference) VALUES ($1, $2, $3) RETURNING owner, amount, reference, created_at"
	stmt, err = driver.Db.Prepare(statement)
	if err != nil {
		return account, err
	}
	defer stmt.Close()
	row = stmt.QueryRow(withdrawparam.Owner, withdrawparam.Amount, withdrawparam.Reference)
	err = row.Scan(&account.Owner, &account.Amount, &account.Reference, &account.CreatedAt)
	if err != nil {
		fmt.Println(err)
	}

	return account, err
}

func (driver *DBClient) SendMoney(w http.ResponseWriter, r *http.Request) {
	postBody, _ := ioutil.ReadAll(r.Body)
	var transferparam TransferParams
	var senderParam WithdrawMoneyParams
	var receiverparam DepositMoneyParams
	err := json.Unmarshal(postBody, &transferparam)
	if err != nil {
		fmt.Println(err)
	}
	senderParam.Amount = transferparam.Amount
	senderParam.Owner = transferparam.Sender
	senderParam.Reference = transferparam.Reference
	receiverparam.Owner = transferparam.Receiver
	receiverparam.Amount = transferparam.Amount
	transferInfo, err := driver.SendMoneQuery(senderParam, receiverparam)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(transferInfo)
		w.Write(response)
	}

}

func (driver *DBClient) SendMoneQuery(senderparam WithdrawMoneyParams, receiverparam DepositMoneyParams) (TransferInfo, error) {
	var account TransferInfo
	// var sender DepositInfo
	// var receiver WithdrawInfo
	receiverAccount, err := driver.DepositMoneyQuery(receiverparam)
	if err != nil {
		fmt.Println(err)
	}
	senderAccount, err := driver.WithdrawMoneyQuery(senderparam)
	if err != nil {
		fmt.Println(err)
	}
	statement := "INSERT INTO transactions (amount, sender, reference, balance, receiver) VALUES ($1, $2, $3, $4, $5) RETURNING amount, sender, balance, receiver, reference, created_at"
	stmt, err := driver.Db.Prepare(statement)
	if err != nil {
		return account, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(senderAccount.Amount, senderAccount.Owner, senderAccount.Reference, senderAccount.Balance, receiverAccount.Owner)
	err = row.Scan(&account.Amount, &account.Sender, &account.Balance, &account.Receiver, &account.Reference, &senderAccount.CreatedAt)
	if err != nil {
		fmt.Println(err)
	}
	return account, err

}
