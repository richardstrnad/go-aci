package models

import "encoding/json"

type Login struct {
	AaaUser AaaUser `json:"aaaUser"`
}

type AaaUser struct {
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	Name string `json:"name"`
	PWD  string `json:"pwd"`
}

func (l *Login) ToJSON() ([]byte, error) {
	return json.Marshal(l)
}

func NewLogin(username, password string) *Login {
	return &Login{
		AaaUser{
			Attributes{
				Name: username,
				PWD:  password,
			},
		},
	}
}
