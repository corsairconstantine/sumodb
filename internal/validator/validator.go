package validator

import (
	"regexp"
	"strconv"
	"strings"
)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}

func ValidTournament(tournament string) bool {
	if strings.TrimSpace(tournament) == "" {
		return false
	}
	fields := strings.Fields(tournament)
	year, err := strconv.Atoi(fields[0])
	if err != nil {
		return false
	}

	if year < 1900 || year > 2050 {
		return false
	}

	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Sep", "Oct", "Nov", "Dec"}
	if !In(fields[len(fields)-1], months...) {
		return false
	}

	return true
}

func ValidDay(day string) bool {
	if day == "Playoff" {
		return true
	}
	isDay_ := strings.HasPrefix(day, "Day ")
	if !isDay_ {
		return false
	}
	n := strings.TrimPrefix(day, "Day ")
	i, err := strconv.Atoi(n)
	if err != nil {
		return false
	}
	if i > 15 {
		return false
	}

	return true
}
