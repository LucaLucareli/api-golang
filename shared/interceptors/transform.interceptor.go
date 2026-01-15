package interceptors

import (
	"net/http"
	"shared/interfaces"

	"github.com/labstack/echo/v4"
)

func TransformInterceptor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			return err
		}

		resp, ok := c.Get("response").(interfaces.HttpResponse)
		if !ok {
			return nil
		}

		return c.JSON(http.StatusOK, resp)
	}
}
