package mongo

import (
    "github.com/revel/revel"
    "gopkg.in/mgo.v2"
    "sync"
    "os"
)

type Mongo struct {
    *revel.Controller
    MongoSession  *mgo.Session
    MongoDatabase *mgo.Database
}

var session *mgo.Session

// Singleton
var dial sync.Once

func GetSession() *mgo.Session {
    dial.Do(func() {
        var err error
        session, err = mgo.Dial("localhost,mongodb")
        if err != nil {
            panic(err)
        }
    })

    return session
}

func (c *Mongo) Bind() revel.Result {
    c.MongoSession = GetSession().Clone()
    c.MongoDatabase = c.MongoSession.DB(os.Getenv("MONGODB_DBNAME"))

    return nil
}

func (c *Mongo) Close() revel.Result {

    if c.MongoSession != nil {
        c.MongoSession.Close()
    }

    return nil
}

func init() {
    revel.InterceptMethod((*Mongo).Bind, revel.BEFORE)
    revel.InterceptMethod((*Mongo).Close, revel.AFTER)
    revel.InterceptMethod((*Mongo).Close, revel.PANIC)
}