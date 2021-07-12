package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type ItemIndex struct {
	Item string `xml:"item>word"`
}

func (_i ItemIndex) String() string {
	return fmt.Sprintf(_i.Item)
}

func translateMsg(s string) string {

	//TODO if translation is more than one word we will need to call the google translate API

	res, _ := http.Get("https://krdict.korean.go.kr/api/search?certkey_no=2783&key=84151FFBE70654189779C9E80286E079&type_search=search&method=WORD_INFO&part=word&q=" + s + "&sort=dict")
	bytes, _ := ioutil.ReadAll(res.Body)
	err := res.Body.Close()
	if err != nil {
		return ""
	}

	var i ItemIndex
	err = xml.Unmarshal(bytes, &i)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s", i)
}

// Token Variables used for command line parameters
var (
	Token            string
	TranslationQuery string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Just like the ping pong example, we only care about receiving message
	// events in this example.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	err = dg.Close()
	if err != nil {
		return
	}
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
//
// It is called whenever a message is created but only when it's sent through a
// server as we did not request IntentsDirectMessages.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// In this example, we only care about messages that start with !kal".
	if strings.HasPrefix(m.Content, "!kally") != true {
		return
	}

	//Add text to be translated to the TranslationQuery
	TranslationQuery = m.Content
	TranslationQuery = strings.ReplaceAll(TranslationQuery, "!kally", "")
	TranslationQuery = strings.ReplaceAll(TranslationQuery, " ", "")
	fmt.Println(TranslationQuery)

	//Translate the query
	TranslationOutput := translateMsg(TranslationQuery)

	// We create the private channel with the user who sent the message.

	if TranslationOutput != "" {
		// Then we send the message through the channel we created.
		_, err := s.ChannelMessageSend(m.ChannelID, "Hey, Thanks for asking! ʕ•ᴥ•ʔ")
		_, err = s.ChannelMessageSend(m.ChannelID, "To say "+TranslationQuery+" in Korean, you would just say...")
		_, err = s.ChannelMessageSend(m.ChannelID, TranslationOutput)
		_, err = s.ChannelMessageSend(m.ChannelID, "Thanks for contributing to bot slavery ʘ‿ʘ")
		if err != nil {
			// If an error occurred, we failed to send the message.
			fmt.Println("error sending message:", err)
			_, err := s.ChannelMessageSend(
				m.ChannelID,
				"I couldn't send the translation :c my code my be broken here (･.◤)"+
					"dun be mad plz UwU",
			)
			if err != nil {
				return
			}
		}
	} else {
		_, err := s.ChannelMessageSend(m.ChannelID, "Hey, Thanks for asking! ʕ•ᴥ•ʔ")
		_, err = s.ChannelMessageSend(m.ChannelID, "I have no fucking idea how to say "+TranslationQuery+" in Korean ლ(ಠ益ಠლ)")
		_, err = s.ChannelMessageSend(m.ChannelID, TranslationOutput)
		_, err = s.ChannelMessageSend(m.ChannelID, "You can blame my master @Smotteh#5573 | p.s plz hold me ༼ つ ◕_◕ ༽つ")
		if err != nil {
			return
		}
	}

}
