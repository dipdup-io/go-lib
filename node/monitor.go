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
	log "github.com/sirupsen/logrus"
)

const (
	filterApplied       = "applied"
	filterRefused       = "refused"
	filterBranchRefused = "branch_refused"
	filterBranchDelayed = "branch_delayed"
)

// Monitor -
type Monitor struct {
	url string

	applied       chan []*Applied
	refused       chan []*FailedMonitor
	branchDelayed chan []*FailedMonitor
	branchRefused chan []*FailedMonitor

	subscribedOnApplied       bool
	subscribedOnRefused       bool
	subscribedOnBranchDelayed bool
	subscribedOnBranchRefused bool

	wg sync.WaitGroup
}

// NewMonitor -
func NewMonitor(url string) *Monitor {
	return &Monitor{
		url:           strings.TrimSuffix(url, "/"),
		applied:       make(chan []*Applied, 4096),
		refused:       make(chan []*FailedMonitor, 4096),
		branchDelayed: make(chan []*FailedMonitor, 4096),
		branchRefused: make(chan []*FailedMonitor, 4096),
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

func (monitor *Monitor) Close() error {
	monitor.wg.Wait()

	close(monitor.applied)
	close(monitor.refused)
	close(monitor.branchDelayed)
	close(monitor.branchRefused)
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

func (monitor *Monitor) pollingMempool(ctx context.Context, filter string) {
	defer monitor.wg.Done()

	if filter == "" {
		filter = filterApplied
	}

	url := fmt.Sprintf("/chains/main/mempool/monitor_operations?%s", filter)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := monitor.process(ctx, filter, url); err != nil {
				log.Error(err)
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
	default:
		return errors.Errorf("unknown filter: %s", filter)
	}
}

func (monitor *Monitor) longPollingApplied(ctx context.Context, url string, ch chan []*Applied) error {
	link := fmt.Sprintf("%s/%s", monitor.url, url)
	req, err := http.NewRequest(http.MethodGet, link, nil)
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

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			for decoder.More() {
				value := make([]*Applied, 0)
				if err := decoder.Decode(&value); err != nil {
					if err == io.EOF || err == io.ErrUnexpectedEOF {
						return nil
					}
					return err
				}
				ch <- value
			}
			time.Sleep(time.Millisecond) // sleeping for CPU usage decreasing
		}
	}
}

func (monitor *Monitor) longPollingFailed(ctx context.Context, url string, ch chan []*FailedMonitor) error {
	link := fmt.Sprintf("%s/%s", monitor.url, url)
	req, err := http.NewRequest(http.MethodGet, link, nil)
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

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			for decoder.More() {
				value := make([]*FailedMonitor, 0)
				if err := decoder.Decode(&value); err != nil {
					if err == io.EOF || err == io.ErrUnexpectedEOF {
						return nil
					}
					return err
				}
				ch <- value
			}
			time.Sleep(time.Millisecond) // sleeping for CPU usage decreasing
		}
	}
}
