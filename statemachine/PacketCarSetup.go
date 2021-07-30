package statemachine

type CarSetupData struct {
  FrontWing   uint8
  RearWing    uint8
  OnThrottle  uint8
  OffThrottle uint8
  FrontCamber float32
  RearCamber  float32
  FrontToe    float32
  RearToe     float32

  FrontSuspension        uint8
  RearSuspension         uint8
  FrontAntiRollBar       uint8
  RearAntiRollBar        uint8
  FrontSuspensionHeight  uint8
  RearSuspensionHeight   uint8
  BrakePressure          uint8
  BrakeBias              uint8
  RearLeftTyrePressure   float32
  RearRightTyrePressure  float32
  FrontLeftTyrePressure  float32
  FrontRightTyrePressure float32

  Ballast  uint8
  FuelLoad float32
}

type PacketCatSetupData struct {
  Header    PacketHeader
  CarSetups [22]CarSetupData
}

func ProcessPacketCarSetups(csm *CommunicationStateMachine, state *StateMachine) {
  buffered := len(csm.UnprocessedBuffer)

  packetCarSetups := PacketCatSetupData{}
  requiredSize := GetMemorySize(packetCarSetups)

  if buffered >= requiredSize {
    ToObject(csm.UnprocessedBuffer[:], &packetCarSetups)

    // println(fmt.Sprintf("data car setups: %+v", packetCarSetups))

    csm.RemoveFirstBytesFromBuffer(requiredSize, state)
  }
}
