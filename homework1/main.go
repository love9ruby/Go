package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "cli is a simple command line tool",
}

func init() {

	csvCmd := &cobra.Command{
		Use:   "csv",
		Short: "read csv file",
		Long:  "read csv file by path",
		Run: func(cmd *cobra.Command, args []string) {
			path, err := cmd.Flags().GetString("path")
			if err != nil {
				println("error:", err)
				return
			}
			println("csv:", path)
			// read csv file
			text, err := os.ReadFile(path)
			if err != nil {
				println("error:", err)
				return
			}
			users := make(map[string]int)
			for _, line := range strings.Split(string(text), "\n")[1:] {
				if line == "" {
					continue
				}
				fields := strings.Split(line, ",")
				if len(fields) < 3 {
					continue
				}
				name := fields[2]
				if _, ok := users[name]; ok {
					users[name]++
					continue
				}
				users[name] = 1
			}
			for name, count := range users {
				println(name, count)
			}
		},
	}
	csvCmd.Flags().StringP("path", "p", "", "csv file path")
	if err := csvCmd.MarkFlagRequired("path"); err != nil {
		panic(err)
	}
	rootCmd.AddCommand(csvCmd)

	wgetCmd := &cobra.Command{
		Use:   "wget",
		Short: "wget website content",
		Long:  "wget website content by url",
		Run: func(cmd *cobra.Command, args []string) {
			url, err := cmd.Flags().GetString("url")
			if err != nil {
				println("error:", err)
				return
			}
			wget(url)
		},
	}
	wgetCmd.Flags().StringP("url", "u", "", "website url")
	if err := wgetCmd.MarkFlagRequired("url"); err != nil {
		panic(err)
	}
	rootCmd.AddCommand(wgetCmd)

	quoteCmd := &cobra.Command{
		Use:   "quote",
		Short: "Get a random quote",
		Long:  `This command prints a random quote from a famous person.`,
		Run: func(cmd *cobra.Command, args []string) {
			quote := getQuote()
			if quote == "" {
				panic("no quotes, please add some quotes first.")
				//return
			}
			fmt.Println(quote)
		},
	}
	rootCmd.AddCommand(quoteCmd)

	greetCmd := &cobra.Command{
		Use:   "greet",
		Short: "Greet someone",
		Long:  `This command greets someone by name.`,
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			fmt.Printf("Hello, %s!\n", name)
		},
	}
	greetCmd.Flags().StringP("name", "n", "World", "Name of the person")
	rootCmd.AddCommand(greetCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		return
	}
}

func getQuote() string {
	quotes := []string{
		"The journey of a thousand miles begins with one step. - Lao Tzu",
		"That which does not kill us makes us stronger. - Friedrich Nietzsche",
		"Life is like riding a bicycle. To keep your balance, you must keep moving. - Albert Einstein",
		"Life is really simple, but we insist on making it complicated. - Confucius",
		"Life is what happens when you're busy making other plans. - John Lennon",
		"Life is 10% what happens to us and 90% how we react to it. - Charles R. Swindoll",
		"Life is a series of natural and spontaneous changes. Don't resist them; that only creates sorrow. Let reality be reality. Let things flow naturally forward in whatever way they like. - Lao Tzu",
		"Life is a long lesson in humility. - James M. Barrie",
		"Life is a question and how we live it is our answer. - Gary Keller",
		"Life is a dream for the wise, a game for the fool, a comedy for the rich, a tragedy for the poor. - Sholom Aleichem",
		"Life is a gift, and it offers us the privilege, opportunity, and responsibility to give something back by becoming more. - Tony Robbins",
		"Life is a daring adventure or nothing at all. - Helen Keller",
		"Life is a great big canvas, and you should throw all the paint on it you can. - Danny Kaye",
		"Life is a mirror and will reflect back to the thinker what he thinks into it. - Ernest Holmes",
		"Life is a process of becoming, a combination of states we have to go through. Where people fail is that they wish to elect a state and remain in it. This is a kind of death. - Anais Nin",
		"Life is a series of natural and spontaneous changes. Don't resist them; that only creates sorrow. Let reality be reality. Let things flow naturally forward in whatever way they like. - Lao Tzu",
		// ... 添加更多名言
	}
	if len(quotes) == 0 {
		println("no quotes, please add some quotes first.")
		return ""
	}
	// init random generator
	source := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(source)
	// 獲取 0 到 n-1 之間的隨機整數
	randomNum := randomGenerator.Intn(len(quotes))
	quote := quotes[randomNum]
	return quote
}

func wget(url string) {
	println("wget:", url)
	// fetch website content
	get, err := http.Get(url)
	if err != nil {
		println("error:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			println("error:", err)
		}
	}(get.Body)
	// print website content
	if _, err = io.Copy(io.Writer(os.Stdout), get.Body); err != nil {
		println("error:", err)
	}
}
