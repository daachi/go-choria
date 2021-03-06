package scoutcmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"

	scoutagent "github.com/choria-io/go-choria/scout/agent/scout"
	scoutclient "github.com/choria-io/go-choria/scout/client/scout"
)

type ResumeCommand struct {
	identities []string
	facts      []string
	classes    []string
	checks     []string
	combined   []string
	json       bool
	cfile      string
	verbose    bool
	log        *logrus.Entry
}

func NewResumeCommand(ids []string, classes []string, facts []string, combined []string, checks []string, json bool, cfile string, verbose bool, log *logrus.Entry) (*ResumeCommand, error) {
	return &ResumeCommand{
		identities: ids,
		classes:    classes,
		checks:     checks,
		facts:      facts,
		combined:   combined,
		json:       json,
		log:        log,
		cfile:      cfile,
		verbose:    verbose,
	}, nil
}

func (t *ResumeCommand) Run(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()

	sc, err := scoutclient.New(scoutclient.ConfigFile(t.cfile), scoutclient.Logger(t.log), scoutclient.Progress())
	if err != nil {
		return err
	}

	var checks = make([]interface{}, len(t.checks))
	for i, c := range t.checks {
		checks[i] = c
	}

	result, err := sc.OptionIdentityFilter(t.identities...).OptionClassFilter(t.classes...).OptionFactFilter(t.facts...).OptionCombinedFilter(t.combined...).Resume().Checks(checks).Do(ctx)
	if err != nil {
		return err
	}

	if t.json {
		return result.RenderResults(os.Stdout, scoutclient.JSONFormat, scoutclient.DisplayDDL, t.verbose, false, t.log)
	}

	if result.Stats().ResponsesCount() == 0 {
		return fmt.Errorf("no responses received")
	}

	mu := sync.Mutex{}
	triggered := 0
	shown := 0
	table := newMarkdownTable("Name", "Triggered", "Skipped", "Failed", "Message")

	result.EachOutput(func(r *scoutclient.ResumeOutput) {
		tr := &scoutagent.ResumeReply{}
		err = r.ParseResumeOutput(tr)
		if err != nil {
			t.log.Errorf("Could not parse output from %s: %s", r.ResultDetails().Sender(), err)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		triggered += len(tr.TransitionedChecks)

		if !t.verbose && r.ResultDetails().OK() && len(tr.FailedChecks) == 0 {
			return
		}

		shown++

		table.Append([]string{r.ResultDetails().Sender(), strings.Join(tr.TransitionedChecks, ", "), strings.Join(tr.SkippedChecks, ", "), strings.Join(tr.FailedChecks, ", "), r.ResultDetails().StatusMessage()})
	})

	if shown == 0 {
		fmt.Printf("Placed %d checks into normal check mode on %d nodes\n", triggered, result.Stats().ResponsesCount())
	} else {
		table.Render()
	}

	fmt.Println()
	return result.RenderResults(os.Stdout, scoutclient.TXTFooter, scoutclient.DisplayDDL, t.verbose, false, t.log)
}
