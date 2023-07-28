package httpd

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	host    string
	cookies []*http.Cookie

	method string
	url    string
	body   *bytes.Buffer

	status  int
	jsonRsp *json.Decoder
	strRsp  string
}

func NewClient(host string) *Client {
	return &Client{host: host}
}

func (c *Client) Get(path string, params map[string]string) error {
	query := url.Values{}
	for k, v := range params {
		query.Add(k, v)
	}

	c.method = "GET"
	c.url = c.host + path + "?" + query.Encode()
	c.body = nil

	return c.do()
}

func (c *Client) Post(path string, body any) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	c.method = "POST"
	c.url = c.host + path
	c.body = bytes.NewBuffer(b)

	return c.do()
}

func (c *Client) do() error {
	req, err := http.NewRequest(c.method, c.url, c.body)
	if err != nil {
		return err
	}

	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	if c.method == "POST" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	c.status = resp.StatusCode
	if c.status == http.StatusOK {
		c.cookies = resp.Cookies()
	}

	if resp.Header.Get("Content-Type") == "application/json" {
		c.jsonRsp = json.NewDecoder(resp.Body)
	} else if resp.Header.Get("Content-Type") == "text/plain" {
		body, _ := ioutil.ReadAll(resp.Body)
		c.strRsp = string(body)
	}

	return nil
}
