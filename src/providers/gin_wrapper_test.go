package providers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chff7cb/swissbank/providers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ginWrapperTestSuite struct {
	suite.Suite
	ctx             *gin.Context
	recorder        *httptest.ResponseRecorder
	wrapperProvider providers.GinWrapperProvider
	wrapper         providers.GinWrapper
}

type dummyData struct {
	Foo string `json:"foo"`
}

func TestGinWrapperTestSuite(t *testing.T) {
	suite.Run(t, new(ginWrapperTestSuite))
}

func (s *ginWrapperTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	s.wrapperProvider = providers.NewGinWrapperProvider()
	s.recorder = httptest.NewRecorder()
	s.ctx, _ = gin.CreateTestContext(s.recorder)
	s.wrapper = s.wrapperProvider.Wrap(s.ctx)
}

func (s *ginWrapperTestSuite) TestGinWrapperImpl_TestShouldBindJSON() {
	inputData := dummyData{"bar"}

	data, err := json.Marshal(inputData)
	assert.Equal(s.T(), nil, err)

	s.ctx.Request = httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(data))

	outputData := dummyData{}
	err = s.wrapper.ShouldBindJSON(&outputData)

	assert.Equal(s.T(), nil, err)
	assert.Equal(s.T(), inputData.Foo, outputData.Foo)
}

func (s *ginWrapperTestSuite) TestGinWrapperImpl_JSON() {
	outputData := dummyData{"bar"}

	encodedData, err := json.Marshal(outputData)
	assert.Equal(s.T(), nil, err)

	s.wrapper.JSON(http.StatusOK, outputData)

	assert.Equal(s.T(), string(encodedData), s.recorder.Body.String())
}

func (s *ginWrapperTestSuite) TestGinWrapperImpl_Param() {
	requestParam := gin.Param{
		Key:   "foo",
		Value: "bar",
	}
	s.ctx.Params = []gin.Param{requestParam}

	assert.Equal(s.T(), requestParam.Value, s.wrapper.Param(requestParam.Key))
}

func (s *ginWrapperTestSuite) TestGinWrapperImpl_Inner() {
	assert.Equal(s.T(), s.ctx, s.wrapper.Inner())
}
