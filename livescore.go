package app

import (
	"strings"

	"github.com/yhat/scrape"

	"golang.org/x/net/html"
)

const (
	classTag        = "class"
	classContentTag = "content"
	classRowGray    = "row-gray"
	classMinElmt    = "min"
	classTrightElmt = "tright"
	classPlyElmt    = "ply"
	classScoreLink  = "scorelink"
)

//LivescoreParser parse livescore
func LivescoreParser(root *html.Node) []Match {
	var matches []Match

	contentElmt, contentOK := scrape.Find(root, scrape.ByClass(classContentTag))
	if contentOK {
		//find all row-gray
		rowGrayMatcher := func(n *html.Node) bool {
			classes := strings.Fields(scrape.Attr(n, "class"))
			for _, c := range classes {
				if c == classRowGray {
					parentClasses := strings.Fields(scrape.Attr(n.Parent, "class"))
					for _, pc := range parentClasses {
						if pc == classContentTag {
							return true
						}
					}
				}
			}
			return false
		}
		rows := scrape.FindAll(contentElmt, rowGrayMatcher)

		matchChann := make(chan Match)
		for _, rowElmt := range rows {
			go func(rowElmt *html.Node) {
				var time string
				var homeTeam string
				var awayTeam string
				var score string

				timeElmt, timeElmtOK := scrape.Find(rowElmt, scrape.ByClass(classMinElmt))
				if timeElmtOK {
					time = scrape.Text(timeElmt)
				}

				scoreElmt, scoreElmtOK := scrape.Find(rowElmt, scrape.ByClass(classScoreLink))
				if scoreElmtOK {
					score = scrape.Text(scoreElmt)
				}

				teamElmts := scrape.FindAll(rowElmt, scrape.ByClass(classPlyElmt))
				for i := 0; i < len(teamElmts); i++ {
					teamElmt := teamElmts[i]
					if i%2 == 0 {
						homeTeam = scrape.Text(teamElmt)
					} else {
						awayTeam = scrape.Text(teamElmt)
					}
				}
				match := Match{
					HomeTeam: homeTeam,
					AwayTeam: awayTeam,
					Score:    score,
					Time:     time,
				}

				matchChann <- match
			}(rowElmt)
		}

		for i := 0; i < len(rows); i++ {
			select {
			case m := <-matchChann:
				matches = append(matches, m)
			}
		}
		close(matchChann)
	}
	return matches
}
