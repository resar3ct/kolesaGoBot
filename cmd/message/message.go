package message

import (
	"encoding/json"
	"io/ioutil"
	"kolesaGoBot/internal/models"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"gopkg.in/telebot.v3"
)

type Message struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type UpgradeBot struct {
	Bot   *telebot.Bot
	Users *models.UserModel
}

func (bot *UpgradeBot) RunServer(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			newMessage := Message{}
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal(b, &newMessage)
			if err != nil {
				log.Println(err)
			}
			bot.SendToAllUsers(newMessage.Title, newMessage.Body)
			log.Println(newMessage.Title, newMessage.Body)
		})

	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Println("starting server at", addr)
	server.ListenAndServe()
}

func (bot *UpgradeBot) SendToAllUsers(title string, body string) {
	users, _ := bot.Users.GetAllUsers()

	for _, v := range users {
		sendTextToTelegramChat(int(v.ChatId), title)
		log.Println(v.ChatId)
	}
}

func sendTextToTelegramChat(chatId int, text string) (string, error) {

	log.Printf("Sending %s to chat_id: %d", text, chatId)
	var telegramApi string = "https://api.telegram.org/bot" + "5620278688:AAFnD6ix0Z90APRkBa8FavEbvengQBPg_04" + "/sendMessage"
	response, err := http.PostForm(
		telegramApi,
		url.Values{
			"chat_id": {strconv.Itoa(chatId)},
			"text":    {text},
		})

	if err != nil {
		log.Printf("error when posting text to the chat: %s", err.Error())
		return "", err
	}
	defer response.Body.Close()

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
	if errRead != nil {
		log.Printf("error in parsing telegram answer %s", errRead.Error())
		return "", err
	}
	bodyString := string(bodyBytes)
	log.Printf("Body of Telegram Response: %s", bodyString)

	return bodyString, nil
}
