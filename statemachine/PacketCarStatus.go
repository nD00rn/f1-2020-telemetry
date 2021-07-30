package statemachine

type CarStatusData struct {
  TractionControl   uint8
  AntiLockBrakes    uint8
  FuelMix           uint8
  FrontBrakeBias    uint8
  PitLimiterStatus  uint8
  FuelInTank        float32
  FuelCapacity      float32
  FuelRemainingLaps float32
  MaxRpm            uint16
  IdleRpm           uint16
  MaxGears          uint8
  DrsAllowed        uint8

  DrsActivationDistance uint16
  TyresWear             [4]uint8
  ActualTypeCompound    uint8
  VisualTyreCompound    uint8
  TyresAgeLaps          uint8
  TyresDamage           [4]uint8
  FrontLeftWingDamage   uint8
  FrontRightWingDamage  uint8
  RearWingDamage        uint8

  DrsFault                uint8
  EngineDamage            uint8
  GearBoxDamage           uint8
  VehicleFiaFlags         int8
  ErsStoreEnergy          float32
  ErsDeployMode           uint8
  ErsHarvestedThisLapMguk float32
  ErsHarvestedThisLapMguh float32
  ErsDeployedThisLap      float32
}

type PacketCarStatusData struct {
  Header        PacketHeader
  CarStatusData [22]CarStatusData
}

func ProcessPacketCarStatus(state *StateMachine) {
  buffered := len(state.UnprocessedBuffer)

  packetCarStatusData := PacketCarStatusData{}
  requiredSize := GetMemorySize(packetCarStatusData)

  if buffered >= requiredSize {
    ToObject(state.UnprocessedBuffer[:], &packetCarStatusData)

    state.CarStatus = packetCarStatusData

    state.RemoveFirstBytesFromBuffer(requiredSize)
  }
}
