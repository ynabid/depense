package db

import (
	"container/list"
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"time"
)

var db *sql.DB

type Depense struct {
	Id          int64
	Date        int64
	Description string
	Montant     float64
	UserId      int64
	CategoryId  int64
}
type DepenseV struct {
	Id          int64
	Date        int64
	Description string
	Montant     float64
	UserId      int64
	Category    string
}

type DepenseCategory struct {
	Category Category
	Montant  float64
}

func OpenMysqlDB() (*sql.DB, error) {
	var err error
	if db != nil {
		return db, nil
	}
	db, err = sql.Open("mysql", "root:Ya3811469@tcp(localhost:3306)/mm")
	return db, err
}

func (d *Depense) Insert() (sql.Result, error) {
	return db.Exec(
		"INSERT INTO depense (date,description,montant,user_id,category_id) VALUES (?,?,?,?,?)",
		d.Date,
		d.Description,
		d.Montant,
		d.UserId,
		d.CategoryId,
	)
}

func DepenseByCategory(from, to int64) ([]DepenseCategory, error) {
	l := list.New()
	rows, err := db.Query(
		"SELECT category.id,category.name,sum(depense.montant) FROM depense LEFT JOIN category ON depense.category_id = category.id WHERE date >= ? and date <= ? GROUP BY category_id ORDER BY name",
		from,
		to,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var d DepenseCategory
		rows.Scan(&d.Category.Id, &d.Category.Name, &d.Montant)
		l.PushBack(d)
	}
	depenseList := make([]DepenseCategory, l.Len())
	i := 0
	for e := l.Front(); e != nil; e = e.Next() {
		depenseList[i] = e.Value.(DepenseCategory)
		i++
	}
	rows.Close()
	return depenseList, nil
}

func DepenseList(from, to int64) ([]DepenseV, error) {
	var depenseList []DepenseV
	l := list.New()
	rows, err := db.Query(
		"SELECT depense.id,date,montant,description,user_id,name FROM depense LEFT JOIN category ON depense.category_id = category.id WHERE date >= ? and date <= ? ORDER BY date DESC, id DESC",
		from,
		to,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var d DepenseV
		rows.Scan(&d.Id, &d.Date, &d.Montant, &d.Description, &d.UserId, &d.Category)
		l.PushBack(d)
	}
	depenseList = make([]DepenseV, l.Len())
	i := 0
	for e := l.Front(); e != nil; e = e.Next() {
		depenseList[i] = e.Value.(DepenseV)
		i++
	}
	rows.Close()
	return depenseList, nil
}

func Read(id int64) (*Depense, error) {
	rows, err := db.Query(
		"SELECT date,montant,description,user_id,category_id from depense WHERE id=?",
		id,
	)
	rows.Close()
	if err != nil {
		return nil, err
	}
	rows.Next()
	var d Depense
	rows.Scan(&d.Date, &d.Montant, &d.Description, &d.UserId, &d.CategoryId)
	return &d, nil
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
	rows.Close()
	return f, nil
}
