package api

import (
	"encoding/json"
	"github.com/globalsign/mgo/bson"
	"github.com/tag-service/model"
	"github.com/tag-service/test"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Query all tags, it should return two tags
func TestQueryAllSuccess(t *testing.T) {

	t.Logf("Given I create two tags")
	{
		controller := NewTagHandler(Repository)
		router := controller.CreateRouter()
		//
		tag1 := model.CreateTagRequest{Name: "Dinner", Colour: "Red"}
		tag2 := model.CreateTagRequest{Name: "Flight", Colour: "Black"}
		id1 := test.CreateTag(tag1, router, t, test.Token2, test.OrgID2)
		id2 := test.CreateTag(tag2, router, t, test.Token2, test.OrgID2)
		expectedTags := map[string]model.Tag{id1: model.Tag{Id: id1, Name: tag1.Name, Colour: tag1.Colour, AccountId: test.AccountID2, OrganisationId: test.OrgID2},
			id2: model.Tag{Id: id2, Name: tag2.Name, Colour: tag2.Colour, AccountId: test.AccountID2, OrganisationId: test.OrgID2}}

		t.Logf("\tWhen Sending Get All tags request to endpoint:  \"%s\"", "\\tags")
		{
			req, err := test.HttpRequest(nil, "/tags", http.MethodGet, test.Token2, test.OrgID2)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			test.Ok(err, t)

			// Assert response code status
			test.CheckStatus(w, t, http.StatusOK)

			var response model.GetAllTagResponse
			json.NewDecoder(w.Body).Decode(&response)

			t.Log(response)
			 //check both tags can be retrieved
			checkRetrievedTag(expectedTags, id1, response, t)
			checkRetrievedTag(expectedTags, id2, response, t)
		}
	}
}

// Query all tags, it should return none.
func TestQueryAllNoTagFound(t *testing.T) {
	t.Logf("Given no tags were created for the given user")
	{
		controller := NewTagHandler(Repository)
		router := controller.CreateRouter()

		t.Logf("\tWhen Sending Get All tags request to endpoint:  \"%s\"", "\\tags")
		{
			req, _ := test.HttpRequest(nil, "/tags", http.MethodGet, test.Token3, test.OrgID1)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert response code status
			test.CheckStatus(w, t, http.StatusOK)

			var response model.GetAllTagResponse
			json.NewDecoder(w.Body).Decode(&response)

			// check body response empty
			if len(response.Tags) == 0 {
				t.Logf("\t\t \"%d\" tags have been found. %v", len(response.Tags), test.CheckMark)
			} else {
				t.Errorf("\t\t \"%d\" tags have been found. %v", len(response.Tags), test.BallotX)
			}
		}
	}
}

// Query for given tag id it should return the actual tag
func TestQueryForGivenTagIdSuccess(t *testing.T) {

	t.Logf("Given I create a tag")
	{
		controller := NewTagHandler(Repository)
		router := controller.CreateRouter()

		tag := model.CreateTagRequest{Name: "Dinner", Colour: "Red"}
		id := test.CreateTag(tag, router, t, test.Token2, test.OrgID2)

		t.Logf("\tWhen Sending Get tag request to endpoint:  \"%s\"", "\\tags\\"+id)
		{
			req, err := test.HttpRequest(nil, "/tags/"+id, http.MethodGet, test.Token2, test.OrgID2)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			test.Ok(err, t)

			// Assert response code status
			test.CheckStatus(w, t, http.StatusOK)

			var response model.Tag
			json.NewDecoder(w.Body).Decode(&response)

			// accountId for token two
			accountId := "5a3922ac86da0c1d779a776"
			expectedResponse := model.Tag{Id: id, Name: tag.Name, Colour: tag.Colour, AccountId: accountId, OrganisationId: test.OrgID2}

			// check body response
			if expectedResponse == response {
				t.Logf("\t\tTag should have been found:  \"%s\". %v", expectedResponse, test.CheckMark)
			} else {
				t.Errorf("\t\tTag should have been found:  \"%s\". %v", response, test.BallotX)
			}
		}
	}
}

// Creates two tag and query them.
func TestQueryForGivenTagIdNotFound(t *testing.T) {

	t.Logf("Given no tag is created")
	{
		controller := NewTagHandler(Repository)
		router := controller.CreateRouter()
		id := bson.NewObjectId().Hex()
		t.Logf("\t\ttWhen Sending Get tag request to endpoint:  \"%s\"", "\\tags\\"+id)
		{
			req, err := test.HttpRequest(nil, "/tags/"+id, http.MethodGet, test.Token2, test.OrgID2)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			test.Ok(err, t)

			// Assert response code status
			test.CheckStatus(w, t, http.StatusNotFound)

		}
	}
}

// helper function
func checkRetrievedTag(expectedTags map[string]model.Tag, id string, response model.GetAllTagResponse, t *testing.T) {
	if expectedTags[id] == getCreatedTags(id, response.Tags) {
		t.Logf("\t\tTag [\"%s\"] has been retrieved successfully:  \"%s\". %v", id, getCreatedTags(id, response.Tags), test.CheckMark)
	} else {
		t.Logf("\t\tTag [\"%s\"] has been retrieved successfully:  \"%s\". %v", id, getCreatedTags(id, response.Tags), test.BallotX)
	}
}

// helper function
func getCreatedTags(id string, tags []model.Tag) model.Tag {
	for _, tag := range tags {
		if tag.Id == id {
			return tag
		}
	}
	return model.Tag{}
}
