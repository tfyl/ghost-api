package ghost_api

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type ApiStruct struct {
	Success     bool `json:"success"`
	Error       error
	Endpoint    string     `json:"endpoint"`
	Validator   Validator  `json:"validator"`
	FormInfo    string     `json:"form_info"`
	Proxy       string     `json:"proxy"`
	UserAgent   string     `json:"user_agent"`
	CookieValue string     `json:"abck_value"`
	PixelValue  string	   `json:"pixel_value"`
	BmSz        string     `json:"bm_sz"`
}


type Validator struct {
	LengthCheck int `json:"length_check"`
	CharCheck   string `json:"char_check"`
}

type Client struct {
	Endpoint  string
	AuthToken string
	http.Client
}

type Getter struct {
	*Client
	C           chan ApiStruct
	FormInfo    string
	ApiEndpoint string
	Validator   Validator
}

func MakeClient(endpoint, token string) (client *Client) {
	return &Client{
		Endpoint:  endpoint,
		AuthToken: token,
		Client:    http.Client{},
	}
}

func MakeGetter(endpoint, token, forminfo, apiEndpoint string, validator Validator) (client *Getter) {
	return &Getter{
		Client: &Client{
			Endpoint:  endpoint,
			AuthToken: token,
			Client:    http.Client{},
		},
		C:           make(chan ApiStruct, 99999),
		FormInfo:    forminfo,
		ApiEndpoint: apiEndpoint,
		Validator:   validator,
	}
}

func (C *Client) Set(form, api string, validator Validator) *Getter {
	return &Getter{
		Client:      C,
		FormInfo:    form,
		ApiEndpoint: api,
		C:           make(chan ApiStruct, 99999),
		Validator:   validator,
	}
}

func (g *Getter) Get(proxy string) {

	data := ApiStruct{
		Endpoint:    g.ApiEndpoint,
		FormInfo:    g.FormInfo,
		Proxy:       proxy,
		UserAgent:   "",
		CookieValue: "",
		Validator:   g.Validator,
	}

	byteBuffer, err := json.Marshal(data)

	if err != nil {
		g.C <- ApiStruct{Success: false, Error: err}
		return
	}
	body := bytes.NewBuffer(byteBuffer)
	req, err := http.NewRequest("POST", g.Endpoint, body)
	if err != nil {
		g.C <- ApiStruct{Success: false, Error: err}
		return
	}

	req.Header.Add("Auth", g.AuthToken)

	do, err := g.Do(req)
	if err != nil {
		g.C <- ApiStruct{Success: false, Error: err}
		return
	}

	err = json.NewDecoder(do.Body).Decode(&data)
	if err != nil {
		g.C <- ApiStruct{Success: false, Error: err}
		return
	}
	g.C <- data

}
