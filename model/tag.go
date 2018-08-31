package model

import (
	"github.com/globalsign/mgo/bson"
)

// for persistence
type TagDAO struct {
	Id             bson.ObjectId `json:"id" bson:"_id,omitempty"`
	AccountId      string        `json:"accountId" bson:"accountId",omitempty`
	OrganisationId string        `json:"organisationId" bson:"organisationId",omitempty`
	Name           string        `json:"name" bson:"name"`
	Colour         string        `json:"colour" bson: "colour"`
}

type CreateTagResponse struct {
	Id string `json:Id`
}

type CreateTagRequest struct {
	Name   string `json:"name" binding:"required"`
	Colour string `json:"colour" binding:"required"`
}

type GetAllTagResponse struct {
	Tags []Tag `json:tags`
}

type Tag struct {
	Id             string `json:"id"`
	AccountId      string `json:"AccountId" `
	OrganisationId string `json:"organisationId" `
	Name           string `json:"name"`
	Colour         string `json:"colour"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:code`
}

type EmptyBody struct{}

func Convert(tags []TagDAO) GetAllTagResponse {
	response := make([]Tag, 0)
	for _, tag := range tags {
		response = append(response, Tag{Id: tag.Id.Hex(), Name: tag.Name, Colour: tag.Colour, AccountId: tag.AccountId})
	}
	return GetAllTagResponse{response}
}

func ConvertToTag(tag TagDAO) Tag {
	return Tag{Id: tag.Id.Hex(), Name: tag.Name, Colour: tag.Colour, AccountId: tag.AccountId, OrganisationId: tag.OrganisationId}
}
