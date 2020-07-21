package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	utils "testTask"
	"testTask/model"
	"testing"
)

type Request struct {
	jsonBody model.Request
	headers  map[string]string
}

type RequestResponce struct {
	request  Request
	responce map[string]string
	status   int
}

var SimpleTest = RequestResponce{
	request: Request{
		jsonBody: model.Request{"win", "10.15", utils.RandString(10)},
		headers:  map[string]string{"Source-Type": "server"}},
	responce: map[string]string{"Message": "Transaction Complete"},
}

//Clear DB, make simple transaction, check result
func TestSinpleBalanceTransaction(t *testing.T) {
	r := SimpleTest.request
	http.Get("http://localhost:8080/nullbalance")
	reqBody, err := json.Marshal(r.jsonBody)
	if err != nil {
		fmt.Println(err)
	}
	client := http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:8080/transaction", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Source-Type", "server")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	assert.Equal(t, 200, resp.StatusCode)
	response := make(map[string]string, 1)
	respBody, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		fmt.Println(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, SimpleTest.responce, response)
	w, err := http.Get("http://localhost:8080/checkbalance")
	if err != nil {
		fmt.Println(err)
	}
	assert.Equal(t, http.StatusOK, w.StatusCode)
	var responseFromBalance map[string]float64
	respBalance, err := ioutil.ReadAll(w.Body)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(respBalance, &responseFromBalance)
	assert.Nil(t, err)
	assert.Equal(t, 10.15, responseFromBalance["balance"])
}

//Different fail cases
var TestCases = []RequestResponce{
	{request: Request{
		jsonBody: model.Request{"lost", "qwerty", utils.RandString(10)},
		headers:  map[string]string{"Source-Type": "server"}},
		responce: map[string]string{"Message": "Invalid Amount"},
		status:   400},
	{request: Request{
		jsonBody: model.Request{"wrong state here", "124", utils.RandString(10)},
		headers:  map[string]string{"Source-Type": "server"}},
		responce: map[string]string{"Message": "Wrong Transaction Status"},
		status:   400},
	{request: Request{
		jsonBody: model.Request{"win", "130", utils.RandString(10)},
		headers:  map[string]string{"Source-Type": "wrong source"}},
		responce: map[string]string{"Message": "Wrong Source"},
		status:   400},
	{request: Request{
		jsonBody: model.Request{"lost", "-100", utils.RandString(10)},
		headers:  map[string]string{"Source-Type": "server"}},
		responce: map[string]string{"Message": "Negative Amount"},
		status:   400},
}

func TestCasesProsess(t *testing.T) {
	for _, testCase := range TestCases {
		reqBody, err := json.Marshal(testCase.request.jsonBody)
		req, err := http.NewRequest("POST", "http://localhost:8080/transaction", bytes.NewBuffer(reqBody))
		req.Header.Set("Source-Type", testCase.request.headers["Source-Type"])
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		response := make(map[string]string, 1)
		respBody, err := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(respBody, &response)
		if err != nil {
			fmt.Println(err)
		}
		assert.Nil(t, err)
		assert.Equal(t, response, testCase.responce)
		assert.Equal(t, resp.StatusCode, testCase.status)
	}
}
