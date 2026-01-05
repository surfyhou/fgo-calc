package model

type ServantDetail struct {
	Name     string           `json:"name"`
	Traits   []int            `json:"traits"`
	Cost     int              `json:"cost"`
	Img      string           `json:"img"`
	TraitSet map[int]struct{} `json:"-"`
}

type Servant struct {
	Id   int                      `json:"id"`
	Name string                   `json:"name"`
	Diff map[string]ServantDetail `json:"diff"`
}
