package main

import "fmt"

func seriesTypeToText(id int) string {
	seriesTypes := []string{"bo1", "bo3", "bo5"}
	return seriesTypes[id]
}

func teamsToVS(radiantTeam, direTeam int) string {
	if radiantTeam > direTeam {
		return fmt.Sprintf("#t%dvs%drsltdtk", direTeam, radiantTeam)
	} else {
		return fmt.Sprintf("#t%dvs%drsltdtk", radiantTeam, direTeam)
	}
}

func makeMatchText(matchData oDotaMatchData) string {
	radiantName := matchData.RadiantTeam.Name
	direName := matchData.DireTeam.Name
	radiantID := matchData.RadiantTeam.TeamID
	direID := matchData.DireTeam.TeamID
	radiantScore := matchData.RadiantScore
	direScore := matchData.DireScore
	hours, minutes, seconds := secondsToTime(matchData.Duration)
	leagueName := matchData.League.Name
	seriesID := matchData.SeriesID
	seriesType := matchData.SeriesType

	if radiantName == "" {
		radiantName = matchData.RadiantName
	}
	if direName == "" {
		direName = matchData.DireName
	}

	if radiantName == "" {
		radiantName = "Radiant"
	}
	if direName == "" {
		direName = "Dire"
	}

	if leagueName == "" {
		leaguesData := getLeaguesData()
		leagueName = getLeagueName(matchData.Leagueid, leaguesData)
	}

	seriesText := seriesTypeToText(seriesType)
	vsText := teamsToVS(radiantID, direID)

	if matchData.RadiantWin {
		if seriesID != 0 {
			return fmt.Sprintf("🏆 [rsltdtk|%s] %d - %d %s\n%d:%02d:%02d @ %s, %s\n#dota2 #s%d@rsltdtk %s",
				radiantName, radiantScore, direScore, direName, hours, minutes, seconds, leagueName, seriesText, seriesID, vsText)
		} else {
			return fmt.Sprintf("🏆 [rsltdtk|%s] %d - %d %s\n%d:%02d:%02d @ %s\n#dota2 %s",
				radiantName, radiantScore, direScore, direName, hours, minutes, seconds, leagueName, vsText)
		}
	} else {
		if seriesID != 0 {
			return fmt.Sprintf("%s %d - %d 🏆 [rsltdtk|%s]\n%d:%02d:%02d @ %s, %s\n#dota2 #s%d@rsltdtk %s",
				radiantName, radiantScore, direScore, direName, hours, minutes, seconds, leagueName, seriesText, seriesID, vsText)
		} else {
			return fmt.Sprintf("%s %d - %d 🏆 [rsltdtk|%s]\n%d:%02d:%02d @ %s\n#dota2 %s",
				radiantName, radiantScore, direScore, direName, hours, minutes, seconds, leagueName, vsText)
		}
	}
}
