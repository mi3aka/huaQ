package fofa

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Fofa struct {
	email []byte
	key   []byte
	*http.Client
}

type result struct {
	Domain string `json:"domain,omitempty"`
	Host   string `json:"host,omitempty"`
	IP     string `json:"ip,omitempty"`
	Port   string `json:"port,omitempty"`
}

type Results []result

const (
	defaultAPIUrl = "https://fofa.info/api/v1/search/all?"
)

func NewFofaClient(email, key []byte) *Fofa {
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &Fofa{
		email: email,
		key:   key,
		Client: &http.Client{
			Transport: transCfg,
		},
	}
}

func (fofa *Fofa) Get(u string) ([]byte, error) {
	body, err := fofa.Client.Get(u)
	if err != nil {
		return nil, err
	}
	defer body.Body.Close()
	content, err := ioutil.ReadAll(body.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (fofa *Fofa) Query(query []byte, fields []byte, page uint) ([]byte, error) {
	q := []byte(base64.StdEncoding.EncodeToString(query))
	q = bytes.Join([][]byte{[]byte(defaultAPIUrl),
		[]byte("email="), fofa.email,
		[]byte("&key="), fofa.key,
		[]byte("&qbase64="), q,
		[]byte("&fields="), fields,
		[]byte("&page="), []byte(strconv.Itoa(int(page))),
		[]byte("&size=10000"),
	}, []byte(""))
	fmt.Printf("%s\n", q)
	content, err := fofa.Get(string(q))
	if err != nil {
		return nil, err
	}
	errmsg, err := jsonparser.GetString(content, "errmsg")
	if err == nil {
		err = errors.New(errmsg)
	} else {
		err = nil
	}
	return content, err
}

func (fofa *Fofa) WriteFile(content []byte, fields []byte, filename string) error {
	results := gjson.GetBytes(content, "results")
	fieldsArray := strings.Split(string(fields), ",")
	content = []byte("[]")
	for _, result := range results.Array() {
		tmp_json := "{}"
		for index, field := range fieldsArray {
			value := result.Get(strconv.Itoa(index)).String()
			tmp_json, _ = sjson.Set(tmp_json, field, value)
		}
		content, _ = sjson.SetBytes(content, "-1", gjson.Parse(tmp_json).Value())
	}
	err := ioutil.WriteFile(filename, content, 0644)
	if err != nil {
		return err
	}
	return nil
}
