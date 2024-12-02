package unit

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	domain "Medods/auth_service/inner_layer/domain"
	"Medods/auth_service/rest/adapter"
	"Medods/auth_service/rest/route"

	service "Medods/auth_service/inner_layer/service/auth"

	"Medods/auth_service/rest/controller"
	domainErrors "github.com/santaasus/errors-handler"
)

type mockRepository struct {
}

type requestConfig struct {
	method  string
	url     string
	reqBody any
	header  map[string]string
}

func (mockRepository) GetUserByGuid(guid string) (*domain.User, error) {
	user := &domain.User{
		ID:   1,
		Guid: guid,
		Hash: "$2a$10$Tc7JS8JeVDHChmD/quVfWOIC2PpYgC9nM2Ji6WLTLzLvVptEszeIO",
	}

	return user, nil
}

func (mockRepository) CreateUser(newUser *domain.NewUser) (*domain.User, error) {
	user := &domain.User{
		ID:   1,
		Guid: newUser.Guid,
		Hash: newUser.Hash,
		IP:   newUser.IP,
	}

	return user, nil
}

func (mockRepository) UpdateUser(params map[string]any, userId int) error {
	if len(params) == 0 {
		return &domainErrors.AppError{
			Err:  errors.New("UpdateUser: params are empty"),
			Type: domainErrors.ValidationError,
		}
	}

	return nil
}

func (mockRepository) DeleteUserByHash(hash string) error {
	return nil
}

// Change to the root path for correct reading files like os.Readfile
func changeWd() {
	path, _ := os.Getwd()
	rootPath := strings.Split(path, "/Medods")[0]
	if len(rootPath) > 0 {
		os.Chdir(rootPath + "/Medods")
	}
}

func getAdapter() *adapter.BaseAdapter {
	return &adapter.BaseAdapter{
		Repository: mockRepository{},
	}
}

func (c *requestConfig) doRequest(gin *gin.Engine) (*httptest.ResponseRecorder, error) {
	jsonValue, err := json.Marshal(c.reqBody)
	if err != nil {
		return nil, err
	}

	newRequest, err := http.NewRequest(c.method, c.url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}

	if len(c.header) > 0 {
		for key, value := range c.header {
			newRequest.Header.Add(key, value)
		}
	}

	writer := httptest.NewRecorder()
	gin.ServeHTTP(writer, newRequest)

	return writer, nil
}

func TestGetTokens(t *testing.T) {
	changeWd()

	engine := gin.Default()

	baseAdapter := getAdapter()

	controller := baseAdapter.AuthAdapter()

	group := engine.Group(route.AUTH_GROUP)
	{
		group.GET(route.TOKENS_PATH, controller.GetTokens)
	}

	reqConfig := requestConfig{
		method:  "GET",
		url:     route.AUTH_GROUP + route.TOKENS_PATH + "?guid=745a8c08-9483-4a5a-b9ba-69ebc2204d17",
		reqBody: nil,
	}

	writer, err := reqConfig.doRequest(engine)
	if err != nil {
		t.Error(err)
		return
	}

	// Check response status
	assert.Equal(t, http.StatusOK, writer.Code)

	var expected service.SecurityData
	err = json.Unmarshal(writer.Body.Bytes(), &expected)
	if err != nil {
		t.Errorf("body len: %v, error for unmarshal: %v", len(writer.Body.Bytes()), err)
		return
	}

	assert.NotEmpty(t, expected.JWTAccessToken)
}

func TestRefreshTokens(t *testing.T) {
	changeWd()

	engine := gin.Default()

	baseAdapter := getAdapter()

	group := engine.Group(route.AUTH_GROUP)
	{
		group.POST(route.REFRESH_TOKENS_PATH, baseAdapter.AuthAdapter().RefreshToken)
	}

	req := &controller.AccessTokenRequest{
		RefreshToken: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpZCI6Ijc0NWE4YzA4LTk0ODMtNGE1YS1iOWJhLTY5ZWJjMjIwNGQxNyIsInR5cGUiOiJyZWZyZXNoIiwicGF5bG9hZCI6eyJpcCI6Ijo6MSJ9LCJleHAiOjE3MzQ0Mjg5MDd9.NlqX3d8pmqH0t7dZODdFthxnch264WeDxu9w7L7ACv1VhGP1G9VjfNSP87bhSDNnbUet-jqNrhBHhcsy4rln1w",
	}

	reqConfig := requestConfig{
		method:  "POST",
		url:     route.AUTH_GROUP + route.REFRESH_TOKENS_PATH,
		reqBody: req,
	}

	writer, err := reqConfig.doRequest(engine)
	if err != nil {
		t.Error(err)
		return
	}

	// Check response status
	assert.Equal(t, http.StatusOK, writer.Code)

	var expected service.SecurityData
	err = json.Unmarshal(writer.Body.Bytes(), &expected)
	if err != nil {
		t.Errorf("body len: %v, error for unmarshal: %v", len(writer.Body.Bytes()), err)
		return
	}

	assert.NotEmpty(t, expected.JWTAccessToken)
}
