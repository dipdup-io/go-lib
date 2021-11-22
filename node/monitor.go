package node

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"
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
	refused       chan []*Applied
	branchDelayed chan []*Applied
	branchRefused chan []*Applied

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
		refused:       make(chan []*Applied, 4096),
		branchDelayed: make(chan []*Applied, 4096),
		branchRefused: make(chan []*Applied, 4096),
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
func (monitor *Monitor) BranchRefused() <-chan []*Applied {
	return monitor.branchRefused
}

// BranchDelayed -
func (monitor *Monitor) BranchDelayed() <-chan []*Applied {
	return monitor.branchDelayed
}

// Refused -
func (monitor *Monitor) Refused() <-chan []*Applied {
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
			ch, err := monitor.selectChannel(filter)
			if err != nil {
				log.Error(err)
				continue
			}

			if err := monitor.longPolling(ctx, url, ch); err != nil {
				log.Error(err)
			}
		}
	}
}

func (monitor *Monitor) selectChannel(filter string) (interface{}, error) {
	switch filter {
	case filterApplied:
		return monitor.applied, nil
	case filterBranchDelayed:
		return monitor.branchDelayed, nil
	case filterBranchRefused:
		return monitor.branchRefused, nil
	case filterRefused:
		return monitor.refused, nil
	default:
		return nil, errors.Errorf("unknown filter: %s", filter)
	}
}

func (monitor *Monitor) longPolling(ctx context.Context, url string, ch interface{}) error {
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
	return monitor.parseLongPollingResponse(ctx, resp, ch)
}

func (monitor *Monitor) parseLongPollingResponse(ctx context.Context, resp *http.Response, ch interface{}) error {
	if resp == nil {
		return errors.New("nil response on mempool long polling request")
	}
	if ch == nil {
		return errors.New("nil output channel during mempool long polling request")
	}

	typ := reflect.TypeOf(ch)
	if typ.Kind() != reflect.Chan {
		return errors.Errorf("invalid channel type: %T", ch)
	}

	decoder := json.NewDecoder(resp.Body)
	cases := []reflect.SelectCase{
		{
			Dir:  reflect.SelectSend,
			Chan: reflect.ValueOf(ch),
		},
		{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ctx.Done()),
		},
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			value := reflect.New(typ.Elem())
			if err := decoder.Decode(value.Interface()); err != nil {
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					return nil
				}
				return err
			}
			cases[0].Send = value.Elem()
			if chosen, _, _ := reflect.Select(cases); chosen == 1 {
				return ctx.Err()
			}
		}
	}
}
