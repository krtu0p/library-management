package integration_test

import (
	"base-gin/domain/dto"
	"base-gin/server"
	"testing"
	"base-gin/domain/dao"
	"base-gin/util"
	"fmt"
	"strings"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"github.com/gin-gonic/gin"
	"base-gin/repository"


	"github.com/stretchr/testify/assert"

)

func TestAccount_Login_Success(t *testing.T) {
	req := dto.AccountLoginReq{
		Username: "admin",
		Password: password,
	}

	w := doTest("POST", server.RootAccount+server.PathLogin, req, "")
	assert.Equal(t, 200, w.Code)
}

func TestAccount_GetProfile_Success(t *testing.T) {
	accessToken := createAuthAccessToken(dummyAdmin.Account.Username)

	w := doTest("GET", server.RootAccount, nil, accessToken)
	assert.Equal(t, 200, w.Code)

	resp := w.Body.String()
	assert.Contains(t, resp, dummyAdmin.Fullname)
}

func TestAccount_GetProfile_ErrorAccessToken(t *testing.T) {
	w := doTest("GET", server.RootAccount, nil, "")
	assert.Equal(t, 401, w.Code)

	w = doTest("GET", server.RootAccount, nil, "accessToken")
	assert.Equal(t, 401, w.Code)
}

func TestAccount_Create_Success(t *testing.T) {
	req := dto.AccountCreateReq{
		Username: util.RandomStringAlpha(10),
		Password: password,
	}

	w := doTest("POST", server.RootAccount, req, "")
	fmt.Println("Request URL:", server.RootAccount)
	fmt.Println("Response Body:", w.Body.String())
	assert.Equal(t, 201, w.Code)
}
func TestAccount_Delete_Success(t *testing.T) {
	o := dao.Account{
		Username: util.RandomStringAlpha(10),
		Password: password,
	}
	err := accountRepo.Create(&o)
	if err != nil {
		t.Fatal(err)
	}

	accessToken := createAuthAccessToken(dummyAdmin.Account.Username)

	w := doTest("DELETE", fmt.Sprintf("%s/%d", server.RootAccount, o.ID), nil, accessToken)
	if w.Code != 200 {
		t.Errorf("expected status code 200, but got %d", w.Code)
	}

	item, err := accountRepo.GetByID(o.ID)
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		t.Fatal(err)
	}
	if item != nil {
		t.Errorf("expected account to be deleted, but got %+v", item)
	}
}

func TestAccount_GetByID_Success(t *testing.T) {

	o := dao.Account{
		Username: util.RandomStringAlpha(10),
		Password: password,
	}


	accountRepo := repository.NewAccountRepository(db)
	err := accountRepo.Create(&o)
	if err != nil {
		t.Fatal(err)
	}

	url := fmt.Sprintf("/accounts/%d", o.ID)
	req := httptest.NewRequest("GET", url, nil)


	w := httptest.NewRecorder()
	router := gin.Default()
	router.GET("/accounts/:id", func(c *gin.Context) {
		id, err := util.StringToUint(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		account, err := accountRepo.GetByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}
		c.JSON(http.StatusOK, account)
	})
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseAccount dao.Account
	err = json.Unmarshal(w.Body.Bytes(), &responseAccount)
	assert.NoError(t, err)
	assert.Equal(t, o.ID, responseAccount.ID)
	assert.Equal(t, o.Username, responseAccount.Username)
	assert.Equal(t, o.Password, responseAccount.Password)

	err = accountRepo.Delete(o.ID)
	assert.NoError(t, err)
}

func TestAccount_Update_Success(t *testing.T) {
	// Create a new account
	account := dao.Account{
		Username: "example",
		Password: "password",
	}
	err := accountRepo.Create(&account)
	if err != nil {
		t.Fatal(err)
	}

	// Update the account
	updatedAccount := dao.Account{
		ID:       account.ID,
		Username: "updated-example",
		Password: "updated-password",
	}
	err = accountRepo.Update(&updatedAccount)
	if err != nil {
		t.Fatal(err)
	}

	// Get the updated account
	updatedAccountFromDB, err := accountRepo.GetByID(account.ID)
	if err != nil {
		t.Fatal(err)
	}


	fmt.Println("Updated account:", updatedAccountFromDB)

	if updatedAccountFromDB.Username != updatedAccount.Username {
		t.Errorf("expected username to be %s, but got %s", updatedAccount.Username, updatedAccountFromDB.Username)
	}
	if updatedAccountFromDB.Password != updatedAccount.Password {
		t.Errorf("expected password to be %s, but got %s", updatedAccount.Password, updatedAccountFromDB.Password)
	}
}
type MockAccountRepo struct {
	GetListFn func(filter interface{}) ([]dao.Account, error)
	CreateFn  func(newItem *dao.Account) error
	GetByIDFn func(id uint) (*dao.Account, error)
	UpdateFn  func(newItem *dao.Account) error
}


func (m *MockAccountRepo) GetList(filter interface{}) ([]dao.Account, error) {
    return m.GetListFn(filter)
}

func (m *MockAccountRepo) Create(newItem *dao.Account) error {
    return m.CreateFn(newItem)
}

func (m *MockAccountRepo) GetByID(id uint) (*dao.Account, error) {
    return m.GetByIDFn(id)
}

func (m *MockAccountRepo) Update(newItem *dao.Account) error {
    return m.UpdateFn(newItem)
}

func TestAccountGetList_Success(t *testing.T) {
	accounts := []dao.Account{
		{Username: "test1", Password: "password1"},
		{Username: "test2", Password: "password2"},
	}

	// Initialize the mock repository
	accountRepo := &MockAccountRepo{
		GetListFn: func(filter interface{}) ([]dao.Account, error) {
			return accounts, nil
		},
		CreateFn: func(newItem *dao.Account) error {
			return nil
		},
		GetByIDFn: func(id uint) (*dao.Account, error) {
			return nil, nil
		},
		UpdateFn: func(newItem *dao.Account) error {
			return nil
		},
	}

	// Use the accountRepo to call the GetList method
	result, err := accountRepo.GetList(nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != len(accounts) {
		t.Errorf("expected %d accounts, but got %d", len(accounts), len(result))
	}
}

