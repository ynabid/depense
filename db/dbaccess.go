package db

import (
	"container/list"
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"time"
)

type Depense struct {
	Id          int64
	Date        int64
	Description string
	Montant     float64
}

func openMysqlDB() (*sql.DB, error) {
	return sql.Open("mysql", "root:Ya3811469@tcp(localhost:3306)/mm")
}

func (d *Depense) Insert() (sql.Result, error) {
	db, err := openMysqlDB()
	if err != nil {
		return nil, err
	}
	return db.Exec(
		"INSERT INTO depense (date,description,montant) VALUES (?,?,?)",
		d.Date,
		d.Description,
		d.Montant,
	)
}
func ParseMonth(date string) (from, to int64, err error) {
	if date != "" && len(date) == 7 {
		t, err := time.Parse("2006-01-02", date+"-01")
		if err != nil {
			return 0, 0, err
		}
		from := t.Unix()
		year, err := strconv.ParseInt(date[0:4], 10, 64)
		if err != nil {
			return 0, 0, err
		}
		month, err := strconv.ParseInt(date[5:7], 10, 64)
		if err != nil {
			return 0, 0, err
		}
		if month == 12 {
			year++
			month = 1
		} else {
			month++
		}

		t, err = time.Parse(
			"2006-1-02",
			strconv.FormatInt(year, 10)+"-"+strconv.FormatInt(month, 10)+"-01",
		)
		if err != nil {
			return 0, 0, err
		}
		to := t.Unix() - 1
		return from, to, nil
	}
	return 0, 0, errors.New("Value of month is incorrect")
}

func DepenseMonth(m string) (float64, error) {
	var f float64
	from, to, err := ParseMonth(m)
	if err != nil {
		return 0.0, err
	}
	db, err := openMysqlDB()
	rows, err := db.Query(
		"SELECT sum(montant) FROM depense WHERE date >= ? and date <= ?",
		from,
		to,
	)
	if err != nil {
		return 0, err
	}
	rows.Next()
	rows.Scan(&f)
	return f, nil
}

func DepenseList(from, to int64) ([]Depense, error) {
	var depenseList []Depense
	l := list.New()
	db, err := openMysqlDB()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(
		"SELECT id,date,montant,description from depense WHERE date >= ? and date <= ? ORDER BY date DESC, id DESC",
		from,
		to,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var d Depense
		rows.Scan(&d.Id, &d.Date, &d.Montant, &d.Description)
		l.PushBack(d)
	}
	depenseList = make([]Depense, l.Len())
	i := 0
	for e := l.Front(); e != nil; e = e.Next() {
		depenseList[i] = e.Value.(Depense)
		i++
	}
	defer rows.Close()
	return depenseList, nil
}

func Read(id int64) (*Depense, error) {
	db, err := openMysqlDB()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(
		"SELECT date,montant,description from depense WHERE id=?",
		id,
	)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	rows.Next()
	var d Depense
	rows.Scan(&d.Date, &d.Montant, &d.Description)
	return &d, nil
}
