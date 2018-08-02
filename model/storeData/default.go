package storeData


type Base struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type BaseInt struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type BaseFloat32 struct {
	Name  string  `json:"name"`
	Value float32 `json:"value"`
}

type BaseFloat64 struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

