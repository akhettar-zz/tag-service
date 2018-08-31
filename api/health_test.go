package api

import (
	"github.com/tag-service/test"
	"net/http"
	"net/http/httptest"
	"testing"
)

/**
  Integration test for the health end point Health Endpoint written in BDD style
*/
func TestHealthRoute(t *testing.T) {

	t.Logf("Given the need to use the health endpoint to query container status")
	{
		t.Logf("\tWhen checking \"%s\" for status code \"%d\"", "\\health", http.StatusOK)
		{
			controller := NewTagHandler(Repository)
			router := controller.CreateRouter()

			w := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/health", nil)
			router.ServeHTTP(w, req)

			if err != nil {
				t.Fatal("\t\tShould be able to make the Get call.",
					test.BallotX, err)
			}
			t.Log("\t\tShould be able to make the Get call.", test.CheckMark)
			if w.Code == http.StatusOK {
				t.Logf("\t\tShould receive a \"%d\" status. %v", http.StatusOK, test.CheckMark)
			} else {
				t.Errorf("\t\tShould receive a \"%d\" status. %v %v", http.StatusOK, test.BallotX, w.Code)
			}
		}
	}

}
