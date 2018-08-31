package test

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/tag-service/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	// AccountID1 embedded in Token1
	AccountID1 = "5a3938ac86da0c1d5549a776"

	// AccountID2 embedded in Token2
	AccountID2 = "5a3922ac86da0c1d779a776"

	// OrgID1 Organisation Id1 used in integration test.
	OrgID1 = "494010c0-5e76-11e8-9c2d-fa7ae01bbebc"

	// OrgID2 Organisation Id2 used in integration test.
	OrgID2 = "49401340-5e76-11e8-9c2d-fa7ae01bbebc"

	// Token1 used in integration test.
	Token1 = "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6Ik5ERXlOekkzUVVZelFUWXhSVGcwUVRRNE5EUTRPVE5HTXpJMVFVTTVPRVZHTlVVM05qVTJSUSJ9.eyJodHRwOi8vYXN0by5jby51ay91c2VySWQiOiI1YTM5MzhhYzg2ZGEwYzFkNTU0OWE3NzYiLCJuaWNrbmFtZSI6Iis0NDc4Mzc4NzYyMjUiLCJuYW1lIjoiKzQ0NzgzNzg3NjIyNSIsInBpY3R1cmUiOiJodHRwczovL2Nkbi5hdXRoMC5jb20vYXZhdGFycy80LnBuZyIsInVwZGF0ZWRfYXQiOiIyMDE4LTAxLTI0VDE0OjI0OjE5LjQ3MloiLCJpc3MiOiJodHRwczovL2FzdG8tZGV2LmV1LmF1dGgwLmNvbS8iLCJzdWIiOiJzbXN8NWEzOTM4YWM4NmRhMGMxZDU1NDlhNzc2IiwiYXVkIjoiajl0dDFDcUowTm5DTTdhZGVoV2kzQTJ1c3JvcVFwZnMiLCJpYXQiOjE1MTY4MDM4NTksImV4cCI6MTUxOTM5NTg1OX0.RtEhjO0GXzR-vQ--ciIVdPW4Jao2rPWvi2gCcF6K7PHn2slubBbXZhbmRZBaJEwpYvKW06-bpgwFqvZTHg5MXNKvbaLyA6jV9jGHbDXRki7BTcCT75tNqWy-x7Hv8eNYJu8Y9-RDSs1OgK4VOd2CpnQvWwSGZIbP_N8vg8Rj-AVo08ih6L2Dg11q4URV8RTpx_J9YE9jxJgP9u-DzyzYMxsfMBPx4zu1OL8o16g1hvSFGiqrgf3-NjjsakW645fHjxQqjiJvT0XsjBy2EmL19pCTjGrpauYf1mgwmixOBuu1drNIf233LxnEOGCdDBh3vaJnqBZCiWWsYe2pJJThRA"

	// Token2 used in integration test
	Token2 = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJodHRwOi8vYXN0by5jby51ay91c2VySWQiOiI1YTM5MjJhYzg2ZGEwYzFkNzc5YTc3NiIsIm5pY2tuYW1lIjoiKzQ0NzgzNzg3NjIyNSIsIm5hbWUiOiIrNDQ3ODM3ODc2MjI1IiwicGljdHVyZSI6Imh0dHBzOi8vY2RuLmF1dGgwLmNvbS9hdmF0YXJzLzQucG5nIiwidXBkYXRlZF9hdCI6IjIwMTgtMDEtMjRUMTQ6MjQ6MTkuNDcyWiIsImlzcyI6Imh0dHBzOi8vYXN0by1kZXYuZXUuYXV0aDAuY29tLyIsInN1YiI6InNtc3w1YTM5MzhhYzg2ZGEwYzFkNTU0OWE3NzYiLCJhdWQiOiJqOXR0MUNxSjBObkNNN2FkZWhXaTNBMnVzcm9xUXBmcyIsImlhdCI6MTUxNjgwMzg1OSwiZXhwIjoxNTE5Mzk1ODU5fQ.3EbaQN6zy7GS0s0-rAUfjKjKYwG8fvLqn8YzcyEZKAY"

	// Token3 used in integration test
	Token3 = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJodHRwOi8vYXN0by5jby51ay91c2VySWQiOiIxRDMzNDRhYzg2ZGEwYzFkNzc5YTg4MSIsIm5pY2tuYW1lIjoiKzQ0NzgzNzg3NjIyNSIsIm5hbWUiOiIrNDQ3ODM3ODc2MjI1IiwicGljdHVyZSI6Imh0dHBzOi8vY2RuLmF1dGgwLmNvbS9hdmF0YXJzLzQucG5nIiwidXBkYXRlZF9hdCI6IjIwMTgtMDEtMjRUMTQ6MjQ6MTkuNDcyWiIsImlzcyI6Imh0dHBzOi8vYXN0by1kZXYuZXUuYXV0aDAuY29tLyIsInN1YiI6InNtc3w1YTM5MzhhYzg2ZGEwYzFkNTU0OWE3NzYiLCJhdWQiOiJqOXR0MUNxSjBObkNNN2FkZWhXaTNBMnVzcm9xUXBmcyIsImlhdCI6MTUxNjgwMzg1OSwiZXhwIjoxNTE5Mzk1ODU5fQ.CWtn73aOMEGS29rAwR57oArBRQURl6AaV60q3TKc-Jc"

	// InvalidToken - accountId not present in the token
	InvalidToken = "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.EkN-DOsnsuRjRO6BxXemmJDm3HbxrbRzXglbN2S4sOkopdU4IsDxTI8jO19W_A4K8ZPJijNLis4EZsHeY559a4DFOd50_OqgHGuERTqYZyuhtF39yxJPAjUESwxk2J5k_4zM3O-vtd1Ghyo4IbqKKSy6J9mTniYJPenn5-HIirE"

	// DummyToken to test token parsing error
	DummyToken = "DymmyToken"

	// CheckMark used for unit test highlight.
	CheckMark = "\u2713"

	// BallotX used for unit test highlight.
	BallotX = "\u2717"
)

