package issue

import "net/http"

func SubmitIssueHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("submitted"))

}
