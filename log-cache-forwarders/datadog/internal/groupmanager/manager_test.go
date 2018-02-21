package groupmanager_test

import (
	"context"
	"sync"
	"time"

	"code.cloudfoundry.org/loggregator-tools/log-cache-forwarders/datadog/internal/groupmanager"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manager", func() {
	var (
		stubGroupProvider *stubGroupProvider
		spyGroupUpdater   *spyGroupUpdater
		t                 chan time.Time
	)

	BeforeEach(func() {
		spyGroupUpdater = newSpyGroupUpdater()
		stubGroupProvider = newStubGroupProvider()
		t = make(chan time.Time, 1)

		groupmanager.Start("group-name", t, stubGroupProvider, spyGroupUpdater)
	})

	It("fetches meta and adds to the 0group", func() {
		stubGroupProvider.sourceIDs = []string{
			"source-id-1",
			"source-id-2",
		}

		t <- time.Now()

		Eventually(spyGroupUpdater.AddRequests).Should(ConsistOf(
			addRequest{name: "group-name", sourceID: "source-id-1"},
			addRequest{name: "group-name", sourceID: "source-id-2"},
		))
	})
})

func newStubGroupProvider() *stubGroupProvider {
	return &stubGroupProvider{}
}

type stubGroupProvider struct {
	sourceIDs []string
}

func (s *stubGroupProvider) SourceIDs() []string {
	return s.sourceIDs
}

func newSpyGroupUpdater() *spyGroupUpdater {
	return &spyGroupUpdater{}
}

type spyGroupUpdater struct {
	mu          sync.Mutex
	addRequests []addRequest
}

func (s *spyGroupUpdater) AddToGroup(ctx context.Context, name, sourceID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.addRequests = append(s.addRequests, addRequest{name: name, sourceID: sourceID})
	return nil
}

func (s *spyGroupUpdater) AddRequests() []addRequest {
	s.mu.Lock()
	defer s.mu.Unlock()

	r := make([]addRequest, len(s.addRequests))
	copy(r, s.addRequests)

	return r
}

type addRequest struct {
	name     string
	sourceID string
}
