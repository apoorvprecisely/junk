package api

import (
	"github.com/labstack/echo"
	"lucene"
	"net/http"
)

func Search(c echo.Context) error {
	field := c.Param("field")
	query := c.Param("query")
	results := lucene.Find(field, query)
	return c.JSON(http.StatusOK, results)
}

func Write(c echo.Context) error {
	call := new(lucene.Call)
	if err := c.Bind(call); err != nil {
		return err
	}
	lucene.Write(*call)
	return c.JSON(http.StatusCreated, call)
}
