package model

type PortReq struct {
	Name        string     `json:"name"`
	Coordinates [2]float32 `json:"coordinates"`
	City        string     `json:"city"`
	Province    string     `json:"province"`
	Country     string     `json:"country"`
	Alias       []string   `json:"alias"`
	Regions     []string   `json:"regions"`
	Timezone    string     `json:"timezone"`
	Unlocs      []string   `json:"unlocs"`
	Code        string     `json:"code"`
}

type Port struct {
	Id          string
	Name        string
	Coordinates [2]float32
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
