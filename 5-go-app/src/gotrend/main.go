package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"

	"github.com/go-pg/pg/v9"
)

// Define HTML template and vars
var tpl = template.Must(template.ParseFiles("index.html"))
var googleOauthConfig *oauth2.Config
var oauthStateString = "pseudo-random"
var likedID = ""

const (
	host     = "localhost"
	port     = 5432
	user     = "jesse.lhuillier"
	password = ""
	dbname   = "gotrend"
)

const googleID = "912232579439-r4e91kgmt2v0tsbcifm3nahb2j4ofnjc.apps.googleusercontent.com"
const googlePass = "_WIUihTgElRldWPslWO4IVz0"

//User table
type User struct {
	ID        string `pg:",pk"`
	Firstname string
	Lastname  string
	Email     string
}

func (u *User) String() string {
	return fmt.Sprintf("User<%s %s %s %s>", u.ID, u.Lastname, u.Email, u.Firstname)
}

//Video table
type Video struct {
	ID           string `pg:",pk"`
	Title        string
	Channelid    string
	Channeltitle string
	Description  string
	Publishedat  string
	Thumbnailurl string
}

func (v *Video) String() string {
	return fmt.Sprintf("Video<%s %s %s %s %s %s %s>", v.ID, v.Title, v.Channelid, v.Channeltitle, v.Description, v.Publishedat, v.Thumbnailurl)
}

//Like table
type Like struct {
	IDvideo string `pg:",pk"`
	IDuser  string
}

func (l *Like) String() string {
	return fmt.Sprintf("Likes<%s %s>", l.IDuser, l.IDvideo)
}

// Userinfo decode table
type Userinfo struct {
	ID            string
	email         string
	verifiedEmail string
	picture       string
	hd            string
}

func main() {
	// handle file server foer static files
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// handle func for dynamic pages
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleGoogleLogin)
	http.HandleFunc("/callback", handleGoogleCallback)
	http.HandleFunc("/items", itemsHandler)
	http.HandleFunc("/like", likevideo)
	http.HandleFunc("/unlike", unlikevideo)
	fmt.Println(http.ListenAndServe(":8080", nil))

}

// init OAuth2 config
func init() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL: "http://localhost:8080/callback",
		//RedirectURL:  "http://gotrend.appspot.com/callback",
		ClientID:     googleID,
		ClientSecret: googlePass,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

// Handle HTTP response for youtube API
func itemsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("LoginCookie")

	if err != nil {
		fmt.Println("error in reading cookie : " + err.Error() + "\n")
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	} else {
		fmt.Println("cookie has : " + c.Value + "\n")
	}
	const developerKey = "AIzaSyBS0BET1LFxBbhIUF554i1jLj3LmHgBav8"
	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	// create client for API call
	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	// Make the API call to YouTube.
	call := service.Videos.List("id,snippet,contentDetails").
		Chart("mostPopular").
		RegionCode("FR").
		MaxResults(20)
	response, err := call.Do()

	// Match the API response with HTML template
	err = tpl.Execute(w, response)
	if err != nil {
		fmt.Println(err)
	}
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("LoginCookie")

	if err != nil {
		fmt.Println("error in reading cookie : " + err.Error() + "\n")
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	} else {
		fmt.Println("cookie has : " + c.Value + "\n")
		http.Redirect(w, r, "/items", http.StatusTemporaryRedirect)
		insertuserid(c.Value)
	}
}

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	content, userID, err := getUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// set up cookies
	c := http.Cookie{
		Name:   "LoginCookie",
		Value:  userID,
		MaxAge: 3600}

	http.SetCookie(w, &c)
	fmt.Println(content)

	http.Redirect(w, r, "/items", http.StatusTemporaryRedirect)
}

func getUserInfo(state string, code string) ([]byte, string, error) {
	null := ""

	if state != oauthStateString {
		return nil, null, fmt.Errorf("invalid oauth state")
	}

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, null, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, null, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, null, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	var userinfos Userinfo
	json.Unmarshal([]byte(contents), &userinfos)
	userID := userinfos.ID
	return contents, userID, nil
}

