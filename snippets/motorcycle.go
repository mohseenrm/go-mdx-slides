package main

import "fmt"

const usixteenbitmax float64 = 65535
const kmphMultiple float64 = 1.60934

type motorcycle struct {
    throttleResponse uint16
    brakePressure uint16
    rearWheelPSI float64
    topSpeedMph float64
}

func (m motorcycle) kmph() float64 {
	return float64(m.throttleResponse) * (m.topSpeedMph / usixteenbitmax) * kmphMultiple
}

func main() {
	panigale := motorcycle{
		throttleResponse: 60000,
		brakePressure: 934,
		rearWheelPSI: 36.2,
		topSpeedMph: 202}

	fmt.Println(panigale.throttleResponse)
	fmt.Println(panigale.kmph())
}
