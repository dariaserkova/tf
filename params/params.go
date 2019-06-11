package params

import (
	"io/ioutil"
	"log"
	"os"
)

const (
	Vkwriterfile = "vkwriter.bot"
	users        = "/usertg/"
	channels     = "/channels/"
	PubNames     = "/pubNames/"
	username     = "/usernametg/"
	pubSubTg     = "/pubSubTg/"
	feedSubTg    = "/feedSubTg/"
	userSubTg    = "/userSubTg/"
	lastPost     = "/vkpublastpost/"
	Feed         = "/feeds/"
	best         = "/best/"
	viralhash    = "/viralhash/"
	links        = "/links/"
	video        = "/video/"
	ShortUrl     = "https://www.googleapis.com/urlshortener/v1/url?key="
	TgApi        = "https://api.telegram.org"
	Example      = "\nExample: \nhttps://www.reddit.com/r/gifs/top/\nhttps://vk.com/evil_incorparate\n\nMore examples: http://telegra.ph/telefeedbot-05-12\n\nЛучшее из вконтакте - собрано здесь: @memefeed\n\nТоп пабликов и фидов по подписчикам: /top"
	SomeErr      = "🇬🇧 Something going wrong. Try later.. 🇷🇺 Ошибка, мать её!"
	Hello        = "🇬🇧 Send me a link on domain/rss.\n\n🇷🇺 Отправь мне ссылку на домен или rss.\n\n" + Example
	Psst         = "🇬🇧 As soon as there are new articles here - i will send them, but with some delay (1-3 hour)\n\n🇷🇺 Я отправлю новые посты, но с некоторой задержкой (1-3 часа)"
	NotFound     = "🇬🇧 Rss feed not found\nPls send me direct link on rss\n\n🇷🇺 Rss поток не найден\nПришли, пожалуйста, прямую ссылку на rss\n"
	NewChannel   = "🇬🇧 Add @telefeedbot as admin in channel\nSend me link on channel, example: https://t.me/channel\n\n 🇷🇺 Добавь @telefeedbot как админа в канал\nПришли ссылку на канал в формате: https://t.me/channel\n"
	SubsHelp     = "🇬🇧 Commands:\nAdd url:\n@channelname url(s)\nDelete url(s):\n@channelname delete url"
	Rate         = "Please rate me here ❤️❤️❤️:\nhttps://storebot.me/bot/telefeedbot\n\nSupport(если что-то не работает кликай сюда): https://t.me/joinchat/AAAAAEMFJOGkHNVp8qKQ1g"
	TopLinks     = `
	Top rss feeds (by subscribers):

0 https://www.reddit.com/r/gifs/top/.rss

1 http://itc.ua/feed/

2 http://pikabu.ru/xmlfeeds.php?cmd=popular

3 https://web.stagram.com/rss/n/p

4 https://wylsa.com/feed/

5 http://news.liga.net/all/rss.xml

6 http://www.opennet.ru/opennews/opennews_all.rss

7 http://feeds.feedburner.com/macdigger/

8 http://droider.ru/feed/

9 https://xakep.ru/feed/

More examples: http://telegra.ph/telefeedbot-05-12
	
`
)

var (
	host               = os.Getenv("BOLT_HOST")
	Api                = "http://" + host + "/bolt"
	BaseUri            = Api + "/"
	Publics            = Api + PubNames
	Feeds              = Api + Feed
	Links              = Api + links
	Video              = Api + video
	Users              = Api + users
	Channels           = Api + channels
	Bests              = Api + best
	ViralHash          = Api + viralhash
	UserName           = Api + username
	Subs               = Api + pubSubTg
	FeedSubs           = Api + feedSubTg
	UserSubs           = Api + userSubTg
	LastPost           = Api + lastPost
	Telefeedfile       = "./telefeed.bot"
	ChannelsFatherfile = "./channelsfather.bot"
	Tokens             = [...]string{"token1", "token2"}
	Stores             = [...]string{"@telefeedcontent1", "@telefeedcontent2", "@telefeedcontent3"}
	StoreIds           = [...]int64{-1001140338639, -1001144965998, -1001122084977, -1001121449455, -1001147806509, -1001069985583, -1001128095164}
	GooglKeys          = [...]string{"AIzaSyCTmUsTGqjl7iWJLiJisrejgTNamp7bfIA", "AIzaSyAZaQiAkSmYFZLMjUOxtKOj3R29TPs81X0"}
)

func init() {
	/*
		f, err := os.OpenFile("./testlogfile.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.SetOutput(f)
		}
		defer f.Close()*/

	log.SetOutput(ioutil.Discard)
}
