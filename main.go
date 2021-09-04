package ghost_api

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type ApiStruct struct {
	Success     bool   `json:"success"`
	Error       error
	Endpoint    string `json:"endpoint"`
	FormInfo    string `json:"form_info"`
	Proxy       string `json:"proxy"`
	UserAgent   string `json:"user_agent"`
	CookieValue string `json:"cookie"`
}


type Client struct {
	Endpoint string
	AuthToken string
	http.Client
}

type Getter struct {
	*Client
	FormInfo string
	ApiEndpoint string
}

func MakeClient(endpoint,token string) (client *Client)  {
	return &Client{
		Endpoint:  endpoint,
		AuthToken: token,
		Client:    http.Client{},
	}
}

func MakeGetter(endpoint,token,forminfo,apiEndpoint string) (client *Getter)  {
	return &Getter{
		Client: &Client{
			Endpoint:  endpoint,
			AuthToken: token,
			Client:    http.Client{},
		},
		FormInfo:    forminfo,
		ApiEndpoint: apiEndpoint,
	}
}

func (C *Client) Set(form,api string) *Getter {
	return &Getter{
		Client: C,
		FormInfo:    form,
		ApiEndpoint: api,
	}
}

func (g *Getter) Get(proxy string) chan ApiStruct{
	c := make(chan ApiStruct,1)
	go func() {
		data := ApiStruct{
			Endpoint:    g.ApiEndpoint,
			FormInfo:    g.FormInfo,
			Proxy:       proxy,
			UserAgent:   "",
			CookieValue: "",
		}


		byteBuffer , err := json.Marshal(data)
		if err != nil{
			c <- ApiStruct{Success: false,Error: err}
			return
		}
		body := bytes.NewBuffer(byteBuffer)
		req, err := http.NewRequest("POST",g.Endpoint,body)
		if err != nil{
			c <- ApiStruct{Success: false,Error: err}
			return
		}
		req.Header.Add("Auth",g.AuthToken)


		do, err := g.Do(req)
		if err != nil{
			c <- ApiStruct{Success: false,Error: err}
			return
		}

		err = json.NewDecoder(do.Body).Decode(&data)
		if err != nil{
			c <- ApiStruct{Success: false,Error: err}
			return
		}
		c <- data


	}()
	return c
}