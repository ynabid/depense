package handlers

import (
	"errors"
	"net/url"
	"testing"
)

func parsePerMonthTest(m string, fromE, toE int64) error {
	v := url.Values{}
	v.Add("month", m)
	from, to, err := parsePerMonth(v)
	if err != nil {
		return err
	}
	if from != fromE {
		return errors.New("From value is incorrect")
	}
	if to != toE {
		return errors.New("To value is incorrect")
	}
	return nil

}
func TestParsePerMonth(t *testing.T) {
	err := parsePerMonthTest("2016-05", 1462060800, 1464739200)
	if err != nil {
		t.Error("2016-05 : " + err.Error())
	}
	err = parsePerMonthTest("2016-12", 1480550400, 1483228800)
	if err != nil {
		t.Error("2016-12 : " + err.Error())
	}
}
