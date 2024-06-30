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
	ZipCode     string `json:"cep"      validate:"required"`
	PublicPlace string `json:"logradouro"  validate:"required"`
	Complement  string `json:"complemento"   validate:""`
	District    string `json:"bairro"     validate:"required"`
	Place       string `json:"localidade"        validate:"required"`
	Uf          string `json:"uf"           validate:"required"`
	IbgeCode    string `json:"ibge"     validate:"required"`
}

type GetCep struct {
	ZipCode string `json:"zipCode"      validate:"required"`
}
