package statemachine

import (
    "bytes"
    "encoding/binary"
    "fmt"
    "math"
)

type StateMachine struct {
    UnprocessedBuffer []byte

    LapData               PacketLapData
    TelemetryData         PacketCarTelemetryData
    ParticipantData       PacketParticipantsData
    EventData             PacketEventData
    SessionData           PacketSessionData
    SpeedTraps            [22]SpeedTrap
    CarStatus             PacketCarStatusData
    DistanceHistory       map[uint32]float32
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
}

func CreateStateMachine() StateMachine {
    sm := StateMachine{
        UnprocessedBuffer: []byte{},
        DistanceHistory:   map[uint32]float32{},
    }
    sm.ResetTimers()
    return sm
}

func (s *StateMachine) Process(input []byte) {
    s.UnprocessedBuffer = append(s.UnprocessedBuffer, input...)

    ProcessPacketHeader(s)
}

func (s *StateMachine) RemoveFirstBytesFromBuffer(reducingAmount int) {
    s.UnprocessedBuffer = s.UnprocessedBuffer[reducingAmount:]

    // Clean unprocessable left over data
    // Make sure we have any data we could possibly process
    if len(s.UnprocessedBuffer) == 0 {
        // fmt.Println("buffer is empty, have to wait for new data to be entered")
        return
    }

    indexOfNextStart := bytes.Index(s.UnprocessedBuffer, []byte{0xe4, 0x07})

    if indexOfNextStart > 0 {
        fmt.Printf("Index is %d\n", indexOfNextStart)
        s.UnprocessedBuffer = s.UnprocessedBuffer[indexOfNextStart:]
        s.Process([]byte{})
    } else if indexOfNextStart == -1 {
        // fmt.Println("no start data packet available, clearing buffer")
        s.UnprocessedBuffer = []byte{}
    }
}

func (s *StateMachine) AvailableData() int {
    return len(s.UnprocessedBuffer)
}

func GetMemorySize(input interface{}) int {
    return binary.Size(input)
}

func ToObject(buffer []byte, input interface{}) {
    _ = binary.Read(bytes.NewBuffer(buffer), binary.LittleEndian, input)
}

func (s *StateMachine) GetTimeForDistance(
    totalDistanceTraveled uint32,
    currentTime float32,
    remainingTries uint8,
) float32 {
    if remainingTries == 0 {
        return float32(0)
    }

    value, exists := s.DistanceHistory[totalDistanceTraveled]
    if exists {
        return currentTime - value
    } else {
        return s.GetTimeForDistance(totalDistanceTraveled-1, currentTime, remainingTries-1)
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

func (s *StateMachine) ResetHistory() {
    s.DistanceHistory = map[uint32]float32{}
}
