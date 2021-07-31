package main

import (
	"cloud.google.com/go/translate"
	"context"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/language"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func translateText(targetLanguage, text string) (string, error) {
	ctx := context.Background()

	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return "", fmt.Errorf("language.Parse: %v", err)
	}

	client, err := translate.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer func(client *translate.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)

	resp, err := client.Translate(ctx, []string{text}, lang, nil)
	if err != nil {
		return "", fmt.Errorf("translate: %v", err)
	}
	if len(resp) == 0 {
		return "", fmt.Errorf("translate returned empty response to text: %s", text)
	}
	fmt.Println(resp[0].Text)
	return resp[0].Text, nil
}

// Token Variables used for command line parameters
var (
	Token string
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
// TODO make system for commands as the current method of making multiple commands is a bit redundant using lots of duplicate code
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	//Add text to be translated to the TranslationQuery
	TranslationQuery := m.Content

	// COMMAND - !kally
	// for this command, we only care about messages that start with !kally".
	if strings.HasPrefix(m.Content, "!kally") && !strings.HasPrefix(m.Content, "!kally -q") {
		fmt.Println("I'm working UwU")

		TranslationQuery = strings.ReplaceAll(TranslationQuery, "!kally", "")
		//TranslationQuery = strings.ReplaceAll(TranslationQuery, " ", "")
		fmt.Println(TranslationQuery)

		//Translate the query
		TranslationOutput, _ := translateText("ko", TranslationQuery)

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
						"You can blame my master @Smotteh#5573 | p.s plz hold me ༼ つ ◕_◕ ༽つ",
				)
				if err != nil {
					return
				}
			}
		}
	} else if strings.HasPrefix(m.Content, "!kally -q") { // COMMAND - !kally -q <-- this command is used for quick translations
		TranslationQuery = strings.ReplaceAll(TranslationQuery, "!kally -q", "")

		//Translate the query
		TranslationOutput, _ := translateText("ko", TranslationQuery)

		if TranslationOutput != "" {
			// Then we send the message through the channel from which we received the message
			_, err := s.ChannelMessageSend(m.ChannelID, TranslationOutput)
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
			_, err := s.ChannelMessageSend(m.ChannelID, "I have no idea how to say "+TranslationQuery+" in Korean ლ(ಠ益ಠლ)")
			if err != nil {
				return
			}
		}
	} else {
		return
	}

}
