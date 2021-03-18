package vatpl

// VATResponse - JSON response object
type VATResponse struct {
	// error response
	Message string `json:"message"`
	Code    string `json:"code"`
	// normal response
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
		Entries []struct {
			Identifier string `json:"identifier"`
			Subjects   []struct {
				AuthorizedClerks        []interface{} `json:"authorizedClerks"`
				Regon                   string        `json:"regon"`
				RestorationDate         string        `json:"restorationDate"`
				WorkingAddress          string        `json:"workingAddress"`
				HasVirtualAccounts      bool          `json:"hasVirtualAccounts"`
				StatusVat               string        `json:"statusVat"`
				Krs                     string        `json:"krs"`
				RestorationBasis        string        `json:"restorationBasis"`
				AccountNumbers          []string      `json:"accountNumbers"`
				RegistrationDenialBasis string        `json:"registrationDenialBasis"`
				Nip                     string        `json:"nip"`
				RemovalDate             string        `json:"removalDate"`
				Partners                []interface{} `json:"partners"`
				Name                    string        `json:"name"`
				RegistrationLegalDate   string        `json:"registrationLegalDate"`
				RemovalBasis            string        `json:"removalBasis"`
				Pesel                   string        `json:"pesel"`
				Representatives         []struct {
					FirstName   string `json:"firstName"`
					LastName    string `json:"lastName"`
					Nip         string `json:"nip"`
					CompanyName string `json:"companyName"`
				} `json:"representatives"`
				ResidenceAddress       string `json:"residenceAddress"`
				RegistrationDenialDate string `json:"registrationDenialDate"`
			} `json:"subjects,omitempty"`
			Error struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"error,omitempty"`
		} `json:"entries"`
	} `json:"result"`
}