func likevideo(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("LoginCookie")

	if err != nil {
		fmt.Println("error in reading cookie : " + err.Error() + "\n")
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	}
	userID := c.Value
	fmt.Println("liked video")
	body, err := ioutil.ReadAll(r.Body)
	likedID = string(body)
	insertlike(userID, likedID)
	insertvideo(likedID)
	fmt.Println(err)
}

func unlikevideo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("unliked video")
	body, err := ioutil.ReadAll(r.Body)
	likedID = string(body)
	deletevideo(likedID)
	fmt.Println(err)
}

func insertvideo(videoID string) {
	const developerKey = "AIzaSyBS0BET1LFxBbhIUF554i1jLj3LmHgBav8"
	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	// create client for API call
	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	// Make the API call to YouTube.
	call := service.Videos.List("id,snippet").
		Id(videoID).
		Fields("items(id,snippet(channelId,title,channelTitle,description,publishedAt,thumbnails(default)))")
	response, err := call.Do()
	if err != nil {
		fmt.Println(err)
	}

	db := pg.Connect(&pg.Options{
		User:     user,
		Password: password,
		Database: dbname,
		Addr:     "localhost:5432",
	})

	resp := response.Items[0]
	_, err = db.Model(video(resp)).
		OnConflict("(id) DO UPDATE").
		Insert()
	fmt.Println(err)
}

func deletevideo(videoID string) {
	db := pg.Connect(&pg.Options{
		User:     user,
		Password: password,
		Database: dbname,
		Addr:     "localhost:5432",
	})
	like := &Like{IDvideo: likedID}
	err := db.Delete(like)
	fmt.Println(err)
}

func insertuserid(userID string) {
	db := pg.Connect(&pg.Options{
		User:     user,
		Password: password,
		Database: dbname,
		Addr:     "localhost:5432",
	})
	user := new(User)
	user.ID = userID
	_, err := db.Model(user).
		Where("user.ID = ?user.ID").
		OnConflict("(id) DO UPDATE").
		Insert()
	fmt.Println(err)
}

func insertlike(userID string, videoID string) {
	db := pg.Connect(&pg.Options{
		User:     user,
		Password: password,
		Database: dbname,
		Addr:     "localhost:5432",
	})
	like := new(Like)
	like.IDuser = userID
	like.IDvideo = videoID
	_, err := db.Model(like).
		Insert()
	fmt.Println(err)
}

// func dbselect() {
// 	db := pg.Connect(&pg.Options{
// 		User:     user,
// 		Password: password,
// 		Database: dbname,
// 		Addr:     "localhost:5432",
// 	})
// 	db.AddQueryHook(dbLogger{})
// 	video := &Video{ID: "REweKmw8lmg"}
// 	r := db.Select(video)
// 	fmt.Println(r)
// 	fmt.Println(video)
// }

func fetchlikes() {
	video := new(Video)
	err := db.Model(video).
		ColumnExpr("video.*").
		ColumnExpr(" v.ID AS Video__id, v.Title AS Video__title, v.Channelid AS Video__channelid, v.Channeltitle AS Video__channeltitle, v.Description AS Video__description, v.Publishedat AS Video__publishedat, v.Thumbnailurl AS Video__thumbnailurl").
		Join("JOIN likes AS l ON l.id =_id")
}

func video(resp *youtube.Video) *Video {
	video := new(Video)
	video.ID = resp.Id
	video.Title = resp.Snippet.Title
	video.Description = resp.Snippet.Description
	video.Channelid = resp.Snippet.ChannelId
	video.Channeltitle = resp.Snippet.ChannelTitle
	video.Publishedat = resp.Snippet.PublishedAt
	video.Thumbnailurl = resp.Snippet.Thumbnails.Default.Url
	return video
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	fmt.Println(q.FormattedQuery())
	return nil
}
