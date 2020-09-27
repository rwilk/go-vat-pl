package vatpl

// VATResponse - JSON response object
type VATResponse struct {
	Result struct {
		RequestDateTime string `json:"requestDateTime"`
		RequestID       string `json:"requestId"`
		Subject         struct {
			AccountNumbers          []string      `json:"accountNumbers"`
			AuthorizedClerks        []interface{} `json:"authorizedClerks"`
			HasVirtualAccounts      bool          `json:"hasVirtualAccounts"`
			Krs                     interface{}   `json:"krs"`
			Name                    string        `json:"name"`
			Nip                     string        `json:"nip"`
			Partners                []interface{} `json:"partners"`
			Pesel                   interface{}   `json:"pesel"`
			RegistrationDenialBasis interface{}   `json:"registrationDenialBasis"`
			RegistrationDenialDate  interface{}   `json:"registrationDenialDate"`
			RegistrationLegalDate   string        `json:"registrationLegalDate"`
			Regon                   string        `json:"regon"`
			RemovalBasis            interface{}   `json:"removalBasis"`
			RemovalDate             interface{}   `json:"removalDate"`
			Representatives         []interface{} `json:"representatives"`
			ResidenceAddress        string        `json:"residenceAddress"`
			RestorationBasis        interface{}   `json:"restorationBasis"`
			RestorationDate         interface{}   `json:"restorationDate"`
			StatusVat               string        `json:"statusVat"`
			WorkingAddress          interface{}   `json:"workingAddress"`
		} `json:"subject"`
	} `json:"result"`
}
