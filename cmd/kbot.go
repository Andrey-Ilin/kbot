package cmd

import (
	"fmt"
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
	url := fmt.Sprintf("https://wttr.in/%s?format=%%C+%%t+%%w", city)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var weather string
	_, err = fmt.Fscanf(resp.Body, "%s", &weather)
	if err != nil {
		return "", err
	}

	return weather, nil
}

var kbotCmd = &cobra.Command{
	Use:     "kbot",
	Aliases: []string{"start"},
	Short:   "A brief description of your command",
	Long:    `A longer description...`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("kbot %s started", appVersion)

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
			log.Printf(m.Message().Payload, m.Text())
			msg := m.Text()
			var err error

			if strings.HasPrefix(strings.ToLower(msg), "weather ") {
				city := strings.TrimPrefix(strings.ToLower(msg), "weather ")
				weather, err := getWeather(city)
				if err != nil {
					err = m.Send(fmt.Sprintf("Error getting weather for %s: %v", city, err))
				} else {
					err = m.Send(fmt.Sprintf("Weather in %s: %s", city, weather))
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
