package statemachine

import (
  "fmt"
)

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

func ProcessPacketLobbyInfo(state *StateMachine) {
  buffered := len(state.UnprocessedBuffer)

  packetLobbyInfoData := PacketLobbyInfoData{}
  requiredSize := GetMemorySize(packetLobbyInfoData)

  if buffered >= requiredSize {
    ToObject(state.UnprocessedBuffer[:], &packetLobbyInfoData)

    println(fmt.Sprintf("data lobby info: %+v", packetLobbyInfoData))

    state.RemoveFirstBytesFromBuffer(requiredSize)
  }
}