// RequestBody helper
func RequestBody(req interface{}) *bytes.Buffer {
	jsonBytes, err := json.Marshal(req)
	if err != nil {
		panic("Failed to marshall json request")
	}
	return bytes.NewBuffer(jsonBytes)
}

// Ok assert helper
func Ok(err error, t *testing.T) {
	if err != nil {
		t.Fatal("\t\tShould be able to make the Post call.", BallotX, err)
	}
}

// CreateTag helper
func CreateTag(body model.CreateTagRequest, router *gin.Engine, t *testing.T, token string, orgId string) string {
	req, err := HttpRequest(body, "/tags", http.MethodPost, token, orgId)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// check call success
	Ok(err, t)
	var response model.CreateTagResponse
	json.NewDecoder(w.Body).Decode(&response)
	return response.Id
}

// HttpRequest helper
func HttpRequest(jsonReq interface{}, endpoint string, method string, token string, orgId string) (*http.Request, error) {
	req, err := http.NewRequest(method, endpoint, RequestBody(jsonReq))
	if err != nil {
		panic("Failed to marshall json request")
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)
	req.Header.Add("Organisation-ID", orgId)
	return req, err
}

// CheckStatus helper
func CheckStatus(w *httptest.ResponseRecorder, t *testing.T, status int) {
	if w.Code == status {
		t.Logf("\t\tShould receive a \"%d\" status. %v", status, CheckMark)
	} else {
		t.Errorf("\t\tShould receive a \"%d\" status. %v %v", status, BallotX, w.Code)
	}
}

// CheckResponseMessage helper.
func CheckResponseMessage(response model.ErrorResponse, expectedResponse model.ErrorResponse, t *testing.T, w *httptest.ResponseRecorder) {
	if response.Message == expectedResponse.Message {
		t.Logf("\t\tThe body response should  contain a message \"%s\" . %v", expectedResponse.Message, CheckMark)
	} else {
		t.Errorf("\t\tThe body response should contain a message \"%s\". %v %v", response.Message, BallotX, w.Code)
	}
}
