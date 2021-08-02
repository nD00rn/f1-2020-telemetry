package statemachine

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

func ProcessPacketHeader(csm *CommunicationStateMachine, state *StateMachine) {
    buffered := len(csm.UnprocessedBuffer)

    packetHeader := PacketHeader{}
    requiredSize := GetMemorySize(packetHeader)

    if buffered >= requiredSize {
        header := packetHeader

        ToObject(csm.UnprocessedBuffer[:], &header)

        if header.PackageFormat != 2020 {
            csm.RemoveFirstBytesFromBuffer(requiredSize, state)
            return
        }

        switch header.PacketId {
        case 0:
            ProcessPacketMotion(csm, state)
            break

        case 1:
            processPacketSession(csm, state)
            break

        case 2:
            ProcessPacketLapData(csm, state)
            break

        case 3:
            ProcessPacketEvent(csm, state)
            break

        case 4:
            ProcessPacketParticipants(csm, state)
            break

        case 5:
            ProcessPacketCarSetups(csm, state)
            break

        case 6:
            ProcessPacketCarTelemetry(csm, state)
            break

        case 7:
            ProcessPacketCarStatus(csm, state)
            break

        case 8:
            ProcessPacketFinalClassification(csm, state)
            break

        case 9:
            ProcessPacketLobbyInfo(csm, state)
            break

        default:
            csm.RemoveFirstBytesFromBuffer(1, state)
            break
        }
    }
}
