package api

import (
	"github.com/BetaProjectWave/jwt-go-plugin"
	gintrace "github.com/DataDog/dd-trace-go/contrib/gin-gonic/gin"
	"github.com/DataDog/dd-trace-go/tracer"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/globalsign/mgo/bson"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "github.com/tag-service/docs"
	"github.com/tag-service/logger"
	"github.com/tag-service/model"
	"github.com/tag-service/repository"
	cfg "github.com/tag-service/vault"
	"net/http"
	"strings"
)

const (
	ServiceName              = "tag-service"
	DatabaseName             = "tag-db"
	DatabaseCollection       = "tags"
	JSONMimeType             = "application/json; charset=utf-8"
	ContentType              = "Content-Type"
	AccountIDField           = "accountId"
	OrganisationIDField      = "Organisation-ID"
	OrganisationId           = "organisationId"
	AccountId                = "accountId"
	TagId                    = "id"
	DataDogAgentHostEnv      = "DD_AGENT_HOST"
	DataDogAgentHostFallback = "localhost"
	DataDogServiceNameEnv    = "DD_SERVICE_NAME"
	DataDogTracerPort        = "8126"
)

type TagHandler struct {
	repo repository.Repository
}

func NewTagHandler(repo repository.Repository) *TagHandler {
	return &TagHandler{repo}
}

// @Summary Creates new tag
// @ID create-tag
// @Description Creates new tag with metadata send
// @Accept  json
// @Produce  json
// @Param new-tag body model.CreateTagRequest true "New tag"
// @Success 201 {object} model.CreateTagResponse "Tag created"
// @Failure 400 {object} model.ErrorResponse "Bad request"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /tags [post]
func (handler *TagHandler) CreateTag(c *gin.Context) {
	var req model.CreateTagRequest

	// This will infer what binder to use depending on the content-type header.
	if errB := c.ShouldBindWith(&req, binding.JSON); errB != nil {
		logger.Error.Println(errB.Error())
		setErrorResponse("Failed to parse Json request", http.StatusBadRequest, c)
		return
	}

	// return bad request if tag is empty
	if tagNameEmptyString(req) {
		setErrorResponse("tag name may not be empty", http.StatusBadRequest, c)
		return
	}

	accountId := c.Request.Header.Get(AccountIDField)
	organisationId := c.Request.Header.Get(OrganisationIDField)
	logger.Info.Printf("Received request to create tag with name \"%v\" for accountId \"%v\" and organisationId \"%v", req.Name, accountId, organisationId)

	tag := model.TagDAO{Id: bson.NewObjectId(), Name: req.Name, Colour: req.Colour, AccountId: accountId, OrganisationId: organisationId}
	logger.Info.Printf("Tag \"%v\" successfully created", tag.Id.Hex())
	err := handler.repo.Insert(DatabaseName, DatabaseCollection, &tag)
	if err != nil {
		logger.Error.Println(err.Error())
		setErrorResponse("Insert failed", http.StatusInternalServerError, c)
		return
	}

	// if all good create success response
	c.Writer.Header().Set(ContentType, JSONMimeType)
	c.JSON(http.StatusCreated, model.CreateTagResponse{Id: tag.Id.Hex()})
}

// @Summary Get tags
// @ID get-tags
// @Accept  json
// @Produce  json
// @Success 200 {object} model.GetAllTagResponse	"ok"
// @Failure 400 {object} model.ErrorResponse "Bad request"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /tags [get]
func (handler *TagHandler) GetAllTags(c *gin.Context) {
	accountId := c.Request.Header.Get(AccountIDField)
	organisationId := c.Request.Header.Get(OrganisationIDField)
	logger.Info.Printf("Received retrieve all tags request for accountId \"%s\" and organisationId \"%v", accountId, organisationId)
	results, err := handler.repo.FindAll(DatabaseName, DatabaseCollection, bson.M{AccountId: accountId, OrganisationId: organisationId})
	if err != nil {
		logger.Error.Println("Failed to retrieve data from the database")
		setErrorResponse("Failed to retrieve data from the database", http.StatusInternalServerError, c)
		return
	}
	c.JSON(http.StatusOK, model.Convert(results))
}

