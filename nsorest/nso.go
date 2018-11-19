package nsorest

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
)

type NSO struct {
	http.Client
	Username string
	Password string
	BaseURI  string
}

type Device struct {
	IPAddress net.IP
	Hostname  string
	CPEName   string
}

func (nso *NSO) send(method string, endpoint string, body io.Reader) (error, *http.Response) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", nso.BaseURI, endpoint), body)

	if err != nil {
		return err, nil
	}

	req.Header.Set("Accept", "application/vnd.yang.datastore+json, application/vnd.yang.data+json, application/vnd.yang.api+json, application/vnd.yang.collection+json")
	req.SetBasicAuth(nso.Username, nso.Password)

	if body != nil {
		req.Header.Set("Content-Type", "application/vnd.yang.data+json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err, nil
	}

	return nil, resp

}

func (nso *NSO) get(endpoint string) (error, *http.Response) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	method := "GET"
	req, err := http.NewRequest(
		method,
		fmt.Sprintf("%s/%s", nso.BaseURI, endpoint),
		nil,
	)

	if err != nil {
		return err, nil
	}

	req.Header.Set("Accept", "application/vnd.yang.datastore+json, application/vnd.yang.data+json, application/vnd.yang.api+json, application/vnd.yang.collection+json")
	req.SetBasicAuth(nso.Username, nso.Password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err, nil
	}

	return nil, resp

}
