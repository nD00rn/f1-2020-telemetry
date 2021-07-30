package rest

type PlayerCar struct {
    Speed uint16 `json:"speed"`
    Gear  int8   `json:"gear"`
    Lap uint8 `json:"lap"`
    LapPercent float32 `json:"lap_percent"`
}
