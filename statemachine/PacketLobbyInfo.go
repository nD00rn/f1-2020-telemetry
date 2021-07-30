package statemachine

type LobbyInfoData struct {
  AiControlled uint8
  TeamId       uint8
  Nationality  uint8
  Name         [48]byte

  Ready uint8
}

type PacketLobbyInfoData struct {
  Header        PacketHeader
  NumPlayers    uint8
  LobbyInfoData [22]LobbyInfoData
}

func ProcessPacketLobbyInfo(csm *CommunicationStateMachine, state *StateMachine) {
  buffered := len(csm.UnprocessedBuffer)

  packetLobbyInfoData := PacketLobbyInfoData{}
  requiredSize := GetMemorySize(packetLobbyInfoData)

  if buffered >= requiredSize {
    ToObject(csm.UnprocessedBuffer[:], &packetLobbyInfoData)

    csm.RemoveFirstBytesFromBuffer(requiredSize, state)
  }
}
