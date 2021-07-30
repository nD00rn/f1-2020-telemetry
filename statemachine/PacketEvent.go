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
    // fmt.Printf("wanted bytes: %+v | %x %s\n", wantedBytes, wantedBytes, wantedBytes)
    // fmt.Printf("event type is %s\n", wantedBytes)

    text := string(wantedBytes)

    switch text {
    case "FTLP":
      fastestLap := FastestLap{}
      if csm.AvailableData() >= GetMemorySize(fastestLap)+sizeHeaderAndEventType {
        ToObject(csm.UnprocessedBuffer[sizeHeaderAndEventType:], &fastestLap)
        // fmt.Printf("event: fastestLap triggered for id:%d, time:%.2f\n", fastestLap.VehicleIndex, fastestLap.LapTime)
        // state.FastestLap = fastestLap
        csm.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType + GetMemorySize(fastestLap), state)
      }
      break

    case "RTMT":
      retirement := Retirement{}
      if csm.AvailableData() >= GetMemorySize(retirement)+sizeHeaderAndEventType {
        ToObject(csm.UnprocessedBuffer[sizeHeaderAndEventType:], &retirement)
        // fmt.Printf("event: retirement triggered for id:%d\n", retirement.VehicleIndex)
        csm.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType + GetMemorySize(retirement), state)
      }
      break

    case "TMPT":
      teamMateInPits := TeamMateInPits{}
      if csm.AvailableData() >= GetMemorySize(teamMateInPits)+sizeHeaderAndEventType {
        ToObject(csm.UnprocessedBuffer[sizeHeaderAndEventType:], &teamMateInPits)
        // fmt.Printf("event: teamMateInPits triggered for id:%d\n", teamMateInPits.VehicleIndex)
        csm.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType + GetMemorySize(teamMateInPits), state)
      }
      break

    case "RCWN":
      raceWinner := RaceWinner{}
      if csm.AvailableData() >= GetMemorySize(raceWinner)+sizeHeaderAndEventType {
        ToObject(csm.UnprocessedBuffer[sizeHeaderAndEventType:], &raceWinner)
        // fmt.Printf("event: raceWinner triggered for id:%d\n", raceWinner.VehicleIndex)
        csm.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType + GetMemorySize(raceWinner), state)
      }
      break

    case "PENA":
      penalty := Penalty{}
      if csm.AvailableData() >= GetMemorySize(penalty)+sizeHeaderAndEventType {
        ToObject(csm.UnprocessedBuffer[sizeHeaderAndEventType:], &penalty)
        // fmt.Printf("event: penalty triggered for id:%d, object:%+v\n", penalty.VehicleIndex, penalty)
        csm.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType + GetMemorySize(penalty), state)
      }
      break

    case "SPTP":
      trap := SpeedTrap{}
      sizeTrap := GetMemorySize(trap)
      if csm.AvailableData() >= sizeTrap+sizeHeaderAndEventType {
        ToObject(csm.UnprocessedBuffer[sizeHeaderAndEventType:], &trap)
        // fmt.Printf("event: speed trap triggered for id:%d, speed %.2f\n", trap.VehicleIndex, trap.Speed)
        state.SpeedTraps[trap.VehicleIndex] = trap
        csm.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType + sizeTrap, state)
      } else {
        // fmt.Printf(
        //   "not enough data for speed trap. trap-mem:%d, size header and event type:%d, available %d\n",
        //   sizeTrap,
        //   sizeHeaderAndEventType,
        //   state.AvailableData(),
        // )
      }
      break

    case "CHQF":
      // fmt.Println("event: chequered flag, no additional data.")
      csm.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType, state)
      break

    case "DRSE":
      // fmt.Println("event: DRS enabled, no additional data.")
      csm.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType, state)
      break

    case "DRSD":
      // fmt.Println("event: DRS disabled, no additional data.")
      csm.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType, state)
      break

    case "SEND":
      // fmt.Println("event: end of session, no additional data.")
      csm.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType, state)
      break

    case "SSTA":
      // fmt.Println("event: start session, no additional data.")
      csm.RemoveFirstBytesFromBuffer(sizeHeaderAndEventType, state)
      state.ResetTimers()
      break
    }
  }
}
