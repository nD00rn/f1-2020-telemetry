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

func ProcessPacketParticipants(csm *CommunicationStateMachine, state *StateMachine) {
  buffered := len(csm.UnprocessedBuffer)

  packetParticipants := PacketParticipantsData{}
  requiredSize := GetMemorySize(packetParticipants)

  if buffered >= requiredSize {
    ToObject(csm.UnprocessedBuffer[:], &packetParticipants)

    state.ParticipantData = packetParticipants
    csm.RemoveFirstBytesFromBuffer(requiredSize, state)
  }
}
