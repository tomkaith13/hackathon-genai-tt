package issue

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func SubmitIssueHandler(w http.ResponseWriter, r *http.Request) {

	var issueRequest IssueRequest

	err := json.NewDecoder(r.Body).Decode(&issueRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("%+v", issueRequest)

	// TODO call the vertex ai model here!!

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
