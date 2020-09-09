package echobox

import (
	"net/http"
	"net/url"
	"path/filepath"
	"simple-auth/pkg/box"

	"github.com/labstack/echo/v4"
)

func Static(root string) echo.HandlerFunc {
	return StaticBox(root, box.Global)
}

func StaticBox(root string, readBox box.Box) echo.HandlerFunc {
	if root == "" {
		root = "."
	}

	return func(c echo.Context) error {
		p, err := url.PathUnescape(c.Param("*"))
		if err != nil {
			return err
		}

		name := filepath.Join(root, filepath.Clean("/"+p))

		fi, err := readBox.Stat(name)
		if err != nil {
			return c.HTML(http.StatusNotFound, "Not found")
		}

		r, err := readBox.Read(name)
		if err != nil {
			return c.HTML(http.StatusNotFound, "Not found")
		}

		http.ServeContent(c.Response(), c.Request(), fi.Name(), fi.ModTime(), r)
		return nil
	}
}
