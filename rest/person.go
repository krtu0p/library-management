package rest

import (
	"base-gin/domain/dto"
	"base-gin/exception"
	"base-gin/server"
	"base-gin/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PersonHandler struct {
	hr      *server.Handler
	service *service.PersonService
}

func NewPersonHandler(
	hr *server.Handler,
	personService *service.PersonService,
) *PersonHandler {
	return &PersonHandler{hr: hr, service: personService}
}

func (h *PersonHandler) Route(app *gin.Engine) {
	grp := app.Group(server.RootPerson)
	grp.GET("", h.getList)
	grp.GET("/:id", h.getByID)
	grp.PUT("/:id", h.hr.AuthAccess(), h.update)
	grp.POST("", h.create)
	grp.DELETE("/:id", h.hr.AuthAccess(), h.delete)
}

// getList godoc
//
//	@Summary Get a list of person
//	@Description Get a list of person.
//	@Produce json
//	@Param q query string false "Person's name"
//	@Param s query int false "Data offset"
//	@Param l query int false "Data limit"
//	@Success 200 {object} dto.SuccessResponse[[]dto.PersonDetailResp]
//	@Failure 422 {object} dto.ErrorResponse
//	@Failure 404 {object} dto.ErrorResponse
//	@Failure 422 {object} dto.ErrorResponse
//	@Failure 500 {object} dto.ErrorResponse
//	@Router /persons [get]
func (h *PersonHandler) getList(c *gin.Context) {
	var req dto.Filter
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(h.hr.BindingError(err))
		return
	}

	data, err := h.service.GetList(&req)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrUserNotFound):
			c.JSON(http.StatusNotFound, h.hr.ErrorResponse(exception.ErrDataNotFound.Error()))
		default:
			h.hr.ErrorInternalServer(c, err)
		}

		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse[[]dto.PersonDetailResp]{
		Success: true,
		Message: "Daftar anggota",
		Data:    data,
	})
}

// getByID godoc
//
//	@Summary Get a person's detail
//	@Description Get a person's detail.
//	@Produce json
//	@Param id path int true "Person's ID"
//	@Success 200 {object} dto.SuccessResponse[dto.PersonDetailResp]
//	@Failure 400 {object} dto.ErrorResponse
//	@Failure 404 {object} dto.ErrorResponse
//	@Failure 500 {object} dto.ErrorResponse
//	@Router /persons/{id} [get]
func (h *PersonHandler) getByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, h.hr.ErrorResponse("ID tidak valid"))
		return
	}

	data, err := h.service.GetByID(uint(id))
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrUserNotFound):
			c.JSON(http.StatusNotFound, h.hr.ErrorResponse(exception.ErrDataNotFound.Error()))
		default:
			h.hr.ErrorInternalServer(c, err)
		}

		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse[dto.PersonDetailResp]{
		Success: true,
		Message: "Detail anggota",
		Data:    data,
	})
}

// update godoc
//
//	@Summary Update a person's detail
//	@Description Update a person's detail.
//	@Accept json
//	@Produce json
//	@Security BearerAuth
//	@Param id path int true "Person's ID"
//	@Param detail body dto.PersonUpdateReq true "Person's detail"
//	@Success 200 {object} dto.SuccessResponse[any]
//	@Failure 400 {object} dto.ErrorResponse
//	@Failure 401 {object} dto.ErrorResponse
//	@Failure 403 {object} dto.ErrorResponse
//	@Failure 404 {object} dto.ErrorResponse
//	@Failure 500 {object} dto.ErrorResponse
//	@Router /persons/{id} [put]
func (h *PersonHandler) update(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, h.hr.ErrorResponse("ID tidak valid"))
		return
	}

	var req dto.PersonUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(h.hr.BindingError(err))
		return
	}
	req.ID = uint(id)

	err = h.service.Update(&req)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrDateParsing):
			c.JSON(http.StatusBadRequest, h.hr.ErrorResponse(err.Error()))
		case errors.Is(err, exception.ErrUserNotFound):
			c.JSON(http.StatusNotFound, h.hr.ErrorResponse(err.Error()))
		default:
			h.hr.ErrorInternalServer(c, err)
		}

		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse[any]{
		Success: true,
		Message: "Data berhasil disimpan",
	})
}


//	@Summary Create a person
//	@Description Create a person.
//	@Accept json
//	@Produce json
//	@Param detail body dto.PersonCreateReq true "Person's detail"
//	@Success 201 {object} dto.SuccessResponse[any]
//	@Failure 401 {object} dto.ErrorResponse
//	@Failure 403 {object} dto.ErrorResponse
//	@Failure 422 {object} dto.ErrorResponse
//	@Failure 500 {object} dto.ErrorResponse
//	@Router /persons [post]
func (h *PersonHandler) create(c *gin.Context) {
	var req dto.PersonCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(h.hr.BindingError(err))
		return
	}

	err := h.service.Create(&req)
	if err != nil {
		h.hr.ErrorInternalServer(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse[any]{
		Success: true,
		Message: "Data berhasil disimpan",
	})
}


// @Summary Delete a Person by ID
// @Description Deletes a Person entity by ID
// @ID person-repository-delete
// @Accept json
// @Produce json
// @Param id path uint true "Person ID"
// @Success 200 {object} interface{} "Person deleted successfully"
// @Failure 404 {object} interface{} "Person not found"
// @Failure 500 {object} interface{} "Internal server error"
// @Router /persons/{id} [delete]
func (h *PersonHandler) delete(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, h.hr.ErrorResponse("ID tidak valid"))
		return
	}

	err = h.service.Delete(uint(id))
	if err != nil {
		h.hr.ErrorInternalServer(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse[any]{
		Success: true,
		Message: "Data berhasil dihapus",
	})
}