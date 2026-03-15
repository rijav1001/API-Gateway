package loadbalancer

import (
	"net/url"
	"sync"
)

type LoadBalancer struct {
	mu      sync.Mutex
	targets []*url.URL
	current int
}

func New(targets []string) (*LoadBalancer, error) {
	var urls []*url.URL
	for _, t := range targets {
		u, err := url.Parse(t)
		if err != nil {
			return nil, err
		}
		urls = append(urls, u)
	}

	return &LoadBalancer{targets: urls}, nil
}

func (lb *LoadBalancer) Next() *url.URL {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	target := lb.targets[lb.current % len(lb.targets)]
	lb.current++
	return target
}

func (lb *LoadBalancer) AddTarget(target string) error {
	u, err := url.Parse(target)
	if err != nil {
		return err
	}
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.targets = append(lb.targets, u)
	return nil
}

func (lb *LoadBalancer) Targets() []*url.URL {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	return lb.targets
}