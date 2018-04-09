package groupmanager

import (
	"context"
	"log"
	"time"
)

// Manager syncs the source IDs from the GroupProvider with the GroupUpdater.
// It does so at the configured interval.
type Manager struct {
	groupName string
	ticker    <-chan time.Time
	gp        GroupProvider
	gu        GroupUpdater
}

// GroupProvider returns the desired SourceIDs.
type GroupProvider interface {
	SourceIDs() []string
}

// GroupUpdater is used to add (and keep alive) the source IDs for a group.
type GroupUpdater interface {
	// SetShardGroup adds source IDs to the LogCache sub-groups.
	SetShardGroup(ctx context.Context, name string, sourceIDs ...string) error
}

// Start creates and starts a Manager.
func Start(groupName string, ticker <-chan time.Time, gp GroupProvider, gu GroupUpdater) {
	m := &Manager{
		groupName: groupName,
		ticker:    ticker,
		gp:        gp,
		gu:        gu,
	}

	go m.run()
}

func (m *Manager) run() {
	sourceIDs := m.gp.SourceIDs()

	m.updateSourceIDs(sourceIDs)

	for range m.ticker {
		sourceIDs := m.gp.SourceIDs()

		m.updateSourceIDs(sourceIDs)
	}
}

func (m *Manager) updateSourceIDs(sourceIDs []string) {
	start := time.Now()
	if err := m.gu.SetShardGroup(context.Background(), m.groupName, sourceIDs...); err != nil {
		log.Printf("failed to set shard group: %s", err)
	}
	log.Printf("Setting %d source ids took %s", len(sourceIDs), time.Since(start).String())
}