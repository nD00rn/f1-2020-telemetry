package statemachine

import (
  "fmt"
)

type PacketEventData struct {
  Header          PacketHeader
  EventStringCode [4]byte
  EventDetails    EventDataDetails
}

type EventDataDetails interface{}

type FastestLap struct {
  VehicleIndex uint8
  LapTime      float32
}

type Retirement struct {
  VehicleIndex uint8
}

type TeamMateInPits struct {
  VehicleIndex uint8
}

type RaceWinner struct {
  VehicleIndex uint8
}

type Penalty struct {
  PenaltyType       uint8
  InfringementType  uint8
  VehicleIndex      uint8
  OtherVehicleIndex uint8
  Time              uint8
  LapNum            uint8
  PlacesGained      uint8
}

type SpeedTrap struct {
  VehicleIndex uint8
  Speed        float32
}

func ProcessPacketEvent(state *StateMachine) {
  event := PacketEventData{}
  sizeHeader := GetMemorySize(event.Header)
  sizeEventType := GetMemorySize(event.EventStringCode)
  sizeHeaderAndEventType := sizeHeader + sizeEventType

  if state.AvailableData() >= sizeHeaderAndEventType {
    wantedBytes := state.UnprocessedBuffer[sizeHeader:sizeHeaderAndEventType]
    // fmt.Printf("wanted bytes: %+v | %x %s\n", wantedBytes, wantedBytes, wantedBytes)

    fmt.Printf("event type is %s\n", wantedBytes)

    text := string(wantedBytes)

    switch text {
    case "FTLP":
      fastestLap := FastestLap{}
      if state.AvailableData() >= GetMemorySize(fastestLap)+sizeHeaderAndEventType {
        ToObject(state.UnprocessedBuffer[sizeHeaderAndEventType:], &fastestLap)
        fmt.Printf("event: fastestLap triggered for id:%d, time:%.2f\n", fastestLap.VehicleIndex, fastestLap.LapTime)
        // state.FastestLap = fastestLap
        state.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType + GetMemorySize(fastestLap))
      }
      break

    case "RTMT":
      retirement := Retirement{}
      if state.AvailableData() >= GetMemorySize(retirement)+sizeHeaderAndEventType {
        ToObject(state.UnprocessedBuffer[sizeHeaderAndEventType:], &retirement)
        fmt.Printf("event: retirement triggered for id:%d\n", retirement.VehicleIndex)
        state.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType + GetMemorySize(retirement))
      }
      break

    case "TMPT":
      teamMateInPits := TeamMateInPits{}
      if state.AvailableData() >= GetMemorySize(teamMateInPits)+sizeHeaderAndEventType {
        ToObject(state.UnprocessedBuffer[sizeHeaderAndEventType:], &teamMateInPits)
        fmt.Printf("event: teamMateInPits triggered for id:%d\n", teamMateInPits.VehicleIndex)
        state.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType + GetMemorySize(teamMateInPits))
      }
      break

    case "RCWN":
      raceWinner := RaceWinner{}
      if state.AvailableData() >= GetMemorySize(raceWinner)+sizeHeaderAndEventType {
        ToObject(state.UnprocessedBuffer[sizeHeaderAndEventType:], &raceWinner)
        fmt.Printf("event: raceWinner triggered for id:%d\n", raceWinner.VehicleIndex)
        state.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType + GetMemorySize(raceWinner))
      }
      break

    case "PENA":
      penalty := Penalty{}
      if state.AvailableData() >= GetMemorySize(penalty)+sizeHeaderAndEventType {
        ToObject(state.UnprocessedBuffer[sizeHeaderAndEventType:], &penalty)
        fmt.Printf("event: penalty triggered for id:%d, object:%+v\n", penalty.VehicleIndex, penalty)
        state.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType + GetMemorySize(penalty))
      }
      break

    case "SPTP":
      trap := SpeedTrap{}
      sizeTrap := GetMemorySize(trap)
      if state.AvailableData() >= sizeTrap+sizeHeaderAndEventType {
        ToObject(state.UnprocessedBuffer[sizeHeaderAndEventType:], &trap)
        fmt.Printf("event: speed trap triggered for id:%d, speed %.2f\n", trap.VehicleIndex, trap.Speed)
        state.SpeedTraps[trap.VehicleIndex] = trap
        state.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType + sizeTrap)
      } else {
        fmt.Printf(
          "not enough data for speed trap. trap-mem:%d, size header and event type:%d, available %d\n",
          sizeTrap,
          sizeHeaderAndEventType,
          state.AvailableData(),
        )
      }
      break

    case "CHQF":
      fmt.Println("event: chequered flag, no additional data.")
      state.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType)
      break

    case "DRSE":
      fmt.Println("event: DRS enabled, no additional data.")
      state.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType)
      break

    case "DRSD":
      fmt.Println("event: DRS disabled, no additional data.")
      state.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType)
      break

    case "SEND":
      fmt.Println("event: end of session, no additional data.")
      state.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType)
      break

    case "SSTA":
      fmt.Println("event: start session, no additional data.")
      state.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType)
      state.ResetTimers()

      break
    }
  }
}
