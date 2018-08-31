package repository

import (
	"github.com/BetaProjectWave/kube-vault-plugin"
	"github.com/globalsign/mgo/bson"
	"github.com/tag-service/model"
	"github.com/tag-service/test"
	"testing"
	"os"
)

const AccountId = "48590485"

func TestMongoRepository_Insert(t *testing.T) {

	t.Logf("Given the tag service is up and running")
	{
		t.Logf("\tWhen Sending Create TagDAO request to endpoint:  \"%s\"", "\\tags")
		{
			tagId := bson.NewObjectId()
			tag := model.TagDAO{Id: tagId, Name: "Lunch", Colour: "Red", AccountId: AccountId}
			err := RepositoryUnderTest.Insert("tags-db", "tags", tag)
			if err == nil {
				t.Logf("\t\tThe insert should have been successful %v", test.CheckMark)
			} else {
				t.Errorf("\t\tThe insert should have been successful %v", test.BallotX)
			}
		}
	}
}

func TestMongoRepository_Delete(t *testing.T) {
	t.Logf("Given the tag service is up and running")
	{
		t.Logf("\tWhen Sending Delete TagDAO request to endpoint:  \"%s\"", "\\tags\\id")
		{
			tagId := CreateTag(t)
			err := RepositoryUnderTest.Delete("tags-db", "tags", tagId)
			if err == nil {
				t.Logf("\t\tThe delete should have been successful %v", test.CheckMark)
			} else {
				t.Errorf("\t\tThe delete should have been successful %v", test.BallotX)
			}
		}

	}
}

func TestMongoRepository_FindAll(t *testing.T) {
	t.Logf("Given the tag service is up and running")
	{
		t.Logf("\tWhen Sending Find all TagDAO request to endpoint:  \"%s\"", "\\tags")
		{
			CreateTag(t)
			_, err := RepositoryUnderTest.FindAll("tags-db", "tags", bson.M{})
			if err == nil {
				t.Logf("\t\tThe find all should have been successful %v", test.CheckMark)
			} else {
				t.Errorf("\t\tThe find all should have been successful %v", test.BallotX)
			}
		}

	}
}

func TestMongoRepository_Find(t *testing.T) {
	t.Logf("Given the tag service is up and running")
	{
		t.Logf("\tWhen Sending Find TagDAO request to endpoint:  \"%s\"", "\\tags\\id")
		{
			tagId := CreateTag(t)
			_, err := RepositoryUnderTest.Find("tags-db", "tags", tagId)
			if err == nil {
				t.Logf("\t\tThe find should have been successful %v", test.CheckMark)
			} else {
				t.Errorf("\t\tThe find should have been successful %v", test.BallotX)
			}
		}

	}
}

func TestNewRepository(t *testing.T) {
	config := vault.Config{Address: "https://domain.com"}
	os.Setenv(ENVIRONMENT, "dev")
	go NewRepository(config)
}


// Helper to create a Tag
func CreateTag(t *testing.T) bson.ObjectId {
	tagId := bson.NewObjectId()
	tag := model.TagDAO{Id: tagId, Name: "Lunch", Colour: "Red", AccountId: AccountId}
	err := RepositoryUnderTest.Insert("tags-db", "tags", tag)
	if err == nil {
		t.Logf("\t\tThe insert should have been successful %v", test.CheckMark)
	} else {
		t.Errorf("\t\tThe insert should have been successful %v", test.BallotX)
	}
	return tagId
}
