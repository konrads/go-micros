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
	ID                string
	Name              string
	Alias             []string
	Constellation     string
	Coordinates       []float32
	Distance          float32
	ApparentMagnitude float32
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
		ID:                id,
		Name:              s.Name,
		Alias:             s.Alias,
		Constellation:     s.Constellation,
		Coordinates:       s.Coordinates,
		Distance:          s.Distance,
		ApparentMagnitude: s.ApparentMagnitude,
	}
}

func (s *Star) ToStarReq() StarReq {
	return StarReq{
		Name:              s.Name,
		Alias:             s.Alias,
		Constellation:     s.Constellation,
		Coordinates:       s.Coordinates,
		Distance:          s.Distance,
		ApparentMagnitude: s.ApparentMagnitude,
	}
}
