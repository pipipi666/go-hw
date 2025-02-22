package hw09structvalidator

import (
	"errors"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

type SimpleChecker func(ruleKey string, ruleValue string, value reflect.Value) (validationErr error, err error)

var (
	ErrStrNotEqualLen = errors.New("invalid number of symbols")
	ErrStrRegMatch    = errors.New("doesn't match regexp")
	ErrStrNotInArr    = errors.New("not in array of strings")
	ErrIntNotInArr    = errors.New("not in array of integers")
	ErrIntMax         = errors.New("value exceeds max")
	ErrIntMin         = errors.New("value is less than min")
)

var (
	ErrReg       = errors.New("invalid regexp")
	ErrAtoi      = errors.New("invalid int value")
	ErrType      = errors.New("invalid type")
	ErrValidator = errors.New("invalid validator")
)

func (v ValidationErrors) Error() string {
	res := []string{}

	for _, s := range v {
		res = append(res, s.Err.Error())
	}

	return strings.Join(res, ", ")
}

func checkStr(ruleKey string, ruleValue string, value reflect.Value) (error, error) {
	parsedValue := value.String()

	switch ruleKey {
	case "len":
		strLen, err := strconv.Atoi(ruleValue)
		if err != nil {
			return nil, ErrAtoi
		}

		if len(parsedValue) != strLen {
			return ErrStrNotEqualLen, nil
		}
	case "regexp":
		re, err := regexp.Compile(ruleValue)
		if err != nil {
			return nil, ErrReg
		}

		if !re.Match([]byte(parsedValue)) {
			return ErrStrRegMatch, nil
		}
	case "in":
		values := strings.Split(ruleValue, ",")

		if !slices.Contains(values, parsedValue) {
			return ErrStrNotInArr, nil
		}
	default:
		return nil, ErrValidator
	}

	return nil, nil
}

func checkInt(ruleKey string, ruleValue string, value reflect.Value) (error, error) {
	if !value.CanInt() {
		return nil, ErrType
	}

	parsedValue := int(value.Int())

	switch ruleKey {
	case "min":
		minVal, err := strconv.Atoi(ruleValue)
		if err != nil {
			return nil, ErrAtoi
		}

		if parsedValue < minVal {
			return ErrIntMin, nil
		}
	case "max":
		maxVal, err := strconv.Atoi(ruleValue)
		if err != nil {
			return nil, ErrAtoi
		}

		if parsedValue > maxVal {
			return ErrIntMax, nil
		}
	case "in":
		values := strings.Split(ruleValue, ",")

		for _, v := range values {
			parsedV, err := strconv.Atoi(v)
			if err != nil {
				return nil, ErrAtoi
			}

			if parsedV == parsedValue {
				return nil, nil
			}
		}

		return ErrIntNotInArr, nil
	default:
		return nil, ErrValidator
	}

	return nil, nil
}

func checkType(value reflect.Value, validators []string) ([]error, error) {
	//nolint:exhaustive
	switch value.Kind() {
	case reflect.String:
		return checkSimpleType(value, validators, checkStr)
	case reflect.Int:
		return checkSimpleType(value, validators, checkInt)
	case reflect.Slice:
		simpleTypeErrors := []error{}

		for i := 0; i < value.Len(); i++ {
			validErr, err := checkType(value.Index(i), validators)
			if err != nil {
				return nil, err
			}

			simpleTypeErrors = append(simpleTypeErrors, validErr...)
		}

		return simpleTypeErrors, nil
	}

	return nil, ErrType
}

func checkSimpleType(value reflect.Value, validators []string, cb SimpleChecker) ([]error, error) {
	simpleTypeErrors := []error{}

	for _, validator := range validators {
		entry := strings.Split(validator, ":")

		if len(entry) != 2 {
			return nil, ErrValidator
		}

		ruleKey := entry[0]
		ruleValue := entry[1]

		validErr, err := cb(ruleKey, ruleValue, value)
		if err != nil {
			return nil, err
		}

		if validErr != nil {
			simpleTypeErrors = append(simpleTypeErrors, validErr)
		}
	}

	return simpleTypeErrors, nil
}

func Validate(v interface{}) (ValidationErrors, error) {
	st := reflect.TypeOf(v)

	if st.Kind() != reflect.Struct {
		return nil, ErrType
	}

	resErrors := ValidationErrors{}

	for i := range st.NumField() {
		tag := st.Field(i).Tag.Get("validate")

		if tag == "" {
			continue
		}

		key := st.Field(i).Name
		value := reflect.Indirect(reflect.ValueOf(v)).FieldByName(key)
		validators := strings.Split(tag, "|")
		errValid, err := checkType(value, validators)
		if err != nil {
			return nil, err
		}

		for _, v := range errValid {
			resErrors = append(resErrors, ValidationError{Field: key, Err: v})
		}
	}

	return resErrors, nil
}
