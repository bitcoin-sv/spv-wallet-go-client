package errors

import "errors"

var (
	// ErrMissingXpriv is returned when the xpriv is missing.
	ErrMissingXpriv = errors.New("xpriv is missing")
	// ErrContactPubKeyInvalid is returned when the contact's PubKey is invalid.
	ErrContactPubKeyInvalid = errors.New("contact's PubKey is invalid")
	// ErrMetadataFilterMaxDepthExceeded is returned when the maximum depth of nesting in metadata map is exceeded.
	ErrMetadataFilterMaxDepthExceeded = errors.New("maximum depth of nesting in metadata map exceeded")
	// ErrMetadataWrongTypeInArray is returned when the wrong type is in the array.
	ErrMetadataWrongTypeInArray = errors.New("wrong type in array")
	// ErrFilterQueryBuilder is returned when the filter query builder fails to build the operation.
	ErrFilterQueryBuilder = errors.New("filter query builder - build operation failure")
	// ErrUnrecognizedAPIResponse indicates that the response received from the SPV Wallet API
	// does not match the expected format or structure.
	ErrUnrecognizedAPIResponse = errors.New("unrecognized response from API")
	// ErrSyncMerkleRootsTimeout is returned when the SyncMerkleRoots operation times out.
	ErrSyncMerkleRootsTimeout = errors.New("SyncMerkleRoots operation timed out")
	// ErrStaleLastEvaluatedKey is returned when the last evaluated key has not changed between requests,
	ErrStaleLastEvaluatedKey = errors.New("the last evaluated key has not changed between requests, indicating a possible loop or synchronization issue.")

	// ErrFailedToFetchMerkleRootsFromAPI is returned when the API fails to fetch merkle roots.
	ErrFailedToFetchMerkleRootsFromAPI = errors.New("failed to fetch merkle roots from API")
)
