package echobox

import (
	"net/http"
	"net/http/httptest"
	"simple-auth/pkg/box"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func testGetFile(t *testing.T, root, reqPath string) *httptest.ResponseRecorder {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, reqPath, nil)
	res := httptest.NewRecorder()

	c := e.NewContext(req, res)
	c.SetPath("/*")
	c.SetParamNames("*")
	c.SetParamValues(reqPath)

	assert.NoError(t, Static(root)(c))

	return res
}

func TestEchoBoxGetFile(t *testing.T) {
	res := testGetFile(t, "./", "/route.go")
	assert.Equal(t, 200, res.Result().StatusCode)
	assert.Greater(t, res.Result().ContentLength, int64(0))

	assert.Contains(t, res.Body.String(), "package echobox")
}

func TestEchoBoxNotFOund(t *testing.T) {
	res := testGetFile(t, "./", "/blabla")
	assert.Equal(t, http.StatusNotFound, res.Result().StatusCode)
}

func TestReadFileAtRoot(t *testing.T) {
	res := testGetFile(t, "./", "/etc/passwd")
	assert.Equal(t, http.StatusNotFound, res.Result().StatusCode)
}

func TestReadingFileUpDir(t *testing.T) {
	res := testGetFile(t, "./", "/../../../../../../../../../../../etc/passwd")
	assert.Equal(t, http.StatusNotFound, res.Result().StatusCode)
}

func TestReadFileInBox(t *testing.T) {
	box := box.NewBox()

	e := echo.New()
	{
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		res := httptest.NewRecorder()

		c := e.NewContext(req, res)
		c.SetPath("/*")
		c.SetParamNames("*")
		c.SetParamValues("/test")

		assert.NoError(t, StaticBox("", box)(c))
		assert.Equal(t, 404, res.Result().StatusCode)
	}

	box.AddBytes("test", []byte{1, 2, 3})

	{
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		res := httptest.NewRecorder()

		c := e.NewContext(req, res)
		c.SetPath("/*")
		c.SetParamNames("*")
		c.SetParamValues("/test")

		assert.NoError(t, StaticBox("", box)(c))
		assert.Equal(t, 200, res.Result().StatusCode)
		assert.Equal(t, int64(3), res.Result().ContentLength)
		assert.Equal(t, []byte{1, 2, 3}, res.Body.Bytes())
	}
}
