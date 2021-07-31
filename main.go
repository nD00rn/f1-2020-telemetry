package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"net"
	"os"
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
	restOptions, terminalOptions := processOptions()
	fmt.Printf("terminal width is %d\n", terminalOptions.width)

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

	// Set up websocket
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
		for idx := 0; idx < 23; idx++ {
			playerIdOrder[idx] = 255
		}

		// get order of ids by car position
		for playerIndex := 0; playerIndex < 22; playerIndex++ {
			carPos := sm.LapData.LapData[playerIndex].CarPosition
			playerIdOrder[carPos] = playerIndex
		}

		textBuf := ""

		drawMarshalZones(sm, terminalOptions.width, &textBuf, true)

		// Drawing progress lines
		for carPosition := 1; carPosition <= 22; carPosition++ {
			playerIndex := playerIdOrder[carPosition]

			if playerIndex == 255 {
				continue
			}

			if uint8(playerIndex) == sm.LapData.Header.PlayerCarIndex || uint8(playerIndex) == sm.LapData.Header.SecondPlayerCarIndex {
				drawTrackProcess(sm, uint8(playerIndex), terminalOptions.width, &textBuf)
			}
		}
		lastPersonDeltaTime := float32(0)
		ticksPerformed++
		timeBetweenPlayers := fmt.Sprintf(
			"%s%s%s%s%s",
			BgWhite, FgRed,
			floatToTimeStamp(sm.TimeBetweenPlayers()),
			BgReset, FgReset,
		)

		fastestLapTime := sm.FastestLapTime
		if fastestLapTime == math.MaxFloat32 {
			fastestLapTime = float32(0)
		}

		fastestPerson := []byte{'X', 'X', 'X'}
		if sm.FastestLapPlayerIndex != 255 {
			fastestPerson = sm.ParticipantData.Participants[sm.FastestLapPlayerIndex].Name[:][0:3]
		}
		textBuf += fmt.Sprintf(
			"%s: %s  |  %s %s%s[%s | %s]%s%s%s\n",
			"Delta to other player",
			timeBetweenPlayers,
			"Fastest lap time",
			BgWhite,
			FgRed,
			floatToTimeStamp(fastestLapTime),
			fastestPerson,
			FgReset,
			BgReset,
			"                |",
		)

		textBuf += fmt.Sprintln(
			" P  L  NAME | LAST LAP | BEST S1 | BEST S2 | BEST S3 ||  LEADER  |   NEXT |          |",
		)
		for carPosition := 1; carPosition <= 22; carPosition++ {
			playerIndex := playerIdOrder[carPosition]

			if playerIndex == 255 {
				continue
			}

			playerLineColour := BgReset

			lap := sm.LapData.LapData[playerIndex]
			participant := sm.ParticipantData.Participants[playerIndex]

			myDeltaToLeader := float32(0)
			deltaToNext := float32(0)
			totalDistance := uint32(sm.LapData.LapData[playerIndex].TotalDistance)
			sessionTime := sm.LapData.Header.SessionTime

			// Set sessionTime stamp first person crossed this path
			if carPosition == 1 {
				_, exists := csm.DistanceHistory[totalDistance]
				if !exists {
					csm.DistanceHistory[totalDistance] = sessionTime
				}
			} else {
				myDeltaToLeader = csm.GetTimeForDistance(totalDistance, sessionTime, 50)
				deltaToNext = myDeltaToLeader - lastPersonDeltaTime
			}
			lastPersonDeltaTime = myDeltaToLeader

			if uint8(playerIndex) == sm.LapData.Header.PlayerCarIndex {
				sm.TimeToLeaderPlayerOne = myDeltaToLeader
			}
			if uint8(playerIndex) == sm.LapData.Header.SecondPlayerCarIndex {
				sm.TimeToLeaderPlayerTwo = myDeltaToLeader
			}

			if lap.ResultStatus == 0 {
				// You are an invalid player, ignore
				continue
			}

			if uint8(playerIndex) == sm.LapData.Header.PlayerCarIndex {
				playerLineColour = BgYellow
			}
			if uint8(playerIndex) == sm.LapData.Header.SecondPlayerCarIndex {
				// if uint8(playerIndex) != sm.LapData.Header.PlayerCarIndex {
				playerLineColour = BgBlue
			}
			textBuf += fmt.Sprint(playerLineColour)

			if sm.LapData.Header.SessionTime < 1 {
				sm.ResetTimers()
				csm.ResetHistory()
			}

			if lap.BestLapTime < sm.FastestLapTime && lap.BestLapTime > 0 {
				sm.FastestLapTime = lap.BestLapTime
				sm.FastestLapPlayerIndex = playerIndex
			}
			if lap.BestOverallSectorOneTimeInMs < sm.FastestS1Time && lap.BestOverallSectorOneTimeInMs > 0 {
				sm.FastestS1Time = lap.BestOverallSectorOneTimeInMs
				sm.FastestS1PlayerIndex = playerIndex
			}
			if lap.BestOverallSectorTwoTimeInMs < sm.FastestS2Time && lap.BestOverallSectorTwoTimeInMs > 0 {
				sm.FastestS2Time = lap.BestOverallSectorTwoTimeInMs
				sm.FastestS2PlayerIndex = playerIndex
			}
			if lap.BestOverallSectorThreeTimeInMs < sm.FastestS3Time && lap.BestOverallSectorThreeTimeInMs > 0 {
				sm.FastestS3Time = lap.BestOverallSectorThreeTimeInMs
				sm.FastestS3PlayerIndex = playerIndex
			}

			bestLapTime := fmt.Sprintf("%8.3f", lap.BestLapTime)
			bestS1Time := fmt.Sprintf("%7.3f", float32(lap.BestOverallSectorOneTimeInMs)/1000)
			bestS2Time := fmt.Sprintf("%7.3f", float32(lap.BestOverallSectorTwoTimeInMs)/1000)
			bestS3Time := fmt.Sprintf("%7.3f", float32(lap.BestOverallSectorThreeTimeInMs)/1000)
			lastLapTime := floatToTimeStamp(lap.LastLapTime)

			if sm.FastestLapPlayerIndex == playerIndex {
				bestLapTime = BgRed + bestLapTime + playerLineColour
			}

			if sm.FastestS1PlayerIndex == playerIndex {
				bestS1Time = BgRed + bestS1Time + playerLineColour
			}

			if sm.FastestS2PlayerIndex == playerIndex {
				bestS2Time = BgRed + bestS2Time + playerLineColour
			}

			if sm.FastestS3PlayerIndex == playerIndex {
				bestS3Time = BgRed + bestS3Time + playerLineColour
			}

			additionalInformation := "   "
			ersText := ""
			switch sm.CarStatus.CarStatusData[playerIndex].ErsDeployMode {
			case 0:
				ersText = " "
				break
			case 1:
				ersText = " "
				break
			case 2:
				ersText = fmt.Sprintf("%s%s%s", BgGreen, "E", playerLineColour)
				break
			case 3:
				ersText = fmt.Sprintf("%s%s%s", BgGreen, "E", playerLineColour)
				break
			}
			penaltyTime := "  "
			playerPenaltyTime := sm.LapData.LapData[playerIndex].Penalties
			switch playerPenaltyTime {
			case 0:
				penaltyTime = "  "
				break
			default:
				penaltyTime = fmt.Sprintf("%2d", playerPenaltyTime)
				break
			}

			name := getParticipantName(participant)
			if sm.TelemetryData.CarTelemetryData[playerIndex].Drs == 1 {
				additionalInformation = fmt.Sprintf(
					"%s%s%s",
					BgGreen,
					"DRS",
					playerLineColour,
				)
			}
			if sm.LapData.LapData[playerIndex].PitStatus != 0 {
				additionalInformation = fmt.Sprintf(
					"%s%s%s",
					BgRed,
					"PIT",
					playerLineColour,
				)
			}

			additionalInformation = fmt.Sprintf(
				"%s %s %s",
				ersText,
				additionalInformation,
				penaltyTime,
			)

			lapNumber := lap.CurrentLapNum
			lapNumberString := fmt.Sprintf("%2d", lapNumber)
			if sm.LapData.LapData[playerIndex].CurrentLapTime < 3 {
				lapNumberString = fmt.Sprintf(
					"%s%s%s",
					FgRed,
					lapNumberString,
					FgReset,
				)
			}

			textBuf += fmt.Sprintf(
				"%2d %s   %3s | %s | %s | %s | %s || %s | %6.2f | %s |",
				lap.CarPosition,
				lapNumberString,
				name,
				lastLapTime,
				bestS1Time,
				bestS2Time,
				bestS3Time,
				floatToTimeStamp(myDeltaToLeader),
				deltaToNext,
				additionalInformation,
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

func getParticipantName(participant statemachine.ParticipantData) string {
	bytes := participant.Name[0:4]
	for _, b := range bytes {
		if b < 'A' || b > 'Z' {
			return string(bytes)
		}
	}
	return string(participant.Name[:])[0:3]
}

func drawMarshalZones(
	sm statemachine.StateMachine,
	terminalWidth uint,
	textBuf *string,
	showOtherPlayers bool,
) {
	numZones := sm.SessionData.NumMarshalZones
	zones := sm.SessionData.MarshalZones

	*textBuf += "[FLG]|"
	usableTerminalWidth := terminalWidth - 7

	for currentPixel := uint(0); currentPixel < usableTerminalWidth; currentPixel++ {
		zoneToUse := 0
		trackProgress := float32(1) / float32(usableTerminalWidth) * float32(currentPixel)

		trackLength := sm.SessionData.TrackLength
		trackText := " "

		if showOtherPlayers {
			for playerIndex := uint8(0); playerIndex < 22; playerIndex++ {
				if playerIndex == sm.LapData.Header.PlayerCarIndex || playerIndex == sm.LapData.Header.SecondPlayerCarIndex {
					continue
				}
				if sm.LapData.LapData[playerIndex].ResultStatus == 0 {
					continue
				}
				lapPercentage := sm.LapData.LapData[playerIndex].LapDistance / float32(trackLength)

				if isCorrectTrackPixel(lapPercentage, currentPixel, usableTerminalWidth) {
					trackText = ">"
					break
				}
			}
		}

		// Determine which zone data to use
		for i, _ := range zones {
			if zones[0].ZoneStart == 0 {
				break
			}

			if isCorrectZone(i, numZones, trackProgress, zones) {
				zoneToUse = i
				break
			}
		}

		trackColour := BgReset

		switch zones[zoneToUse].ZoneFlag {
		case 1:
			trackText = " "
			trackColour = BgGreen
			break
		case 2:
			trackText = " "
			trackColour = BgBlue
			break
		case 3:
			trackText = " "
			trackColour = BgYellow
			break
		case 4:
			trackText = " "
			trackColour = BgRed
			break
		}

		*textBuf += trackColour + trackText
	}
	*textBuf += BgReset + "|\n"
}

func getCorrectTrackPixel(progress float32, maxPixels uint) uint {
	for current := uint(0); current < maxPixels; current++ {
		if isCorrectTrackPixel(progress, current, maxPixels) {
			return current
		}
	}

	return 0
}

func isCorrectTrackPixel(progress float32, currentPixel uint, maxPixels uint) bool {
	current := float32(currentPixel) / float32(maxPixels)
	next := float32(currentPixel+1) / float32(maxPixels)

	return progress >= current && progress <= next
}

func isCorrectZone(zoneIndex int, numZones uint8, progress float32, zones [21]statemachine.MarshalZone) bool {
	return uint8(zoneIndex+1) == numZones || progress > zones[zoneIndex].ZoneStart && progress < zones[zoneIndex+1].ZoneStart
}

func drawTrackProcess(
	sm statemachine.StateMachine,
	playerIndex uint8,
	terminalWidth uint,
	textBuf *string,
) {
	if playerIndex == 255 {
		return
	}
	name := string(sm.ParticipantData.Participants[playerIndex].Name[:])[0:3]

	lap := sm.LapData.LapData[playerIndex]
	trackLength := sm.SessionData.TrackLength

	lapDistance := lap.LapDistance
	if lapDistance < 0 {
		lapDistance = 0
	}
	lapPercentage := lapDistance / float32(trackLength)
	if lapPercentage >= 1 {
		lapPercentage = 0.99
	}

	availableBlockSpace := terminalWidth - 7
	current := getCorrectTrackPixel(lapPercentage, availableBlockSpace)

	passed := current
	remaining := availableBlockSpace - current - 1

	if remaining < 0 {
		remaining = 0
	}
	if passed < 0 || passed > terminalWidth {
		passed = 0
	}

	fmt.Printf("rem: %d, passed %d\n", remaining, passed)

	arrowHead := ">"
	*textBuf += fmt.Sprintf("[%-3s]|"+strings.Repeat(" ", int(passed))+arrowHead, name)
	*textBuf += fmt.Sprintf(strings.Repeat(" ", int(remaining)) + "|\n")
}

func processOptions() (
	rest.Options,
	TerminalOptions,
) {
	restOptions := rest.DefaultOptions()
	terminalOptions := defaultTerminalOptions()

	flag.UintVar(
		&restOptions.Port,
		"restport",
		8000,
		"port to allow REST communication",
	)
	flag.UintVar(
		&terminalOptions.width,
		"terminalwidth",
		90,
		"terminal width used to generate table",
	)

	flag.Parse()

	return restOptions, terminalOptions
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
