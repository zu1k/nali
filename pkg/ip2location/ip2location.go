package ip2location

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/ip2location/ip2location-go/v9"
)

// IP2Location
type IP2Location struct {
	db *ip2location.DB
}

// new IP2Location from database file
func NewIP2Location(filePath string) (*IP2Location, error) {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，请自行下载 IP2Location 库，并保存在", filePath)
		return nil, err
	} else {
		db, err := ip2location.OpenDB(filePath)

		if err != nil {
			log.Fatal(err)
		}
		return &IP2Location{db: db}, nil
	}
}

func (x IP2Location) Find(query string, params ...string) (result fmt.Stringer, err error) {
	ip := net.ParseIP(query)
	if ip == nil {
		return nil, errors.New("Query should be valid IP")
	}
	record, err := x.db.Get_all(ip.String())

	if err != nil {
		return
	}

	result = Result{
		Country: record.Country_long,
		Region:  record.Region,
		City:    record.City,
	}
	return
}

type Result struct {
	Country string
	Region  string
	City    string
}

func (r Result) String() string {
	return fmt.Sprintf("%s %s %s", r.Country, r.Region, r.City)
}
