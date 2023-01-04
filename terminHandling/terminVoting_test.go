package terminHandling

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTerminVotingHandler(t *testing.T) {
	setup()
	defer after()
	request, _ := http.NewRequest(http.MethodPost, "/terminVoting", nil)
	response := httptest.NewRecorder()
	http.HandlerFunc(TerminVotingHandler).ServeHTTP(response, request)
	body, _ := io.ReadAll(response.Result().Body)
	fmt.Println(string(body))

	//form = url.Values{}
	//// User zu Terminfindung einladen-Button
	//form.Add("inviteUserSubmit", "Terminfindung1")
	//form.Add("username", "Anna")
	//request.PostForm = form
}
