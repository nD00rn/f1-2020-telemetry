package terminal

import (
    "fmt"
    "math"

    "github.com/nD00rn/f1-2020-telemetry/statemachine"
    "github.com/nD00rn/f1-2020-telemetry/util"
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

type Options struct {
    Width uint
}

func DefaultOptions() Options {
    return Options{
        Width: 84,
    }
}

func CreateBuffer(
    sm statemachine.StateMachine,
    csm statemachine.CommunicationStateMachine,
    terminalOptions Options,
) string {
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

    drawMarshalZones(sm, terminalOptions.Width, &textBuf, true)

    // Drawing progress lines
    for carPosition := 1; carPosition <= 22; carPosition++ {
        playerIndex := playerIdOrder[carPosition]

        if playerIndex == 255 {
            continue
        }

        if uint8(playerIndex) == sm.LapData.Header.PlayerCarIndex || uint8(playerIndex) == sm.LapData.Header.SecondPlayerCarIndex {
            drawTrackProcess(sm, uint8(playerIndex), terminalOptions.Width, &textBuf)
        }
    }
    lastPersonDeltaTime := float32(0)
    timeBetweenPlayers := fmt.Sprintf(
        "%s%s%s%s%s",
        BgWhite, FgRed,
        util.FloatToTimeStamp(sm.TimeBetweenPlayers()),
        BgReset, FgReset,
    )

    fastestS1 := uint16(math.MaxUint16)
    fastestS2 := uint16(math.MaxUint16)
    fastestS3 := uint16(math.MaxUint16)
    fastestLap := float32(math.MaxFloat32)
    fastestS1Idx := 255
    fastestS2Idx := 255
    fastestS3Idx := 255
    fastestLapIdx := 255

    for i, lapData := range sm.LapData.LapData {
        if lapData.ResultStatus == 0 {
            continue
        }

        bestS1 := lapData.BestOverallSectorOneTimeInMs
        bestS2 := lapData.BestOverallSectorTwoTimeInMs
        bestS3 := lapData.BestOverallSectorThreeTimeInMs
        bestLap := lapData.BestLapTime

        if bestS1 < fastestS1 && bestS1 > 1 {
            fastestS1 = bestS1
            fastestS1Idx = i
        }
        if bestS2 < fastestS2 && bestS2 > 1 {
            fastestS2 = bestS2
            fastestS2Idx = i
        }
        if bestS3 < fastestS3 && bestS3 > 1 {
            fastestS3 = bestS3
            fastestS3Idx = i
        }
        if bestLap < fastestLap && bestLap > float32(1) {
            fastestLap = bestLap
            fastestLapIdx = i
        }
    }

    if fastestLap == math.MaxFloat32 {
        fastestLap = float32(0)
    }

    fastestPerson := []byte{'X', 'X', 'X'}
    if fastestLapIdx != 255 {
        fastestPerson = sm.ParticipantData.Participants[fastestLapIdx].Name[:][0:3]
    }
    textBuf += fmt.Sprintf(
        "%s: %s  |  %s %s%s[%s | %s]%s%s%s\n",
        "Delta to other player",
        timeBetweenPlayers,
        "Fastest lap time",
        BgWhite,
        FgRed,
        util.FloatToTimeStamp(fastestLap),
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
            playerLineColour = BgBlue
        }
        textBuf += fmt.Sprint(playerLineColour)

        if sm.LapData.Header.SessionTime < 1 || sm.LastSessionUid != sm.LapData.Header.SessionUid {
            sm.LastSessionUid = sm.LapData.Header.SessionUid
            sm.ResetTimers()
            csm.ResetHistory()
        }

        bestLapTime := fmt.Sprintf("%8.3f", lap.BestLapTime)
        bestS1Time := fmt.Sprintf("%7.3f", float32(lap.BestOverallSectorOneTimeInMs)/1000)
        bestS2Time := fmt.Sprintf("%7.3f", float32(lap.BestOverallSectorTwoTimeInMs)/1000)
        bestS3Time := fmt.Sprintf("%7.3f", float32(lap.BestOverallSectorThreeTimeInMs)/1000)
        lastLapTime := util.FloatToTimeStamp(lap.LastLapTime)

        if fastestLapIdx == playerIndex {
            bestLapTime = BgRed + bestLapTime + playerLineColour
        }

        if fastestS1Idx == playerIndex {
            bestS1Time = BgRed + bestS1Time + playerLineColour
        }

        if fastestS2Idx == playerIndex {
            bestS2Time = BgRed + bestS2Time + playerLineColour
        }

        if fastestS3Idx == playerIndex {
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
            deltaToNext = float32(0)
            myDeltaToLeader = float32(0)
            additionalInformation = fmt.Sprintf(
                "%s%s%s",
                BgRed,
                "PIT",
                playerLineColour,
            )
        }

        switch sm.LapData.LapData[playerIndex].ResultStatus {
        case 3:
            deltaToNext = float32(0)
            myDeltaToLeader = float32(0)
            additionalInformation = fmt.Sprintf(
                "%s%s%s",
                BgYellow,
                "FIN",
                playerLineColour,
            )
            break
        case 4:
            deltaToNext = float32(0)
            myDeltaToLeader = float32(0)
            additionalInformation = fmt.Sprintf(
                "%s%s%s",
                BgRed,
                "DSQ",
                playerLineColour,
            )
            break
        case 5:
            deltaToNext = float32(0)
            myDeltaToLeader = float32(0)
            additionalInformation = fmt.Sprintf(
                "%s%s%s",
                BgRed,
                "NCL",
                playerLineColour,
            )
            break
        case 6:
            deltaToNext = float32(0)
            myDeltaToLeader = float32(0)
            additionalInformation = fmt.Sprintf(
                "%s%s%s",
                BgRed,
                "OUT",
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
            util.FloatToTimeStamp(myDeltaToLeader),
            deltaToNext,
            additionalInformation,
        )

        textBuf += fmt.Sprint(FgReset + BgReset + "\n")
    }

    lastStreamBuffer = textBuf
    return textBuf
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
