package statemachine

import (
  "fmt"
)

type PacketHeader struct {
  PackageFormat        uint16
  GameMajorVersion     uint8
  GameMinorVersion     uint8
  PacketVersion        uint8
  PacketId             uint8
  SessionUid           uint64
  SessionTime          float32
  FrameIdentifier      uint32
  PlayerCarIndex       uint8
  SecondPlayerCarIndex uint8
}

func ProcessPacketHeader(state *StateMachine) {
  buffered := len(state.UnprocessedBuffer)

  packetHeader := PacketHeader{}
  requiredSize := GetMemorySize(packetHeader)

  if buffered >= requiredSize {
    header := packetHeader

    ToObject(state.UnprocessedBuffer[:], &header)

    if header.PackageFormat != 2020 {
      fmt.Println("improbable game version, discarding packet")
      state.RemoveFirstBytesFromBuffer(requiredSize)
      return
    }

    switch header.PacketId {
    case 0:
      ProcessPacketMotion(state)
      break

    case 1:
      processPacketSession(state)
      break

    case 2:
      ProcessPacketLapData(state)
      break

    case 3:
      ProcessPacketEvent(state)
      break

    case 4:
      ProcessPacketParticipants(state)
      break

    case 5:
      ProcessPacketCarSetups(state)
      break

    case 6:
      ProcessPacketCarTelemetry(state)
      break

    case 7:
      ProcessPacketCarStatus(state)
      break

    case 8:
      ProcessPacketFinalClassification(state)
      break

    case 9:
      ProcessPacketLobbyInfo(state)
      break

    default:
      fmt.Println("[switch] no idea on how to process this packet")
      break
    }
  }
}
