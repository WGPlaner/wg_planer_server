package integrations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/wgplaner/wg_planer_server/controllers"
	"github.com/wgplaner/wg_planer_server/models"
	"github.com/wgplaner/wg_planer_server/modules/setting"
	"github.com/wgplaner/wg_planer_server/restapi"
	"github.com/wgplaner/wg_planer_server/restapi/operations"

	"github.com/stretchr/testify/assert"
	"gopkg.in/testfixtures.v2"
)

var server *restapi.Server

var (
	AuthValid   = "1234567890fakefirebaseid0001"
	AuthInvalid = "invalid"
	AuthEmpty   = ""
)

// Main Test Function
// Called first for integrations tests
func TestMain(m *testing.M) {
	initIntegrationTest()
	defer server.Shutdown()

	var sqlHelper = &testfixtures.SQLite{}

	err := models.InitFixtures(
		sqlHelper,
		path.Join(setting.AppWorkPath, "models/fixtures/"),
	)

	if err != nil {
		fmt.Printf("Error initializing test database: %v\n", err)
		os.Exit(1)
	}

	// Run the tests
	os.Exit(m.Run())
}

func initIntegrationTest() {
	var api *operations.WgplanerAPI

	if wgPlanerRoot := os.Getenv("WGPLANER_ROOT"); wgPlanerRoot == "" {
		log.Fatalln("Environment variable $WGPLANER_ROOT not set. " +
			"It is required for integration tests!")

	} else {
		// Set path so that config and data directory are found
		setting.AppWorkPath = wgPlanerRoot
		setting.AppPath = path.Join(wgPlanerRoot, "wgplaner")
	}

	setting.NewConfigContext()
	setting.AppConfig.Auth.IgnoreFirebase = true

	api = operations.NewWgplanerAPI(setting.LoadSwaggerSpec(restapi.SwaggerJSON))
	server = restapi.NewServer(api)
	server.Port = setting.AppConfig.Server.Port

	controllers.GlobalInit()
	controllers.InitializeControllers(api)

	// Set handler
	server.SetHandler(api.Serve(nil))
}

func prepareTestEnv(t testing.TB) {
	assert.NoError(t, models.LoadFixtures())
}

type TestResponseWriter struct {
	HeaderCode int
	Writer     io.Writer
	Headers    http.Header
}

func (w *TestResponseWriter) Header() http.Header {
	return w.Headers
}

func (w *TestResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *TestResponseWriter) WriteHeader(n int) {
	w.HeaderCode = n
}

type TestResponse struct {
	HeaderCode int
	Body       []byte
	Headers    http.Header
}

func NewRequest(t testing.TB, method, auth string, urlStr string) *http.Request {
	return NewRequestWithBody(t, method, auth, urlStr, nil)
}

func NewRequestWithJSON(t testing.TB, method, auth string, urlStr string, v interface{}) *http.Request {
	jsonBytes, err := json.Marshal(v)
	assert.NoError(t, err)
	req := NewRequestWithBody(t, method, auth, urlStr, bytes.NewBuffer(jsonBytes))
	req.Header.Add("Content-Type", "application/json")
	return req
}

func NewRequestWithBody(t testing.TB, method, auth string, urlStr string, body io.Reader) *http.Request {
	request, err := http.NewRequest(method, urlStr, body)
	assert.NoError(t, err)
	request.Header.Add("Authorization", auth)
	request.RequestURI = urlStr
	return request
}

const NoExpectedStatus = -1

func MakeRequest(t testing.TB, req *http.Request, expectedStatus int) *TestResponse {
	buffer := bytes.NewBuffer(nil)
	respWriter := &TestResponseWriter{
		Writer:  buffer,
		Headers: make(map[string][]string),
	}

	fmt.Println(server)
	server.GetHandler().ServeHTTP(respWriter, req)

	if expectedStatus != NoExpectedStatus {
		assert.EqualValues(t, expectedStatus, respWriter.HeaderCode,
			"Request: %s %s", req.Method, req.URL.String())
	}
	return &TestResponse{
		HeaderCode: respWriter.HeaderCode,
		Body:       buffer.Bytes(),
		Headers:    respWriter.Headers,
	}
}

func DecodeJSON(t testing.TB, resp *TestResponse, v interface{}) bool {
	decoder := json.NewDecoder(bytes.NewBuffer(resp.Body))
	return assert.NoError(t, decoder.Decode(v))
}
