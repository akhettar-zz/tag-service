package api

import (
	"encoding/json"
	"github.com/tag-service/model"
	"github.com/tag-service/test"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Create TagDAO endpoint
func TestCreateTagSuccess(t *testing.T) {

	t.Logf("Given the tag service is up and running")
	{
		t.Logf("\tWhen Sending Create TagDAO request to endpoint:  \"%s\"", "\\tags")
		{
			handler := NewTagHandler(Repository)
			router := handler.CreateRouter()

			body := model.CreateTagRequest{Name: "Dinner", Colour: "Red"}
			req, err := test.HttpRequest(body, "/tags", http.MethodPost, test.Token1, test.OrgID1)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			// check call success
			test.Ok(err, t)

			if w.Code == http.StatusCreated {
				t.Logf("\t\tShould receive a \"%d\" status. %v", http.StatusOK, test.CheckMark)
			} else {
				t.Errorf("\t\tShould receive a \"%d\" status. %v %v", http.StatusOK, test.BallotX, w.Code)
			}
			var response model.CreateTagResponse
			json.NewDecoder(w.Body).Decode(&response)
			if response.Id != "" {
				t.Logf("\t\tThe response should contain tag id. %v %v", response.Id, test.CheckMark)
			} else {
				t.Errorf("\t\tThe response should contain tag id. %v %v", response.Id, test.BallotX)
			}
		}
	}
}

func TestCreateRouteWithInvalidCreateTagRequest(t *testing.T) {

	t.Logf("Given the tag service is up and running")
	{
		t.Logf("\tWhen Sending an Invalid Create request to endpoint : \"%s\"", "\\tags")
		{
			handler := NewTagHandler(Repository)
			router := handler.CreateRouter()

			w := httptest.NewRecorder()
			body := model.TagDAO{Name: "Dinner"}
			req, err := test.HttpRequest(body, "/tags", http.MethodPost, test.Token1, test.OrgID1)
			req.Header.Add("Authorization", "DummyToken")
			router.ServeHTTP(w, req)

			// check call success
			test.Ok(err, t)

			// Assert response code status
			test.CheckStatus(w, t, http.StatusBadRequest)

			// assert content-type in the case of an error
			contentType := w.Header().Get(ContentType)
			t.Logf(contentType)
			if contentType == JSONMimeType {
				t.Logf("\t\tThe content type header should be \"%s\". %v", contentType, test.CheckMark)
			} else {
				t.Errorf("\t\tThe content type header should be \"%s\". %v", contentType, test.BallotX)
			}

			var response model.ErrorResponse
			json.NewDecoder(w.Body).Decode(&response)
			expectedResponse := model.ErrorResponse{Code: http.StatusBadRequest, Message: "Failed to parse Json request"}

			// check body response
			test.CheckResponseMessage(response, expectedResponse, t, w)
		}
	}
}


func TestCreateRouteWithInvalidOrganisationId(t *testing.T) {
	t.Logf("Given the tag service is up and running")
	{
		t.Logf("\tWhen Sending Create TagDAO request to endpoint : \"%s\" with an empty OrganisationId", "\\tags")
		{
			controller := NewTagHandler(Repository)
			router := controller.CreateRouter()

			w := httptest.NewRecorder()
			body := model.CreateTagRequest{Name: "Dinner", Colour: "Red"}
			req, err := test.HttpRequest(body, "/tags", http.MethodPost, test.Token1, "")
			router.ServeHTTP(w, req)

			// check call success
			test.Ok(err, t)

			// Assert response code status
			test.CheckStatus(w, t, http.StatusUnauthorized)

			var response model.ErrorResponse
			json.NewDecoder(w.Body).Decode(&response)
			expectedResponse := model.ErrorResponse{Code: http.StatusUnauthorized, Message: "Authorization failed"}

			// check body response
			test.CheckResponseMessage(response, expectedResponse, t, w)
		}
	}
}

func TestCreateRouteWithDummyToken(t *testing.T) {

	t.Logf("Given the tag service is up and running")
	{
		t.Logf("\tWhen Sending Create TagDAO request to endpoint : \"%s\" with a dummy token", "\\tags")
		{
			handler := NewTagHandler(Repository)
			router := handler.CreateRouter()

			w := httptest.NewRecorder()
			body := model.TagDAO{Name: "Dinner", Colour: "Red", AccountId: "48590485"}
			req, err := test.HttpRequest(body, "/tags", http.MethodPost, test.DummyToken, test.OrgID1)
			req.Header.Add("Authorization", "DummyToken")
			router.ServeHTTP(w, req)

			// check call success
			test.Ok(err, t)

			// Assert response code status
			test.CheckStatus(w, t, http.StatusUnauthorized)

			var response model.ErrorResponse
			json.NewDecoder(w.Body).Decode(&response)
			expectedResponse := model.ErrorResponse{Code: http.StatusUnauthorized, Message: "invalid auth header"}

			// check body response
			test.CheckResponseMessage(response, expectedResponse, t, w)
		}
	}
}

func TestCreateRouteWithInvalidToken(t *testing.T) {
	t.Logf("Given the tag service is up and running")
	{
		t.Logf("\tWhen Sending Create TagDAO request to endpoint : \"%s\" with a dummy token", "\\tags")
		{
			controller := NewTagHandler(Repository)
			router := controller.CreateRouter()

			w := httptest.NewRecorder()
			body := model.CreateTagRequest{Name: "Dinner", Colour: "Red"}
			req, err := test.HttpRequest(body, "/tags", http.MethodPost, test.InvalidToken, test.OrgID1)
			router.ServeHTTP(w, req)

			// check call success
			test.Ok(err, t)

			// Assert response code status
			test.CheckStatus(w, t, http.StatusUnauthorized)

			var response model.ErrorResponse
			json.NewDecoder(w.Body).Decode(&response)
			expectedResponse := model.ErrorResponse{Code: http.StatusUnauthorized, Message: "Authorization failed"}

			// check body response
			test.CheckResponseMessage(response, expectedResponse, t, w)
		}
	}
}

func TestCreateTagWithEmptyTagName(t *testing.T) {

	t.Logf("Given the tag service is up and running")
	{
		t.Logf("\tWhen Sending Create request to endpoint : \"%s\" with an empty tag name", "\\tags")
		{
			handler := NewTagHandler(Repository)
			router := handler.CreateRouter()

			body := model.CreateTagRequest{Name: " ", Colour: "Red"}
			req, err := test.HttpRequest(body, "/tags", http.MethodPost, test.Token1, test.OrgID1)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check call success
			test.Ok(err, t)

			// Assert response code status
			test.CheckStatus(w, t, http.StatusBadRequest)

			var response model.ErrorResponse
			json.NewDecoder(w.Body).Decode(&response)
			expectedResponse := model.ErrorResponse{Code: http. StatusBadRequest, Message: "tag name may not be empty"}

			// check body response
			test.CheckResponseMessage(response, expectedResponse, t, w)

		}
	}
}
