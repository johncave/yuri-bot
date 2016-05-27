package main

import (
	"flag"
	"fmt"
	"strings"
	"database/sql"
	"regexp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/bwmarrin/discordgo"
)

var (
	Email    string
	Password string
	Token    string
	BotID    string
)

func init() {

	flag.StringVar(&Email, "e", "", "Account Email")
	flag.StringVar(&Password, "p", "", "Account Password")
	flag.StringVar(&Token, "t", "", "Account Token")
	flag.Parse()
}

var userid int

func main() {



	// Create a new Discord session using the provided login information.
	// Use discordgo.New(Token) to just use a token for login.
	dg, err := discordgo.New(Email, Password, Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)

	// Open the websocket and begin listening.
	dg.Open()

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	//First, see if this image contains the keyword ~yuri
	if(strings.HasPrefix(strings.ToLower(m.Content), "~yuri")){
	//The image does contain the keyword. Connect to the database and get ready to do some work.
		s.ChannelTyping(m.ChannelID)

		db, err := sql.Open("mysql", "user:password@/database")
		if err != nil {
        		//panic(err.Error()) // proper error handling instead of panic in your app
			s.ChannelMessageSend(m.ChannelID, "Error: Could not connect to database. Please contact my author and remember I am a work in progress.")
    		}

		//Remove the keyword from the inputted
		var comm = strings.TrimPrefix(m.Content, "~yuri ")
		var words = strings.Split(comm, " ")
		
		fmt.Printf("Ran command \"%s\" for user %s\n", comm, m.Author.Username)
		//fmt.Printf("Command: %s\n", words[0])
		if(strings.ToLower(words[0]) == "id"){
			
			//Open up the images table for reading.
			re := regexp.MustCompile( "[^0-9]" )
			image, err := db.Query("SELECT title, description, file_name, width, height from image WHERE approved = 'yes' AND id = ?", re.ReplaceAllString(words[1], ""))
			if err != nil {
        			//panic(err.Error()) // proper error handling instead of panic in your app
				s.ChannelMessageSend(m.ChannelID, "Error: DB-table CNS. Please contact my author and remember I am a work in progress.")
    			}
			defer image.Close()
			var title string
			var description string
			var file_name string
			var width string
			var height string
			for image.Next(){
				err := image.Scan(&title, &description, &file_name, &width, &height)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "Sorry. This image does not appear to exist.")
				}

			}

			if(title == ""){
				s.ChannelMessageSend(m.ChannelID, "Sorry. This image does not appear to exist.")
			} else{
				var urlname = strings.Replace(file_name, " ", "%20", -1)
				s.ChannelMessageSend(m.ChannelID, "**"+title+"**\n"+description+"\n"+width+"W x "+height+"H\nhttps://img.yuriplease.com/full/0/"+urlname)
			}
		} else if(strings.ToLower(words[0]) == "random"){
			
			//Open up a random table of the database for reading.
			image, err := db.Query("SELECT title, description, file_name, width, height from image WHERE approved = 'yes' AND id >= (SELECT FLOOR( MAX(id) * RAND()) FROM `image` ) ORDER BY id LIMIT 1")
			if err != nil {
        			//panic(err.Error()) // proper error handling instead of panic in your app
				s.ChannelMessageSend(m.ChannelID, "Error: DB-rand-table CNS. Please contact my author and remember I am a work in progress.")
    			}
			defer image.Close()
			var title string
			var description string
			var file_name string
			var width string
			var height string
			for image.Next(){
				err := image.Scan(&title, &description, &file_name, &width, &height)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "Sorry. This image does not appear to exist.")
				}

			}

				var urlname = strings.Replace(file_name, " ", "%20", -1)
				s.ChannelMessageSend(m.ChannelID, "**"+title+"**\n"+description+"\n"+width+"W x "+height+"H\nhttps://img.yuriplease.com/full/0/"+urlname)
		} else if(strings.ToLower(words[0]) == "help"){
			s.ChannelMessageSend(m.ChannelID, "Command list: \n *id <number>* post image number <number> into channel.\n*random* post a random image to the chat\n*help* print this help message\n*more* inform the bot author that they need to make more functions.") 
		} else{
			s.ChannelMessageSend(m.ChannelID, "Command not recognised. Type \"~yuri help\" for a list of commands.") 
		}

	} else{
		//fmt.Print("Message does not start with keyword. Ignore.")
	}
	// Print message to stdout.
	//fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)
	//fmt.Printf("%s says %s", m.Author.Username, m.Content)
}











