package statemachine

type CarPacketMotion struct {
  WorldPositionX    float32
  WorldPositionY    float32
  WorldPositionZ    float32
  WorldVelocityX    float32
  WorldVelocityY    float32
  WorldVelocityZ    float32
  WorldForwardDirX  int16
  WorldForwardDirY  int16
  WorldForwardDirZ  int16
  WorldRightDirX    int16
  WorldRightDirY    int16
  WorldRightDirZ    int16
  ForceLateral      float32
  ForceLongitudinal float32
  ForceVertical     float32
  Yaw               float32
  Pitch             float32
  Roll              float32
}

type PacketMotion struct {
  Header PacketHeader

  CarMotionData [22]CarPacketMotion

  // Extra player car ONLY data
  SuspensionPosition     [4]float32
  SuspensionVelocity     [4]float32
  SuspensionAcceleration [4]float32
  WheelSpeed             [4]float32
  WheelSlip              [4]float32
  LocalVelocityX         float32
  LocalVelocityY         float32
  LocalVelocityZ         float32
  AngularVelocityX       float32
  AngularVelocityY       float32
  AngularVelocityZ       float32
  AngularAccelerationX   float32
  AngularAccelerationY   float32
  AngularAccelerationZ   float32
  FrontWheelsAngle       float32
}

func ProcessPacketMotion(csm *CommunicationStateMachine, state *StateMachine) {
  buffered := len(csm.UnprocessedBuffer)

  packetMotion := PacketMotion{}
  requiredSize := GetMemorySize(packetMotion)

  if buffered >= requiredSize {
    ToObject(csm.UnprocessedBuffer[:], &packetMotion)

    csm.RemoveFirstBytesFromBuffer(requiredSize, state)
  }
}
