package controllers

import (
    "github.com/revel/revel"
    "fmt"
    "time"
    "gopkg.in/mgo.v2/bson"
    "net/url"
    "github.com/ChimeraCoder/anaconda"
    "strconv"
    "io/ioutil"
    "encoding/json"
    "os"
    "github.com/aladagemre/tweetnotes/app/modules/mongo"
)

type App struct {
    *revel.Controller
    mongo.Mongo
}

type Tweet struct {
    Id bson.ObjectId "_id,omitempty"
    User_id string
    Screen_name string
    Text string "text,omitempty"
    Created_at time.Time
    Note string
    Id_str string "id_str"
}

func (t Tweet) String() string {
    return fmt.Sprintf("[%s]\t[%s]\t%s: %s", t.Created_at, t.Id_str, t.Screen_name, t.Text)
}

func (c App) Index() revel.Result {
    return c.Render()
}


// Displays tweets of a specific user.
func (c App) Tweets(screen_name string) revel.Result {
    collection := c.MongoDatabase.C("tweetnotes_tweetmongo")

    var tweets []Tweet
    err := collection.Find(bson.M{"screen_name": screen_name}).Sort("-created_at").All(&tweets)
    if err != nil {
        panic(err)
    }
    //fmt.Println("Results All: ", tweets)
    return c.Render(screen_name, tweets)
}

// Displays tweets of a specific user in json format.
// json file can be saved and then imported.
func (c App) TweetsJson(screen_name string) revel.Result {
    collection := c.MongoDatabase.C("tweetnotes_tweetmongo")
    var tweets []Tweet
    err := collection.Find(bson.M{"screen_name": screen_name}).Sort("-created_at").All(&tweets)
    if err != nil {
        panic(err)
    }
    return c.RenderJson(tweets)
}


// Handler for json import operation.
// Displays the number of imported tweets.
func (c App) Import() revel.Result {
    myfile, _, err := c.Request.FormFile("file")
    if err != nil {
        panic(err)
    }
    defer myfile.Close()

    byt, err := ioutil.ReadAll(myfile)
    var tweets []Tweet

    if err := json.Unmarshal(byt, &tweets); err != nil {
        panic(err)
    }

    screen_name := tweets[0].Screen_name
    fmt.Sprintln("Read %d tweets.", len(tweets))
    fmt.Println("Importing...")
    // TODO: Handle upsert.
    successful, failed := AddTweets(c, screen_name, tweets)
    fmt.Println("Imported tweets count:", successful)
    return c.Render(successful, failed, screen_name)
}

// Updates the tweet list of the specified user.
func (c App) Update(screen_name string) revel.Result {
    collection := c.MongoDatabase.C("tweetnotes_tweetmongo")
    var last_tweet Tweet
    err := collection.Find(bson.M{"screen_name": screen_name}).Sort("-created_at").Limit(1).One(&last_tweet)
    if err != nil {
        // panic(err)
        // No tweets yet.
        tweets := GetAllTweetsSince(screen_name, "")
        fmt.Println(tweets)
        AddTweets(c, screen_name, tweets)
    } else {
        // we already have some tweets.
        fmt.Println(last_tweet.Id_str)
        tweets := GetAllTweetsSince(screen_name, last_tweet.Id_str)
        fmt.Println(tweets)
        AddTweets(c, screen_name, tweets)
    }

    return c.Redirect("/tweets/%s", screen_name)
}

// Adds given tweets to the database.
func AddTweets(app App, screen_name string, tweets []Tweet) (int, int) {
    //c := GetCollection()
    c := app.MongoDatabase.C("tweetnotes_tweetmongo")

    fmt.Println("Adding new tweets...")
    successful := 0
    for i := 0; i<len(tweets); i++ {
        err := c.Insert(tweets[i])
        if err != nil {
            fmt.Println(i)
            fmt.Println(err)
        } else {
            successful++
        }
    }
    failed := len(tweets) - successful
    fmt.Println("Adding finished.")
    return successful, failed
}

