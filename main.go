package main

import (
    "fmt"
    "math"
    "net"
    "os"
    "os/exec"
    "strconv"
    "strings"

    "github.com/nD00rn/f1-2020-telemetry/statemachine"
)

var FgReset = "\033[39m"
var BgReset = "\033[49m"
var FgRed = "\033[31m"
var BgRed = "\033[41m"
var FgGreen = "\033[32m"
var BgGreen = "\033[42m"
var FgYellow = "\033[33m"
var BgYellow = "\033[43m"
var FgBlue = "\033[34m"
var BgBlue = "\033[44m"
var FgPurple = "\033[35m"
var BgPurple = "\033[45m"
var FgCyan = "\033[36m"
var BgCyan = "\033[46m"
var FgGray = "\033[37m"
var BgGray = "\033[47m"
var FgWhite = "\033[97m"
var BgWhite = "\033[47m"

func main() {
    cmd := exec.Command("stty", "size")
    cmd.Stdin = os.Stdin
    out, err := cmd.Output()
    dimensions := string(out)
    dimensions = strings.ReplaceAll(dimensions, "\n", "")
    width, _ := strconv.Atoi(strings.Split(dimensions, " ")[1])

    fmt.Printf("terminal width is %d\n", width)

    addr := net.UDPAddr{
        IP:   net.ParseIP("0.0.0.0"),
        Port: 20777,
    }

    conn, err := net.ListenUDP("udp", &addr)
    if err != nil {
        panic(err)
    }
    defer conn.Close()

    sm := statemachine.CreateStateMachine()

    ticksPerformed := 0

    for {
        buffer := make([]byte, 4096)
        n, _, err := conn.ReadFrom(buffer)
        if err != nil {
            panic(err)
        }

        sm.Process(buffer[:n])

        playerIdOrder := [23]int{}

        // get order of ids by car position
        for i := 0; i < 22; i++ {
            if i == 0 {
                continue
            }

            carPos := sm.LapData.LapData[i].CarPosition
            playerIdOrder[carPos] = i
        }

        for j := 1; j <= 22; j++ {
            i := playerIdOrder[j]

            if i == 0 {
                continue
            }

            if uint8(i) == sm.LapData.Header.PlayerCarIndex || uint8(i) == sm.LapData.Header.SecondPlayerCarIndex {
                drawTrackProcess(sm, uint8(i), width)
            }
        }
        lastPersonDeltaTime := float32(0)
        ticksPerformed++
        // fmt.Printf(
        //     "t%.2f | %d | %d %d %d | %f [%02d , %02d]\n",
        //     sm.LapData.Header.SessionTime,
        //     ticksPerformed,
        //     sm.FastestS1Time,
        //     sm.FastestS2Time,
        //     sm.FastestS3Time,
        //     sm.FastestLapTime,
        //     sm.LapData.Header.PlayerCarIndex,
        //     sm.LapData.Header.SecondPlayerCarIndex,
        // )

        // Player one and player two specific data
        // indexP1 := sm.LapData.Header.PlayerCarIndex
        // indexP2 := sm.LapData.Header.SecondPlayerCarIndex

        // Hmm, what do we want to show
        // Speed ? Gear ?
        // Sector 1, Sector 2

        // currentTimePlayerOne := sm.LapData.LapData[indexP1].CurrentLapTime
        // sectorOnePlayerOne := float32(sm.LapData.LapData[indexP1].SectorOneTimeInMs) / 1000
        // sectorTwoPlayerOne := float32(sm.LapData.LapData[indexP1].SectorTwoTimeInMs) / 1000
        //
        // currentTimePlayerTwo := float32(0)
        // sectorOnePlayerTwo := float32(0)
        // sectorTwoPlayerTwo := float32(0)
        // if indexP2 != 255 {
        //     currentTimePlayerTwo = sm.LapData.LapData[indexP2].CurrentLapTime
        //     sectorOnePlayerTwo = float32(sm.LapData.LapData[indexP2].SectorOneTimeInMs) / 1000
        //     sectorTwoPlayerTwo = float32(sm.LapData.LapData[indexP2].SectorTwoTimeInMs) / 1000
        // }
        // fmt.Println()
        timeBetweenPlayers := fmt.Sprintf("%s%s %5.2f %s%s", BgWhite, FgRed , sm.TimeBetweenPlayers() , BgReset, FgReset)
        fastestLapTime := sm.FastestLapTime
        if fastestLapTime == math.MaxFloat32 {
            fastestLapTime = float32(0)
        }
        fmt.Printf(
            "%s: %s   |   %s %5.2f\n",
            "delta to other player",
            timeBetweenPlayers,
            "Fastest lap time",
            fastestLapTime,
        )
        // fmt.Println("                    [Player one]   ||    [Player two]")
        // fmt.Printf(
        //     "%-20s :  %6.2f || %6.2f\n",
        //     "Current lap time",
        //     currentTimePlayerOne,
        //     currentTimePlayerTwo,
        // )
        // fmt.Printf(
        //     "%-20s :  %6.2f || %6.2f\n",
        //     "Sector one",
        //     sectorOnePlayerOne,
        //     sectorOnePlayerTwo,
        // )
        // fmt.Printf(
        //     "%-20s :  %6.2f || %6.2f\n",
        //     "Sector two",
        //     sectorTwoPlayerOne,
        //     sectorTwoPlayerTwo,
        // )

        fmt.Println(
            "[NAME]          | BEST S1 | BEST S2 | BEST S3 || LEADER | TO NEXT",
        )
        for j := 1; j <= 22; j++ {
            i := playerIdOrder[j]

            if i == 0 {
                continue
            }

            // telemetry := sm.TelemetryData.CarTelemetryData[i]
            lap := sm.LapData.LapData[i]
            participant := sm.ParticipantData.Participants[i]

            myDeltaToLeader := float32(0)
            deltaToNext := float32(0)
            totalDistance := uint32(sm.LapData.LapData[i].TotalDistance)
            time := sm.LapData.Header.SessionTime

            // Set time stamp first person crossed this path
            if j == 1 {
                sm.DistanceHistory[totalDistance] = time
            } else {
                myDeltaToLeader = sm.GetTimeForDistance(totalDistance, time, 50)
                deltaToNext = myDeltaToLeader - lastPersonDeltaTime
            }
            lastPersonDeltaTime = myDeltaToLeader

            if uint8(i) == sm.LapData.Header.PlayerCarIndex {
                sm.TimeToLeaderPlayerOne = myDeltaToLeader
            }
            if uint8(i) == sm.LapData.Header.SecondPlayerCarIndex {
                sm.TimeToLeaderPlayerTwo = myDeltaToLeader
            }

            if lap.ResultStatus == 0 {
                // You are an invalid player, ignore
                continue
            }

            if uint8(i) == sm.LapData.Header.PlayerCarIndex {
                fmt.Print(BgGray)
            }
            if uint8(i) == sm.LapData.Header.SecondPlayerCarIndex {
                fmt.Print(BgYellow)
            }

            if sm.LapData.Header.SessionTime < 1 {
                sm.ResetTimers()
                sm.ResetHistory()
            }

            if lap.BestLapTime < sm.FastestLapTime && lap.BestLapTime > 0 {
                sm.FastestLapTime = lap.BestLapTime
                sm.FastestLapPlayerIndex = i
            }
            if lap.BestOverallSectorOneTimeInMs < sm.FastestS1Time && lap.BestOverallSectorOneTimeInMs > 0 {
                sm.FastestS1Time = lap.BestOverallSectorOneTimeInMs
                sm.FastestS1PlayerIndex = i
            }
            if lap.BestOverallSectorTwoTimeInMs < sm.FastestS2Time && lap.BestOverallSectorTwoTimeInMs > 0 {
                sm.FastestS2Time = lap.BestOverallSectorTwoTimeInMs
                sm.FastestS2PlayerIndex = i
            }
            if lap.BestOverallSectorThreeTimeInMs < sm.FastestS3Time && lap.BestOverallSectorThreeTimeInMs > 0 {
                sm.FastestS3Time = lap.BestOverallSectorThreeTimeInMs
                sm.FastestS3PlayerIndex = i
            }

            bestLapTime := fmt.Sprintf("%8.3f", lap.BestLapTime)
            bestS1Time := fmt.Sprintf("%7.3f", float32(lap.BestOverallSectorOneTimeInMs)/1000)
            bestS2Time := fmt.Sprintf("%7.3f", float32(lap.BestOverallSectorTwoTimeInMs)/1000)
            bestS3Time := fmt.Sprintf("%7.3f", float32(lap.BestOverallSectorThreeTimeInMs)/1000)

            if sm.FastestLapPlayerIndex == i {
                bestLapTime = FgPurple + bestLapTime + FgReset
            }

            if sm.FastestS1PlayerIndex == i {
                bestS1Time = FgPurple + bestS1Time + FgReset
            }

            if sm.FastestS2PlayerIndex == i {
                bestS2Time = FgPurple + bestS2Time + FgReset
            }

            if sm.FastestS3PlayerIndex == i {
                bestS3Time = FgPurple + bestS3Time + FgReset
            }

            name := string(participant.Name[:])[0:3]
            if sm.TelemetryData.CarTelemetryData[i].Drs == 1 {
                name = fmt.Sprintf("%s%s%s",
                    FgGreen,
                    name,
                    FgReset,
                )
            }
            fmt.Printf(
                "[ %3s ] p%-2d l%-2d | %s | %s | %s || %6.2f | %6.2f",
                name,
                lap.CarPosition,
                lap.CurrentLapNum,
                bestS1Time,
                bestS2Time,
                bestS3Time,
                myDeltaToLeader,
                deltaToNext,
            )

            fmt.Print(FgReset + BgReset + "\n")
        }
    }
}

func drawTrackProcess(sm statemachine.StateMachine, playerIndex uint8, terminalWidth int) {
    if playerIndex == 255 {
        // fmt.Printf("[%-3s]|" + strings.Repeat("=", 0) + ">\n", name)
        return
    }
    name := string(sm.ParticipantData.Participants[playerIndex].Name[:])[0:3]

    lap := sm.LapData.LapData[playerIndex]
    trackLength := sm.SessionData.TrackLength

    lapPercentage := lap.LapDistance / float32(trackLength)
    blocks := int(float32(terminalWidth-6) * lapPercentage)
    if blocks < 0 {
        blocks = 0
    }
    fmt.Printf("[%-3s]|"+strings.Repeat("=", blocks)+">\n", name)
}
