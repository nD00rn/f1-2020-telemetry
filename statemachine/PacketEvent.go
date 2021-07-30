package statemachine

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

func ProcessPacketEvent(csm *CommunicationStateMachine, state *StateMachine) {
  event := PacketEventData{}
  sizeHeader := GetMemorySize(event.Header)
  sizeEventType := GetMemorySize(event.EventStringCode)
  sizeHeaderAndEventType := sizeHeader + sizeEventType

  if csm.AvailableData() >= sizeHeaderAndEventType {
    wantedBytes := csm.UnprocessedBuffer[sizeHeader:sizeHeaderAndEventType]

    text := string(wantedBytes)

    switch text {
    case "FTLP":
      fastestLap := FastestLap{}
      if csm.AvailableData() >= GetMemorySize(fastestLap)+sizeHeaderAndEventType {
        ToObject(csm.UnprocessedBuffer[sizeHeaderAndEventType:], &fastestLap)
        csm.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType + GetMemorySize(fastestLap), state)
      }
      break

    case "RTMT":
      retirement := Retirement{}
      if csm.AvailableData() >= GetMemorySize(retirement)+sizeHeaderAndEventType {
        ToObject(csm.UnprocessedBuffer[sizeHeaderAndEventType:], &retirement)
        csm.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType + GetMemorySize(retirement), state)
      }
      break

    case "TMPT":
      teamMateInPits := TeamMateInPits{}
      if csm.AvailableData() >= GetMemorySize(teamMateInPits)+sizeHeaderAndEventType {
        ToObject(csm.UnprocessedBuffer[sizeHeaderAndEventType:], &teamMateInPits)
        csm.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType + GetMemorySize(teamMateInPits), state)
      }
      break

    case "RCWN":
      raceWinner := RaceWinner{}
      if csm.AvailableData() >= GetMemorySize(raceWinner)+sizeHeaderAndEventType {
        ToObject(csm.UnprocessedBuffer[sizeHeaderAndEventType:], &raceWinner)
        csm.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType + GetMemorySize(raceWinner), state)
      }
      break

    case "PENA":
      penalty := Penalty{}
      if csm.AvailableData() >= GetMemorySize(penalty)+sizeHeaderAndEventType {
        ToObject(csm.UnprocessedBuffer[sizeHeaderAndEventType:], &penalty)
        csm.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType + GetMemorySize(penalty), state)
      }
      break

    case "SPTP":
      trap := SpeedTrap{}
      sizeTrap := GetMemorySize(trap)
      if csm.AvailableData() >= sizeTrap+sizeHeaderAndEventType {
        ToObject(csm.UnprocessedBuffer[sizeHeaderAndEventType:], &trap)
        state.SpeedTraps[trap.VehicleIndex] = trap
        csm.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType + sizeTrap, state)
      }
      break

    case "CHQF":
      csm.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType, state)
      break

    case "DRSE":
      csm.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType, state)
      break

    case "DRSD":
      csm.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType, state)
      break

    case "SEND":
      csm.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType, state)
      break

    case "SSTA":
      csm.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType, state)
      state.ResetTimers()
      break
    }
  }
}
