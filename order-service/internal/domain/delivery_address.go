package domain

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrInvalidAddressID    = errors.New("invalid address ID")
	ErrInvalidFullName     = errors.New("invalid full name")
	ErrInvalidStreetAddress = errors.New("invalid street address")
	ErrInvalidCity         = errors.New("invalid city")
	ErrInvalidState        = errors.New("invalid state")
	ErrInvalidPostalCode   = errors.New("invalid postal code")
	ErrInvalidCountry      = errors.New("invalid country")
	ErrInvalidPhone        = errors.New("invalid phone number")
)

type DeliveryAddress struct {
	ID            string
	UserID        string
	FullName      string
	StreetAddress string
	Apartment     string // Optional
	City          string
	State         string
	PostalCode    string
	Country       string
	Phone         string
	IsDefault     bool
}

// Regular expressions for validation
var (
	phoneRegex    = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	postalRegex   = regexp.MustCompile(`^[0-9A-Z]{3,10}$`)
	nameRegex     = regexp.MustCompile(`^[a-zA-Z\s\-']{2,100}$`)
)

func NewDeliveryAddress(
	userID string,
	fullName string,
	streetAddress string,
	apartment string,
	city string,
	state string,
	postalCode string,
	country string,
	phone string,
	isDefault bool,
) (*DeliveryAddress, error) {
	address := &DeliveryAddress{
		UserID:        userID,
		FullName:      strings.TrimSpace(fullName),
		StreetAddress: strings.TrimSpace(streetAddress),
		Apartment:     strings.TrimSpace(apartment),
		City:          strings.TrimSpace(city),
		State:         strings.TrimSpace(state),
		PostalCode:    strings.TrimSpace(postalCode),
		Country:       strings.TrimSpace(country),
		Phone:         strings.TrimSpace(phone),
		IsDefault:     isDefault,
	}

	if err := address.Validate(); err != nil {
		return nil, err
	}

	return address, nil
}

func (a *DeliveryAddress) Validate() error {
	if a.UserID == "" {
		return ErrInvalidUserID
	}

	if !nameRegex.MatchString(a.FullName) {
		return ErrInvalidFullName
	}

	if len(a.StreetAddress) < 5 {
		return ErrInvalidStreetAddress
	}

	if len(a.City) < 2 {
		return ErrInvalidCity
	}

	if len(a.State) < 2 {
		return ErrInvalidState
	}

	if !postalRegex.MatchString(a.PostalCode) {
		return ErrInvalidPostalCode
	}

	if len(a.Country) < 2 {
		return ErrInvalidCountry
	}

	if !phoneRegex.MatchString(a.Phone) {
		return ErrInvalidPhone
	}

	return nil
}

func (a *DeliveryAddress) Update(
	fullName string,
	streetAddress string,
	apartment string,
	city string,
	state string,
	postalCode string,
	country string,
	phone string,
) error {
	a.FullName = strings.TrimSpace(fullName)
	a.StreetAddress = strings.TrimSpace(streetAddress)
	a.Apartment = strings.TrimSpace(apartment)
	a.City = strings.TrimSpace(city)
	a.State = strings.TrimSpace(state)
	a.PostalCode = strings.TrimSpace(postalCode)
	a.Country = strings.TrimSpace(country)
	a.Phone = strings.TrimSpace(phone)

	return a.Validate()
}

func (a *DeliveryAddress) SetDefault(isDefault bool) {
	a.IsDefault = isDefault
}

func (a *DeliveryAddress) FormatFull() string {
	parts := []string{
		a.StreetAddress,
	}

	if a.Apartment != "" {
		parts = append(parts, a.Apartment)
	}

	parts = append(parts,
		a.City,
		a.State,
		a.PostalCode,
		a.Country,
	)

	return strings.Join(parts, ", ")
} 