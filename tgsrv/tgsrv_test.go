package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/recoilme/tf/httputils"
	"github.com/dariaserkova/tf/params"
	"github.com/recoilme/tf/vkapi"

	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
)

func TestAverage(t *testing.T) {
	groupName := "myakotkapub"
	group := pubDbGet(groupName)
	//if group == nil {
	log.Println("group1", group)
}

func TestVkWallUpd(t *testing.T) {
	log.Println("vkWallUpd")
	//vkapi.vkWallUpd()
}

func TestLinkExtract(t *testing.T) {
	var defHeaders = make(map[string]string)
	defHeaders["User-Agent"] = "script::recoilme:v1"
	defHeaders["Authorization"] = "Client-ID 4191ffe3736cfcb"

	b := httputils.HttpGet("https://www.reddit.com/.rss?feed=32f7ac01a37b80c88037018e186bb2581de14d55&user=recoilme", defHeaders)
	if b == nil {
		return
	}
	s := string(b)
	//log.Println(s)
	var rss string
	//s := `<link rel="alternate" type="application/rss+xml" href="https://vc.ru/feed">`
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		log.Fatal(err)
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "link" {
			var isRss bool
			for _, a := range n.Attr {
				if a.Key == "type" {
					if a.Val == "application/rss+xml" || a.Val == "application/atom+xml" {
						isRss = true
						break
					}
				}
			}
			if isRss {
				for _, a := range n.Attr {
					if a.Key == "href" {
						rss = a.Val
						break
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
			if rss != "" {
				break
			}
		}
	}
	f(doc)
	log.Println("rss", rss)
}

func TestMainDomain(t *testing.T) {
	var result string
	result, _ = publicsuffix.EffectiveTLDPlusOne("www.reddit.com")
	if result != "reddit.com" {
		t.Error("Expected reddit.com, got ", result)
	}
	result, _ = publicsuffix.EffectiveTLDPlusOne("m.vk.com")
	if result != "vk.com" {
		t.Error("Expected vk.com, got ", result)
	}
	result, _ = publicsuffix.EffectiveTLDPlusOne("adsfffsf")
	log.Println("result:", result)
	if result != "" {
		t.Error("Expected reddit.com, got ", result)
	}
	result, _ = publicsuffix.EffectiveTLDPlusOne("en.reddit.com")
	if result != "reddit.com" {
		t.Error("Expected reddit.com, got ", result)
	}
}

func TestRss(t *testing.T) {
	url := "https://www.reddit.com/.rss?feed=32f7ac01a37b80c88037018e186bb2581de14d55&user=recoilme"
	//url := "https://vc.ru/feed"
	link := getFeedLink(url)
	if link != url {
		t.Error("Expected got", link)
	}
}

func TestRssExtract(t *testing.T) {
	url := "https://vc.ru/"
	urlexpect := "https://vc.ru/feed"
	link := rssExtract(url)
	if link != urlexpect {
		t.Error("Expected got", link)
	}
}

func TestSubs(t *testing.T) {
	subs := usersub("", 1263310, true)
	for k, v := range subs {
		log.Println("k", k, "v", v)
	}
}

func TestMigrate(t *testing.T) {
	migrate := false
	if migrate == false {
		return
	}
	//get publics

	domains := vkdomains()
	for i := range domains {
		domain := domains[i]
		log.Println(domain.ScreenName)
		users := domUsers(domains[i])
		for user := range users {
			log.Println(user, domain.ScreenName)
			usersub(params.Publics+domain.ScreenName, user, false)
		}
		time.Sleep(1 * time.Second)
		//break
	}

	//get feeds
	rsss := rssdomains()
	for hash, feedlink := range rsss {

		log.Println(hash, "::", feedlink)
		users := domUsersRss(hash)
		for user := range users {
			//log.Println(user, ":", feedlink)
			usersub(params.Feeds+GetMD5Hash(feedlink), user, false)
		}
		//time.Sleep(1 * time.Second)
	}
}

func domUsersRss(hash string) (users map[int64]bool) {
	mask := params.FeedSubs + "%s"
	url := fmt.Sprintf(mask, hash)
	log.Println(url)
	b := httputils.HttpGet(url, nil)
	if b != nil {
		json.Unmarshal(b, &users)
	}
	return users
}

func rssdomains() map[string]string {
	domains := make(map[string]string)
	url := params.BaseUri + "feeds/Last?cnt=1000000&order=desc&vals=false"
	log.Println("rssdomains", url)
	resp, err := http.Post(url, "application/json", nil)
	if err == nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			var keys = make([]string, 0)
			json.Unmarshal(body, &keys)

			log.Println("keys", keys)
			for _, key := range keys {
				log.Println("key", key)
				b := httputils.HttpGet(params.Feeds+key, nil)
				if b != nil {
					domains[key] = string(b)
				}
			}
		}
	}
	return domains
}

