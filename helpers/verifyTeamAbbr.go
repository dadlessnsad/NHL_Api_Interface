package helpers

var AllAbbreviations = []string{
	"ANA", // Anaheim Ducks
	"ARI", // Arizona Coyotes
	"BOS", // Boston Bruins
	"BUF", // Buffalo Sabres
	"CAR", // Carolina Hurricanes
	"CBJ", // Columbus Blue Jackets
	"CGY", // Calgary Flames
	"CHI", // Chicago Blackhawks
	"COL", // Colorado Avalanche
	"DAL", // Dallas Stars
	"DET", // Detroit Red Wings
	"EDM", // Edmonton Oilers
	"FLA", // Florida Panthers
	"LAK", // Los Angeles Kings
	"MIN", // Minnesota Wild
	"MTL", // Montreal Canadiens
	"NJD", // New Jersey Devils
	"NSH", // Nashville Predators
	"NYI", // New York Islanders
	"NYR", // New York Rangers
	"OTT", // Ottawa Senators
	"PHI", // Philadelphia Flyers
	"PIT", // Pittsburgh Penguins
	"SEA", // Seattle Kraken
	"SJS", // San Jose Sharks
	"STL", // St. Louis Blues
	"TBL", // Tampa Bay Lightning
	"TOR", // Toronto Maple Leafs
	"VAN", // Vancouver Canucks
	"VGK", // Vegas Golden Knights
	"WPG", // Winnipeg Jets
	"WSH", // Washington Capitals
}

func GetAllTeamAbbrs() []string {
	return AllAbbreviations
}

func VerifyTeamAbbr(teamAbbr string) bool {
	for _, abbr := range AllAbbreviations {
		if abbr == teamAbbr {
			return true
		}
	}
	return false
}
