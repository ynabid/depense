package db

import (
	"container/list"
	"database/sql"
)

type Category struct {
	Id   int64
	Name string
}

func (category *Category) Insert() (sql.Result, error) {
	return db.Exec(
		"INSERT INTO category (Name) VALUES(?)",
		category.Name,
	)
}
func ReadAll() ([]Category, error) {
	l := list.New()
	rows, err := db.Query(
		"SELECT * FROM category ORDER BY name",
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var c Category
		rows.Scan(&c.Id, &c.Name)
		l.PushBack(c)
	}
	categoryList := make([]Category, l.Len())
	i := 0
	for e := l.Front(); e != nil; e = e.Next() {
		categoryList[i] = e.Value.(Category)
		i++
	}
	rows.Close()
	return categoryList, nil
}
