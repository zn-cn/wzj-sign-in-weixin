package db

import (
	"config"
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
)

// mgo controller
type MgoDBCntlr struct {
	sess *mgo.Session
	db   *mgo.Database
}

var (
	DBNAME     = config.Conf.DB.DBName
	globalSess *mgo.Session
	mongoURL   string
)

const (
	MongoCopyType  = "1"
	MongoCloneType = "2"
)

func init() {
	dbConf := config.Conf.DB
	if dbConf.User != "" && dbConf.PW != "" {
		mongoURL = fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", dbConf.User, dbConf.PW, dbConf.Host, dbConf.Port, dbConf.AdminDBName)
	} else {
		mongoURL = fmt.Sprintf("mongodb://%s:%s", dbConf.Host, dbConf.Port)
	}

	var err error
	globalSess, err = GetDBSession()
	if err != nil {
		panic(err)
	}
}

/****************************************** db session manage ****************************************/

// GetSession get the db session
func GetDBSession() (*mgo.Session, error) {
	globalMgoSession, err := mgo.DialWithTimeout(mongoURL, 10*time.Second)
	if err != nil {
		return nil, err
	}
	globalMgoSession.SetMode(mgo.Monotonic, true)
	//default is 4096
	globalMgoSession.SetPoolLimit(1000)
	return globalMgoSession, nil
}

func NewCloneMgoDBCntlr() *MgoDBCntlr {
	sess := globalSess.Clone()
	return &MgoDBCntlr{
		sess: sess,
		db:   sess.DB(DBNAME),
	}
}

func NewCopyMgoDBCntlr() *MgoDBCntlr {
	sess := globalSess.Copy()
	return &MgoDBCntlr{
		sess: sess,
		db:   sess.DB(DBNAME),
	}
}

func (this *MgoDBCntlr) Close() {
	this.sess.Close()
}

func (this *MgoDBCntlr) GetDB() *mgo.Database {
	return this.db
}

func (this *MgoDBCntlr) GetTable(tableName string) *mgo.Collection {
	return this.db.C(tableName)
}
