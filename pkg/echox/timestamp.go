package echox

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Timestamp time.Time

// UnmarshalParam echo api @see https://echo.labstack.com/guide/request
func (t *Timestamp) UnmarshalParam(src string) error {
	if src != "" {
		m, err := strconv.ParseInt(src, 10, 64)
		if err != nil {
			return err
		}

		ts := time.Unix(0, m*int64(time.Millisecond)).Local()
		*t = Timestamp(ts)
	}
	return nil
}

// MarshalJSON echo api json response
func (t *Timestamp) MarshalJSON() ([]byte, error) {
	if t != nil {
		ts := time.Time(*t)
		return []byte(fmt.Sprintf(`%d`, ts.UnixNano()/int64(time.Millisecond))), nil
	}
	return nil, nil
}

func (t *Timestamp) UnmarshalJSON(p []byte) error {
	data := string(p)
	if data == "null" {
		return nil
	}

	if p != nil {
		i, err := strconv.ParseInt(strings.Replace(data, `"`, "", -1), 10, 64)
		if err != nil {
			return err
		}

		*t = Timestamp(time.Unix(0, int64(time.Millisecond)*i))
	}
	return nil
}

// for sql log, print readable format
func (t Timestamp) String() string {
	ts := time.Time(t)
	return ts.Format("2006-01-02T15:04:05")
}

// insert into database conversion
func (t Timestamp) Value() (driver.Value, error) {
	return time.Time(t), nil
}

// read from database conversion
func (t *Timestamp) Scan(src interface{}) error {
	switch v := src.(type) {
	case *Timestamp:
		*t = *v
	case Timestamp:
		*t = v
	case time.Time:
		*t = Timestamp(v)
	case *time.Time:
		*t = Timestamp(*v)
	case string:
		v = strings.TrimSpace(v)
		return t.parse(v)
	case *string:
		str := strings.TrimSpace(*v)
		return t.parse(str)
	case int, int32, int64, uint, uint32, uint64:
		i := reflect.ValueOf(v).Int()
		*t = Timestamp(time.Unix(0, int64(time.Millisecond)*int64(i)))
	case *int, *int32, *int64, *uint, *uint32, *uint64:
		i := reflect.ValueOf(v).Elem().Int()
		*t = Timestamp(time.Unix(0, int64(time.Millisecond)*int64(i)))
	}
	return nil
}

func (t *Timestamp) parse(v string) error {
	switch {
	case regexp.MustCompile(`^\d+$`).MatchString(v):
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			*t = Timestamp(time.Unix(0, int64(time.Millisecond)*i))
		}
	case regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`).MatchString(v):
		if tt, err := time.Parse("2006-01-02 15:04:05", v); err == nil {
			*t = Timestamp(tt)
		}
	case regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}$`).MatchString(v):
		if tt, err := time.Parse("2006-01-02T15:04:05", v); err == nil {
			*t = Timestamp(tt)
		}
	case regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}[\+\-]\d{4}$`).MatchString(v):
		if tt, err := time.Parse("2006-01-02T15:04:05-0700", v); err == nil {
			*t = Timestamp(tt)
		}
	}
	return nil
}
