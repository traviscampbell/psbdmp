package psbdmp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Dump holds the id, date added, and the associated tags in some cases of a given dump
type Dump struct {
	ID   string `json:"id"`
	Tags string `json:"tags"`
	Time string `json:"time"`
}

// Dumplings represents server response containing a list dump IDs
type Dumplings struct {
	Search    string `json:"search"`
	Count     int    `json:"count"`
	Dumps     []Dump `json:"data"`
	Error     int    `json:"error"`
	ErrorInfo string `json:"error_info"`
}

// DumpClient facilitates interaction with the psbdmp service
type DumpClient struct {
	c       *http.Client
	ua      string
	baseURL *url.URL
}

// NewDumpClient returns a dumpclient instance initialized with sane defaults
func NewDumpClient() *DumpClient {
	bURL, _ := url.Parse("https://psbdmp.ws")
	
	return &DumpClient{
		c:       &http.Client{Timeout: 9 * time.Second},
		ua:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36",
		baseURL: bURL,
	}
}

// SetHTTPClient allows a custom HTTP Client to be used
func (dc *DumpClient) SetHTTPClient(client *http.Client) *DumpClient {
	dc.c = client
	return dc
}

// SetUserAgent allows a custom User-Agent to be applied
func (dc *DumpClient) SetUserAgent(useragent string) *DumpClient {
	dc.ua = useragent
	return dc
}

// GetByDate returns a list of the dumps that we're posted between the given dates
func (dc *DumpClient) GetByDate(from, to time.Time) ([]Dump, error) {
	u, _ := dc.baseURL.Parse("api/dump/getbydate")

	format := "02.01.2006"
	payload := fmt.Sprintf("from=%s&to=%s", from.Format(format), to.Format(format))

	d := Dumplings{}
	if err := dc.pitterPatter("POST", u.String(), bytes.NewBufferString(payload), &d); err != nil {
		return nil, err
	}

	if d.Error != 0 {
		return nil, errors.New(d.ErrorInfo)
	}

	return d.Dumps, nil
}

// Search performs a generalized search across the dumps for the provided key word
func (dc *DumpClient) Search(keyword string) ([]Dump, error) {
	u, _ := dc.baseURL.Parse("api/search/" + keyword)

	d := Dumplings{}
	if err := dc.pitterPatter("GET", u.String(), nil, &d); err != nil {
		return nil, err
	}

	if d.Error != 0 {
		return nil, errors.New(d.ErrorInfo)
	}

	return d.Dumps, nil
}

// SearchByDomain searches for dumps containing the provided domain
func (dc *DumpClient) SearchByDomain(domain string) ([]Dump, error) {
	u, _ := dc.baseURL.Parse("api/search/domain/" + domain)

	d := Dumplings{}
	if err := dc.pitterPatter("GET", u.String(), nil, &d); err != nil {
		return nil, err
	}

	if d.Error != 0 {
		return nil, errors.New(d.ErrorInfo)
	}

	return d.Dumps, nil
}

// SearchByEmail searches for dumps containing the provided email
func (dc *DumpClient) SearchByEmail(email string) ([]Dump, error) {
	u, _ := dc.baseURL.Parse("api/search/email/" + email)

	d := Dumplings{}
	if err := dc.pitterPatter("GET", u.String(), nil, &d); err != nil {
		return nil, err
	}

	if d.Error != 0 {
		return nil, errors.New(d.ErrorInfo)
	}

	return d.Dumps, nil
}

// GetDumpContent returns the full dump as a string
func (dc *DumpClient) GetDumpContent(id string) (string, error) {
	var d struct {
		ID        string `json:"id"`
		Data      string `json:"data"`
		Time      string `json:"time"`
		Error     int    `json:"error"`
		ErrorInfo string `json:"error_info"`
	}

	u, _ := dc.baseURL.Parse("api/dump/get/" + id)

	if err := dc.pitterPatter("GET", u.String(), nil, &d); err != nil {
		return "", err
	}

	if d.Error != 0 {
		return "", errors.New(d.ErrorInfo)
	}

	return d.Data, nil
}

// pitterPatter lets get at 'er. Makes the request, closes the body, and unmarshals the response
// into the provided result object which is expected to be given a pointer value.
func (dc *DumpClient) pitterPatter(meth, u string, body io.Reader, result interface{}) error {
	req, err := http.NewRequest(meth, u, body)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", dc.ua)
	req.Header.Set("Accept", "application/json")
	
	resp, err := dc.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf := bytes.Buffer{}
	if _, err := io.Copy(&buf, resp.Body); err != nil {
		return err
	}

	err = json.Unmarshal(buf.Bytes(), result)
	return err
}
