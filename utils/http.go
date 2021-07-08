package utils

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type setRequest func(*http.Request)

func DoHTTP(method string, url string, body io.Reader, setFunc setRequest, out interface{}) error {

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	if setFunc != nil {
		setFunc(request)
	}

	tr := &http.Transport{ //x509: certificate signed by unknown authority
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: 15 * time.Second}

	resp, err := client.Do(request)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	defer client.CloseIdleConnections()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("read body failed:%s", err)
		return err
	}

	if out != nil {
		if s, ok := out.(*string); ok {
			*s = string(buf)
		} else {
			if err := json.Unmarshal(buf, obj); err != nil {
				log.Printf("unmarshal failed:%s", err)
				return err
			}
		}
	}

	return nil

}
