package statemachine

import (
    "bytes"
    "encoding/binary"
    "fmt"
    "math"
)

type CommunicationStateMachine struct {
    UnprocessedBuffer []byte
    DistanceHistory   map[uint32]float32
}

type StateMachine struct {
    LapData               PacketLapData
    TelemetryData         PacketCarTelemetryData
    ParticipantData       PacketParticipantsData
    EventData             PacketEventData
    SessionData           PacketSessionData
    SpeedTraps            [22]SpeedTrap
    CarStatus             PacketCarStatusData
    FastestLapPlayerIndex int
    FastestLapTime        float32
    FastestS1PlayerIndex  int
    FastestS1Time         uint16
    FastestS2PlayerIndex  int
    FastestS2Time         uint16
    FastestS3PlayerIndex  int
    FastestS3Time         uint16
    TimeToLeaderPlayerOne float32
    TimeToLeaderPlayerTwo float32
    PlayerOneIndex        uint8
    PlayerTwoIndex        uint8
}

func CreateCommunicationStateMachine() CommunicationStateMachine {
    csm := CommunicationStateMachine{
        UnprocessedBuffer: []byte{},
        DistanceHistory:   map[uint32]float32{},
    }

    return csm
}

func CreateStateMachine() StateMachine {
    sm := StateMachine{
        // UnprocessedBuffer: []byte{},
        // DistanceHistory:   map[uint32]float32{},
    }
    sm.ResetTimers()
    return sm
}

func (csm *CommunicationStateMachine) Process(input []byte, sm *StateMachine) {
    csm.UnprocessedBuffer = append(csm.UnprocessedBuffer, input...)

    ProcessPacketHeader(csm, sm)
}

func (csm *CommunicationStateMachine) RemoveFirstBytesFromBuffer(reducingAmount int, sm *StateMachine) {
    csm.UnprocessedBuffer = csm.UnprocessedBuffer[reducingAmount:]

    // Clean unprocessable left over data
    // Make sure we have any data we could possibly process
    if len(csm.UnprocessedBuffer) == 0 {
        // fmt.Println("buffer is empty, have to wait for new data to be entered")
        return
    }

    indexOfNextStart := bytes.Index(csm.UnprocessedBuffer, []byte{0xe4, 0x07})

    if indexOfNextStart > 0 {
        fmt.Printf("Index is %d\n", indexOfNextStart)
        csm.UnprocessedBuffer = csm.UnprocessedBuffer[indexOfNextStart:]
        csm.Process([]byte{}, sm)
    } else if indexOfNextStart == -1 {
        // fmt.Println("no start data packet available, clearing buffer")
        csm.UnprocessedBuffer = []byte{}
    }
}

func (csm *CommunicationStateMachine) AvailableData() int {
    return len(csm.UnprocessedBuffer)
}

func GetMemorySize(input interface{}) int {
    return binary.Size(input)
}

func ToObject(buffer []byte, input interface{}) {
    _ = binary.Read(bytes.NewBuffer(buffer), binary.LittleEndian, input)
}

func (csm *CommunicationStateMachine) GetTimeForDistance(
    totalDistanceTraveled uint32,
    currentTime float32,
    remainingTries uint8,
) float32 {
    if remainingTries == 0 {
        return float32(0)
    }

    value, exists := csm.DistanceHistory[totalDistanceTraveled]
    if exists {
        return currentTime - value
    } else {
        return csm.GetTimeForDistance(totalDistanceTraveled-1, currentTime, remainingTries-1)
    }
}

func (s *StateMachine) ResetTimers() {
    s.FastestLapTime = math.MaxFloat32
    s.FastestS1Time = math.MaxUint16
    s.FastestS2Time = math.MaxUint16
    s.FastestS3Time = math.MaxUint16

    s.TimeToLeaderPlayerOne = float32(0)
    s.TimeToLeaderPlayerTwo = float32(0)
}

func (s *StateMachine) TimeBetweenPlayers() float32 {
    return s.TimeToLeaderPlayerOne - s.TimeToLeaderPlayerTwo
}

func (csm *CommunicationStateMachine) ResetHistory() {
    csm.DistanceHistory = map[uint32]float32{}
}
