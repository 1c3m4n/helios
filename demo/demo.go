package main

import (
	"code.google.com/p/gorest"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"time"
)

type NewFeed struct {
	EmailHash string
	Url       string
	FeedName  string
}

type User struct {
	UserName         string
	Email            string
	EmailHash        string
	PasswordHash     string
	RegistrationDate time.Time
}

type UserFeed struct {
	EmailHash        string
	FeedName         string
	FeedHash         string
	LastFeedItemHash string
}

type UserUnreadFeed struct {
	EmailHash string
	FeedHash  string
	ItemUrl   string
	TimeStamp time.Time
}

type Feed struct {
	FeedHash    string
	FeedName    string
	RssUrl      string
	LastUpdated time.Time
}

type FeedItem struct {
	FeedHash     string
	FeedItemHash string
	ItemUrl      string
	TimeStamp    time.Time
}

func ExecuteWithCollection(database, collection string, f func(*mgo.Collection) error) error {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(database).C(collection)
	return f(c)
}

func GetUser(emailhash string) (result *User) {
	result = &User{}
	ExecuteWithCollection("rss", "users", func(c *mgo.Collection) error {
		return c.Find(bson.M{"emailhash": emailhash}).One(result)
	})
	return
}

func CreateUser(user User) {
	ExecuteWithCollection("rss", "users", func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"emailhash": user.EmailHash}, user)
		return err
	})
}

func (serv ReaderService) GetUserFeeds(emailhash string) (userFeeds []UserFeed) {

	ExecuteWithCollection("rss", "userfeeds", func(c *mgo.Collection) error {
		return c.Find(bson.M{"emailhash": emailhash}).All(&userFeeds)
	})
	serv.ResponseBuilder().AddHeader("Access-Control-Allow-Origin", "http://localhost:8080")
	fmt.Printf("Found %d UserFeeds for : %s \n", len(userFeeds), emailhash)
	return
}

func CreateUserFeed(userfeed UserFeed) {
	ExecuteWithCollection("rss", "userfeeds", func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"emailhash": userfeed.EmailHash, "feedhash": userfeed.FeedHash}, userfeed)
		return err
	})
	return
}

func (serv ReaderService) GetUserUnreadFeeds(emailhash string) (userunreadfeeds []UserUnreadFeed) {
	ExecuteWithCollection("rss", "userunreadfeeds", func(c *mgo.Collection) error {
		return c.Find(bson.M{"emailhash": emailhash}).All(&userunreadfeeds)
	})
	return
}

func CreateUserUnreadFeed(userunreadfeed UserUnreadFeed) {
	ExecuteWithCollection("rss", "userunreadfeeds", func(c *mgo.Collection) error {
		return c.Insert(userunreadfeed)
	})
	return
}

func CreateFeed(feed Feed) {
	ExecuteWithCollection("rss", "feeds", func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"feedhash": feed.FeedHash}, feed)
		return err
	})
}

func CreateFeedItem(feeditem FeedItem) {
	ExecuteWithCollection("rss", "feeditems", func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"feeditemhash": feeditem.FeedItemHash}, feeditem)
		return err
	})
}

func HashEncode(input string) (output string) {
	h := sha256.New()
	h.Write([]byte(input))
	output = base64.URLEncoding.EncodeToString(h.Sum(nil))
	return
}

func (serv ReaderService) AddUserFeed(newfeed NewFeed) {
	serv.ResponseBuilder().AddHeader("Access-Control-Allow-Origin", "http://localhost:8080")
	resp, err := http.Head(newfeed.Url)
	fmt.Printf("Response from %s : %s \n", newfeed.Url, resp.Status)
	if resp.Status == "200 OK" {
		feedhash := HashEncode(newfeed.Url)
		feed := Feed{feedhash, "", newfeed.Url, time.Now()}
		CreateFeed(feed)
		userfeed := UserFeed{newfeed.EmailHash, newfeed.FeedName, feedhash, ""}
		CreateUserFeed(userfeed)
		fmt.Printf("Adding URL(%s) for EmailHash: %s \n", newfeed.Url, newfeed.EmailHash)
	} else {
		fmt.Printf("<ERROR> adding URL(%s) for EmailHash: %s , Error: %s \n", newfeed.Url, newfeed.EmailHash, err)
	}
}

func (serv ReaderService) DoOptions(varArgs ...string) {
	rb := serv.ResponseBuilder()
	rb.AddHeader("Access-Control-Allow-Origin", "http://localhost:8080")
	rb.AddHeader("Access-Control-Allow-Method", "GET, PUT, POST, DELETE")
	rb.AddHeader("Access-Control-Allow-Headers", "accept, origin, x-requested-with, content-type")

	fmt.Printf("Received an Options request, responding with some options. \n")
}

type ReaderService struct {
	gorest.RestService `root:"/UserFeed" consumes:"application/json" produces:"application/json"`

	doOptions    gorest.EndPoint `method:"OPTIONS"	path:"/{...:string}"`
	addUserFeed  gorest.EndPoint `method:"POST" 	    path:"/"							postdata:"NewFeed"`
	getUserFeeds gorest.EndPoint `method:"GET" 		path:"/{emailhash:string}" 			output:"[]UserFeed"`
	//getUserUnreadFeeds gorest.EndPoint `method:"GET" 		path:"/UserUnreadFeed/{emailhash:string}" 	output:"[]UserUnreadFeed"`
}

func main() {
	gorest.RegisterService(new(ReaderService))
	http.Handle("/", gorest.Handle())
	http.ListenAndServe(":9090", nil)
}
