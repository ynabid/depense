package main

import (
	"testing"
)

func TestDB(t *testing.T) {
	d := Depense{Date: 123, Montant: 123, Description: "ceci est un test"}
	r, err := d.insert()
	if err != nil {
		t.Error(err.Error())
	}
	id, err := r.LastInsertId()
	dr, err := read(id)
	if err != nil {
		t.Error(err.Error())
	}
	if dr.Date == 123 && dr.Montant == 123 && dr.Description == "ceci est un test" {

	} else {
		t.Error("Error when reading")
	}
}
