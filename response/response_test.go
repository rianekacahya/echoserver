package response

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/rianekacahya/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestError(t *testing.T) {

	e := echo.New()

	e.GET("/ping", func (c echo.Context) error {
		err := errors.New(errors.BADREQUEST, errors.Message("Error"))
		return Error(c, err)
	})

	ts := httptest.NewServer(e)
	defer ts.Close()
	var cl http.Client
	req, _ := http.NewRequest("GET", ts.URL+"/ping", nil)
	resp, _ := cl.Do(req)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, fmt.Sprintf("%s\n", `{"message":"Error"}`), string(body))
}

func TestRender(t *testing.T) {

	e := echo.New()

	e.GET("/pong", func (c echo.Context) error {
		return Render(c, http.StatusOK, nil)
	})

	ts := httptest.NewServer(e)
	defer ts.Close()
	var cl http.Client
	req, _ := http.NewRequest("GET", ts.URL+"/pong", nil)
	resp, _ := cl.Do(req)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, fmt.Sprintf("%s\n", `{"message":"success"}`), string(body))
}