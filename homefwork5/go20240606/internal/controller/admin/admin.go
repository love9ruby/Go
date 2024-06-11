package admin

import (
	"go20240606/internal/services/users"
	_ "go20240606/pkg/pg"
	"net/http"

	"github.com/labstack/echo/v4"
)

var UserService users.IUserService

func init() {
	UserService = users.NewUserService()
}

// GetAdminStatistics godoc
//
//	@Summary		Get admin statistics
//	@Description	Get statistics for an admin user with the provided password
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			password	path		string		true	"Admin password"
//	@Success		200			{object}	[]pg.Event	"Statistics events"
//	@Failure		500			{object}	error		"Internal server error"
//	@Router			/admin/statistic/{password} [get]
func GetAdminStatistics(c echo.Context) error {
	// get pwd on url path
	password := c.Param("password")
	// get statistics
	events, err := UserService.GetUserStatistics(password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, events)
}

// CreateAdmin godoc
//
//	@Summary		Create admin user
//	@Description	Create a new admin user with the provided email
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			admin	body		CreateAdminRequest	true	"Admin details"
//	@Success		200		{object}	map[string]string	"Created admin user"
//	@Failure		400		{object}	error				"Bad request"
//	@Failure		500		{object}	error				"Internal server error"
//	@Router			/admin/register [post]
func CreateAdmin(c echo.Context) error {
	// Create a new event for a user
	req := new(CreateAdminRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	pwd, err := UserService.CreateUser(req.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, map[string]string{"password": pwd, "email": req.Email})
}
