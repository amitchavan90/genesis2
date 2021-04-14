package sms

import "fmt"

// Mock - a mock SMS Client
type Mock struct {
}

// NewMock creates a new Mock SMS Client
func NewMock() *Mock {
	return &Mock{}
}

// SendMessage prints an SMS message to the console (mocking SMS)
func (m *Mock) SendMessage(mobileNumber string, content string) error {
	fmt.Printf("\n\n---SMS---\nto: %s\n\n%s\n\n---SMS---\n\n", mobileNumber, content)
	return nil
}

// SendToken prints a verification token to the console (mocking SMS)
func (m *Mock) SendToken(mobileNumber string, token string) error {
	content := fmt.Sprintf("%s is your Fetch verification code.", token)
	fmt.Printf("\n\n---SMS---\nto: %s\n\n%s\n\n---SMS---\n\n", mobileNumber, content)
	return nil
}
