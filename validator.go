package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// Name of the struct tag used in examples.
const tagName = "validate"

// Regular expression to validate email address.
var mailRe = regexp.MustCompile(`\A[\w+\-.]+@[a-z\d\-]+(\.[a-z]+)*\.[a-z]+\z`)

// Generic data validator.
type Validator interface {
	// Validate method performs validation and returns result and optional error.
	Validate(interface{}) (bool, error)
}

// DefaultValidator does not perform any validations.
type DefaultValidator struct {
}

func (v DefaultValidator) Validate(val interface{}) (bool, error) {
	return true, nil
}

// StringValidator validates string presence and/or its length.
type StringValidator struct {
	Min int
	Max int
}

func (v StringValidator) Validate(val interface{}) (bool, error) {
	l := len([]rune(val.(string)))

	if l < v.Min {
		if v.Min == 1 {
			return false, fmt.Errorf("不能为空")
		}
		return false, fmt.Errorf("至少%v位", v.Min)
	}

	if v.Max >= v.Min && l > v.Max {
		return false, fmt.Errorf("最大长度为%v位", v.Max)
	}

	return true, nil
}

// StringValidator validates string presence and/or its length.
type RegexValidator struct {
	Pattern string
}

func (v RegexValidator) Validate(val interface{}) (bool, error) {
	l := val.(string)
	if l != "" {
		if matched, err := regexp.MatchString(v.Pattern, l); err != nil {
			return false, err
		} else {
			if matched {
				return true, nil
			} else {
				return false, fmt.Errorf("不符合规则")
			}
		}
	}
	return true, nil
}

// NumberValidator performs numerical value validation.
// Its limited to int type for simplicity.
type NumberValidator struct {
	Min float64
	Max float64
}

func (v NumberValidator) Validate(val interface{}) (bool, error) {
	var num float64
	switch v := val.(type) {
	case float64:
		num = v
	case int64:
		num = float64(v)
	case int32:
		num = float64(v)
	case int:
		num = float64(v)
	}

	if num < v.Min {
		return false, fmt.Errorf("必须大于%v", v.Min)
	}

	if v.Max >= v.Min && num > v.Max {
		return false, fmt.Errorf("不能超过%.2f", v.Max)
	}

	return true, nil
}

// EmailValidator checks if string is a valid email address.
type EmailValidator struct {
}

func (v EmailValidator) Validate(val interface{}) (bool, error) {
	if !mailRe.MatchString(val.(string)) {
		return false, fmt.Errorf("is not a valid email address")
	}
	return true, nil
}

// Returns validator struct corresponding to validation type
func getValidatorFromTag(tag string) Validator {
	args := strings.Split(tag, ",")
	switch args[0] {
	case "number":
		validator := NumberValidator{}
		arg1 := args[1:]
		if len(arg1) == 2 {
			fmt.Sscanf(strings.Join(args[1:], ","), "min=%g,max=%g", &validator.Min, &validator.Max)
		} else if len(arg1) == 1 {
			if strings.Contains(arg1[0], "min") {
				fmt.Sscanf(strings.Join(args[1:], ","), "min=%g", &validator.Min)
				validator.Max = 0
			} else if strings.Contains(arg1[0], "max") {
				fmt.Sscanf(strings.Join(args[1:], ","), "max=%g", &validator.Max)
				validator.Min = 0
			} else {
				fmt.Printf("Error validate formart")
			}
		} else {
			fmt.Printf("Error validate formart")
		}
		return validator
	case "string":
		validator := StringValidator{}
		arg1 := args[1:]

		if len(arg1) == 2 {
			fmt.Sscanf(strings.Join(args[1:], ","), "min=%d,max=%d", &validator.Min, &validator.Max)
		} else if len(arg1) == 1 {
			if strings.Contains(arg1[0], "min") {
				fmt.Sscanf(strings.Join(args[1:], ","), "min=%d", &validator.Min)
				validator.Max = 0
			} else if strings.Contains(arg1[0], "max") {
				fmt.Sscanf(strings.Join(args[1:], ","), "max=%d", &validator.Max)
				validator.Min = 0
			} else {
				fmt.Printf("Error validate formart")
			}
		} else {
			fmt.Printf("Error validate formart")
		}
		return validator
	case "email":
		return EmailValidator{}
	case "regex":
		validator := RegexValidator{}
		fmt.Sscanf(strings.Join(args[1:], ","), "pattern=%s", &validator.Pattern)
		return validator
	}
	return DefaultValidator{}
}

// Performs actual data validation using validator definitions on the struct
func ValidateStruct(s interface{}) []error {
	errs := []error{}

	// ValueOf returns a Value representing the run-time data
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		// Get the field tag value
		tag := v.Type().Field(i).Tag.Get(tagName)

		// Skip if tag is not defined or ignored
		if tag == "" || tag == "-" {
			continue
		}

		// Get a validator that corresponds to a tag
		validator := getValidatorFromTag(tag)

		// Perform validation
		valid, err := validator.Validate(v.Field(i).Interface())

		// Append error to results
		if !valid && err != nil {
			name := ""
			tagjson := v.Type().Field(i).Tag.Get("json")
			tagMsg := v.Type().Field(i).Tag.Get("msg")
			if tagMsg != "" {
				name = tagMsg
			} else if tagjson != "" {
				name = tagjson
			} else {
				name = v.Type().Field(i).Name
			}
			errs = append(errs, fmt.Errorf("%s%s", name, err.Error()))
		}
	}

	return errs
}
