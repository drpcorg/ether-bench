package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/INFURA/go-ethlibs/node"
	"github.com/jessevdk/go-flags"
	ethspam "github.com/p2p-org/ethspam/lib"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type Options struct {
	Target     string `long:"target" short:"t" description:"target eth host" required:"false"`
	SourceHost string `long:"source" short:"s" description:"source eth host" default:"https://eth.drpc.org"`
	Scenario   string `long:"scenario" short:"c" description:"scenario file" default:"stages.json"`
	Result     string `long:"result" short:"r" description:"result file" default:"result.json"`
}

func main() {
	options := Options{}
	_, err := flags.Parse(&options)
	if err != nil {
		return
	}
	fmt.Println("Start")

	jsonString, err := os.ReadFile(options.Scenario)
	if err != nil {
		exit(1, "stages reading error: %v", err)
	}

	var stages []Stage
	err = json.Unmarshal(jsonString, &stages)
	if err != nil {
		exit(1, "stages parsing error: %v", err)
	}

	stagesResults := make([]StageResult, 0)

	for _, stage := range stages {
		fmt.Printf("Stage %v\n", stage.Name)
		var stageMetrics vegeta.Metrics
		stepResults := make([]StepResult, 0)
		target := stage.Target
		if target == "" {
			target = options.Target
		}
		for n, step := range stage.Steps {
			fmt.Printf("Step %v, rate - %v\n", n+1, step.Rate)
			rate := vegeta.Rate{Freq: step.Rate, Per: time.Second}
			duration := time.Duration(step.Duration) * time.Second
			ctx, cancel := context.WithCancel(context.Background())
			targeter := NewEthSpamTargeter(ctx,
				target,
				stage.Profile,
				options.SourceHost,)
			attacker := vegeta.NewAttacker()

			var metrics vegeta.Metrics
			for res := range attacker.Attack(targeter, rate, duration, stage.Name) {
				processEthErrors(res)
				metrics.Add(res)
				stageMetrics.Add(res)
			}
			attacker.Stop()
			metrics.Close()
			cancel()
			stepResults = append(stepResults, StepResult{StepSummary: metrics})
			if metrics.Success < 0.85 {
				break
			}
			fmt.Printf("Step %v - succ: %v - mean: %v - p90: %v \n", n+1, metrics.Success, metrics.Latencies.Mean, metrics.Latencies.P90)
			<-time.After(10 * time.Second)
		}
		stageMetrics.Close()
		stageResult := StageResult{
			Name:         stage.Name,
			StageSummary: stageMetrics,
			Steps:        stepResults,
		}
		stagesResults = append(stagesResults, stageResult)
		fmt.Printf("Stage %v - succ: %v - mean: %v - p90: %v \n", stage.Name, stageMetrics.Success, stageMetrics.Latencies.Mean, stageMetrics.Latencies.P90)
	}
	fmt.Println("Finish")

	res, _ := json.Marshal(stagesResults)
	os.WriteFile(options.Result, res, 0644)
}

func processEthErrors(data *vegeta.Result) {
	if data.Code != 200 {
		return
	}
	var result RpcResponse
	err := json.Unmarshal(data.Body, &result)
	if err != nil {
		fmt.Printf("Error unmarshalling response JSON: %v\n", err)
		data.Code = 500
		data.Error = err.Error()
		return
	}

	if result.Error != nil {
		if strings.Contains(result.Error.Message, "reverted") || strings.Contains(result.Error.Data, "Reverted") {
			// not an error
			return
		}

		data.Code = 500
		data.Error = result.Error.Message
	}
}

func NewEthSpamTargeter(ctx context.Context, host string, queryParams map[string]int64, parentHost string) vegeta.Targeter {
	if host == "" {
		exit(1, "host is required")
	}

	header := http.Header{}
	header.Add("Content-Type", "application/json; charset=utf8")

	generator, err := ethspam.MakeQueriesGenerator(queryParams)
	if err != nil {
		exit(1, "failed to install defaults: %s", err)
	}

	client, err := node.NewClient(ctx, parentHost)
	if err != nil {
		exit(1, "failed to make a new client: %s", err)
	}
	mkState := ethspam.StateProducer{
		Client: client,
	}

	stateChannel := make(chan ethspam.State, 1)

	go func() {
		randSrc := rand.NewSource(time.Now().UnixNano())
		state := ethspam.LiveState{
			IdGen:   &ethspam.IdGenerator{},
			RandSrc: randSrc,
		}
		defer close(stateChannel)
		for {
			newState, err := mkState.Refresh(&state)
			if err != nil {
				// It can happen in some testnets that most of the blocks
				// are empty(no transaction included), don't refresh the
				// QueriesGenerator state without new inclusion.
				if err == ethspam.ErrEmptyBlock {
					select {
					case <-time.After(5 * time.Second):
					case <-ctx.Done():
						return
					}
					continue
				}
				fmt.Printf("failed to refresh state: %s", err)
				<-time.After(1 * time.Second)
				continue
			}
			select {
			case stateChannel <- newState:
			case <-ctx.Done():
				return
			}

			select {
			case <-time.After(15 * time.Second):
			case <-ctx.Done():
				return
			}
		}
	}()

	state := <-stateChannel

	queries := make(chan string)

	go func() {
		defer close(queries)
		for {
			// Update state when a new one is emitted
			select {
			case state = <-stateChannel:
			case <-ctx.Done():
				return
			default:
			}
			if q, err := generator.Query(state); err == io.EOF {
				return
			} else if err != nil {
				exit(2, "failed to write generated query: %s", err)
			} else {
				select {
				case queries <- q.GetBody():
				case <-ctx.Done():
					return
				}
				
			}
		}
	}()

	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		req, ok := <-queries

		if !ok {
			return vegeta.ErrNoTargets
		}

		tgt.URL = host
		tgt.Method = "POST"
		tgt.Header = header

		tgt.Body = []byte(req)

		return nil
	}
}

func exit(code int, format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(code)
}
