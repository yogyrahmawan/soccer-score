package main

import (
	"log"
	"os"
	"strings"

	"github.com/Syfaro/telegram-bot-api"
	"github.com/yogyrahmawan/soccer-score"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("your_api")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	err = bot.UpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	//create session
	session, err := app.Session()
	if err != nil {
		log.Fatal("Cannot create session, err = ", err.Error())
		os.Exit(1)
	}
	defer session.Close()

	for update := range bot.Updates {
		msgText := strings.ToLower(strings.Replace(update.Message.Text, "/", "", -1))
		var msg tgbotapi.MessageConfig
		switch {
		case msgText == "help":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "/list - see available leagues")
		case msgText == "start":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to soccer live score bot\nPress /list to see available leagues")
		case msgText == "list":
			leagueList, err := app.LeagueList(session)
			if err != nil {
				log.Printf("internal error, err=%v", err)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "sorry,we encounter error in our system")
			} else {
				msgTextLeagueList := ""
				for _, v := range leagueList {
					msgTextLeagueList += "/" + v.Key + "\n"
				}

				msg = tgbotapi.NewMessage(update.Message.Chat.ID, msgTextLeagueList)
			}
		default:
			leagueMapper, err := app.GetLeagueMapperByKey(session, msgText)
			if err != nil {
				log.Printf("internal error, err=%v", err)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "sorry, we can not get any data")
			} else {
				node, err := app.GetParseableHTML(leagueMapper.URL)
				if err != nil {
					log.Printf("internal error, err=%v", err)
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "sorry,we can not get any data")
				} else {
					matches := app.LivescoreParser(node)
					parsedMatch := ""
					for _, v := range matches {
						if len(strings.Split(v.Time, ":")) > 1 {
							parsedMatch += v.Time + " UTC \t"
						} else {
							parsedMatch += v.Time + "\t"
						}

						parsedMatch += v.HomeTeam + "\t" + v.Score + "\t" + v.AwayTeam + "\n"
					}
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, parsedMatch)
				}
			}
		}
		bot.SendMessage(msg)

	}
}
