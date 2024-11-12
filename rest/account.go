package rest

import (
	"base-gin/domain/dto"
	"base-gin/domain/dao"
	"base-gin/exception"
	"base-gin/server"
	"base-gin/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AccountHandler struct {
	hr            *server.Handler
	service       *service.AccountService
	personService *service.PersonService
}

func NewAccountHandler(
	hr *server.Handler,
	accountService *service.AccountService,
	personService *service.PersonService,
) *AccountHandler {
	return &AccountHandler{
		hr:            hr,
		service:       accountService,
		personService: personService,
	}
}

func (h *AccountHandler) Route(app *gin.Engine) {
	grp := app.Group(server.RootAccount)
	grp.POST(server.PathLogin, h.login)
	grp.GET("", h.hr.AuthAccess(), h.getProfile)
	grp.POST("", h.create)
	grp.DELETE("/:id", h.hr.AuthAccess(), h.delete)
	grp.GET("/:id", h.getProfile)
	grp.GET("/profile", h.hr.AuthAccess(), h.getByID)
}

// login godoc
//
//	@Summary Account login
//	@Description Account login using username & password combination.
//	@Accept json
//	@Produce json
//	@Param cred body dto.AccountLoginReq true "Credential"
//	@Success 200 {object} dto.SuccessResponse[dto.AccountLoginResp]
//	@Failure 400 {object} dto.ErrorResponse
//	@Failure 422 {object} dto.ErrorResponse
//	@Failure 500 {object} dto.ErrorResponse
//	@Router /accounts/login [post]
func (h *AccountHandler) login(c *gin.Context) {
	var req dto.AccountLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(h.hr.BindingError(err))
		return
	}

	data, err := h.service.Login(req)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrUserNotFound),
			errors.Is(err, exception.ErrUserLoginFailed):
			c.JSON(http.StatusBadRequest, h.hr.ErrorResponse(exception.ErrUserLoginFailed.Error()))
		default:
			h.hr.ErrorInternalServer(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse[dto.AccountLoginResp]{
		Success: true,
		Message: "Login berhasil",
		Data:    data,
	})
}

// getProfile godoc
//
//	@Summary Get account's profile
//	@Description Get profile of logged-in account.
//	@Produce json
//	@Security BearerAuth
//	@Success 200 {object} dto.SuccessResponse[dto.AccountProfileResp]
//	@Failure 401 {object} dto.ErrorResponse
//	@Failure 403 {object} dto.ErrorResponse
//	@Failure 404 {object} dto.ErrorResponse
//	@Failure 500 {object} dto.ErrorResponse
//	@Router /accounts [get]
func (h *AccountHandler) getProfile(c *gin.Context) {
	accountID, _ := c.Get(server.ParamTokenUserID)

	data, err := h.personService.GetAccountProfile((accountID).(uint))
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrUserNotFound):
			c.JSON(http.StatusNotFound, h.hr.ErrorResponse(err.Error()))
		default:
			h.hr.ErrorInternalServer(c, err)
		}

		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse[dto.AccountProfileResp]{
		Success: true,
		Message: "Profile pengguna",
		Data:    data,
	})
}

// create godoc
//
//	@Summary Create a new account
//	@Description Create a new account using the provided credentials
//	@Accept json
//	@Produce json
//	@Param cred body dto.AccountCreateReq true "Account creation request"
//	@Success 201 {object} dto.SuccessResponse[dto.AccountCreateResp]
//	@Failure 400 {object} dto.ErrorResponse
//	@Failure 422 {object} dto.ErrorResponse
//	@Failure 500 {object} dto.ErrorResponse
//	@Router /accounts/ [post]
//	@Security BearerAuth
func (h *AccountHandler) create(c *gin.Context) {
	var req dto.AccountCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(h.hr.BindingError(err))
		return
	}

	account, err := h.service.Create(req)
	if err != nil {
		h.hr.ErrorInternalServer(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse[any]{
		Success: true,
		Message: "Data berhasil disimpan",
		Data:    account,
	})
}

