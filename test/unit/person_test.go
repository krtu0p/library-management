package unit_test

import (
	"base-gin/domain"
	"base-gin/domain/dao"
	"base-gin/domain/dto"
	"base-gin/repository"
	"base-gin/util"
	"testing"
	"time"
	"math/rand"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestPerson_Update_Success(t *testing.T) {
	birthDate, _ := time.Parse("2006-01-02", "1993-09-13")
	gender := domain.GenderFemale
	params := dto.PersonUpdateReq{
		ID:           dummyMember.ID,
		Fullname:     util.RandomStringAlpha(4) + " " + util.RandomStringAlpha(6) + " " + util.RandomStringAlpha(6),
		Gender:       string(gender),
		BirthDateStr: birthDate.Format("2006-01-02"),
		BirthDate:    birthDate,
	}

	err := personRepo.Update(&params)
	assert.Nil(t, err)

	item, _ := personRepo.GetByID(dummyMember.ID)
	assert.Equal(t, params.Fullname, item.Fullname)
	assert.EqualValues(t, params.Gender, string(*item.Gender))
	assert.EqualValues(t, params.BirthDateStr, item.BirthDate.Format("2006-01-02"))
}

func TestPerson_Create_Success(t *testing.T) {
    // Mock data untuk create
    birthDate, _ := time.Parse("2006-01-02", "1995-05-12")
    gender := domain.GenderMale
    createReq := dto.PersonCreateReq{
        Fullname:     util.RandomStringAlpha(4) + " " + util.RandomStringAlpha(6),
        Gender:       string(gender),
        BirthDateStr: birthDate.Format("2006-01-02"),
        BirthDate:    birthDate,
    }
    createPerson := createReq.ToEntity()
    err := personRepo.Create(&createPerson)
    assert.Nil(t, err)


    // Ambil data yang baru saja dibuat
    createdPerson, _ := personRepo.GetByID(createPerson.ID)
    assert.Equal(t, createReq.Fullname, createdPerson.Fullname)
    assert.Equal(t, createReq.Gender, string(*createdPerson.Gender))
    assert.Equal(t, createReq.BirthDateStr, createdPerson.BirthDate.Format("2006-01-02"))
}

func TestPerson_GetByID_Success(t *testing.T) {
    // Mock ID untuk pengujian
    id := uint(1)

    // Ambil data berdasarkan ID
    person, err := personRepo.GetByID(id)
    assert.Nil(t, err)
    assert.NotNil(t, person, "Data person harusnya tidak nil")

    // Validasi data
    assert.Equal(t, id, person.ID)                      // Pastikan tipe data sesuai
    assert.NotEmpty(t, person.Fullname, "Fullname tidak boleh kosong")
}

func genRandomInt(t *testing.T) int {
	t.Helper()
	return rand.Int()
}

func TestPerson_GetList_Success(t *testing.T) {
	db.Exec("SET FOREIGN_KEY_CHECKS=0")
	db.Exec("TRUNCATE TABLE borrowings")
	db.Exec("TRUNCATE TABLE persons")
	db.Exec("SET FOREIGN_KEY_CHECKS=1")
	

	personRepo := repository.NewPersonRepository(db)

	// Membuat data mock
	male := domain.GenderMale
	female := domain.GenderFemale
	
	// Generate a unique id for each person
	id1 := uint(genRandomInt(t))
	id2 := uint(genRandomInt(t))

	mockData := []dao.Person{
		{Model: gorm.Model{ID: id1}, Fullname: "John Doe", Gender: &male},
		{Model: gorm.Model{ID: id2}, Fullname: "Jane Smith", Gender: &female},
	}
	// Menambahkan data mock ke database
	for _, person := range mockData {
		err := personRepo.Create(&person)
		assert.Nil(t, err)
	}

	// Memanggil GetList tanpa filter
	persons, err := personRepo.GetList(nil)
	assert.Nil(t, err)
	assert.NotNil(t, persons, "Hasil GetList tidak boleh nil")
	assert.Equal(t, len(mockData), len(persons), "Jumlah data tidak sesuai")

	// Validasi data yang dikembalikan
	for i, person := range persons {
		assert.Equal(t, mockData[i].Fullname, person.Fullname)
		assert.Equal(t, *mockData[i].Gender, *person.Gender)
	}
}


func TestPerson_Delete_Success(t *testing.T) {
    // ID yang akan dihapus
    id := 1

    // Hapus data menggunakan repository
    err := personRepo.Delete(uint(id))
    assert.Nil(t, err)

    // Coba ambil data yang sudah dihapus
    deletedPerson, err := personRepo.GetByID(uint(id))
    assert.NotNil(t, err)
    assert.Nil(t, deletedPerson)
}
