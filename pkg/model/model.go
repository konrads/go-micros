package model

type StarReq struct {
	Name              string    `json:"name"`
	Alias             []string  `json:"alias"`
	Constellation     string    `json:"constellation"`
	Coordinates       []float32 `json:"coordinates"`
	Distance          float32   `json:"distance"`
	ApparentMagnitude float32   `json:"apparentMagnitude"`
}

type Star struct {
	Id                string
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

func (s *StarReq) ToStar(id string) Star {
	return Star{
		Id:                id,
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

// FIXME: another way to parse json to structs?
func StarReqFromJson(kv map[string]interface{}) StarReq {
	var name, constellation string
	var distance, apparentMagnitude float32
	if v, ok := kv["name"].(string); ok {
		name = v
	}
	if v, ok := kv["constellation"].(string); ok {
		constellation = v
	}
	if v, ok := kv["distance"].(float64); ok {
		distance = float32(v)
	}
	if v, ok := kv["apparentMagnitude"].(float64); ok {
		apparentMagnitude = float32(v)
	}
	return StarReq{
		Name:              name,
		Alias:             toStringSlice(kv["alias"].([]interface{})),
		Constellation:     constellation,
		Coordinates:       toFloat32Slice(kv["coordinates"].([]interface{})),
		Distance:          distance,
		ApparentMagnitude: apparentMagnitude,
	}
}

func toFloat32Slice(org []interface{}) []float32 {
	slice := []float32{}
	for _, x := range org {
		slice = append(slice, float32(x.(float64)))
	}
	return slice
}

func toStringSlice(org []interface{}) []string {
	slice := []string{}
	for _, x := range org {
		slice = append(slice, x.(string))
	}
	return slice
}
