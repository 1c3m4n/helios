package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"os"
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

type UserUnreadFeedItem struct {
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

func GetFeeds() (Feeds []Feed) {

	ExecuteWithCollection("rss", "feeds", func(c *mgo.Collection) error {
		return c.Find(bson.M{}).All(&Feeds)
	})
	return
}

func HashEncode(input string) (output string) {
	h := sha256.New()
	h.Write([]byte(input))
	output = base64.URLEncoding.EncodeToString(h.Sum(nil))
	return
}

func AddUserFeed(newfeed NewFeed) {
	resp, err := http.Head(newfeed.Url)
	fmt.Printf("Response from %s : %s \n", newfeed.Url, resp.Status)
	if resp.Status == "200 OK" {

		//feedhash := HashEncode(newfeed.Url)
		//feed := Feed{feedhash, "", newfeed.Url, time.Now()}
		//userfeed := UserFeed{newfeed.EmailHash, newfeed.FeedName, feedhash, ""}
		fmt.Printf("Adding URL(%s) for EmailHash: %s \n", newfeed.Url, newfeed.EmailHash)
		fmt.Println(resp.Header)
	} else {
		fmt.Printf("<ERROR> adding URL(%s) for EmailHash: %s , Error: %s \n", newfeed.Url, newfeed.EmailHash, err)
	}
}

func main() {
	os.Clearenv()
	fmt.Printf("Sarting NewsAgent")
	feeds := GetFeeds()
	for key := range feeds {
		fmt.Printf("Processing %s - %s (LAST:%s) URI:%s \n", feeds[key].FeedHash, feeds[key].FeedName, feeds[key].LastUpdated, feeds[key].RssUrl)
		resp, err := http.Get(feeds[key].RssUrl)
		if resp.Status == "200 OK" {
			//fmt.Println()
			//fmt.Println(resp.Header)
			fmt.Println()
			defer resp.Body.Close()
			contents, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(string(contents))
			fmt.Println()
		} else {
			fmt.Printf("<ERROR> Could not access %s, Error: %s \n", feeds[key].RssUrl, err)
		}
	}
	fmt.Printf("Stopping NewsAgent")
}
