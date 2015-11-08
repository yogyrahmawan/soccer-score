package app

import (
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"labix.org/v2/mgo"
)

const (
	defaultURI = "mongodb://127.0.0.1/soccerscore"
)

var (
	ses *mgo.Session
)

func getDbURI() string {
	if mgoURI := os.Getenv("MONGODB_URI"); mgoURI != "" {
		return mgoURI
	}
	return defaultURI
}

// getDbName extracts database name from MongoDB URI (e.g. mongodb://localhost/test?replicaSet=test),
// see http://docs.mongodb.org/manual/reference/connection-string/.
func getDbName() string {
	uriParts := strings.SplitN(getDbURI(), "/", 4)

	// get the last part (maybe a name with or without options)
	name := uriParts[len(uriParts)-1]

	// get the first part of name (omit options if any)
	nameParts := strings.SplitN(name, "?", 2)
	return nameParts[0]
}

//Session get session from mongodb
func Session() (*mgo.Session, error) {
	var err error

	if ses == nil {
		mgoURI := getDbURI()
		if ses, err = mgo.Dial(mgoURI); err != nil {
			log.Fatal("can't connect to mongo:", err)
			return nil, err
		}
		// sets mongodb mode
		ses.SetMode(mgo.Monotonic, true)

		// sets mongodb safety mode
		ses.SetSafe(&mgo.Safe{})
	}
	return ses.Copy(), nil
}

//DB get specific collection
func DB(ses *mgo.Session) *mgo.Database {
	return ses.DB(getDbName())
}

//LeagueMappers get sites collection
func LeagueMappers(ses *mgo.Session) *mgo.Collection {
	return DB(ses).C("league_mappers")
}