// delete godoc
//
//	@Summary Delete an account
//	@Description Delete an account using the provided ID
//	@Accept json
//	@Produce json
//	@Param id path uint true "Account ID"
//	@Success 200 {object} dto.SuccessResponse[any]
//	@Failure 400 {object} dto.ErrorResponse
//	@Failure 404 {object} dto.ErrorResponse
//	@Failure 500 {object} dto.ErrorResponse
//	@Router /accounts/{id} [delete]
// 	@Security BearerAuth
func (h *AccountHandler) delete(c *gin.Context){
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

// get godoc
//
//	@Summary Get logged-in account profile
//	@Description Get the profile information of the logged-in account
//	@Accept json
//	@Produce json
//  @Security BearerAuth
//	@Success 200 {object} dto.SuccessResponse[dto.AccountResp]
//	@Failure 400 {object} dto.ErrorResponse
//	@Failure 404 {object} dto.ErrorResponse
//	@Failure 500 {object} dto.ErrorResponse
//	@Router /accounts/profile [get]
func (h *AccountHandler) GetProfile(c *gin.Context) {
	account, err := h.service.GetAccountByToken(c)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, h.hr.ErrorResponse("Account not found"))
		} else {
			h.hr.ErrorInternalServer(c, err)
		}
		return
	}

	dataResp := dto.AccountResp{
		ID:       account.ID,
		Username: account.Username,
	}

	c.JSON(http.StatusOK, dto.SuccessResponse[dto.AccountResp]{
		Success: true,
		Message: "Account profile found",
		Data:    dataResp,
	})
}


// @Summary Get account by ID
// @Description Retrieves an account by its ID
// @ID get-account-by-id
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "Account ID"
// @Success 200 {object} dto.SuccessResponse[dto.AccountResp] "Account found"
// @Failure 400 {object} dto.ErrorResponse "Invalid ID"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Router /accounts/{id} [get]
func (h *AccountHandler) getByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, h.hr.ErrorResponse("ID tidak valid"))
		return
	}

	data, err := h.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, h.hr.ErrorResponse("Data tidak ditemukan"))
		} else {
			h.hr.ErrorInternalServer(c, err)
		}
	}
	
	personData, err := h.personService.GetAccountProfile(data.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, h.hr.ErrorResponse("Data tidak ditemukan"))
		} else {
			h.hr.ErrorInternalServer(c, err)
		}
		return
	}
	
	accountResp := dto.AccountResp{
		ID:       data.ID,
		Username: data.Username,
		Fullname: personData.Fullname,
		Gender:   personData.Gender,
		Age:      personData.Age,
	}
	
	c.JSON(http.StatusOK, dto.SuccessResponse[dto.AccountResp]{
		Success: true,
		Message: "Account found",
		Data:    accountResp,
	})
}

// @Summary Get a list of accounts
// @Description Retrieves a list of accounts based on the provided filter
// @ID get-accounts
// @Tags default
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param filter query dto.Filter true "Filter for the accounts"
// @Success 200 {array} dao.Account "List of accounts"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Router /accounts [get]
func (h *AccountHandler) GetList(c *gin.Context) {
	params := &dto.Filter{}
	if err := c.ShouldBindQuery(params); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	accounts, err := h.getAccountList(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, accounts)
}


func (h *AccountHandler) getAccountList(params *dto.Filter) ([]dao.Account, error) {
	accounts, err := h.service.GetList(params)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, exception.ErrDataNotFound
		default:
			return nil, err
		}
	}
	return accounts, nil
}



// @Summary Update an account
// @Description Updates an account based on the provided ID and request body
// @ID update-account
// @Tags default
// @Accept json
// @Produce json
// @Param id path uint true "Account ID"
// @Param account body dto.AccountUpdateReq true "Account data"
// @Success 200 {object} dao.Account "Updated account"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 404 {object} dto.ErrorResponse "Not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /accounts/{id} [put]
// @Security BearerAuth
func (h *AccountHandler) Update(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, h.hr.ErrorResponse("ID tidak valid"))
		return
	}

	var req dto.AccountUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(h.hr.BindingError(err))
		return
	}
	req.ID = uint(id)

	account, err := h.service.Update(&req)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, h.hr.ErrorResponse("Data tidak ditemukan"))
		default:
			h.hr.ErrorInternalServer(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse[dto.AccountResp]{
			Success: true,
			Message: "Data berhasil disimpan",
			Data: dto.AccountResp{
				ID:       account.ID,
				Username: account.Username,
			},
		})
}
