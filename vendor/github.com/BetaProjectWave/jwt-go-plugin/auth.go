package jwt

import (
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	jwt2 "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	// AuthHeaderEmptyError thrown when an empty Authorization header is received
	AuthHeaderEmptyError = errors.New("auth header empty")

	// InvalidAuthHeaderError thrown when an invalid Authorization header is received
	InvalidAuthHeaderError = errors.New("invalid auth header")

	// ErrInvalidSigningAlgorithm for invalid signing algorithm
	ErrInvalidSigningAlgorithm = errors.New("invalid signing algorithm")
)

const (

	// AccountIDClaim used to retrieve the account id from the claim
	AccountIDClaim = "http://wso2.org/claims/enduser"

	// AccountIDField header to set on the request
	AccountIDField = "accountId"

	// OrganisationIDHeader header to retrieve from request
	OrganisationIDHeader = "Organisation-ID"

	// AuthenticateHeader the Gin authenticate header
	AuthenticateHeader = "WWW-Authenticate"

	// Application token claim Id
	TokenTypeClaimId = "http://wso2.org/claims/usertype"

	// Application token value
	UserTokenType = "APPLICATION_USER"

	// WS02TokenHeader the WS02 JWT token header name
	WS02TokenHeader = "X-JWT-Assertion"
)

// WaveJWTMiddleware middleware
type WaveJWTMiddleware struct {
	// Callback function that should perform the authorization of the authenticated user. Called
	// only after an authentication success. Must return true on success, false on failure.
	// Optional, default to success.
	TokenValidator func(c *gin.Context) bool

	// User can define own Unauthorized func.
	Unauthorized func(*gin.Context, int, string)

	Timeout time.Duration

	TokenLookup string

	TimeFunc func() time.Time

	// Realm name to display to the user. Required.
	Realm string

	// signing algorithm - possible values are HS256, HS384, HS512
	// Optional, default is HS256.
	SigningAlgorithm string

	// Secret key used for signing. Required.
	Key []byte

	// Issuer
	Issuer string
}

// ErrorResponse to be returned to the client
type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:code`
}

// MiddlewareInit initialize jwt configs.
func (mw *WaveJWTMiddleware) MiddlewareInit() {

	if mw.TokenLookup == "" {
		mw.TokenLookup = "header:" + WS02TokenHeader
	}

	if mw.Timeout == 0 {
		mw.Timeout = time.Hour
	}

	if mw.TimeFunc == nil {
		mw.TimeFunc = time.Now
	}

	if mw.Unauthorized == nil {
		mw.Unauthorized = func(c *gin.Context, code int, message string) {
			c.JSON(code, ErrorResponse{Code: code, Message: message})
		}
	}

	if mw.SigningAlgorithm == "" {
		mw.SigningAlgorithm = "HS256"
	}

	if mw.Realm == "" {
		mw.Realm = "gin jwt"
	}

}

func (mw *WaveJWTMiddleware) parseToken(c *gin.Context) (*jwt2.Token, error) {
	var token string
	var err error

	parts := strings.Split(mw.TokenLookup, ":")
	switch parts[0] {
	case "header":
		token, err = mw.jwtFromHeader(c, parts[1])
	}

	if err != nil {
		return nil, err
	}

	return jwt2.Parse(token, func(t *jwt2.Token) (interface{}, error) {
		if jwt2.GetSigningMethod(mw.SigningAlgorithm) != t.Method {
			return nil, ErrInvalidSigningAlgorithm
		}

		// save token string if vaild
		c.Set("JWT_TOKEN", token)

		return mw.Key, nil
	})
}

func (mw *WaveJWTMiddleware) middlewareImpl(c *gin.Context) {

	// Parse the given token
	token, err := mw.parseToken(c)

	if err != nil {
		mw.unauthorized(c, http.StatusUnauthorized, err.Error())
		return
	}

	claims := token.Claims.(jwt2.MapClaims)

	// Verify the issuer
	if ok := claims.VerifyIssuer(mw.Issuer, true); !ok {
		Error.Println("The JWT token issuer is not valid")
		mw.unauthorized(c, http.StatusUnauthorized, "The JWT token issuer is not valid")
		return
	}

	// Assert OrganisationId and accountId only for user Token Type
	if isUserToken(claims) {

		if organisationID := c.Request.Header.Get(OrganisationIDHeader); organisationID == "" {
			Warning.Println("No Organisation-ID header param found. Authorization failed")
			mw.unauthorized(c, http.StatusUnauthorized, "Authorization failed")
			return
		}

		// check accountID
		accountID := extractAccountID(claims)
		if accountID == "" {
			Error.Println("No AccountID claim found. Authorization failed")
			mw.unauthorized(c, http.StatusUnauthorized, "Authorization failed")
			return
		}
		c.Request.Header.Add(AccountIDField, accountID)
	}
	c.Next()
}

func (mw *WaveJWTMiddleware) jwtFromHeader(c *gin.Context, key string) (string, error) {
	authHeader := c.Request.Header.Get(key)

	if authHeader == "" {
		return "", AuthHeaderEmptyError
	}
	return authHeader, nil
}

func extractAccountID(claims jwt2.MapClaims) string {
	id, _ := claims[AccountIDClaim].(string)
	return id
}

func (mw *WaveJWTMiddleware) unauthorized(c *gin.Context, code int, message string) {
	if mw.Realm == "" {
		mw.Realm = "gin jwt"
	}

	c.Header(AuthenticateHeader, "JWT realm="+mw.Realm)
	c.Abort()

	mw.Unauthorized(c, code, message)
	return
}

// MiddlewareFunc makes GinJWTMiddleware implement the Middleware interface.
func (mw *WaveJWTMiddleware) MiddlewareFunc() gin.HandlerFunc {
	// initialise
	mw.MiddlewareInit()
	return func(c *gin.Context) {
		mw.middlewareImpl(c)
		return
	}
}

// GinJWTMiddleware create an instance of WaveJWTMiddleware
func GinJWTMiddleware(secret string, issuer string) *WaveJWTMiddleware {

	key, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		log.Fatal("failed to decode signing string")
	}
	authMiddleware := &WaveJWTMiddleware{
		Timeout: time.Hour,

		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, ErrorResponse{Code: code, Message: message})
		},

		// Token header
		TokenLookup:      "header:" + WS02TokenHeader,
		TimeFunc:         time.Now,
		Key:              key,
		SigningAlgorithm: "HS256",
		Issuer:           issuer,
	}
	return authMiddleware
}

func isUserToken(claims jwt2.MapClaims) bool {
	tokenType, _ := claims[TokenTypeClaimId].(string)
	return tokenType == UserTokenType
}
