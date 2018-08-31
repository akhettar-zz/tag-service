package api

import (
	"github.com/globalsign/mgo/bson"
	"github.com/tag-service/model"
	"github.com/tag-service/test"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Creates two tag and query them.
func TestDeleteTagSuccess(t *testing.T) {

	t.Logf("Given I create a tag")
	{
		controller := NewTagHandler(Repository)
		router := controller.CreateRouter()

		tag := model.CreateTagRequest{Name: "Dinner", Colour: "Red"}
		id := test.CreateTag(tag, router, t, test.Token2, test.OrgID1)

		t.Logf("\t\ttWhen Sending Delete tag request to endpoint:  \"%s\"", "\\tags\\"+id)
		{
			req, err := test.HttpRequest(nil, "/tags/"+id, http.MethodDelete, test.Token2, test.OrgID1)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			test.Ok(err, t)

			// Assert response code status
			test.CheckStatus(w, t, http.StatusNoContent)
		}
	}
}

// Attempt to delete tag and you get empty response with not found
func TestDeleteTagIdNotFound(t *testing.T) {

	t.Logf("Given no tag is created")
	{
		controller := NewTagHandler(Repository)
		router := controller.CreateRouter()
		randomId := bson.NewObjectId().Hex()
		t.Logf("\t\ttWhen Sending Get tag request to endpoint:  \"%s\"", "\\tags\\"+randomId)
		{
			req, err := test.HttpRequest(nil, "/tags/"+randomId, http.MethodDelete, test.Token2, test.OrgID1)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			test.Ok(err, t)

			// Assert response code status
			test.CheckStatus(w, t, http.StatusNotFound)

		}
	}
}

// Attempt to delete a tag not belonging to the organisation result in conflict error.
func TestDeleteTagNotBelongingToTheUser(t *testing.T) {

	t.Logf("Given I create a tag")
	{
		controller := NewTagHandler(Repository)
		router := controller.CreateRouter()

		tag := model.CreateTagRequest{Name: "Dinner", Colour: "Red"}
		id := test.CreateTag(tag, router, t, test.Token2, test.OrgID1)

		t.Logf("\t\tWhen Sending Delete tag request to endpoint:  \"%s\" with a different organisationId", "\\tags\\"+id)
		{
			req, err := test.HttpRequest(nil, "/tags/"+id, http.MethodDelete, test.Token1, test.OrgID2)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			test.Ok(err, t)

			// Assert response code status
			test.CheckStatus(w, t, http.StatusConflict)
		}
	}
}
