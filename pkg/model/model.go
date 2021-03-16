package model

type PortReq struct {
	Name        string    `json:"name"`
	Coordinates []float32 `json:"coordinates"`
	City        string    `json:"city"`
	Province    string    `json:"province"`
	Country     string    `json:"country"`
	Alias       []string  `json:"alias"`
	Regions     []string  `json:"regions"`
	Timezone    string    `json:"timezone"`
	Unlocs      []string  `json:"unlocs"`
	Code        string    `json:"code"`
}

type Port struct {
	Id          string
	Name        string
	Coordinates []float32
	City        string
	Province    string
	Country     string
	Alias       []string
	Regions     []string
	Timezone    string
	Unlocs      []string
	Code        string
}

// type Ports struct {
// 	Ports []Port
// }

type Pong struct {
	Message string `json:"message"`
}

func (pr *PortReq) ToPort(id string) Port {
	return Port{
		Id:          id,
		Name:        pr.Name,
		Coordinates: pr.Coordinates,
		City:        pr.City,
		Province:    pr.Province,
		Country:     pr.Country,
		Alias:       pr.Alias,
		Regions:     pr.Regions,
		Timezone:    pr.Timezone,
		Unlocs:      pr.Unlocs,
		Code:        pr.Code,
	}
}

func (pr *Port) ToPortReq() PortReq {
	return PortReq{
		Name:        pr.Name,
		Coordinates: pr.Coordinates,
		City:        pr.City,
		Province:    pr.Province,
		Country:     pr.Country,
		Alias:       pr.Alias,
		Regions:     pr.Regions,
		Timezone:    pr.Timezone,
		Unlocs:      pr.Unlocs,
		Code:        pr.Code,
	}
}

// FIXME: another way to parse json to structs?
func PortReqFromJson(kv map[string]interface{}) PortReq {
	var name, city, province, country, timezone, code string
	if v, ok := kv["name"].(string); ok {
		name = v
	}
	if v, ok := kv["city"].(string); ok {
		city = v
	}
	if v, ok := kv["province"].(string); ok {
		province = v
	}
	if v, ok := kv["country"].(string); ok {
		country = v
	}
	if v, ok := kv["timezone"].(string); ok {
		timezone = v
	}
	if v, ok := kv["code"].(string); ok {
		code = v
	}
	return PortReq{
		Name:        name,
		Coordinates: toFloat32Slice(kv["coordinates"].([]interface{})),
		City:        city,
		Province:    province,
		Country:     country,
		Alias:       toStringSlice(kv["alias"].([]interface{})),
		Regions:     toStringSlice(kv["regions"].([]interface{})),
		Timezone:    timezone,
		Unlocs:      toStringSlice(kv["unlocs"].([]interface{})),
		Code:        code,
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
