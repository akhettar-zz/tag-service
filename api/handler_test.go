package api

import (
	"encoding/json"
	"errors"
	"github.com/globalsign/mgo/bson"
	"github.com/golang/mock/gomock"
	"github.com/tag-service/mocks"
	"github.com/tag-service/model"
	"github.com/tag-service/test"
	"net/http"
	"net/http/httptest"
	"testing"
)

/**
Unit tests around the negative edges of testing against the DB failures.
*/
func TestTagHandler_CreateTag_Failure_to_Insert(t *testing.T) {

	t.Logf("Given the tag service is up and running")
	{
		t.Logf("\tWhen Sending Create TagDAO request to endpoint:  \"%s\"", "\\tags")
		{
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockRepo := mocks.NewMockRepository(mockCtrl)

			expectedErrorMessage := "Insert failed"
			body := model.CreateTagRequest{Name: "Dinner", Colour: "Red"}
			err := errors.New(expectedErrorMessage)

			mockRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(err).Times(1)

			handler := NewTagHandler(mockRepo)
			router := handler.CreateRouter()

			req, err := test.HttpRequest(body, "/tags", http.MethodPost, test.Token1, test.OrgID1)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check call success
			test.Ok(err, t)

			// Assert response code status
			test.CheckStatus(w, t, http.StatusInternalServerError)

			var response model.ErrorResponse
			json.NewDecoder(w.Body).Decode(&response)

			expectedResponse := model.ErrorResponse{Code: http.StatusInternalServerError, Message: expectedErrorMessage}

			// check body response matches the expected response
			test.CheckResponseMessage(response, expectedResponse, t, w)
		}
	}

}

/**
Unit tests around the negative edges of testing against the DB failures.
*/
func TestTagHandler_DeleteTag_Failure_User_id_mismatch(t *testing.T) {
	t.Logf("Given the tag service is up and running")
	{
		t.Logf("\tWhen Sending Delete TagDAO request with dummy accountId to the endpoint:  \"%s\"", "\\tags")
		{
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockRepo := mocks.NewMockRepository(mockCtrl)

			expectedErrorMessage := "The given tag id does not belong to the organisation"
			err := errors.New(expectedErrorMessage)
			body := model.CreateTagRequest{Name: "Dinner", Colour: "Red"}

			tag := model.TagDAO{Id: bson.NewObjectId(), Name: body.Name, Colour: body.Colour, OrganisationId: test.OrgID1}

			mockRepo.EXPECT().Find(gomock.Any(), gomock.Any(), gomock.Any()).Return(tag, nil).Times(1)

			handler := NewTagHandler(mockRepo)
			router := handler.CreateRouter()

			req, err := test.HttpRequest(nil, "/tags/"+bson.NewObjectId().Hex(), http.MethodDelete, test.Token1, test.OrgID2)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check call success
			test.Ok(err, t)

			// Assert response code status
			test.CheckStatus(w, t, http.StatusConflict)

			var response model.ErrorResponse
			json.NewDecoder(w.Body).Decode(&response)

			expectedResponse := model.ErrorResponse{Code: http.StatusConflict, Message: expectedErrorMessage}

			// check body response matches the expected response
			test.CheckResponseMessage(response, expectedResponse, t, w)
		}
	}
}

func TestTagHandler_DeleteTag_Tag_not_found(t *testing.T) {
	t.Logf("Given the tag service is up and running")
	{
		t.Logf("\tWhen Sending Delete TagDAO request to endpoint:  \"%s\"", "\\tags")
		{
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockRepo := mocks.NewMockRepository(mockCtrl)

			expectedErrorMessage := "Tag not found"
			body := model.CreateTagRequest{Name: "Dinner", Colour: "Red"}
			err := errors.New(expectedErrorMessage)

			tag := model.TagDAO{Id: bson.NewObjectId(), Name: body.Name, Colour: body.Colour, AccountId: test.AccountID2, OrganisationId: test.OrgID1}

			findCall := mockRepo.EXPECT().Find(gomock.Any(), gomock.Any(), gomock.Any()).Return(tag, nil).Times(1)
			mockRepo.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(err).Times(1).After(findCall)

			handler := NewTagHandler(mockRepo)
			router := handler.CreateRouter()

			req, err := test.HttpRequest(nil, "/tags/"+bson.NewObjectId().Hex(), http.MethodDelete, test.Token2, test.OrgID1)
			t.Logf("\tThen I shsould get an error response stating that tag not found  \"%s\"", "\\tags")
			{
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				// check call success
				test.Ok(err, t)

				// Assert response code status
				test.CheckStatus(w, t, http.StatusNotFound)

			}
		}
	}
}

func TestTagHandler_GetAllTags_Failed_to_Fetch_Due_to_DB_error(t *testing.T) {

	t.Logf("Given Tag service is up and running")
	{
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockRepo := mocks.NewMockRepository(mockCtrl)
		controller := NewTagHandler(mockRepo)
		router := controller.CreateRouter()
		expectedErrorMessage := "Failed to retrieve data from the database"
		err := errors.New(expectedErrorMessage)

		// set mockrepo expectations
		mockRepo.EXPECT().FindAll(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, err).Times(1)

		t.Logf("\tWhen Sending Get All tags request to endpoint:  \"%s\"", "\\tags")
		{
			req, err := test.HttpRequest(nil, "/tags", http.MethodGet, test.Token2, test.OrgID1)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check call success
			test.Ok(err, t)

			// Assert response code status
			test.CheckStatus(w, t, http.StatusInternalServerError)

			var response model.ErrorResponse
			json.NewDecoder(w.Body).Decode(&response)

			expectedResponse := model.ErrorResponse{Code: http.StatusInternalServerError, Message: expectedErrorMessage}

			// check body response matches the expected response
			test.CheckResponseMessage(response, expectedResponse, t, w)
		}
	}
}
