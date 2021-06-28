package usingpin

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// CheckValidBDMobileNumber checks and match if the given mobile value is valid or not.
func CheckValidBDMobileNumber(value interface{}) error {
	value, isNil := validation.Indirect(value)
	if isNil {
		return nil
	}

	re := regexp.MustCompile("(^([+]{1}[8]{2}|0088)?(01){1}[3-9]{1}\\d{8})$")

	isString, str, isBytes, bs := validation.StringOrBytes(value)
	if isString && (str == "" || re.MatchString(str)) {
		return nil
	} else if isBytes && (len(bs) == 0 || re.Match(bs)) {
		return nil
	}
	return validation.NewError("8001", "Invalid mobile phone number format.")
}
