package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"
)

type (
	testCasesStruct struct {
		Cases []caseStruct
	}
	caseStruct struct {
		Username string
		Date     string
		Meal     string
		State    bool
	}
	response struct {
		HTTPStatusCode int         `json:"http_status_code"` // http response status code
		StatusText     string      `json:"status"`           // user-level status message
		ErrorText      string      `json:"error"`            // application-level error message, for debugging
		Data           interface{} `json:"data"`             // application-level data
	}
	testSuiteFile string
)

var (
	testCases testCasesStruct
	url       string
)

func startServer(configFile string) {
	url = "http://localhost:8083/orders"
	go func() {
		restServer(configFile)
	}()
	time.Sleep(5 * time.Second)
}

func timeElapsed(fn string) func() {
	start := time.Now()
	return func() {
		log.Printf("%s took %v\n", fn, time.Since(start))
	}
}

func (suite testSuiteFile) readTestCases() (testCases testCasesStruct) {
	rawData, err := ioutil.ReadFile(string(suite))
	if err != nil {
		log.Fatal(err)
	}
	if err = json.Unmarshal(rawData, &testCases); err != nil {
		log.Fatal(err)
	}
	return
}

func SendPostRequest(cases caseStruct) (*http.Response, error) {
	jsontest, err := json.Marshal(cases)
	if err != nil {
		panic(err)
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsontest))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}
func SendDeleteRequest(cases caseStruct) (*http.Response, error) {
	jsontest, err := json.Marshal(cases)
	if err != nil {
		panic(err)
	}
	request, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsontest))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func TestMain(t *testing.T) {
	configFile = "config/test.json"
	startServer(configFile)
}

func TestSuitesFile(t *testing.T) {
	var suiteFile testSuiteFile
	suiteFile = "test_suite.json"
	testCases = suiteFile.readTestCases()
}

