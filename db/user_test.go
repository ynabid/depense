package db

import (
	"testing"
)

type TestCase struct {
	user   User
	passed bool
}

var testCases = []TestCase{
	{
		User{
			1,
			"yn.abid@gmail.com",
			"3811469",
			"Yassine",
			"ABID",
			"472953600",
		},
		true,
	},
	{
		User{
			2,
			"jai.chaimae@gmail.com",
			"chaimae",
			"Chaimae",
			"JAI",
			"472953600",
		},
		true,
	},
	{
		User{
			3,
			"jzai.chaimae@gmail.com",
			"chZaimae",
			"ChaZimae",
			"JAI",
			"472953600",
		},
		false,
	},
}

func TestUserInsert(t *testing.T) {
	for _, testCase := range testCases {
		s, err := Login(testCase.user.Email, testCase.user.Password)
		if err != nil {
			if !testCase.passed {
				return
			}
			t.Error(err.Error())
			return
		}
		if len(s.Id) < 20 {
			t.Error("Error in session ID : " + s.Id)
			return
		}

		_, err = IsAuthentified(s.UserId, s.Id)
		if err != nil {
			t.Error(err.Error())
			return
		}
		_, err = s.Delete()
		if err != nil {
			t.Error(err.Error())
			return
		}
	}
}
