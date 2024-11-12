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

type BorrowingHandler struct {
	hr *server.Handler
	service *service.BorrowingService
}

func NewBorrowingHandler (
	hr *server.Handler,
	borrowingService *service.BorrowingService,
) *BorrowingHandler {
	return &BorrowingHandler{hr: hr, service: borrowingService}
}

func (h *BorrowingHandler) Route(app *gin.Engine) {
	grp := app.Group(server.RootBorrowing)
	grp.GET("", h.getList)
	grp.GET("/:id", h.getByID)
	grp.POST("", h.hr.AuthAccess(), h.create)
	grp.PUT("/:id", h.hr.AuthAccess(), h.update)
	grp.DELETE("/:id", h.hr.AuthAccess(), h.delete)
}

// create godoc
//
//	@Summary Create a borrowing
//	@Description Create a borrowing.
//	@Accept json
//	@Produce json
//	@Security BearerAuth
//	@Param detail body dto.BorrowingUpdateReq true "Borrowing's detail"
//	@Success 201 {object} dto.SuccessResponse[any]
//	@Failure 401 {object} dto.ErrorResponse
//	@Failure 403 {object} dto.ErrorResponse
//	@Failure 422 {object} dto.ErrorResponse
//	@Failure 500 {object} dto.ErrorResponse
//	@Router /borrowings [post]
func (h *BorrowingHandler) create(c *gin.Context) {
	var req dto.BorrowingCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.hr.BindingError(err)
		return
	}

	err := h.service.Create(&req)
	if err != nil {
		h.hr.ErrorInternalServer(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse[any]{
		Success: true,
		Message: "Data berhasil disimpan",
	})
}

// getList godoc
//
//	@Summary Get a list of borrowings
//	@Description Get a list of borrowings.
//	@Produce json
//	@Param q query string false "Borrowing's name"
//	@Param s query int false "Data offset"
//	@Param l query int false "Data limit"
//	@Success 200 {object} dto.SuccessResponse[[]dto.BorrowingResp]
//	@Failure 422 {object} dto.ErrorResponse
//	@Failure 404 {object} dto.ErrorResponse
//	@Failure 422 {object} dto.ErrorResponse
//	@Failure 500 {object} dto.ErrorResponse
//	@Router /borrowings [get]
func (h *BorrowingHandler) getList(c *gin.Context) {
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

	c.JSON(http.StatusOK, dto.SuccessResponse[[]dto.BorrowingResp]{
		Success: true,
		Message: "Daftar buku",
		Data:    data,
	})
}

// getByID godoc
//
//	@Summary Get a borrowing's detail
//	@Description Get a borrowing's detail.
//	@Produce json
//	@Param id path int true "Borrowing's ID"
//	@Success 200 {object} dto.SuccessResponse[dto.BorrowingResp]
//	@Failure 400 {object} dto.ErrorResponse
//	@Failure 404 {object} dto.ErrorResponse
//	@Failure 500 {object} dto.ErrorResponse
//	@Router /borrowings/{id} [get]
func (h *BorrowingHandler) getByID(c *gin.Context) {
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

	c.JSON(http.StatusOK, dto.SuccessResponse[dto.BorrowingResp]{
		Success: true,
		Message: "Detail buku",
		Data:    data,
	})
}

// update godoc
//
//	@Summary Update a borrowing's detail
//	@Description Update a borrowing's detail.
//	@Accept json
//	@Produce json
//	@Security BearerAuth
//	@Param id path int true "Borrowing's ID"
//	@Param detail body dto.BorrowingUpdateReq true "Borrowing's detail"
//	@Success 200 {object} dto.SuccessResponse[any]
//	@Failure 400 {object} dto.ErrorResponse
//	@Failure 401 {object} dto.ErrorResponse
//	@Failure 403 {object} dto.ErrorResponse
//	@Failure 404 {object} dto.ErrorResponse
//	@Failure 422 {object} dto.ErrorResponse
//	@Failure 500 {object} dto.ErrorResponse
//	@Router /borrowings/{id} [put]
func (h *BorrowingHandler) update(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, h.hr.ErrorResponse("ID tidak valid"))
		return
	}

	var req dto.BorrowingUpdateReq
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

// delete godoc
//
//	@Summary Delete a borrowing
//	@Description Delete a borrowing.
//	@Produce json
//	@Security BearerAuth
//	@Param id path int true "Borrowing's ID"
//	@Success 200 {object} dto.SuccessResponse[any]
//	@Failure 400 {object} dto.ErrorResponse
//	@Failure 401 {object} dto.ErrorResponse
//	@Failure 403 {object} dto.ErrorResponse
//	@Failure 404 {object} dto.ErrorResponse
//	@Failure 500 {object} dto.ErrorResponse
//	@Router /borrowings/{id} [delete]
func (h *BorrowingHandler) delete(c *gin.Context) {
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
