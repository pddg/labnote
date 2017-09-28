package md

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
	"io/ioutil"
)

const (
	date_format       = "2006/01/02"
	first_line_format = "@({year}/(month}/{day})[tags]"
	date_confirm_msg  = "This file does not have an information of date. Use today's date? (y/n):"
	date_input_msg    = "Please input date (YYYY/MM/DD):"
)

var (
	first_line_regex = regexp.MustCompile("@\\([0-9]{4}/[0-9]{2}/[0-9]{2}\\)(\\[(.+,?)*])?$")
	date_regex = regexp.MustCompile("([0-9]+/?)+")
	tags_regex = regexp.MustCompile("\\[(.+,?)]$")
)

func getInput(message string) string {
	fmt.Print(message)
	var text string
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		text = scanner.Text()
	}
	return text
}

func parseDate(line string) (time.Time, error) {
	line = strings.Trim(line, " \n")
	if d_str := date_regex.FindString(line); len(d_str) != 0 {
		date, err := time.ParseInLocation(date_format, d_str, time.Local)
		if err != nil {
			return time.Now(), err
		}
		return date, nil
	}
	return time.Now(), nil
}

func parseTags(line string) []string {
	var tags []string
	line = strings.Trim(line, " \n")
	if tags_str := tags_regex.FindString(line); len(tags_str) != 0 {
		// 両端の[]を削除
		tags_str := strings.Trim(tags_str, " \n")[1 : len(tags_str)-1]
		tags_list := strings.Split(tags_str, ",")
		for _, tag := range tags_list {
			tags = append(tags, strings.Trim(tag, " "))
		}
	}
	return tags
}

func validateFirstLine(line string) error {
	if !first_line_regex.Match([]byte(line)) {
		return fmt.Errorf("'%s' is invalid. '%s' was expected\n", line, first_line_format)
	}
	return nil
}

type Markdown struct {
	Path      string
	FirstLine string
	Date      time.Time
	Tags      []string
	HasDate   bool
	HasTags   bool
}

func (m *Markdown) SetPath(filename string) error {
	var (
		fp        *os.File
		err       error
		firstline string
	)
	// init this struct
	m.HasDate = false
	m.HasTags = false
	// open
	if fp, err = os.Open(filename); err != nil {
		panic(err)
	}
	defer fp.Close()
	// Get first line
	reader := bufio.NewScanner(fp)
	for reader.Scan() {
		firstline = reader.Text()
		if utf8.RuneCountInString(firstline) != 0 {
			break
		}
	}
	// validate
	if err = validateFirstLine(firstline); err != nil {
		return err
	}
	// set variables
	m.Path = filename
	m.FirstLine = firstline
	return nil
}

func (m *Markdown) SetDate() bool {
	var (
		date       time.Time
		err        error
		user_input = false
		manual = false
	)
	if date, err = parseDate(m.FirstLine); err != nil {
		for {
			switch strings.ToLower(getInput(date_confirm_msg)) {
			case "y":
				user_input = true
			case "n":
				for {
					input := getInput(date_input_msg)
					if input_date, err := time.ParseInLocation(date_format, input, time.Local); err != nil {
						fmt.Printf("Your input ('%s') is invalid. ", input)
					} else {
						date = input_date
						break
					}
				}
				user_input = true
			}
			if user_input {
				manual = true
				break
			}
		}
	}
	m.Date = date
	m.HasDate = true
	return manual
}

func (m *Markdown) SetTags() {
	m.Tags = parseTags(m.FirstLine)
	m.HasTags = true
}

func (m *Markdown) ModFirstLine() error {
	var (
		err error
		body []byte
		f_line string
		f_line_bytes []byte
	)
	if body, err = ioutil.ReadFile(m.Path); err != nil {
		return err
	}
	if !m.HasDate {
		return fmt.Errorf("'%s' has no date information", m.Path)
	}
	f_line = "@(" + m.Date.Format(date_format) + ")"
	if m.HasTags {
		f_line += "[" + strings.Join(m.Tags, ",") + "]"
	}
	f_line += "\n"
	for _, b := range []byte(f_line) {
		f_line_bytes = append(f_line_bytes, b)
	}
	for _, b := range body {
		f_line_bytes = append(f_line_bytes, b)
	}
	if err = ioutil.WriteFile(m.Path, f_line_bytes, 0666); err != nil {
		return err
	}
	m.FirstLine = f_line
	return nil
}
