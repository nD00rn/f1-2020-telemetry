package statemachine

type LapData struct {
  LastLapTime    float32
  CurrentLapTime float32

  SectorOneTimeInMs uint16
  SectorTwoTimeInMs uint16

  BestLapTime                float32
  BestLapNum                 uint8
  BestLapSectorOneTimeInMs   uint16
  BestLapSectorTwoTimeInMs   uint16
  BestLapSectorThreeTimeInMs uint16

  BestOverallSectorOneTimeInMs   uint16
  BestOverallSectorOneLapNum     uint8
  BestOverallSectorTwoTimeInMs   uint16
  BestOverallSectorTwoLapNum     uint8
  BestOverallSectorThreeTimeInMs uint16
  BestOverallSectorThreeLapNum   uint8

  LapDistance       float32
  TotalDistance     float32
  SafetyCarDelta    float32
  CarPosition       uint8
  CurrentLapNum     uint8
  PitStatus         uint8
  Sector            uint8
  CurrentLapInvalid uint8
  Penalties         uint8
  GridPosition      uint8
  DriverStatus      uint8
  ResultStatus      uint8
}

type PacketLapData struct {
  Header  PacketHeader
  LapData [22]LapData
}

func ProcessPacketLapData(csm *CommunicationStateMachine, state *StateMachine) {
  buffered := len(csm.UnprocessedBuffer)

  packetLapData := PacketLapData{}
  requiredSize := GetMemorySize(packetLapData)

  if buffered >= requiredSize {
    ToObject(csm.UnprocessedBuffer[:], &packetLapData)

    state.LapData = packetLapData

    csm.RemoveFirstBytesFromBuffer(requiredSize, state)
  }
}
