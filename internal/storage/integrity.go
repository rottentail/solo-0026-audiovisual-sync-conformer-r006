package storage

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/rottentail/solo-0026-audiovisual-sync-conformer/internal/domain"
)

// shaOrEmpty returns the SHA-256 digest of data, including empty input.
func shaOrEmpty(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// bundleMatchesRecord reports whether the committed bundle bytes still agree
// with the completed job record that published them. The record digest and
// byte count are the durable evidence of the publication, so this check
// holds across process restarts without any in-memory state.
func bundleMatchesRecord(record domain.ConformJob, bundle []byte) bool {
	return record.BundleDigest != "" &&
		record.BundleBytes == len(bundle) &&
		record.BundleDigest == shaOrEmpty(bundle)
}

// receiptMatchesRecord reports whether the published receipt still describes
// the completed job record, keeping the receipt and bundle routes in
// agreement for any reader of the data directory.
func receiptMatchesRecord(record domain.ConformJob, receipt domain.ConformReceipt) bool {
	return receipt.SchemaVersion == 1 &&
		receipt.ConformID == record.ID &&
		receipt.AnalysisID == record.AnalysisID &&
		receipt.TakeID == record.TakeID &&
		receipt.MeasuredOffsetMS == record.MeasuredOffsetMS &&
		receipt.AppliedOffsetMS == record.AppliedOffsetMS &&
		receipt.AppliedSamples == record.AppliedSamples &&
		receipt.AudioDigest == record.AudioOutputDigest &&
		receipt.VideoDigest == record.VideoOutputDigest
}
