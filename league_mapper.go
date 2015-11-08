package app

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

const (
	//SourceLiveScore represent livescore source
	SourceLiveScore = "livescore"
)

//LeagueMapper represent league mapper
type LeagueMapper struct {
	ID        string `json:"id" bson:"_id"`
	Key       string `json:"key" bson:"key"`
	URL       string `json:"url" bson:"url"`
	SourceKey string `json:"source_key" bson:"source_key"`
}

//NewLeagueMapper handle create new league mapper
func NewLeagueMapper(session *mgo.Session, key string, url string, sourceKey string) (*LeagueMapper, error) {
	l := new(LeagueMapper)
	l.ID = bson.NewObjectId().Hex()
	l.Key = key
	l.URL = url
	l.SourceKey = sourceKey

	if err := LeagueMappers(session).Insert(l); err != nil {
		return nil, err
	}

	return l, nil
}

//LeagueList get league list
func LeagueList(session *mgo.Session) ([]LeagueMapper, error) {
	ll := []LeagueMapper{}
	if err := LeagueMappers(session).Find(bson.M{}).Sort("key").All(&ll); err != nil {
		return nil, err
	}

	return ll, nil
}

//RemoveLeagueMappersByID remove leagueMappersByID
func RemoveLeagueMappersByID(session *mgo.Session, ID string) error {
	if err := LeagueMappers(session).RemoveId(ID); err != nil {
		return err
	}

	return nil
}

//GetLeagueMapperByKey get league mapper by specific key
func GetLeagueMapperByKey(session *mgo.Session, key string) (LeagueMapper, error) {
	l := LeagueMapper{}
	if err := LeagueMappers(session).Find(bson.M{"key": key}).One(&l); err != nil {
		return l, err
	}

	return l, nil
}
