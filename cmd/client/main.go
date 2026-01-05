package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Saintrad/todo-server-client/internal/apiclient"
)

func main() {
	baseURL := envOrDefault("TODO_BASE_URL", "http://localhost:8080")

	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	c := apiclient.New(baseURL)

	switch cmd {
	case "list":
		if err := cmdList(c); err != nil {
			fail(err)
		}

	case "create":
		if err := cmdCreate(c, args); err != nil {
			fail(err)
		}

	case "get":
		if err := cmdGet(c, args); err != nil {
			fail(err)
		}

	case "update":
		if err := cmdUpdate(c, args); err != nil {
			fail(err)
		}

	case "delete":
		if err := cmdDelete(c, args); err != nil {
			fail(err)
		}

	default:
		fmt.Fprintln(os.Stderr, "unknown command:", cmd)
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `Usage:
  client list

  client create --title "..." [--category "work"] [--due "2026-01-10"]
  client get <id>
  client update <id> [--title "..."] [--category "..."] [--due "YYYY-MM-DD"] [--done | --undone]
  client delete <id>

Environment:
  TODO_BASE_URL (default http://localhost:8080)
`)
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}

func cmdList(c *apiclient.Client) error {
	tasks, err := c.ListTasks()
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		fmt.Println("(no tasks)")
		return nil
	}
	for _, t := range tasks {
		box := " "
		if t.IsDone {
			box = "x"
		}
		cat := ""
		if t.Category != nil && *t.Category != "" {
			cat = " (" + *t.Category + ")"
		}
		due := ""
		if t.DueDate != nil {
			due = " due:" + t.DueDate.Format("2006-01-02")
		}
		fmt.Printf("%d [%s] %s%s%s\n", t.ID, box, t.Title, cat, due)
	}
	return nil
}

func cmdCreate(c *apiclient.Client, args []string) error {
	fs := flag.NewFlagSet("create", flag.ContinueOnError)
	fs.SetOutput(ioDiscard{}) // suppress default flag error printing; we handle errors ourselves

	title := fs.String("title", "", "task title (required)")
	category := fs.String("category", "", "optional category")
	due := fs.String("due", "", "optional due date in YYYY-MM-DD")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(*title) == "" {
		return fmt.Errorf("--title is required")
	}

	var catPtr *string
	if strings.TrimSpace(*category) != "" {
		v := strings.TrimSpace(*category)
		catPtr = &v
	}

	var duePtr *time.Time
	if strings.TrimSpace(*due) != "" {
		tm, err := time.Parse("2006-01-02", strings.TrimSpace(*due))
		if err != nil {
			return fmt.Errorf("invalid --due (expected YYYY-MM-DD): %w", err)
		}
		duePtr = &tm
	}

	created, err := c.CreateTask(apiclient.CreateTaskRequest{
		Title:    strings.TrimSpace(*title),
		Category: catPtr,
		DueDate:  duePtr,
	})
	if err != nil {
		return err
	}

	fmt.Printf("created task %d: %s\n", created.ID, created.Title)
	return nil
}

func cmdGet(c *apiclient.Client, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: client get <id>")
	}
	id, err := strconv.Atoi(args[0])
	if err != nil || id <= 0 {
		return fmt.Errorf("invalid id: %s", args[0])
	}

	t, err := c.GetTask(id)
	if err != nil {
		return err
	}
	printTask(t)
	return nil
}

func cmdUpdate(c *apiclient.Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: client update <id> [--title ...] [--category ...] [--due ...] [--done|--undone]")
	}
	id, err := strconv.Atoi(args[0])
	if err != nil || id <= 0 {
		return fmt.Errorf("invalid id: %s", args[0])
	}

	fs := flag.NewFlagSet("update", flag.ContinueOnError)
	fs.SetOutput(ioDiscard{})

	title := fs.String("title", "", "new title")
	category := fs.String("category", "", "new category")
	due := fs.String("due", "", "new due date in YYYY-MM-DD")
	done := fs.Bool("done", false, "mark as done")
	undone := fs.Bool("undone", false, "mark as not done")

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	if *done && *undone {
		return fmt.Errorf("use only one of --done or --undone")
	}

	var req apiclient.UpdateTaskRequest
	changed := false

	if fs.Lookup("title").Value.String() != "" {
		v := strings.TrimSpace(*title)
		// Allow setting empty title? Typically no; you can decide policy here.
		// We'll treat empty as invalid if flag provided.
		if v == "" {
			return fmt.Errorf("--title cannot be empty")
		}
		req.Title = &v
		changed = true
	}

	if fs.Lookup("category").Value.String() != "" {
		v := strings.TrimSpace(*category)
		// This sets category to provided string; clearing category requires an explicit design later.
		req.Category = &v
		changed = true
	}

	if strings.TrimSpace(*due) != "" {
		tm, err := time.Parse("2006-01-02", strings.TrimSpace(*due))
		if err != nil {
			return fmt.Errorf("invalid --due (expected YYYY-MM-DD): %w", err)
		}
		req.DueDate = &tm
		changed = true
	}

	if *done {
		v := true
		req.IsDone = &v
		changed = true
	}
	if *undone {
		v := false
		req.IsDone = &v
		changed = true
	}

	if !changed {
		return fmt.Errorf("no update fields provided")
	}

	updated, err := c.UpdateTask(id, req)
	if err != nil {
		return err
	}
	fmt.Printf("updated task %d\n", updated.ID)
	return nil
}

func cmdDelete(c *apiclient.Client, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: client delete <id>")
	}
	id, err := strconv.Atoi(args[0])
	if err != nil || id <= 0 {
		return fmt.Errorf("invalid id: %s", args[0])
	}
	if err := c.DeleteTask(id); err != nil {
		return err
	}
	fmt.Printf("deleted task %d\n", id)
	return nil
}

func printTask(t apiclient.Task) {
	fmt.Printf("ID: %d\n", t.ID)
	fmt.Printf("Title: %s\n", t.Title)
	fmt.Printf("Done: %v\n", t.IsDone)
	if t.Category != nil && *t.Category != "" {
		fmt.Printf("Category: %s\n", *t.Category)
	}
	if t.DueDate != nil {
		fmt.Printf("Due: %s\n", t.DueDate.Format("2006-01-02"))
	}
	fmt.Printf("CreatedAt: %s\n", t.CreatedAt.Format(time.RFC3339))
	if t.UpdatedAt != nil {
		fmt.Printf("UpdatedAt: %s\n", t.UpdatedAt.Format(time.RFC3339))
	}
}

// env helper
func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// ioDiscard suppresses FlagSet default output; avoids double-printing.
type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) { return len(p), nil }
