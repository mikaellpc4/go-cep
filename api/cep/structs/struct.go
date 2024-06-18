package structs

var CepFieldsToJSONMap = map[string]string{
	"ZipCode":     "zipCode",
	"PublicPlace": "publicPlace",
	"Complement":  "complement",
	"District":    "district",
	"Place":       "place",
	"Uf":          "uf",
	"IbgeCode":    "ibgeCode",
}

type Cep struct {
	ZipCode     string `json:"zipCode"      validate:"required"`
	PublicPlace string `json:"publicPlace"  validate:"required"`
	Complement  string `json:"complement"   validate:""`
	District    string `json:"district"     validate:"required"`
	Place       string `json:"place"        validate:"required"`
	Uf          string `json:"uf"           validate:"required"`
	IbgeCode    string `json:"ibgeCode"     validate:"required"`
}

type GetCep struct {
	ZipCode     string `json:"zipCode"      validate:"required"`
}
