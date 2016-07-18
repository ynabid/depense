package db

import (
	"container/list"
	"database/sql"
)

type Account struct {
	Id   int64
	Name string
}
type AccountTR struct {
	Account string
	Credit  float64
	Debit   float64
	Total   float64
}
type AccountBilan struct {
	Account Account
	Amount  float64
}

func (account *Account) Insert() (sql.Result, error) {
	return db.Exec(
		"INSERT INTO account (Name) VALUES(?)",
		account.Name,
	)
}
func ReadAccounts() ([]Account, error) {
	l := list.New()
	rows, err := db.Query(
		"SELECT * FROM account ORDER BY name",
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var c Account
		rows.Scan(&c.Id, &c.Name)
		l.PushBack(c)
	}
	accountList := make([]Account, l.Len())
	i := 0
	for e := l.Front(); e != nil; e = e.Next() {
		accountList[i] = e.Value.(Account)
		i++
	}
	rows.Close()
	return accountList, nil
}
func ReadAccountsBilan(from, to int64) (map[string]float64, error) {
	var maps map[string]float64 = make(map[string]float64)
	rows, err := db.Query(
		`SELECT
			account.name,
			sum(account_tr.amount)
			FROM account_tr 
			JOIN account
			  ON account_tr.account_id = account.id
			WHERE date >= ? AND date <= ?
			GROUP BY account.name
			ORDER BY account.name
			`,
		from,
		to,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var name string
		var amount float64
		rows.Scan(
			&name,
			&amount,
		)
		maps[name] = amount
	}
	rows.Close()
	return maps, nil
}
func ReadAccountsTR(from, to int64) (map[string]map[string]float64, error) {
	maps := make(map[string]map[string]float64)
	rows, err := db.Query(
		`SELECT
			lacc.tr_id,
			lacc.type,
			account.name,
			lacc.amount
			FROM account_tr AS lacc
			JOIN account
			  ON lacc.account_id = account.id
			WHERE date >= ? AND date <= ?
			ORDER BY lacc.tr_id,lacc.type DESC
			`,
		from,
		to,
	)
	if err != nil {
		return nil, err
	}
	var last_tr_id int64 = 0
	var last_account string = ""
	var last_type int64 = -1
	var last_amount float64 = 0

	if maps[""] == nil {
		maps[""] = make(map[string]float64)
	}
	var tr_id, typ int64
	var account string
	var amount float64

	for rows.Next() {
		rows.Scan(
			&tr_id,
			&typ,
			&account,
			&amount,
		)
		if maps[last_account] == nil {
			maps[last_account] = make(map[string]float64)
		}
		if last_tr_id == tr_id {
			maps[last_account][account] += amount
			last_tr_id = 0
			last_account = ""
			last_type = -1
			last_amount = 0
		} else {
			if last_tr_id != 0 {
				if last_type == 0 {
					maps[""][last_account] += last_amount
				} else if last_type == 1 {
					maps[last_account][""] += last_amount
				}
			}
			last_tr_id = tr_id
			last_account = account
			last_type = typ
			last_amount = amount
		}
	}
	if last_tr_id != 0 {
		if last_type == 0 {
			maps[""][last_account] += amount
		} else if last_type == 1 {
			maps[last_account][""] += amount
		}
	}

	rows.Close()
	return maps, nil
}
