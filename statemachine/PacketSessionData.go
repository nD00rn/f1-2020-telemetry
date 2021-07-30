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

func processPacketSession(state *StateMachine) {
  buffered := len(state.UnprocessedBuffer)

  packetSessionData := PacketSessionData{}
  requiredSize := GetMemorySize(packetSessionData)

  if buffered >= requiredSize {
    ToObject(state.UnprocessedBuffer[:], &packetSessionData)

    // println(fmt.Sprintf("data car setups: %+v", packetSessionData))

    state.SessionData = packetSessionData
    state.RemoveFirstBytesFromBuffer(requiredSize)
  }
}

// func convertDataToSessionData(data []byte) (int, *PacketSessionData, error) {
//   sessionData := PacketSessionData{}
//
//   if len(data) < 42 {
//     // fmt.Println("not enough data to fill 42 bytes")
//     return 0, nil, errors.New("not enough data")
//   }
//
//   deserializeBytesSessionData(data[0:24], &sessionData.Header)
//   deserializeByteSessionData(24, data, &sessionData.Weather)
//   deserializeByteSessionData(25, data, &sessionData.TrackTemperature)
//   deserializeByteSessionData(26, data, &sessionData.AirTemperature)
//   deserializeByteSessionData(27, data, &sessionData.TotalLaps)
//   deserializeBytesSessionData(data[28:30], &sessionData.TrackLength)
//   deserializeByteSessionData(30, data, &sessionData.SessionType)
//   deserializeByteSessionData(31, data, &sessionData.TrackId)
//   deserializeByteSessionData(32, data, &sessionData.Formula)
//   deserializeBytesSessionData(data[33:35], &sessionData.SessionTimeLeft)
//   deserializeBytesSessionData(data[35:37], &sessionData.SessionDuration)
//   deserializeByteSessionData(37, data, &sessionData.PitSpeedLimit)
//   deserializeByteSessionData(38, data, &sessionData.GamePaused)
//   deserializeByteSessionData(39, data, &sessionData.IsSpectating)
//   deserializeByteSessionData(40, data, &sessionData.SpectatorCarIndex)
//   deserializeByteSessionData(41, data, &sessionData.SliProNativeSupport)
//   deserializeByteSessionData(42, data, &sessionData.NumMarshalZones)
//
//   zones := make([]MarshalZone, sessionData.NumMarshalZones)
//   sessionData.MarshalZones = zones
//   zoneSize := GetMemorySize(zones)
//   // fmt.Printf("zone size is %d\n", zoneSize)
//   endIndexMarshalZones := 43 + zoneSize
//
//   if len(data) < (43 + zoneSize) {
//     fmt.Println("not enough data to fill zone bytes")
//     return 0, nil, errors.New("not enough data")
//   }
//
//   deserializeBytesSessionData(data[43:endIndexMarshalZones], &sessionData.MarshalZones)
//   deserializeByteSessionData(endIndexMarshalZones, data, &sessionData.SafetyCarStatus)
//   deserializeByteSessionData(endIndexMarshalZones+1, data, &sessionData.NetworkGame)
//   deserializeByteSessionData(endIndexMarshalZones+2, data, &sessionData.NumWeatherForecastSamples)
//
//   weatherSamples := make([]WeatherForecastSample, sessionData.NumWeatherForecastSamples)
//   sessionData.WeatherForecastSamples = weatherSamples
//   weatherSamplesSize := GetMemorySize(weatherSamples)
//   // fmt.Printf("weather sample size is %d\n", weatherSamplesSize)
//   endIndexWeatherSamples := endIndexMarshalZones + 3 + weatherSamplesSize
//
//   if len(data) < endIndexWeatherSamples {
//     // fmt.Println("not enough data to fill weather bytes")
//     return 0, nil, errors.New("not enough data")
//   }
//
//   // fmt.Printf("start of weather data is at %d\n", endIndexMarshalZones+3)
//   // fmt.Printf("end of weather data is at %d\n", endIndexWeatherSamples)
//
//   deserializeBytesSessionData(
//     data[endIndexMarshalZones+3:endIndexWeatherSamples],
//     &sessionData.WeatherForecastSamples,
//   )
//
//   return endIndexWeatherSamples, &sessionData, nil
// }
