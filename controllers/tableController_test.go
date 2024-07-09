
package controllers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
   // "restaurant/database"
    "restaurant/models"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "gorm.io/gorm"
)

// MockDB struct
type MockDB struct {
    mock.Mock
}

// Implement the methods of the DBInterface
func (m *MockDB) Create(value interface{}) *gorm.DB {
    args := m.Called(value)
    return args.Get(0).(*gorm.DB)
}

func (m *MockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
    args := m.Called(dest, conds)
    return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Save(value interface{}) *gorm.DB {
    args := m.Called(value)
    return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Delete(value interface{}) *gorm.DB {
    args := m.Called(value)
    return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Find(dest interface{}, conds ...interface{}) *gorm.DB {
    args := m.Called(dest, conds)
    return args.Get(0).(*gorm.DB)
}

func TestCreateTable(t *testing.T) {
    mockDb := new(MockDB)
    //database.DB = mockDb

    gin.SetMode(gin.TestMode)
    router := gin.Default()
    router.POST("/tables", CreateTable)

    table := models.TablesModel{
        Capacity:     4,
        Availability: true,
    }

    mockDb.On("Create", mock.AnythingOfType("*models.TablesModel")).Return(&gorm.DB{Error: nil})

    body, _ := json.Marshal(table)
    req, _ := http.NewRequest("POST", "/tables", bytes.NewBuffer(body))
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, 201, w.Code)

    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, "Table created successfully", response["message"])

    mockDb.AssertExpectations(t)
}