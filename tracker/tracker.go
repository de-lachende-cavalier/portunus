// tracker contains all the functionality that allows us to correctly track files in various directories,
// in our case that means tracking files in $HOME/.ssh, but the package has been kept generic enough that
// it could be applied to any directory in the fs
package tracker

import (
	"errors"
	"os"
	"path/filepath"
	"time"
)

// struct meant to store values output by time.Date()
type Date struct {
	year  int
	month time.Month
	day   int
}

func newDate(t time.Time) Date {
	y, m, d := t.Date()

	return Date{
		year:  y,
		month: m,
		day:   d,
	}
}

// TODO redo the whole thing after completing locksmith functionality
// Reads the contents of the directory with path=dirPath and checks the last modified time
// of the files within, comparing it to the date=expiration specified.
//
// If there were no errors in its execution, it return a hashmap containing the estimated
// seconds till expiration of the various files within the directory.
//
// expirationDate is in the default time.Time Go format (YEAR-MONTH-DAY).
func EstimateExpiration(dirPath string, expiration string) (map[string]Date, error) {
	var nameDateMap map[string]Date

	d, err := os.Open(filepath.Dir(dirPath))
	if err != nil {
		return nil, errors.New(err.Error())
	}

	inodes, err := d.Readdir(-1)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	d.Close()

	for _, inode := range inodes {
		date := newDate(inode.ModTime())
		nameDateMap[inode.Name()] = date
	}

	date, err := time.Parse("2006-01-02", expiration)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	expirationDate := Date{
		year:  date.Year(),
		month: date.Month(),
		day:   date.Day(),
	}

	for fileName, modDate := range nameDateMap {
	}
}
