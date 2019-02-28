package company

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/apex/log"
)

// Service interface define methods of service
type Service interface {
	findAll() ([]Company, error)
	findByNameAndZipCode(string, string) (Company, error)
	add(Company) error
	InitDatabase(string) error
	loadWebsites(f io.Reader) error
}

type csvLineHandler func([]string)

// companyService struct
type companyService struct {
	repository Repository
}

// NewService returns new Service
func NewService(r Repository) Service {
	return companyService{r}
}

func (s companyService) findAll() ([]Company, error) {
	return s.repository.FindAll()
}

func (s companyService) add(c Company) error {
	return s.repository.Add(c)
}

func isSemicolonSeparated(t string) bool {
	result := strings.Contains(t, ";") && !strings.Contains(t, ",")
	if result {
		log.Debug("Separated by semicolon")
	} else {
		log.Debug("Separated by comma")
	}
	return result
}

func (s companyService) loadWebsites(f io.Reader) error {
	log.Debug("calls [loadWebsites] service")
	return s.iterateFileAndCall(f, s.mergeDataByArray)
}

func (s companyService) InitDatabase(file string) error {
	log.Debug("Start database setup")
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return s.iterateFileAndCall(f, s.addByArray)
}

func (s companyService) addByArray(fields []string) {
	zipcode, _ := strconv.ParseInt(fields[1], 10, 0)
	c := Company{Name: fields[0], Zipcode: zipcode}
	s.add(c)
}

func (s companyService) mergeDataByArray(fields []string) {
	c, err := s.validateAndParseToEntity(fields)
	if err != nil {
		log.WithError(err).Error("Cannot update values")
		return
	}
	info, err := s.repository.MergeWebsite(c)
	if err != nil {
		log.WithError(err).Error("Cannot update values")
		return
	}
	log.Debug("Changed info")
	log.Debug(fmt.Sprint(info))

	log.Error("No fields avaliable to update.")

}

func (s companyService) iterateFileAndCall(f io.Reader, c csvLineHandler) error {
	t, err := ioutil.ReadAll(f)
	check(err)
	text := string(t[:])
	reader := csv.NewReader(strings.NewReader(text))
	if isSemicolonSeparated(text) {
		reader.Comma = ';'
		reader.Comment = '#'
	}
	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		c(row)
	}
}

func (s companyService) validateAndParseToEntity(fields []string) (Company, error) {
	var c Company
	if len(fields) < 3 {
		return c, errors.New("Missing fields")
	}
	zipcode, err := validateZipcode(fields[1])
	if err != nil {
		return c, err
	}
	return Company{Name: fields[0], Zipcode: zipcode, Website: fields[2]}, nil
}

func validateZipcode(zipcode string) (int64, error) {
	if len(zipcode) != 5 {
		return 0, errors.New("Invalid Zipcode lenght")
	}
	return strconv.ParseInt(zipcode, 10, 0)
}

func (s companyService) findByNameAndZipCode(name string, zip string) (Company, error) {
	zipcode, err := validateZipcode(zip)
	if err != nil {
		return Company{}, err
	}
	return s.repository.FindByNameAndZip(name, zipcode)
}

func check(e error) {
	if e != nil {
		log.WithError(e).Error("failed with error - ")
	}
}
