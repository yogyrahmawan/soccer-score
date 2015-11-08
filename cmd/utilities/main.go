package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/yogyrahmawan/soccer-score"

	"labix.org/v2/mgo"
)

var (
	errInvalidURL = errors.New("Invalid URL")
)

func usage() {
	fmt.Print("This is utilities to soccer score. \n\nUsage:\n\n")
	fmt.Print("	./utilities <command> [arguments]\n\n")
	fmt.Print("The commands are:\n\n")
	fmt.Print("	add_league <key> <url> <source> 	Add league, example: add_league premier http://livescore.com livescore\n")
	fmt.Print("	list_league 		List all league\n")
	fmt.Print("	delete_league <id>	Delete league by id\n")
	fmt.Println()
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	commands := map[string]int{"add_league": 3, "list_league": 0, "delete_league": 1}

	valid := false
	for k, v := range commands {
		if os.Args[1] == k {
			if len(os.Args) < v+2 {
				usage()
				os.Exit(1)
			}
		}

		valid = true
		break
	}

	if !valid {
		usage()
		os.Exit(1)
	}

	session, err := app.Session()
	if err != nil {
		log.Fatal("Cannot create session, err = ", err.Error())
		os.Exit(1)
	}
	defer session.Close()

	switch os.Args[1] {
	case "add_league":
		err = commandAddLeague(session, os.Args[2:])
	case "list_league":
		err = commandListLeague(session)
	case "delete_league":
		err = commandDeleteLeague(session, os.Args[2:])
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

func commandAddLeague(session *mgo.Session, args []string) error {
	if _, err := url.Parse(args[1]); err != nil {
		return errInvalidURL
	}
	if _, err := app.NewLeagueMapper(session, strings.Replace(args[0], " ", "", -1), args[1], args[2]); err != nil {
		return err
	}

	return nil
}

func commandListLeague(session *mgo.Session) error {
	ll, err := app.LeagueList(session)
	if err != nil {
		return err
	}

	fmt.Println("ID		Key		URL		Source")
	for _, val := range ll {
		fmt.Printf("%v		%v		%v		%v\n", val.ID, val.Key, val.URL, val.SourceKey)
	}
	return nil
}

func commandDeleteLeague(session *mgo.Session, args []string) error {
	if err := app.RemoveLeagueMappersByID(session, args[0]); err != nil {
		return err
	}
	return nil
}
