package balancer

// DefaultBalancer is an global balancer
var DefaultBalancer = NewWeightedRoundRobin()

// Select gets next selected item.
func Select(_ ...string) interface{} {
	return DefaultBalancer.Select()
}

// Name load balancer name.
func Name() string {
	return DefaultBalancer.Name()
}

// Update reinitialize the balancer items.
func Update(choices []*Choice) bool {
	return DefaultBalancer.Update(choices)
}
