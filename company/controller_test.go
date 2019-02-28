package company

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
)

type serviceMock struct {
	findAllFn              func() ([]Company, error)
	findByNameAndZipCodeFn func(string, string) (Company, error)
	addFn                  func(Company) error
	InitDatabaseFn         func(string) error
	loadWebsitesFn         func(io.Reader) error
}

func (s serviceMock) findByNameAndZipCode(n string, z string) (Company, error) {
	return s.findByNameAndZipCodeFn(n, z)
}

func (s serviceMock) add(c Company) error {
	return s.addFn(c)
}

func (s serviceMock) InitDatabase(st string) error {
	return s.InitDatabaseFn(st)
}

func (s serviceMock) findAll() ([]Company, error) {
	return s.findAllFn()
}

func (s serviceMock) loadWebsites(f io.Reader) error {
	return s.loadWebsitesFn(f)
}

func TestNewController(t *testing.T) {
	cMock := serviceMock{}
	type args struct {
		service Service
	}
	tests := []struct {
		name string
		args args
		want Controller
	}{
		{"Create controller with service", args{cMock}, companyController{cMock}},
		{"Create controller empty", args{}, companyController{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewController(tt.args.service); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewController() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_companyController_GetAll(t *testing.T) {
	ctxMock, _ := gin.CreateTestContext(httptest.NewRecorder())
	sMockError := serviceMock{
		findAllFn: func() ([]Company, error) { return []Company{}, errors.New("mock error") },
	}
	sMockOne := serviceMock{
		findAllFn: func() ([]Company, error) { return []Company{{}, {}}, nil },
	}
	sMockMany := serviceMock{
		findAllFn: func() ([]Company, error) { return []Company{}, nil },
	}
	type fields struct {
		service Service
	}
	type args struct {
		ctx *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Test getAll returns error", fields{sMockError}, args{ctx: ctxMock}},
		{"Test getAll returns one", fields{sMockOne}, args{ctx: ctxMock}},
		{"Test getAll returns many", fields{sMockMany}, args{ctx: ctxMock}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := companyController{
				service: tt.fields.service,
			}
			c.GetAll(tt.args.ctx)
		})
	}
}

func Test_companyController_Find(t *testing.T) {
	ctxMockMany, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctxMockMany.Request, _ = http.NewRequest("GET", "ab.com/test", strings.NewReader(""))
	ctxMockError, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctxMockError.Request, _ = http.NewRequest("GET", "ab.com/test?zipcode=12345", strings.NewReader(""))
	ctxMockOne, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctxMockOne.Request, _ = http.NewRequest("GET", "ab.com/test?zipcode=12345&name=asdf", strings.NewReader(""))
	ctxMockMgoError, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctxMockMgoError.Request, _ = http.NewRequest("GET", "ab.com/test?zipcode=123&name=mgo", strings.NewReader(""))

	sMock := serviceMock{
		findAllFn: func() ([]Company, error) { return []Company{{}, {}}, nil },
		findByNameAndZipCodeFn: func(a string, b string) (Company, error) {
			if a == "mgo" {
				return Company{}, mgo.ErrNotFound
			}
			if a == "error" {
				return Company{}, errors.New("mock error")
			}
			return Company{}, nil
		},
	}
	type fields struct {
		service Service
	}
	type args struct {
		ctx *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Find many", fields{sMock}, args{ctxMockMany}},
		{"Find one", fields{sMock}, args{ctxMockOne}},
		{"Company not found", fields{sMock}, args{ctxMockError}},
		{"mgo.error", fields{sMock}, args{ctxMockMgoError}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := companyController{
				service: tt.fields.service,
			}
			c.Find(tt.args.ctx)
		})
	}
}

func Test_companyController_LoadWebsites(t *testing.T) {
	ctxMockFile, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctxMockFile.Request = &http.Request{Method: "GET", Host: "ab.com/test",
		MultipartForm: &multipart.Form{
			Value: make(map[string][]string),
			File:  make(map[string][]*multipart.FileHeader),
		}}
	ctxMockFile.Request.MultipartForm.Value["data"] = []string{"file"}
	ctxMockFile.Request.MultipartForm.File["data"] = []*multipart.FileHeader{{Filename: "", Header: textproto.MIMEHeader{}, Size: 128}}
	ctxMockNoFile, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctxMockNoFile.Request, _ = http.NewRequest("GET", "ab.com/test", strings.NewReader(""))

	sMock := serviceMock{
		findAllFn: func() ([]Company, error) { return []Company{{}, {}}, nil },
		findByNameAndZipCodeFn: func(a string, b string) (Company, error) {
			if a == "mgo" {
				return Company{}, mgo.ErrNotFound
			}
			if a == "error" {
				return Company{}, errors.New("mock error")
			}
			return Company{}, nil
		},
	}

	type fields struct {
		service Service
	}
	type args struct {
		ctx *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Load website", fields{sMock}, args{ctxMockFile}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := companyController{
				service: tt.fields.service,
			}
			c.LoadWebsites(tt.args.ctx)
		})
	}
}

func Test_companyController_InitDatabase(t *testing.T) {
	ctxMockMany, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctxMockMany.Request, _ = http.NewRequest("GET", "ab.com/test", strings.NewReader(""))

	sMock := serviceMock{
		InitDatabaseFn: func(a string) error {
			if a == "" {
				return nil
			} else {
				return errors.New("mock Error")
			}
		},
		findByNameAndZipCodeFn: func(a string, b string) (Company, error) {
			if a == "mgo" {
				return Company{}, mgo.ErrNotFound
			}
			if a == "error" {
				return Company{}, errors.New("mock error")
			}
			return Company{}, nil
		},
	}
	type fields struct {
		service Service
	}
	type args struct {
		file string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"call", fields{sMock}, args{file: ""}},
		{"call error", fields{sMock}, args{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := companyController{
				service: tt.fields.service,
			}
			c.InitDatabase(tt.args.file)
		})
	}
}
