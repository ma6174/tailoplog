package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func main() {
	mongoHost := "127.0.0.1:27017"
	if len(os.Args) > 1 {
		mongoHost = os.Args[1]
	}
	session, err := mgo.Dial(mongoHost)
	if err != nil {
		log.Fatal(err)
	}
	var last bson.M
	err = session.DB("local").C("oplog.rs").Find(nil).Sort("-$natural").One(&last)
	if err != nil {
		panic(err)
	}
	iter := session.DB("local").C("oplog.rs").Find(bson.M{"ts": map[string]interface{}{"$gt": last["ts"]}}).LogReplay().Tail(time.Second * 5)
	for {
		for {
			var result bson.M
			ok := iter.Next(&result)
			if !ok {
				break
			}
			if result["ns"].(string) == "" {
				continue
			}
			b, err := json.Marshal(result)
			if err != nil {
				panic(err)
			}
			println(string(b))
		}
		if iter.Err() != nil {
			panic(iter.Close())
		}
		if iter.Timeout() {
			continue
		}
		panic("unexpect error")
	}
}
