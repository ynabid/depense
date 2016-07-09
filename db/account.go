package db

import (
	"container/list"
	"database/sql"
)

type Account struct {
	Id   int64
	Name string
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
