package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/hirosassa/zerodriver"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"

	telebot "gopkg.in/telebot.v3"
)

var (
	TeleToken = os.Getenv("TELE_TOKEN")
	MetricsHost = os.Getenv("METRICS_HOST")
)

// Initialize OpenTelemetry
func initMetrics(ctx context.Context) {

	// Create a new OTLP Metric gRPC exporter with the specified endpoint and options
	exporter, _ := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(MetricsHost),
		otlpmetricgrpc.WithInsecure(),
	)

	// Define the resource with attributes that are common to all metrics.
	// labels/tags/resources that are common to all metrics.
	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(fmt.Sprintf("kbot_%s", appVersion)),
	)

	// Create a new MeterProvider with the specified resource and reader
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithReader(
			// collects and exports metric data every 10 seconds.
			sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(10*time.Second)),
		),
	)

	// Set the global MeterProvider to the newly created MeterProvider
	otel.SetMeterProvider(mp)

}

func pmetrics(ctx context.Context, payload string) {
	// Get the global MeterProvider and create a new Meter with the name "kbot_light_signal_counter"
	meter := otel.GetMeterProvider().Meter("kbot_weather_by_city")

	message := "";

	if strings.HasPrefix(strings.ToLower(payload), "weather ") {
		message = strings.TrimPrefix(strings.ToLower(payload), "weather ")
	} else {
		return
	}

	// Get or create an Int64Counter instrument with the name "kbot_weather_by_<message>"
	counter, _ := meter.Int64Counter(fmt.Sprintf("kbot_weather_by_%s", message))

	// Add a value of 1 to the Int64Counter
	counter.Add(ctx, 1)
}

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

		logger := zerodriver.NewProductionLogger();

		kbot, err := telebot.NewBot(telebot.Settings{
			URL:    "",
			Token:  TeleToken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})

		if err != nil {
			logger.Fatal()
			return
		} else {
			logger.Info().Str("Version", appVersion).Msg("kbot started")
		}

		kbot.Handle(telebot.OnText, func(m telebot.Context) error {
			logger.Info().Str("Payload", m.Text()).Msg(m.Message().Payload)
			msg := m.Text()

			pmetrics(context.Background(), msg)

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
	ctx := context.Background()
	initMetrics(ctx)
	rootCmd.AddCommand(kbotCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kbotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kbotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
