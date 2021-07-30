package statemachine

import (
  "fmt"
)

type FinalClassificationData struct {
  Position     uint8
  NumLaps      uint8
  GridPosition uint8
  Points       uint8
  NumPitStops  uint8
  ResultStatus uint8

  BestLapTime      float32
  TotalRaceTime    float64
  PenaltiesTime    uint8
  NumPenalties     uint8
  NumTyreStints    uint8
  TypeStintsActual [8]uint8
  TypeStintsVisual [8]uint8
}

type PacketFinalClassificationData struct {
  Header                  PacketHeader
  NumCars                 uint8
  FinalClassificationData [22]FinalClassificationData
}

func ProcessPacketFinalClassification(state *StateMachine) {
  buffered := len(state.UnprocessedBuffer)

  packetClassificationData := PacketFinalClassificationData{}
  requiredSize := GetMemorySize(packetClassificationData)

  if buffered >= requiredSize {
    ToObject(state.UnprocessedBuffer[:], &packetClassificationData)

    println(fmt.Sprintf("data final classification: %+v", packetClassificationData))

    state.RemoveFirstBytesFromBuffer(requiredSize)
  }
}
