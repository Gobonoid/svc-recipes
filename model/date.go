package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

//DateTime structure is specific for unmarshaling given CSV as it doesn't follow any standard format
type DateTime struct {
	time.Time
}

// Convert the CSV string as internal date
func (date *DateTime) UnmarshalCSV(csvTime string) (err error) {
	return date.Unmarshal(csvTime)
}

// Convert the CSV string as internal date
func (date *DateTime) UnmarshalJSON(jsonTime []byte) (err error) {
	//You have to trim first " and last "
	s := fmt.Sprintf("%s", jsonTime)
	return date.Unmarshal(strings.TrimLeft(strings.TrimRight(s, `"`), `"`))
}

func (date *DateTime) Unmarshal(timeToUnmarshal string) (err error) {
	date.Time, err = time.Parse("02/01/2006 15:04:05", timeToUnmarshal)
	if err != nil {
		return errors.Wrapf(err, "Can't unmarshal time: %s to DateTime", timeToUnmarshal)
	}
	return nil
}
