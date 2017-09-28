package md

import (
	"reflect"
	"testing"
	"time"
)

func TestValidateFirstLine(t *testing.T) {
	const (
		first_line     = "@(2017/09/23)[hoge, fuga]"
		not_close_line = "@(2017/09/23)[hoge, fuga"
		invalid_line   = "# 2017/09/23 hoge"
	)
	if err := validateFirstLine(first_line); err != nil {
		t.Errorf("This function recognize that '%s' is invalid", first_line)
	}
	if err := validateFirstLine(not_close_line); err == nil {
		t.Errorf("This function recognize that '%s' is valid", not_close_line)
	}
	if err := validateFirstLine(invalid_line); err == nil {
		t.Errorf("This function recognize that '%s' is valid", invalid_line)
	}
}

func TestParseDate(t *testing.T) {
	const (
		d_format            = "2006/01/02 Mon 15:04:05"
		first_line          = "@(2017/09/23)[hoge, fuga]"
		invalid_date_format = "@(20017/20/30)[hoge, fuga]"
	)
	test_date := time.Date(2017, 9, 23, 0, 0, 0, 0, time.Local)
	if date, err := parseDate(first_line); err != nil {
		t.Errorf("Could not parse '%s'.\n%s\n", first_line, err)
	} else {
		if !test_date.Equal(date) {
			t.Errorf("Expect:\t%s\nResult:\t%s\n", test_date.Format(d_format), date.Format(d_format))
		}
	}
	if invalid_date, err := parseDate(invalid_date_format); err == nil {
		t.Errorf("Causing error was expected behavior.\nResult:\t%s\n", invalid_date.Format(date_format))
	}
}

func TestParseTags(t *testing.T) {
	const (
		first_line = "@(2017/09/23)[hoge, fuga]"
	)
	valid_tags := []string{"hoge", "fuga"}
	if tags := parseTags(first_line); !reflect.DeepEqual(tags, valid_tags) {
		t.Errorf("Expect:\t%s\nResult:\t%s\n", valid_tags, tags)
	}
}
