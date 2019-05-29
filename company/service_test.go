package company

import (
	"errors"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/globalsign/mgo"
)

type repoMock struct {
	FindAllFn          func() ([]Company, error)
	FindByNameAndZipFn func(string, int64) (Company, error)
	AddFn              func(Company) error
	MergeWebsiteFn     func(Company) (*mgo.ChangeInfo, error)
}

func (r repoMock) FindAll() ([]Company, error) { return r.FindAllFn() }
func (r repoMock) FindByNameAndZip(a string, b int64) (Company, error) {
	return r.FindByNameAndZipFn(a, b)
}
func (r repoMock) Add(c Company) error                             { return r.AddFn(c) }
func (r repoMock) MergeWebsite(c Company) (*mgo.ChangeInfo, error) { return r.MergeWebsiteFn(c) }

func TestNewService(t *testing.T) {
	type args struct {
		r Repository
	}
	tests := []struct {
		name string
		args args
		want Service
	}{
		{"Test new Service", args{}, companyService{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewService(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_companyService_findAll(t *testing.T) {
	type fields struct {
		repository Repository
	}
	tests := []struct {
		name    string
		fields  fields
		want    []Company
		wantErr bool
	}{
		{"Return empty company list", fields{repository: repoMock{FindAllFn: func() ([]Company, error) {
			return []Company{}, nil
		}}}, []Company{}, false},
		{"Return 2 companies", fields{repository: repoMock{FindAllFn: func() ([]Company, error) {
			return []Company{{}, {}}, nil
		}}}, []Company{{}, {}}, false},
		{"Return error", fields{repository: repoMock{FindAllFn: func() ([]Company, error) {
			return []Company{}, errors.New("mockError")
		}}}, []Company{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := companyService{
				repository: tt.fields.repository,
			}
			got, err := s.findAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("companyService.findAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("companyService.findAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_companyService_add(t *testing.T) {
	type fields struct {
		repository Repository
	}
	type args struct {
		c Company
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Add a company",
			fields{repository: repoMock{
				AddFn: func(Company) error { return nil },
			}},
			args{},
			false},
		{"Throw exception on add",
			fields{repository: repoMock{
				AddFn: func(Company) error { return nil },
			}},
			args{},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := companyService{
				repository: tt.fields.repository,
			}
			if err := s.add(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("companyService.add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_isSemicolonSeparated(t *testing.T) {
	type args struct {
		t string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Is semicolon", args{"a;b;c"}, true},
		{"Is not semicolon", args{"a,b,c"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSemicolonSeparated(tt.args.t); got != tt.want {
				t.Errorf("isSemicolonSeparated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_companyService_loadWebsites(t *testing.T) {
	type fields struct {
		repository Repository
	}
	type args struct {
		f io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Load websites comma",
			fields{repoMock{}},
			args{strings.NewReader("a,b,c")},
			false},
		{"Load websites semicolon",
			fields{repoMock{}},
			args{strings.NewReader("a;b;c")},
			false},
		{"throws error",
			fields{repoMock{}},
			args{strings.NewReader("")},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := companyService{
				repository: tt.fields.repository,
			}
			if err := s.loadWebsites(tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("companyService.loadWebsites() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_companyService_InitDatabase(t *testing.T) {
	d1 := []byte("abc,asdf\n")
	ioutil.WriteFile("dat1", d1, 0644)
	type fields struct {
		repository Repository
	}
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Init database empty file",
			fields{
				repoMock{AddFn: func(Company) error { return nil }},
			},
			args{"dat1"},
			false},
		{"Init database with file",
			fields{
				repoMock{AddFn: func(Company) error { return nil }},
			},
			args{"dat1"},
			false},
		{"Init database with invalid file",
			fields{
				repoMock{AddFn: func(Company) error { return nil }},
			},
			args{"dat2"},
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := companyService{
				repository: tt.fields.repository,
			}
			if err := s.InitDatabase(tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("companyService.InitDatabase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_companyService_addByArray(t *testing.T) {
	type fields struct {
		repository Repository
	}
	type args struct {
		fields []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Add by array",
			fields{repoMock{AddFn: func(Company) error { return nil }}},
			args{[]string{"asdf", "123"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := companyService{
				repository: tt.fields.repository,
			}
			s.addByArray(tt.args.fields)
		})
	}
}

func Test_companyService_mergeDataByArray(t *testing.T) {
	type fields struct {
		repository Repository
	}
	type args struct {
		fields []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Should call mock repository",
			fields{repoMock{MergeWebsiteFn: func(Company) (*mgo.ChangeInfo, error) {
				return nil, nil
			}}},
			args{[]string{"adf", "12345", "site"}}},
		{"Should handler error",
			fields{repoMock{MergeWebsiteFn: func(Company) (*mgo.ChangeInfo, error) {
				return nil, errors.New("mock error")
			}}},
			args{[]string{"adf", "12345", "site"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := companyService{
				repository: tt.fields.repository,
			}
			s.mergeDataByArray(tt.args.fields)
		})
	}
}

func Test_companyService_iterateFileAndCall(t *testing.T) {
	type fields struct {
		repository Repository
	}
	type args struct {
		f io.Reader
		c csvLineHandler
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := companyService{
				repository: tt.fields.repository,
			}
			if err := s.iterateFileAndCall(tt.args.f, tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("companyService.iterateFileAndCall() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_companyService_validateAndParseToEntity(t *testing.T) {
	type fields struct {
		repository Repository
	}
	type args struct {
		fields []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Company
		wantErr bool
	}{
		{"should validate fields",
			fields{},
			args{[]string{"asdf", "12345", "site"}},
			Company{Name: "asdf", Zipcode: 12345, Website: "site"},
			false},
		{"should not validate fields when zipcode has not 5 digits",
			fields{},
			args{[]string{"asdf", "1235", "site"}},
			Company{},
			true},
		{"should not validate fields when missing fields",
			fields{},
			args{[]string{"asdf"}},
			Company{},
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := companyService{
				repository: tt.fields.repository,
			}
			got, err := s.validateAndParseToEntity(tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("companyService.validateAndParseToEntity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("companyService.validateAndParseToEntity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateZipcode(t *testing.T) {
	type args struct {
		zipcode string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateZipcode(tt.args.zipcode)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateZipcode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateZipcode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_companyService_findByNameAndZipCode(t *testing.T) {
	type fields struct {
		repository Repository
	}
	type args struct {
		name string
		zip  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Company
		wantErr bool
	}{
		{"Call repo mock find by name and zip",
			fields{repoMock{FindByNameAndZipFn: func(string, int64) (Company, error) {
				return Company{}, nil
			}}},
			args{"name", "12345"},
			Company{}, false},
		{"Throw error when zipinvalid",
			fields{repoMock{FindByNameAndZipFn: func(string, int64) (Company, error) {
				return Company{}, nil
			}}},
			args{"name", "123"},
			Company{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := companyService{
				repository: tt.fields.repository,
			}
			got, err := s.findByNameAndZipCode(tt.args.name, tt.args.zip)
			if (err != nil) != tt.wantErr {
				t.Errorf("companyService.findByNameAndZipCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("companyService.findByNameAndZipCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_check(t *testing.T) {
	type args struct {
		e error
	}
	tests := []struct {
		name string
		args args
	}{
		{"log when error", args{errors.New("mock error")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			check(tt.args.e)
		})
	}
}
