package statemachine

type CarTelemetryData struct {
  Speed                   uint16
  Throttle                float32
  Steer                   float32
  Brake                   float32
  Clutch                  uint8
  Gear                    int8
  EngineRpm               uint16
  Drs                     uint8
  RevLightsPercent        uint8
  BrakeTemperature        [4]uint16
  TyresSurfaceTemperature [4]uint8
  TyresInnerTemperature   [4]uint8
  EngineTemperature       uint16
  TyresPressure           [4]float32
  SurfaceType             [4]uint8
}

type PacketCarTelemetryData struct {
  Header                       PacketHeader
  CarTelemetryData             [22]CarTelemetryData
  ButtonsStatus                uint32
  MfdPanelIndex                uint8
  MfdPanelIndexSecondaryPlayer uint8
  SuggestedGear                int8
}

func ProcessPacketCarTelemetry(state *StateMachine) {
  buffered := len(state.UnprocessedBuffer)

  packetCarTelemetryData := PacketCarTelemetryData{}
  requiredSize := GetMemorySize(packetCarTelemetryData)

  if buffered >= requiredSize {
    ToObject(state.UnprocessedBuffer[:], &packetCarTelemetryData)

    state.TelemetryData = packetCarTelemetryData

    state.RemoveFirstBytesFromBuffer(requiredSize)
  }
}
