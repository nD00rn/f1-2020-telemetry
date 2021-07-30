package statemachine

type ParticipantData struct {
  AiControlled  uint8
  DriverId      uint8
  TeamId        uint8
  RaceNumber    uint8
  Nationality   uint8
  Name          [48]byte
  YourTelemetry uint8
}

type PacketParticipantsData struct {
  Header        PacketHeader
  NumActiveCars uint8
  Participants  [22]ParticipantData
}

func ProcessPacketParticipants(state *StateMachine) {
  buffered := len(state.UnprocessedBuffer)

  packetParticipants := PacketParticipantsData{}
  requiredSize := GetMemorySize(packetParticipants)

  if buffered >= requiredSize {
    ToObject(state.UnprocessedBuffer[:], &packetParticipants)

    // println(fmt.Sprintf("data participants: %+v", packetParticipants))

    state.ParticipantData = packetParticipants
    state.RemoveFirstBytesFromBuffer(requiredSize)
  }
}
