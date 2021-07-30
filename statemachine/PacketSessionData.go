package statemachine

type PacketSessionData struct {
  Header                    PacketHeader
  Weather                   uint8
  TrackTemperature          int8
  AirTemperature            int8
  TotalLaps                 uint8
  TrackLength               uint16
  SessionType               uint8
  TrackId                   int8
  Formula                   uint8
  SessionTimeLeft           uint16
  SessionDuration           uint16
  PitSpeedLimit             uint8
  GamePaused                uint8
  IsSpectating              uint8
  SpectatorCarIndex         uint8
  SliProNativeSupport       uint8
  NumMarshalZones           uint8
  MarshalZones              [21]MarshalZone
  SafetyCarStatus           uint8
  NetworkGame               uint8
  NumWeatherForecastSamples uint8
  WeatherForecastSamples    [20]WeatherForecastSample
}

type WeatherForecastSample struct {
  SessionType      uint8
  TimeOffset       uint8
  Weather          uint8
  TrackTemperature int8
  AirTemperature   int8
}

type MarshalZone struct {
  ZoneStart float32
  ZoneFlag  int8
}

func processPacketSession(csm *CommunicationStateMachine, state *StateMachine) {
  buffered := len(csm.UnprocessedBuffer)

  packetSessionData := PacketSessionData{}
  requiredSize := GetMemorySize(packetSessionData)

  if buffered >= requiredSize {
    ToObject(csm.UnprocessedBuffer[:], &packetSessionData)

    // println(fmt.Sprintf("data car setups: %+v", packetSessionData))

    state.SessionData = packetSessionData
    csm.RemoveFirstBytesFromBuffer(requiredSize, state)
  }
}
