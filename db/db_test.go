package db

import (
	"fmt"
	"os"
	"testing"
)

var depenses = []Depense{
	//2016-05
	{
		0,
		1463184000, //2016-06-14
		"Gaz",
		42,
	},
	{
		0,
		1462147200, //2016-06-02
		"Biomil",
		79.5,
	},
	{
		0,
		1463702400, //2016-06-20
		"Dodot",
		140,
	},

	//2016-06
	{
		0,
		1465862400, //2016-06-14
		"Dodot",
		140,
	},
	{
		0,
		1464739200, //2016-06-01
		"Biomil",
		79.5,
	},
	{
		0,
		1464825600, //2016-06-02
		"Cerelac",
		40,
	},
}

type dmer struct { //Depense Month Expected Result
	Month  string
	Amount float64
}

var dmers = []dmer{
	{
		"2016-05",
		261.5,
	},
	{
		"2016-06",
		259.5,
	},
}

/*
func TestInsertRead(t *testing.T) {
	d := Depense{Date: 123, Montant: 123, Description: "ceci est un test"}
	r, err := d.Insert()
	if err != nil {
		t.Error(err.Error())
	}
	id, err := r.LastInsertId()
	dr, err := Read(id)
	if err != nil {
		t.Error(err.Error())
	}
	if dr.Date == 123 && dr.Montant == 123 && dr.Description == "ceci est un test" {

	} else {
		t.Error("Error when reading")
	}
}
*/

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	teardown()
	os.Exit(retCode)
}
func TestDepenseMonth(t *testing.T) {
	for _, dm := range dmers {
		f, _ := DepenseMonth(dm.Month)
		if f != dm.Amount {
			t.Error(
				fmt.Sprintf(
					"Month %s \t: Actual\t: %f\t Expected\t:%f",
					dm.Month,
					f,
					dm.Amount,
				),
			)
		}
	}
}
func setup() {
	for _, d := range depenses {
		_, err := d.Insert()
		if err != nil {
			panic(err.Error())
		}
	}
}

func teardown() {
	db, err := openMysqlDB()
	if err != nil {
		panic(err.Error())
	}

	db.Exec(
		"DELETE FROM depense WHERE 1",
	)
}
