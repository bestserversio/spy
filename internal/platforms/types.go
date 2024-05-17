package platforms

type Platform struct {
	Flags []string `json:"flags"`

	Banner string `json:"banner"`
	Icon   string `json:"icon"`

	Url         string `json:"url"`
	Name        string `json:"name"`
	NameShort   string `json:"nameShort"`
	Description string `json:"description"`

	JsInternal string `json:"jsInternal"`
	JsExternal string `json:"jsExternal"`
}
