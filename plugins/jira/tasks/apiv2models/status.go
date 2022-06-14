package apiv2models

type Status struct {
	Description string `json:"description"`
	IconURL     string `json:"iconUrl"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Scope       *struct {
		Type    string `json:"type"`
		Project struct {
			ID string `json:"id"`
		} `json:"project"`
	} `json:"scope"`
	Self           string `json:"self"`
	StatusCategory struct {
		ColorName string `json:"colorName"`
		ID        int    `json:"id"`
		Key       string `json:"key"`
		Name      string `json:"name"`
		Self      string `json:"self"`
	} `json:"statusCategory"`
	UntranslatedName string `json:"untranslatedName"`
}
