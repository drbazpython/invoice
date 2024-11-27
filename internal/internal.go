package internal
import (
	"time"
	"os"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
)

//DefineMyLogo ...
func DefineMyLogo() string {
	multiLineString := `

      _      _               
   __| |_ __| |__   __ _ ____
  / _  | __ | _  \ / _  |_  /
 | (_| | |  | |_) | (_| |/ / 
  \__._|_|  |_.__/ \__._/___|
		
  `
	return multiLineString
}

//DefineAppLogo ...
func DefineAppLogo() string {
	multiLineString := `
 ____             _       _     
| __ ) _   _ _ __(_) __ _| |___ 
|  _ \| | | | '__| |/ _' | / __|
| |_) | |_| | |  | | (_| | \__ \
|____/ \__,_|_|  |_|\__,_|_|___/
      
	`
	return multiLineString
}

//DefineLogger ...
func DefineLogger(level log.Level) log.Logger {
	logger := log.NewWithOptions(os.Stderr, log.Options{
    ReportCaller: true,
    ReportTimestamp: true,
    TimeFormat: time.Kitchen,
    Prefix: "🪲 ",
	})
	styles := log.DefaultStyles()
	styles.Levels[log.DebugLevel] = lipgloss.NewStyle().
		SetString("DEBUG").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("204")).
		Foreground(lipgloss.Color("0"))
	styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
		SetString("INFO").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("301")).
		Foreground(lipgloss.Color("0"))
	logger.SetLevel(level)
	logger.SetReportTimestamp(false)
	logger.SetStyles(styles)
	return *logger
}

//MyLipglossStyle ...
func MyLipglossStyle(word string) lipgloss.Style {
	width := 40
	var style = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#FAFAFA")).
    Background(lipgloss.Color("#7D56F4")).
    PaddingTop(1).
	PaddingBottom(1).
    PaddingLeft(4).
	PaddingRight(4).
	Align(lipgloss.Center).
    Width(width).
	Margin(1)

	return style

}

//GetHourlyRate ...
func GetHourlyRate() string{
err := godotenv.Load(".env")
    if err != nil {
    	log.Fatal("Error loading .env file")
    } 
	return os.Getenv("HOURLY_RATE")
	}