package util

import (
    "fmt"
)

func FloatToTimeStamp(input float32) string {
    if input < 0 {
        input *= -1
    }
    minutes := uint8(input / 60)
    seconds := input - float32(minutes*60)
    return fmt.Sprintf("%d:%06.3f", minutes, seconds)
}
