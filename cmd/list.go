package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	gh "github.com/cli/go-gh/v2"
	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/cli/go-gh/v2/pkg/tableprinter"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List discussions of the current repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listDiscussions()
	},
}

var limit int = 30

func init() {
	listCmd.Flags().IntVarP(&limit, "limit", "l", limit, "Limit the number of discussions to display. Maximum is 100.")
}

func listDiscussions() error {
	// Get the current repository context
	repo, err := repository.Current()
	if err != nil {
		return fmt.Errorf("failed to get current repository: %w", err)
	}

	args := []string{"api", fmt.Sprintf("repos/%s/%s/discussions", repo.Owner, repo.Name), "--method", "GET", "-f", fmt.Sprintf("per_page=%d", limit)}
	stdOut, _, err := gh.Exec(args...)
	if err != nil {
		return fmt.Errorf("failed to get discussions: %w", err)
	}

	type Discussion struct {
		Number    int       `json:"number"`
		Title     string    `json:"title"`
		HtmlUrl   string    `json:"html_url"`
		State     string    `json:"state"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	var discussions []Discussion

	// Parse the JSON output
	if err := json.Unmarshal(stdOut.Bytes(), &discussions); err != nil {
		return fmt.Errorf("failed to parse discussions: %w", err)
	}

	// Create a new tab writer for table formatting
	printer := tableprinter.New(os.Stdout, true, 100)
	printer.AddHeader([]string{"ID", "TITLE", "STATE", "CREATED AT"}, tableprinter.WithColor(gray))

	// Print the pull requests in a table format
	for _, d := range discussions {
		var numberColorFunc func(string) string
		if d.State == "open" {
			numberColorFunc = green
		} else if d.State == "closed" {
			numberColorFunc = red
		} else {
			numberColorFunc = gray
		}

		printer.AddField(numberColorFunc(fmt.Sprintf("#%d", d.Number)))
		printer.AddField(d.Title)
		printer.AddField(d.State)
		printer.AddField(d.CreatedAt.Format("2006-01-02 15:04:05"))
		printer.EndRow()
	}

	// Render the table to the terminal
	if err := printer.Render(); err != nil {
		return fmt.Errorf("failed to render table: %w", err)
	}

	return nil
}

func green(text string) string {
	return "\033[32m" + text + "\033[0m"
}

func gray(text string) string {
	return "\033[90m" + text + "\033[0m"
}

func red(text string) string {
	return "\033[31m" + text + "\033[0m"
}
