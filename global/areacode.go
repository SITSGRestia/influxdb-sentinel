package global

type AreaCode struct {
	CN string `json:"cn"`
	ZN string `json:"zn"`
}

type City map[string]AreaCode

