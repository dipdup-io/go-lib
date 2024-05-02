package node

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	filterApplied       = "applied"
	filterRefused       = "refused"
	filterBranchRefused = "branch_refused"
	filterBranchDelayed = "branch_delayed"
	filterOutdated      = "outdated"
)

// Monitor -
type Monitor struct {
	url    string
	client *http.Client

	applied       chan []*Applied
	refused       chan []*FailedMonitor
	branchDelayed chan []*FailedMonitor
	branchRefused chan []*FailedMonitor
	outdated      chan []*FailedMonitor

	subscribedOnApplied       bool
	subscribedOnRefused       bool
	subscribedOnBranchDelayed bool
	subscribedOnBranchRefused bool
	subscribedOnOutdated      bool

	wg sync.WaitGroup
}

// NewMonitor -
func NewMonitor(url string) *Monitor {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	return &Monitor{
		url:           strings.TrimSuffix(url, "/"),
		applied:       make(chan []*Applied, 4096),
		refused:       make(chan []*FailedMonitor, 4096),
		branchDelayed: make(chan []*FailedMonitor, 4096),
		branchRefused: make(chan []*FailedMonitor, 4096),
		outdated:      make(chan []*FailedMonitor, 4096),
		client: &http.Client{
			Transport: t,
			Timeout:   time.Minute,
		},
	}
}

// SubscribeOnMempoolApplied -
func (monitor *Monitor) SubscribeOnMempoolApplied(ctx context.Context) {
	if monitor.subscribedOnApplied {
		return
	}
	monitor.subscribedOnApplied = true

	monitor.wg.Add(1)
	go monitor.pollingMempool(ctx, filterApplied)
}

// SubscribeOnMempoolRefused -
func (monitor *Monitor) SubscribeOnMempoolRefused(ctx context.Context) {
	if monitor.subscribedOnRefused {
		return
	}
	monitor.subscribedOnRefused = true

	monitor.wg.Add(1)
	go monitor.pollingMempool(ctx, filterRefused)
}

// SubscribeOnMempoolBranchRefused -
func (monitor *Monitor) SubscribeOnMempoolBranchRefused(ctx context.Context) {
	if monitor.subscribedOnBranchRefused {
		return
	}
	monitor.subscribedOnBranchRefused = true

	monitor.wg.Add(1)
	go monitor.pollingMempool(ctx, filterBranchRefused)
}

// SubscribeOnMempoolBranchDelayed -
func (monitor *Monitor) SubscribeOnMempoolBranchDelayed(ctx context.Context) {
	if monitor.subscribedOnBranchDelayed {
		return
	}
	monitor.subscribedOnBranchDelayed = true

	monitor.wg.Add(1)
	go monitor.pollingMempool(ctx, filterBranchDelayed)
}

// SubscribeOnMempoolOutdated -
func (monitor *Monitor) SubscribeOnMempoolOutdated(ctx context.Context) {
	if monitor.subscribedOnOutdated {
		return
	}
	monitor.subscribedOnOutdated = true

	monitor.wg.Add(1)
	go monitor.pollingMempool(ctx, filterOutdated)
}

func (monitor *Monitor) Close() error {
	monitor.wg.Wait()

	close(monitor.applied)
	close(monitor.refused)
	close(monitor.branchDelayed)
	close(monitor.branchRefused)
	close(monitor.outdated)
	return nil
}

// Applied -
func (monitor *Monitor) Applied() <-chan []*Applied {
	return monitor.applied
}

// BranchRefused -
func (monitor *Monitor) BranchRefused() <-chan []*FailedMonitor {
	return monitor.branchRefused
}

// BranchDelayed -
func (monitor *Monitor) BranchDelayed() <-chan []*FailedMonitor {
	return monitor.branchDelayed
}

// Refused -
func (monitor *Monitor) Refused() <-chan []*FailedMonitor {
	return monitor.refused
}

// Outdated -
func (monitor *Monitor) Outdated() <-chan []*FailedMonitor {
	return monitor.outdated
}

func (monitor *Monitor) pollingMempool(ctx context.Context, filter string) {
	defer monitor.wg.Done()

	if filter == "" {
		filter = filterApplied
	}

	url := fmt.Sprintf("chains/main/mempool/monitor_operations?version=0&%s", filter)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := monitor.process(ctx, filter, url); err != nil {
				log.Err(err).Msg("")
				continue
			}

		}
	}
}

func (monitor *Monitor) process(ctx context.Context, filter, url string) error {
	switch filter {
	case filterApplied:
		return monitor.longPollingApplied(ctx, url, monitor.applied)
	case filterBranchDelayed:
		return monitor.longPollingFailed(ctx, url, monitor.branchDelayed)
	case filterBranchRefused:
		return monitor.longPollingFailed(ctx, url, monitor.branchRefused)
	case filterRefused:
		return monitor.longPollingFailed(ctx, url, monitor.refused)
	case filterOutdated:
		return monitor.longPollingFailed(ctx, url, monitor.outdated)
	default:
		return errors.Errorf("unknown filter: %s", filter)
	}
}

func (monitor *Monitor) longPollingApplied(ctx context.Context, url string, ch chan []*Applied) error {
	link := fmt.Sprintf("%s/%s", monitor.url, url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		return err
	}
	client := http.Client{
		Timeout: time.Minute,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	return monitor.parseLongPollingAppliedResponse(ctx, resp, ch)
}

func (monitor *Monitor) parseLongPollingAppliedResponse(ctx context.Context, resp *http.Response, ch chan []*Applied) error {
	if resp == nil {
		return errors.New("nil response on mempool long polling request")
	}
	if ch == nil {
		return errors.New("nil output channel during mempool long polling request")
	}

	decoder := json.NewDecoder(resp.Body)

	for decoder.More() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			value := make([]*Applied, 0)
			if err := decoder.Decode(&value); err != nil {
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					return nil
				}
				return err
			}
			ch <- value
		}
	}
	return nil
}

func (monitor *Monitor) longPollingFailed(ctx context.Context, url string, ch chan []*FailedMonitor) error {
	link := fmt.Sprintf("%s/%s", monitor.url, url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		return err
	}

	resp, err := monitor.client.Do(req)
	if err != nil {
		return err
	}
	return monitor.parseLongPollingFailedResponse(ctx, resp, ch)
}

func (monitor *Monitor) parseLongPollingFailedResponse(ctx context.Context, resp *http.Response, ch chan []*FailedMonitor) error {
	if resp == nil {
		return errors.New("nil response on mempool long polling request")
	}
	if ch == nil {
		return errors.New("nil output channel during mempool long polling request")
	}

	decoder := json.NewDecoder(resp.Body)

	for decoder.More() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			value := make([]*FailedMonitor, 0)
			if err := decoder.Decode(&value); err != nil {
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					return nil
				}
				return err
			}
			ch <- value
		}
	}
	return nil
}
