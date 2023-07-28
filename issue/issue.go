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
	vertexAiDomainUrl string = "https://us-central1-aiplatform.googleapis.com/v1/projects/"
	projectId         string = "league-stage-datalake-play"
	modelId           string = "text-bison@001"

	// Double check this works!!
	BearerToken string = "ya29.a0AbVbY6McEVFfIqBJubc3IE71zbkqk4v1nBmitVD8OXoTMhsNJVXGkB5xxiomkoEZNHVkzcqEPHyKl7Fj2cKg_eXKHv0wIGlh58mh1p_fsT3lrhb4K9Xy8GYaAceNBFVKRuqdPLSCVt8i1pS4_lTpoo5s2Z0zo1a8nRU5j3kaCgYKAYQSARESFQFWKvPlJcmFmbmiU2v0nBmRTCumIQ0174"
)

func vertexAIUrlConstructor() string {
	// return vertexAiDomainUrl + projectId + "/locations/us-central1/endpoints/publishers/google/models/" + modelId + ":predict"
	return "https://us-central1-aiplatform.googleapis.com/v1/projects/league-stage-datalake-play/locations/us-central1/publishers/google/models/text-bison@001:predict"
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
	contentBody := `given a user submitted issue or a suggestion, classify the issue as the following: Critical High Low Medium` +
		`input: Hi Team!
	Good day.
	Asking for your assistance to determine why the member received their League Invitation Email at a later date than expected. The member got hired last June 12th 2023 but was only sent an invite July 20th.
	Might you be able to provide a better insight into this as the member is concerned about submitting their claims.
	Thank you!` +
		` output: Low` +
		`input: The app is running slowly after my update` +
		` output: Low` +
		`input: my prescriptions are no longer on the page` +
		` output: Medium` +
		`input: I can't see any info on my wallet` +
		` output: Critical` +
		`input: my date of hire is wrong` +
		` output: High` +
		`input: ` + issueRequest.Issue +
		` output:`
	postBody.Instances = []Instance{
		{
			Content: contentBody,
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
	fmt.Println("\n\npostUrl::", posturl)
	fmt.Println("postBody::", string(postBodyJson))
	vRequest, err := http.NewRequest("POST", posturl, bytes.NewBuffer(postBodyJson))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error:" + err.Error()))
	}
	vRequest.Header.Add("Authorization", bearer)
	vRequest.Header.Add("Content-Type", "application/json")

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
	// dump, err = httputil.DumpResponse(resp, true)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Println("body from vertexAI" + string(body))
	// fmt.Printf("\nVertexAI POST Resp: \n%q\n\n", dump)

	vAIResponse := VertexAIResponse{}
	json.Unmarshal(body, &vAIResponse)
	response := ClassificationResponse{}
	//  Classification Sender
	if len(vAIResponse.Predictions) == 0 {
		// Unknown means we havent got a classification from VertexAI and we need a human in the loop.
		response.Severity = "Unknown"

	} else {
		response.Severity = vAIResponse.Predictions[0].Content
	}

	// response := ClassificationResponse{
	// 	// This assumes that there is always one prediction.
	// 	// TODO: We may wanna clean this to pass an unknown severity to get human eyes on this.
	// 	Severity: vAIResponse.Predictions[0].Content,
	// }
	// response.Severity = "Critical"
	b, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(b)

}
