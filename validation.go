// Package validation is a package to handle various validation procedures
//
//	Author: Elizalde G. Baguinon
//	Created: January 24, 2023
package validation

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	ssd "github.com/shopspring/decimal"
	"golang.org/x/exp/constraints"
)

const EMAIL_PATTERN string = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9]" +
	"(?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

type (
	NumericConstraint interface {
		constraints.Integer | constraints.Float
	}
	StringValidationOptions struct {
		Empty    bool // Allow empty string. Default: false, will raise an error if the string is empty
		Null     bool // Allow null. Default: false, will raise an error if the string is null
		Min      int  // Minimum length. Default: 0
		Max      int  // Maximum length. Default: 0
		NoSpaces bool // Do not allow spaces in the string. Default: false. Setting to true will raise an error if the string has spaces
		Extended []func(value *string) error
	}
	TimeValidationOptions struct {
		Null     bool       // Allow null. Default: false, will raise an error if the time is null
		Empty    bool       // Allow zero time Default: false, will raise an error if the time is zero
		Min      *time.Time // Minimum time. Default: nil
		Max      *time.Time // Maximum time. Default: nil
		DateOnly bool       // Compare dates only. Default: false
		Extended []func(value *time.Time) error
	}
	NumericValidationOptions[T NumericConstraint] struct {
		Null     bool // Allow null. Default: false, will raise an error if the time is null
		Empty    bool // Allow zero time Default: false, will raise an error if the time is zero
		Min      T    // Minimum time. Default: nil
		Max      T    // Maximum time. Default: nil
		Extended []func(value *T) error
	}
	DecimalValidationOptions struct {
		Null     bool         // Allow null. Default: false, will raise an error if the decimal is null
		Empty    bool         // Allow zero decimal. Default: false, will raise an error if the decimal is zero
		Min      *ssd.Decimal // Minimum decimal value. Default: nil
		Max      *ssd.Decimal // Maximum decimal value. Default: nil
		Extended []func(value *ssd.Decimal) error
	}
)

// ValidateEmail validates an e-mail address
func ValidateEmail(email *string) error {
	if email == nil || *email == "" {
		return fmt.Errorf("is an invalid email address")
	}
	re := regexp.MustCompile(EMAIL_PATTERN)
	if !re.MatchString(*email) {
		return fmt.Errorf("is an invalid email address")
	}
	return nil
}

// ValidateString validates an input string against the string validation options
func ValidateString(value *string, opts *StringValidationOptions) error {

	// If options were not set, this string is valid
	// If value is nil and the Null option is false, we raise an error
	// If value is empty and the Empty option is false, we raise an error
	if opts == nil {
		return nil
	}
	if value == nil {
		if !opts.Null {
			return fmt.Errorf("must be provided (nil)")
		}
		return nil
	}
	ln := len([]rune(*value))
	if ln == 0 {
		if !opts.Empty {
			return fmt.Errorf("must be provided (empty)")
		}
		return nil
	}
	if opts.Min > 0 && ln < opts.Min {
		return fmt.Errorf("is shorter than %d characters", opts.Min)
	}
	if opts.Max > 0 && ln > opts.Max {
		return fmt.Errorf("is longer than %d characters", opts.Max)
	}
	if opts.NoSpaces && strings.Contains(*value, " ") {
		return fmt.Errorf("contains spaces")
	}
	for _, f := range opts.Extended {
		if err := f(value); err != nil {
			return err
		}
	}

	return nil
}

// ValidateTime validates an input time against the time validation options
func ValidateTime(value *time.Time, opts *TimeValidationOptions) error {

	// If options were not set, this time is valid
	// If value is nil and the Null option is false, we raise an error
	if opts == nil {
		return nil
	}
	if value == nil {
		if !opts.Null {
			return fmt.Errorf("must be provided (nil)")
		}
		return nil
	}
	if value.IsZero() {
		if !opts.Empty {
			return fmt.Errorf("must be provided (empty)")
		}
		return nil
	}

	if opts.DateOnly {
		dv := *value
		*value = time.Date(dv.Year(), dv.Month(), dv.Day(), 0, 0, 0, 0, dv.Location())
		if opts.Min != nil {
			dc := opts.Min
			*opts.Min = time.Date(dc.Year(), dc.Month(), dc.Day(), 0, 0, 0, 0, dc.Location())
		}
		if opts.Max != nil {
			dc := opts.Max
			*opts.Max = time.Date(dc.Year(), dc.Month(), dc.Day(), 0, 0, 0, 0, dc.Location())
		}
	}

	if opts.Min != nil && value.Before(*opts.Min) {
		return fmt.Errorf("is earlier than %s minimum time", opts.Min)
	}

	if opts.Max != nil && value.After(*opts.Max) {
		return fmt.Errorf("is later than %s maximum time", opts.Max)
	}

	for _, f := range opts.Extended {
		if err := f(value); err != nil {
			return err
		}
	}

	return nil
}

// ValidateNumeric validates a numeric input against numeric validation options
//
// Currently supported data types are:
//   - constraints.Integer (Signed | Unsigned)
//   - constraints.Float (~float32 | ~float64)
//
// This function requires version 1.18+
func ValidateNumeric[T NumericConstraint](value *T, opts *NumericValidationOptions[T]) error {

	// If options were not set, this time is valid
	// If value is nil and the Null option is false, we raise an error
	if opts == nil {
		return nil
	}
	if value == nil {
		if !opts.Null {
			return fmt.Errorf("must be provided (nil)")
		}
		return nil
	}
	if *value == 0 {
		if !opts.Empty {
			return fmt.Errorf("must be provided (empty)")
		}
	}
	if opts.Min > 0 && *value < opts.Min {
		return fmt.Errorf("is lesser than %v minimum value", opts.Min)
	}
	if opts.Max > 0 && *value > opts.Max {
		return fmt.Errorf("is greater than %v maximum value", opts.Max)
	}
	for _, f := range opts.Extended {
		if err := f(value); err != nil {
			return err
		}
	}
	return nil
}

// ValidateDecimal validates a decimal input against decimal validation options
func ValidateDecimal(value *ssd.Decimal, opts *DecimalValidationOptions) error {

	// If options were not set, this decimal is valid
	// If value is nil and the Null option is false, we raise an error
	if opts == nil {
		return nil
	}
	if value == nil {
		if !opts.Null {
			return fmt.Errorf("must be provided (nil)")
		}
		return nil
	}
	if value.IsZero() {
		if !opts.Empty {
			return fmt.Errorf("must be provided (empty)")
		}
	}
	zero := ssd.NewFromInt(0)
	if opts.Min != nil && opts.Min.GreaterThan(zero) && value.LessThan(*opts.Min) {
		return fmt.Errorf("is lesser than %v minimum value", *opts.Min)
	}
	if opts.Max != nil && opts.Max.GreaterThan(zero) && value.GreaterThan(*opts.Max) {
		return fmt.Errorf("is greater than %v maximum value", *opts.Max)
	}
	for _, f := range opts.Extended {
		if err := f(value); err != nil {
			return err
		}
	}
	return nil
}
