package dto

type UserLocation struct {
	PlaceID     int         `json:"place_id"`
	Licence     string      `json:"licence"`
	OsmType     string      `json:"osm_type"`
	OsmID       int         `json:"osm_id"`
	Lat         string      `json:"lat"`
	Lon         string      `json:"lon"`
	Class       string      `json:"class"`
	Type        string      `json:"type"`
	PlaceRank   int         `json:"place_rank"`
	Importance  float64     `json:"importance"`
	AddressType string      `json:"addresstype"`
	Name        string      `json:"name"`
	DisplayName string      `json:"display_name"`
	UserAddress UserAddress `json:"address"`
	BoundingBox []string    `json:"boundingbox"`
}

type UserAddress struct {
	Industrial  string `json:"industrial"`
	Suburb      string `json:"suburb"`
	City        string `json:"city"`
	Regency     string `json:"regency"`
	State       string `json:"state"`
	ISO4        string `json:"ISO3166-2-lvl4"`
	Region      string `json:"region"`
	ISO3        string `json:"ISO3166-2-lvl3"`
	Postcode    string `json:"postcode"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
}
