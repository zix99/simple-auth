package common

import "github.com/labstack/echo"

type Controller interface {
	Mount(group echo.Group)
}
