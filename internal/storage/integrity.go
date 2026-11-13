package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"

	"github.com/rottentail/solo-0026-audiovisual-sync-conformer/internal/domain"
)

type publicationEvidence struct {
	bundleDigest  string
	bundleBytes   int
	receiptDigest string
}

type publicationIndex struct {
	mu      sync.RWMutex
	entries map[string]publicationEvidence
}

func newPublicationIndex() *publicationIndex {
	return &publicationIndex{entries: make(map[string]publicationEvidence)}
}

func (index *publicationIndex) remember(
	id string,
	bundle []byte,
	receipt []byte,
) {
	index.mu.Lock()
	defer index.mu.Unlock()
	index.entries[id] = publicationEvidence{
		bundleDigest:  shaOrEmpty(bundle),
		bundleBytes:   len(bundle),
		receiptDigest: shaOrEmpty(receipt),
	}
}

func (index *publicationIndex) matchesRecord(record domain.ConformJob) bool {
	index.mu.RLock()
	defer index.mu.RUnlock()
	evidence, ok := index.entries[record.ID]
	return ok &&
		evidence.bundleDigest == record.BundleDigest &&
		evidence.bundleBytes == record.BundleBytes
}

func (index *publicationIndex) matchesFiles(
	record domain.ConformJob,
	bundle []byte,
	receipt []byte,
) bool {
	index.mu.RLock()
	defer index.mu.RUnlock()
	evidence, ok := index.entries[record.ID]
	return ok &&
		evidence.bundleDigest == shaOrEmpty(bundle) &&
		evidence.bundleBytes == len(bundle) &&
		evidence.receiptDigest == shaOrEmpty(receipt)
}

func shaOrEmpty(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
