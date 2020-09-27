package vatpl

// VATError represents error with extra Permanent bool
type VATError struct {
	Message   string
	Permanent bool
}

func (v *VATError) Error() string {
	return v.Message
}

func (v *VATError) IsPermanent() bool {
	return v.Permanent
}

// NewVATError create error
func NewVATError(message string, permanent bool) *VATError {
	return &VATError{Message: message, Permanent: permanent}
}