func TestNonExistConfig(t *testing.T) {
	if os.Getenv("GOTEST") == "1" {
		restServer("./noconfig.json")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestNonExistConfig")
	cmd.Env = append(os.Environ(), "GOTEST=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		fmt.Println(e)
		return
	}
	t.Fatalf("Process ran with error %v, want exit status 1", err)
}

func TestRead(t *testing.T) {
	resp, err := http.Get("http://localhost:8083/orders")
	if err != nil {
		t.Errorf(err.Error())
	}
	if resp.StatusCode != 200 {
		t.Errorf("Got status code %v instead of 200", resp.StatusCode)
	}
}

func TestPostResponses(t *testing.T) {
	for id, val := range testCases.Cases {
		t.Run(fmt.Sprintf("case_%v", id+1), testPostResponses(val))
	}
}

func testPostResponses(cases caseStruct) func(*testing.T) {
	return func(t *testing.T) {
		resp, err := SendPostRequest(cases)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		var data response
		decoder := json.NewDecoder(resp.Body)
		if err = decoder.Decode(&data); err != nil {
			t.Fatal(err)
		}
		if !cases.State {
			if resp.StatusCode != 400 {
				t.Errorf("Got status code %v instead of 400", resp.StatusCode)
			}
			return
		}
		if resp.StatusCode != 200 {
			t.Errorf("Got status code %v instead of 200", resp.StatusCode)
		}
		if data.StatusText != "OK" {
			t.Error(data.ErrorText)
		}
	}
}

func TestDeleteResponses(t *testing.T) {
	for id, val := range testCases.Cases {
		t.Run(fmt.Sprintf("case_%v", id+1), testDeleteResponses(val))
	}
}

func testDeleteResponses(cases caseStruct) func(*testing.T) {
	return func(t *testing.T) {
		resp, err := SendDeleteRequest(cases)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		var data response
		decoder := json.NewDecoder(resp.Body)
		if err = decoder.Decode(&data); err != nil {
			t.Fatal(err)
		}
		if !cases.State {
			if resp.StatusCode != 400 {
				t.Errorf("Got status code %v instead of 400", resp.StatusCode)
			}
			return
		}
		if resp.StatusCode != 200 {
			t.Errorf("Got status code %v instead of 200", resp.StatusCode)
		}
		if data.StatusText != "OK" {
			t.Error(data.ErrorText)
		}
	}
}

func TestPostExistingResponse(t *testing.T) {
	for id, val := range testCases.Cases {
		t.Run(fmt.Sprintf("case_%v", id+1), testPostExistingResponse(val))
	}
}

func testPostExistingResponse(cases caseStruct) func(*testing.T) {
	return func(t *testing.T) {
		resp, err := SendPostRequest(cases)
		if err != nil {
			t.Fatal(err)
		}
		resp, err = SendPostRequest(cases)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		var data response
		decoder := json.NewDecoder(resp.Body)
		if !cases.State {
			if resp.StatusCode != 400 {
				t.Errorf("Got status code %v instead of 400", resp.StatusCode)
			}
			return
		}
		if err = decoder.Decode(&data); err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != 400 {
			t.Errorf("Got status code %v instead of 400", resp.StatusCode)
		}
		if data.StatusText != "Bad request" {
			t.Error(data.ErrorText)
		}
		if data.ErrorText != "Order for this date exist" {
			t.Errorf("Showing %v instead of 'Order does not exists'", data.ErrorText)
		}
	}
}

func TestDeleteExistingResponse(t *testing.T) {
	for id, val := range testCases.Cases {
		t.Run(fmt.Sprintf("case_%v", id+1), testDeleteExistingResponse(val))
	}
}

func testDeleteExistingResponse(cases caseStruct) func(*testing.T) {
	return func(t *testing.T) {
		resp, err := SendDeleteRequest(cases)
		if err != nil {
			t.Fatal(err)
		}
		resp, err = SendDeleteRequest(cases)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		var data response
		decoder := json.NewDecoder(resp.Body)
		if err = decoder.Decode(&data); err != nil {
			t.Fatal(err)
		}
		if !cases.State {
			if resp.StatusCode != 400 {
				t.Errorf("Got status code %v instead of 400", resp.StatusCode)
			}
			return
		}
		if resp.StatusCode != 400 {
			t.Errorf("Got status code %v instead of 400", resp.StatusCode)
		}
		if data.StatusText != "Bad request" {
			t.Error(data.ErrorText)
		}
		if data.ErrorText != "Order does not exists" {
			t.Errorf("Showing %v instead of 'Order does not exists'", data.ErrorText)
		}
	}
}

func BenchmarkRead(b *testing.B) {
	resp, err := http.Get("http://localhost:8083/orders")
	if err != nil {
		b.Errorf(err.Error())
	}
	defer resp.Body.Close()
}

func BenchmarkPostResponses(b *testing.B) {
	for id, val := range testCases.Cases {
		b.Run(fmt.Sprintf("case_%v", id+1), benchmarkPostResponses(val))
	}
}

func benchmarkPostResponses(cases caseStruct) func(*testing.B) {
	return func(b *testing.B) {
		resp, err := SendPostRequest(cases)
		if err != nil {
			b.Fatal(err)
		}
		defer resp.Body.Close()
	}
}

func BenchmarkDeleteResponses(b *testing.B) {
	for id, val := range testCases.Cases {
		b.Run(fmt.Sprintf("case_%v", id+1), benchmarkDeleteResponses(val))
	}
}

func benchmarkDeleteResponses(cases caseStruct) func(*testing.B) {
	return func(b *testing.B) {
		resp, err := SendDeleteRequest(cases)
		if err != nil {
			b.Fatal(err)
		}
		defer resp.Body.Close()
	}
}

func BenchmarkPostExistingResponse(b *testing.B) {
	for id, val := range testCases.Cases {
		b.Run(fmt.Sprintf("case_%v", id+1), benchmarkPostExistingResponse(val))
	}
}

func benchmarkPostExistingResponse(cases caseStruct) func(*testing.B) {
	return func(b *testing.B) {
		resp, err := SendPostRequest(cases)
		if err != nil {
			b.Fatal(err)
		}
		resp, err = SendPostRequest(cases)
		if err != nil {
			b.Fatal(err)
		}
		defer resp.Body.Close()
	}
}

func BenchmarkDeleteExistingResponse(b *testing.B) {
	for id, val := range testCases.Cases {
		b.Run(fmt.Sprintf("case_%v", id+1), benchmarkDeleteExistingResponse(val))
	}
}

func benchmarkDeleteExistingResponse(cases caseStruct) func(*testing.B) {
	return func(b *testing.B) {
		resp, err := SendDeleteRequest(cases)
		if err != nil {
			b.Fatal(err)
		}
		resp, err = SendDeleteRequest(cases)
		if err != nil {
			b.Fatal(err)
		}
		defer resp.Body.Close()
	}
}
