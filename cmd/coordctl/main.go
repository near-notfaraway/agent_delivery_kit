package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"example.com/agent_delivery_kit/internal/control"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "usage: coordctl <claim|complete|context|decide|skip|summary> [flags]")
		return 2
	}

	switch args[0] {
	case "claim":
		flags := flag.NewFlagSet("claim", flag.ContinueOnError)
		flags.SetOutput(stderr)
		project := flags.String("project", "", "path to project JSON")
		task := flags.String("task", "", "task ID")
		owner := flags.String("owner", "", "agent owner")
		if err := flags.Parse(args[1:]); err != nil {
			return 2
		}
		if *project == "" || *task == "" || *owner == "" {
			fmt.Fprintln(stderr, "project, task, and owner are required")
			return 2
		}
		if err := control.NewStore(*project).Claim(*task, *owner); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintf(stdout, "claimed %s by %s\n", *task, *owner)
		return 0

	case "complete":
		flags := flag.NewFlagSet("complete", flag.ContinueOnError)
		flags.SetOutput(stderr)
		project := flags.String("project", "", "path to project JSON")
		task := flags.String("task", "", "task ID")
		owner := flags.String("owner", "", "agent owner")
		commit := flags.String("commit", "", "commit or immutable revision")
		validation := flags.String("validation", "", "validation evidence")
		handoff := flags.String("handoff", "", "handoff summary")
		if err := flags.Parse(args[1:]); err != nil {
			return 2
		}
		if *commit == "" || *validation == "" || *handoff == "" {
			fmt.Fprintln(stderr, "commit, validation, and handoff are required")
			return 2
		}
		if *project == "" || *task == "" || *owner == "" {
			fmt.Fprintln(stderr, "project, task, and owner are required")
			return 2
		}
		evidence := control.Completion{Commit: *commit, Validation: []string{*validation}, Handoff: *handoff}
		if err := control.NewStore(*project).Complete(*task, *owner, evidence); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintf(stdout, "completed %s\n", *task)
		return 0

	case "context":
		flags := flag.NewFlagSet("context", flag.ContinueOnError)
		flags.SetOutput(stderr)
		project := flags.String("project", "", "path to project JSON")
		task := flags.String("task", "", "task ID")
		if err := flags.Parse(args[1:]); err != nil {
			return 2
		}
		if *project == "" || *task == "" {
			fmt.Fprintln(stderr, "project and task are required")
			return 2
		}
		context, err := control.NewStore(*project).Context(*task)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		return writeJSON(context, stdout, stderr)

	case "decide":
		flags := flag.NewFlagSet("decide", flag.ContinueOnError)
		flags.SetOutput(stderr)
		project := flags.String("project", "", "path to project JSON")
		id := flags.String("id", "", "decision ID")
		task := flags.String("task", "", "task ID")
		original := flags.String("original", "", "original design")
		problem := flags.String("problem", "", "discovered problem")
		choice := flags.String("choice", "", "chosen resolution")
		needsApproval := flags.Bool("needs-approval", false, "mark decision for human approval")
		if err := flags.Parse(args[1:]); err != nil {
			return 2
		}
		if *project == "" || *id == "" || *task == "" || *original == "" || *problem == "" || *choice == "" {
			fmt.Fprintln(stderr, "project, id, task, original, problem, and choice are required")
			return 2
		}
		decision := control.Decision{ID: *id, TaskID: *task, Original: *original, Problem: *problem, Choice: *choice, NeedsApproval: *needsApproval}
		if err := control.NewStore(*project).AddDecision(decision); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintf(stdout, "recorded decision %s\n", *id)
		return 0

	case "skip":
		flags := flag.NewFlagSet("skip", flag.ContinueOnError)
		flags.SetOutput(stderr)
		project := flags.String("project", "", "path to project JSON")
		task := flags.String("task", "", "task ID")
		reason := flags.String("reason", "", "reason the task is unnecessary")
		if err := flags.Parse(args[1:]); err != nil {
			return 2
		}
		if *project == "" || *task == "" || *reason == "" {
			fmt.Fprintln(stderr, "project, task, and reason are required")
			return 2
		}
		if err := control.NewStore(*project).Skip(*task, *reason); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintf(stdout, "skipped %s\n", *task)
		return 0

	case "summary":
		flags := flag.NewFlagSet("summary", flag.ContinueOnError)
		flags.SetOutput(stderr)
		project := flags.String("project", "", "path to project JSON")
		if err := flags.Parse(args[1:]); err != nil {
			return 2
		}
		if *project == "" {
			fmt.Fprintln(stderr, "project is required")
			return 2
		}
		summary, err := control.NewStore(*project).Summary()
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		return writeJSON(summary, stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown command %q\n", args[0])
		return 2
	}
}

func writeJSON(value any, stdout, stderr io.Writer) int {
	encoder := json.NewEncoder(stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(value); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}
