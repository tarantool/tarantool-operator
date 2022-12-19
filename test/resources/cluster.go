package resources

func (r *FakeCartridge) WithClusterName(name string) *FakeCartridge {
	r.Cluster.Name = name

	return r
}

func (r *FakeCartridge) Bootstrapped() *FakeCartridge {
	r.Cluster.Status.Bootstrapped = true

	return r
}

func (r *FakeCartridge) WithLeader(leaderName string) *FakeCartridge {
	r.Cluster.Status.Leader = leaderName

	return r
}
