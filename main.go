package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/nD00rn/f1-2020-telemetry/rest"
	"github.com/nD00rn/f1-2020-telemetry/statemachine"
	"github.com/nD00rn/f1-2020-telemetry/websocket"
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

var localSm *statemachine.StateMachine

func main() {
	restOptions := processOptions()

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

	csm := statemachine.CreateCommunicationStateMachine()
	sm := statemachine.CreateStateMachine()
	localSm = &sm

	// Start the REST server
	rest.SetUpRestApiRouter(restOptions, &sm)

	go websocket.Broadcast()

	go constantStreamWebSocket()

	// Mainly debug information
	ticksPerformed := 0

	for {
		// the number of incoming bytes lays around the 1500 bytes per request.
		// this should give some space whenever a bigger request comes in.
		buffer := make([]byte, 4096)
		n, _, err := conn.ReadFrom(buffer)
		if err != nil {
			panic(err)
		}

		csm.Process(buffer[:n], &sm)

		sm.PlayerOneIndex = sm.LapData.Header.PlayerCarIndex
		sm.PlayerTwoIndex = sm.LapData.Header.SecondPlayerCarIndex

		playerIdOrder := [23]int{}

		// get order of ids by car position
		for i := 0; i < 22; i++ {
			if i == 0 {
				continue
			}

			carPos := sm.LapData.LapData[i].CarPosition
			playerIdOrder[carPos] = i
		}

        textBuf := ""

        for j := 1; j <= 22; j++ {
			i := playerIdOrder[j]

			if i == 0 {
				continue
			}

			if uint8(i) == sm.LapData.Header.PlayerCarIndex || uint8(i) == sm.LapData.Header.SecondPlayerCarIndex {
				drawTrackProcess(sm, uint8(i), width, textBuf)
			}
		}
		lastPersonDeltaTime := float32(0)
		ticksPerformed++
		timeBetweenPlayers := fmt.Sprintf("%s%s %5.2f %s%s", BgWhite, FgRed, sm.TimeBetweenPlayers(), BgReset, FgReset)

		fastestLapTime := sm.FastestLapTime
		if fastestLapTime == math.MaxFloat32 {
			fastestLapTime = float32(0)
		}

        textBuf += fmt.Sprintf(
			"%s: %s   |   %s [%s | %s]\n",
			"delta to other player",
			timeBetweenPlayers,
			"Fastest lap time",
			floatToTimeStamp(fastestLapTime),
			sm.ParticipantData.Participants[sm.FastestLapPlayerIndex].Name[:][0:3],
		)

        textBuf += fmt.Sprintln(
			"[NAME]          | LAST LAP | BEST S1 | BEST S2 | BEST S3 || LEADER | TO NEXT |",
		)
		for j := 1; j <= 22; j++ {
			i := playerIdOrder[j]

			if i == 0 {
				continue
			}

			lap := sm.LapData.LapData[i]
			participant := sm.ParticipantData.Participants[i]

			myDeltaToLeader := float32(0)
			deltaToNext := float32(0)
			totalDistance := uint32(sm.LapData.LapData[i].TotalDistance)
			sessionTime := sm.LapData.Header.SessionTime

			// Set sessionTime stamp first person crossed this path
			if j == 1 {
				_, exists := csm.DistanceHistory[totalDistance]
				if !exists {
					csm.DistanceHistory[totalDistance] = sessionTime
				}
			} else {
				myDeltaToLeader = csm.GetTimeForDistance(totalDistance, sessionTime, 50)
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
				textBuf += fmt.Sprint(BgGray)
			}
			if uint8(i) == sm.LapData.Header.SecondPlayerCarIndex {
				textBuf += fmt.Sprint(BgYellow)
			}

			if sm.LapData.Header.SessionTime < 1 {
				sm.ResetTimers()
				csm.ResetHistory()
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
			lastLapTime := floatToTimeStamp(lap.LastLapTime)

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
            textBuf += fmt.Sprintf(
				"[ %3s ] p%-2d l%-2d | %s | %s | %s | %s || %6.2f | %6.2f   ",
				name,
				lap.CarPosition,
				lap.CurrentLapNum,
				lastLapTime,
				bestS1Time,
				bestS2Time,
				bestS3Time,
				myDeltaToLeader,
				deltaToNext,
			)

            textBuf += fmt.Sprint(FgReset + BgReset + "\n")
		}

        f, err := os.Create("/tmp/f1.screen.tmp")
        if err != nil {
            panic(err)
        }
		_, _ = fmt.Fprint(f, textBuf)
		_ = f.Close()
		_ = os.Rename("/tmp/f1.screen.tmp", "/tmp/f1.screen.txt")
	}
}

func drawTrackProcess(sm statemachine.StateMachine, playerIndex uint8, terminalWidth int, textBuf string) {
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
    textBuf += fmt.Sprintf("[%-3s]|"+strings.Repeat("=", blocks)+">\n", name)
}

func processOptions() rest.Options {
	restOptions := rest.DefaultOptions()
	flag.UintVar(
		&restOptions.Port,
		"restport",
		8000,
		"port to allow REST communication",
	)

	flag.Parse()

	return restOptions
}

func constantStreamWebSocket() {
	for {
		time.Sleep(250 * time.Millisecond)

		marshal, err := json.Marshal(localSm)
		if err != nil {
			continue
		}

		websocket.BroadcastMessage(string(marshal))
	}
}

func floatToTimeStamp(input float32) string {
	minutes := uint8(input / 60)
	seconds := input - float32(minutes*60)
	return fmt.Sprintf("%d:%06.3f", minutes, seconds)
}
