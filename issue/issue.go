package issue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
)

const (
	vertexAiDomainUrl string = "https://us-central1-aiplatform.googleapis.com/v1/projects"
	projectId         string = "463895878287"
	modelId           string = "7547958208682262528"

	// Double check this works!!
	BearerToken string = "ya29.a0AbVbY6O2_ma5JUdVMZ4INglb-3J0vvDZipkvTawbBu1lDj3e1S548TbMhyTWaaIeeHGOh0Yqgw66vmbdNRneun__SfRiRVjoJctSAQ1Y_H-NcRXJ9XBS35F8vgldDlyyI-6los7Yrb_oeezQiWMayRzChU-Ag8_5m_a32FEaCgYKAV0SARESFQFWKvPlDmiGC8GNsrQECUWk08nnJA0174"
)

func vertexAIUrlConstructor() string {
	// "https://us-central1-aiplatform.googleapis.com/v1/projects/463895878287/locations/us-central1/endpoints/7547958208682262528:predict"

	// return fmt.Sprintf("%s/locations/us-central1/endpoints/%s:predict", vertexAiDomainUrl, projectId, modelId)
	return vertexAiDomainUrl + projectId + "/locations/us-central1/endpoints/" + modelId + ":predict"
}
func SubmitIssueHandler(w http.ResponseWriter, r *http.Request) {

	var issueRequest IssueRequest

	err := json.NewDecoder(r.Body).Decode(&issueRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("%+v", issueRequest)

	// TODO call the vertex ai model here!!
	posturl := vertexAIUrlConstructor()
	postBody := VertexAIRequest{}
	postBody.Instances = []Instance{
		{
			Content: `given a user submitted issue or a suggestion, classify the issue as the following: Critical High Low Medium

			input: Some optional features are not working
			output: Low`,
		},
	}
	postBody.Parameters = Params{
		Temperature:     0.2,
		MaxOutputTokens: 256,
		TopP:            0.8,
		TopK:            40.0,
	}

	bearer := "Bearer " + BearerToken
	postBodyJson, err := json.Marshal(postBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error:" + err.Error()))
		return
	}
	fmt.Println("postUrl::", posturl)
	vRequest, err := http.NewRequest("POST", posturl, bytes.NewBuffer(postBodyJson))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error:" + err.Error()))
	}
	vRequest.Header.Add("Authorization", bearer)

	// http dump request
	dump, err := httputil.DumpRequest(vRequest, true)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	fmt.Printf("\nVertexAI HTTP POST:\n %s\n\n", string(dump))

	client := &http.Client{}
	resp, err := client.Do(vRequest)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}

	//dump response
	dump, err = httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
	fmt.Printf("\nVertexAI POST Resp: \n%q\n\n", dump)

	//  Classification Sender
	response := ClassificationResponse{}
	response.Severity = "Critical"
	b, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(b)

}
