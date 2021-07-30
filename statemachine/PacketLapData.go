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

func ProcessPacketLapData(state *StateMachine) {
  buffered := len(state.UnprocessedBuffer)

  packetLapData := PacketLapData{}
  requiredSize := GetMemorySize(packetLapData)

  if buffered >= requiredSize {
    ToObject(state.UnprocessedBuffer[:], &packetLapData)

    // println(fmt.Sprintf("data lap data: %+v", packetLapData))

    // println(fmt.Sprintf("data lap data: %f %f | %f | S1:%d S2:%d S3:%d | S1:%d S2:%d LL:%f",
    //   packetLapData.LapData[0].LapDistance,
    //   packetLapData.LapData[0].TotalDistance,
    //   packetLapData.LapData[0].CurrentLapTime,
    //   packetLapData.LapData[0].BestLapSectorOneTimeInMs,
    //   packetLapData.LapData[0].BestLapSectorTwoTimeInMs,
    //   packetLapData.LapData[0].BestLapSectorThreeTimeInMs,
    //   packetLapData.LapData[0].SectorOneTimeInMs,
    //   packetLapData.LapData[0].SectorTwoTimeInMs,
    //   packetLapData.LapData[0].LastLapTime,
    //   ))

    state.LapData = packetLapData

    state.RemoveFirstBytesFromBuffer(requiredSize)
  }
}