func domUsers(domain vkapi.Group) (users map[int64]bool) {
	mask := params.Subs + "%d"
	url := fmt.Sprintf(mask, domain.Gid)
	log.Println(url)
	b := httputils.HttpGet(url, nil)
	if b != nil {
		json.Unmarshal(b, &users)
	}
	return users
}

func vkdomains() (domains []vkapi.Group) {
	var domainNames []string
	url := params.BaseUri + "pubNames/Last?cnt=1000000&order=desc&vals=false"
	log.Println("vkdomains", url)
	resp, err := http.Post(url, "application/json", nil)
	if err == nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		err := json.Unmarshal(body, &domainNames)
		if err == nil {
			for i := range domainNames {
				domainName := domainNames[i]
				//log.Println("domainName", domainName)
				b := httputils.HttpGet(params.Publics+domainName, nil)
				if b != nil {
					var domain vkapi.Group
					err := json.Unmarshal(b, &domain)
					if err == nil {
						domains = append(domains, domain)
					}
				}
			}
		} else {
			log.Println(err)
		}
	}
	return
}

func TestInst(t *testing.T) {
	word := "https://www.instagram.com/MillaJovovich/"
	urls, _ := url.Parse(word)
	parts := strings.Split(urls.Path, "/")
	log.Println(parts)
	for _, part := range parts {
		if part != "" {
			log.Println("oart", part)
		}
	}
	//if len(parts) == 2 {
	//findFeed("https://web.stagram.com/rss/n/"+parts[1], msg, delete)
	//}
}

func TestPublicFmt(t *testing.T) {
	domain := "public48287065"
	res := vkapi.GroupsGetById(domain)
	if len(res) == 0 {
		t.Error("public48287065, got ", res)
	}
}

func getUsersSub() (users []string) {
	log.Println("getUsersSub")
	//var users []string
	url := params.UserSubs + "Last?cnt=1000000&order=desc&vals=false"
	log.Println("users", url)
	resp, err := http.Post(url, "application/json", nil)
	if err == nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &users)

	}
	return
}
func TestConvertSub(t *testing.T) {
	users := getUsersSub()
	for _, user := range users {
		log.Println(user)
		id, _ := strconv.ParseInt(user, 10, 64) //strconv.Atoi(user)
		subs := usersub("", id, true)
		newsubs := make(map[string]bool)
		for sub := range subs {
			var newsub string

			newsub = strings.Replace(sub, "http://localhost:5000/bolt", "", -1)
			newsub = strings.Replace(newsub, "http://badtobefat.ru/bolt", "", -1)
			newsub = strings.Replace(newsub, "http://badtobefat.ru:5000/bolt", "", -1)
			//log.Println(newsub)
			newsubs[newsub] = true
		}
		newsubs["/pubNames/telefeedbot"] = true

		bsubs, _ := json.Marshal(newsubs)
		suburl := params.UserSubs + user
		res := httputils.HttpPut(suburl, nil, bsubs)
		if !res {
			log.Println("err")
			log.Println("err", user)
			break
		}
	}
}

func TestTelefeedSub(t *testing.T) {
	users := getUsersSub()
	usersMap := make(map[int64]bool)
	for _, user := range users {
		fmt.Println(user)
		id, _ := strconv.ParseInt(user, 10, 64) //strconv.FormatInt(user, 10)
		if id > 0 {
			usersMap[id] = true
		}
		subs := usersub("", id, true)

		newsubs := make(map[string]bool)
		for sub := range subs {
			newsubs[sub] = true
		}
		newsubs["/pubNames/telefeedbot"] = true

		bsubs, _ := json.Marshal(newsubs)
		suburl := params.UserSubs + user
		res := httputils.HttpPut(suburl, nil, bsubs)
		if !res {
			log.Println("err")
			log.Println("err", user)
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	data, err := json.Marshal(usersMap)
	if err == nil {
		log.Println("pubSubTgAdd data ", string(data))
		result := httputils.HttpPut(params.Subs+"146445941", nil, data)
		fmt.Println("res", result)
	}
}
