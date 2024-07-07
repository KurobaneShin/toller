package types

type Invoice struct {
	OBUID         int     `json:"obuid,omitempty"`
	TotalDistance float64 `json:"total_distance,omitempty"`
	TotalAmount   float64 `json:"total_amount,omitempty"`
}

type Distance struct {
	Value float64 `json:"value,omitempty"`
	OBUID int     `json:"obuid,omitempty"`
	Unix  int64   `json:"unix,omitempty"`
}

type OBUData struct {
	OBUID     int     `json:"obuId,omitempty"`
	Lat       float64 `json:"lat,omitempty"`
	Long      float64 `json:"long,omitempty"`
	RequestId int     `json:"request_id,omitempty"`
}
