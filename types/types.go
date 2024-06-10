package types

type OBUData struct {
	OBUID int     `json:"obuId,omitempty"`
	Lat   float64 `json:"lat,omitempty"`
	Long  float64 `json:"long,omitempty"`
}