// Gets an API object created with auth keys in env.
func GetAPI() *anaconda.TwitterApi {
    CONSUMER_KEY := os.Getenv("CONSUMER_KEY")
    CONSUMER_SECRET := os.Getenv("CONSUMER_SECRET")
    ACCESS_TOKEN := os.Getenv("ACCESS_TOKEN")
    ACCESS_TOKEN_SECRET := os.Getenv("ACCESS_TOKEN_SECRET")

    anaconda.SetConsumerKey(CONSUMER_KEY)
    anaconda.SetConsumerSecret(CONSUMER_SECRET)
    api := anaconda.NewTwitterApi(ACCESS_TOKEN, ACCESS_TOKEN_SECRET)
    //fmt.Println(*api.Credentials)
    return api
}

// Fetches tweets since a tweet and returns the list.
func GetAllTweetsSince(screen_name string, id_str string) []Tweet {
    api := GetAPI()

    // Prepare variables
    v := url.Values{}
    v.Set("screen_name", screen_name)
    v.Set("count", "200")
    if id_str != "" {
        v.Set("since_id", id_str)
    }

    // Make the call
    atweets, err := api.GetUserTimeline(v)
    if err != nil {
        fmt.Println(err);
        fmt.Println(atweets);
        panic(err);
    }

    // Parse the tweets
    all_tweets := make([]Tweet, 200)
    new_tweets := ParseTweets(atweets)
    //fmt.Println(new_tweets)
    if len(new_tweets) == 0 {
        return new_tweets
    }

    // Add new tweets to all tweets list.
    for i := 0; i < len(new_tweets); i++ {
        all_tweets = append(all_tweets, new_tweets[i])
    }

    // Find the oldest.
    last_int, err := strconv.Atoi(all_tweets[len(all_tweets)-1].Id_str)
    oldest := strconv.Itoa(last_int - 1)


    // Fetch the remaining tweets until there is no more tweet.
    for len(new_tweets) > 0 {
        // Set the max_id.
        fmt.Println("Getting tweets before ", oldest)
        v.Set("max_id", oldest)

        // Fetch the remaining tweets
        atweets, err = api.GetUserTimeline(v)
        if err != nil {
            //fmt.Println(err.Error())
            //panic(err);
            fmt.Println(err)
        }

        // Parse the tweets
        new_tweets := ParseTweets(atweets)
        fmt.Println("New tweets:", len(new_tweets))
        if len(new_tweets) == 0 {
            break
        }


        // Add new tweets to all tweets list.
        for i := 0; i < len(new_tweets); i++ {
            all_tweets = append(all_tweets, new_tweets[i])
        }

        // Find new oldest id.
        last_int, err = strconv.Atoi(all_tweets[len(all_tweets)-1].Id_str)
        oldest = strconv.Itoa(last_int - 1)
        fmt.Println("Tweets downloaded so far:", len(all_tweets))
    }
    return all_tweets
}

// Parses the tweet response and constructs an object.
func ParseTweet(atweet anaconda.Tweet) Tweet {
    parsedTime, err := time.Parse(time.RubyDate, atweet.CreatedAt)

    if (err != nil ) {
        panic(err)
    }
    tweet := Tweet{}
    tweet.Id_str = atweet.IdStr
    tweet.Created_at = parsedTime
    tweet.Text = atweet.Text
    tweet.Screen_name = atweet.User.ScreenName
    tweet.User_id = atweet.User.IdStr
    return tweet
}

// Foor loop for calling parse tweet.
func ParseTweets(atweets []anaconda.Tweet) []Tweet {
    tweets := make([]Tweet, len(atweets))
    for i := 0; i < len(atweets); i++ {
        tweets[i] = ParseTweet(atweets[i])
    }
    return tweets
}

// Structure for the note text.
type NoteText struct {
    Text string
}

// Handler for posting a note for a tweet.
func (c App) PostNote(screen_name string, id_str string, note NoteText) revel.Result {
    /* Posts note on a specific tweet. */
    fmt.Println(note.Text)
    collection := c.MongoDatabase.C("tweetnotes_tweetmongo")
    colQuerier := bson.M{"screen_name": screen_name, "id_str": id_str}
    change := bson.M{"$set": bson.M{"note": note.Text}}
    err := collection.Update(colQuerier, change)

    payload := make(map[string]interface{})
    if err != nil {
        payload["error"] = err.Error()
        payload["success"] = false
    } else {
        payload["success"] = true
    }

    return c.RenderJson(payload)
}