package model

type StarReq struct {
	Name              string    `json:"name" validate:"required"`
	Alias             []string  `json:"alias"`
	Constellation     string    `json:"constellation"`
	Coordinates       []float32 `json:"coordinates"`
	Distance          float32   `json:"distance"`
	ApparentMagnitude float32   `json:"apparentMagnitude"`
}

type Star struct {
	ID string
	*StarReq
}

type Pong struct {
	Message string `json:"message"`
}

func DefaultStarReq() *StarReq {
	// only need to assign slice values, which default to nil
	return &StarReq{
		Alias:       []string{},
		Coordinates: []float32{},
	}
}

func (s *StarReq) ToStar(id string) Star {
	return Star{
		ID:      id,
		StarReq: s,
	}
}

func (s *Star) ToStarReq() StarReq {
	return *s.StarReq
}
