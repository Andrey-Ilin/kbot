package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	telebot "gopkg.in/telebot.v3"
)

var (
	TeleToken = os.Getenv("TELE_TOKEN")
)

func getWeather(city string) (string, error) {
	url := fmt.Sprintf("https://wttr.in/%s?format=3", city)
	
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	weatherReport := string(body)
	
	// Add some emoji to make it more visually appealing
	weatherReport = strings.ReplaceAll(weatherReport, "°C", "°C🌡")
	weatherReport = strings.ReplaceAll(weatherReport, "km/h", "km/h💨")
	
	return weatherReport, nil
}

var kbotCmd = &cobra.Command{
	Use:     "start",
	Aliases: []string{"run"},
	Short:   "Start the Telegram bot",
	Long:    `Start the Telegram bot and begin listening for messages.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("kbot %s started\n", appVersion)

		kbot, err := telebot.NewBot(telebot.Settings{
			URL:    "",
			Token:  TeleToken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})

		if err != nil {
			log.Fatalf("Please check TELE_TOKEN env variable, %s", err)
			return
		}

		kbot.Handle(telebot.OnText, func(m telebot.Context) error {
			log.Printf("Received message: %s", m.Text())
			msg := m.Text()
			var err error

			if strings.HasPrefix(strings.ToLower(msg), "weather ") {
				city := strings.TrimPrefix(strings.ToLower(msg), "weather ")
				weather, err := getWeather(city)
				if err != nil {
					err = m.Send(fmt.Sprintf("Error getting weather for %s: %v", city, err))
				} else {
					// Add a title to the weather report
					weatherWithTitle := fmt.Sprintf("<b>Weather in %s</b>\n%s", city, weather)
					err = m.Send(weatherWithTitle, &telebot.SendOptions{ParseMode: telebot.ModeHTML})
				}
			} else {
				switch msg {
				case "hello":
					err = m.Send(fmt.Sprintf("Hello I'm kbot %s", appVersion))
				case "bye":
					err = m.Send(fmt.Sprintf("Bye from kbot %s", appVersion))
				case "whatsup":
					err = m.Send(fmt.Sprintf("All fine %s", appVersion))
				default:
					err = m.Send("I don't understand that command. Try 'weather [city]' to get weather information.")
				}
			}

			return err
		})

		kbot.Start()
	},
}

func init() {
	rootCmd.AddCommand(kbotCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kbotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kbotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