// @Summary Get tag by ID
// @ID get-tag
// @Accept  json
// @Produce  json
// @Param id path string true "Tag ID"
// @Success 200 {object} model.Tag "ok"
// @Failure 404 {object} model.EmptyBody "Tag not found"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /tags/{id} [get]
func (handler *TagHandler) GetTag(c *gin.Context) {
	accountId := c.Request.Header.Get(AccountIDField)
	id := c.Params.ByName("id")
	organisationId := c.Request.Header.Get(OrganisationIDField)
	logger.Info.Printf("Received request to retrieve tag \"%v\" from accountId \"%v\" for organisationId \"%v", id, accountId, organisationId)
	oid := bson.ObjectIdHex(id)
	result, err := handler.repo.Find(DatabaseName, DatabaseCollection, oid)
	if err != nil {
		c.JSON(http.StatusNotFound, model.EmptyBody{})
		return
	}
	c.JSON(http.StatusOK, model.ConvertToTag(result))
}

// @Summary Delete tag by ID
// @ID delete-tag
// @Accept  json
// @Produce  json
// @Param id path string true "Tag ID"
// @Success 204 "Tag deleted"
// @Failure 404 {object} model.ErrorResponse "Tag not found"
// @Failure 409 {object} model.ErrorResponse "The given tag id does not belong to the user"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /tags/{id} [delete]
func (handler *TagHandler) DeleteTag(c *gin.Context) {

	accountId := c.Request.Header.Get(AccountIDField)
	organisationId := c.Request.Header.Get(OrganisationIDField)

	id := c.Params.ByName(TagId)
	logger.Info.Printf("Received request to delete tag for given accountId \"%v\", organisationId \"%v\" and tagId \"%v\"", accountId, organisationId, id)
	oid := bson.ObjectIdHex(id)

	// query the tag
	result, errQ := handler.repo.Find(DatabaseName, DatabaseCollection, oid)
	if errQ != nil {
		c.JSON(http.StatusNotFound, model.EmptyBody{})
		return
	}

	// the tag does not belong to the organisation.
	if result.OrganisationId != organisationId {
		logger.Error.Println("The given tag id does not belong to the organisation")
		setErrorResponse("The given tag id does not belong to the organisation", http.StatusConflict, c)
		return
	}

	// remove the tag
	err := handler.repo.Delete(DatabaseName, DatabaseCollection, oid)
	if err != nil {
		c.JSON(http.StatusNotFound, model.EmptyBody{})
		return
	}
	logger.Info.Printf("Tag successfully deleted \"%v\"", id)
	c.Status(http.StatusNoContent)
}

// @Summary Health endpoint
// @ID health
// @Success 200 {object} model.EmptyBody "ok"
// @Failure 500 {object} model.EmptyBody "Server is down"
// @Router /health [get]
func (handler *TagHandler) Health(c *gin.Context) {
	c.String(http.StatusOK, "Success")
}

// Registers all the routes
func (handler *TagHandler) CreateRouter() *gin.Engine {

	// Start DataDog tracer
	t := tracer.NewTracerTransport(tracer.NewTransport(cfg.GetEnv(DataDogAgentHostEnv, DataDogAgentHostFallback), DataDogTracerPort))
	defer t.ForceFlush()

	// Create router
	router := gin.New()
	router.Use(logger.Logger())
	router.Use(gin.Recovery())
	router.Use(gintrace.MiddlewareTracer(cfg.GetEnv(DataDogServiceNameEnv, ServiceName), t))

	router.GET("/tags", jwt.GinJWTMiddleware().MiddlewareFunc(), handler.GetAllTags)
	router.GET("/tags/:id", jwt.GinJWTMiddleware().MiddlewareFunc(), handler.GetTag)
	router.DELETE("/tags/:id", jwt.GinJWTMiddleware().MiddlewareFunc(), handler.DeleteTag)
	router.GET("/health", handler.Health)
	router.POST("/tags", jwt.GinJWTMiddleware().MiddlewareFunc(), handler.CreateTag)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

// Helper method
func setErrorResponse(msg string, status int, c *gin.Context) {
	c.JSON(status, model.ErrorResponse{Message: msg, Code: status})
}

// checks if the given payload contain tag name with an empty string
func tagNameEmptyString(req model.CreateTagRequest) bool {
	return len(strings.TrimSpace(req.Name)) == 0
}
