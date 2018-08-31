package model

import (
	"github.com/globalsign/mgo/bson"
	"testing"
)

const (
	CheckMark = "\u2713"
	BallotX   = "\u2717"
)

func TestConvertWithDataReturnsGetAllTagResponse(t *testing.T) {
	id := bson.NewObjectId()
	tagDAOList := []TagDAO{{id, "user", "org", "tag text", "blue"}}
	t.Logf("Given a tagDAO array ")
	{
		response := Convert(tagDAOList)
		if response.Tags != nil {
			t.Logf("\t\tThe GetAllResponse is correctly converted:  \"%s\". %v", response, CheckMark)
		} else {
			t.Errorf("\t\tThe GetAllResponse.Tags shouldn't be nil: \"%s\". %v", response, BallotX)
		}
	}
}

func TestConvertEmptyList(t *testing.T) {
	t.Logf("Given a tagDAO array ")
	{
		response := Convert([]TagDAO{})
		if len(response.Tags) == 0 {
			t.Logf("\t\tThe GetAllResponse is correctly converted:  \"%s\". %v", response, CheckMark)
		} else {
			t.Errorf("\t\tThe GetAllResponse.Tags shouldn't be nil: \"%s\". %v", response, BallotX)
		}
	}
}

func TestConvertToTag(t *testing.T) {
	t.Logf("Given a tagDAO")
	{
		id := bson.NewObjectId()
		expectedType := Tag{id.Hex(), "user", "org", "tag text", "blue"}
		response := ConvertToTag(TagDAO{id, "user","org",  "tag text", "blue"})
		if response == expectedType {
			t.Logf("\t\tThe tag converted matches with the tagDAO:  \"%s\". %v", expectedType, CheckMark)
		} else {
			t.Errorf("\t\t The tagDAO and tag converted should match: \"%s\". %v", response, BallotX)
		}
	}
}
