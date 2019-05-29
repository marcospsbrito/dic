package company

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Company entity
type Company struct {
	ID      bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty" example:"12345"`
	Name    string        `json:"name" example:"Company Name"`
	Zipcode int64         `json:"Zipcode,omitempty" example:"123"`
	Website string        `json:"website,omitempty" example:"1" example:"http://localhost"`
}

// Repository interface difines necessary methods
type Repository interface {
	FindAll() ([]Company, error)
	FindByNameAndZip(string, int64) (Company, error)
	Add(Company) error
	MergeWebsite(Company) (*mgo.ChangeInfo, error)
}

type companyRepository struct {
	companies *mgo.Collection
}

// NewRepository function returns a Repository impl
func NewRepository(db *mgo.Database) Repository {
	if db == nil {
		return nil
	}
	db.C("Company").EnsureIndexKey("$text:name")
	return companyRepository{db.C("Company")}
}

func (r companyRepository) FindAll() ([]Company, error) {
	var results []Company
	err := r.companies.Find(nil).All(&results)
	return results, err
}

func (r companyRepository) FindByNameAndZip(name string, zipcode int64) (Company, error) {
	var result Company
	query := getCompanyNameAndZipQuery(name, zipcode)
	err := r.companies.Find(query).One(&result)
	return result, err
}

func (r companyRepository) Add(c Company) error {
	count, err := r.companies.Find(getCompanyNameAndZipQuery(c.Name, c.Zipcode)).Count()
	if err != nil || count > 0 {
		return err
	}
	return r.companies.Insert(c)
}

func (r companyRepository) MergeWebsite(c Company) (*mgo.ChangeInfo, error) {
	query := getCompanyNameOrZipQuery(c.Name, c.Zipcode)
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"website": c.Website}},
		ReturnNew: true,
	}
	return r.companies.Find(query).Apply(change, &c)
}

func getCompanyNameAndZipQuery(name string, zipcode int64) bson.M {
	return bson.M{"$and": []bson.M{
		{"$text": bson.M{"$search": name}},
		{"zipcode": zipcode}}}
}

func getCompanyNameOrZipQuery(name string, zipcode int64) bson.M {
	return bson.M{"$or": []bson.M{
		{"name": name},
		{"zipcode": zipcode}}}
}
