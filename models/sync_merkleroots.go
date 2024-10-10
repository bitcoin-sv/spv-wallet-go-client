package models

// ExclusiveStartKeyPage represents a paginated response for database records using Exclusive Start Key paging
type ExclusiveStartKeyPage[T any] struct {
	// List of records for the response
	Content T
	// Pagination details
	Page ExclusiveStartKeyPageInfo
}

// ExclusiveStartKeyPageInfo represents the pagination information for limiting and sorting database query results
type ExclusiveStartKeyPageInfo struct {
	// Field by which to order the results
	OrderByField *string `json:"orderByField,omitempty"` // Optional ordering field
	// Direction in which to order the results (ASC or DESC)
	SortDirection *string `json:"sortDirection,omitempty"` // Optional sort direction
	// Total count of elements
	TotalElements int `json:"totalElements"`
	// Size of the page or returned data
	Size int `json:"size"`
	// Last evaluated key returned from the database
	LastEvaluatedKey string `json:"lastEvaluatedKey"`
}

// MerkleRoot holds the content of the synced Merkle root response
type MerkleRoot struct {
	MerkleRoot  string `json:"merkleRoot"`
	BlockHeight int    `json:"blockHeight"`
}

// MerkleRootsRepository is an interface responsible for saving synced merkleroots and getting last evaluat key from database
type MerkleRootsRepository interface {
	// GetLastMerkleRoot should return the merkle root with the heighest height from your storage or undefined if empty
	GetLastMerkleRoot() string
	// SaveMerkleRoots should store newly synced merkle roots into your storage;
	// NOTE: items are ordered with ascending order by block height
	SaveMerkleRoots(syncedMerkleRoots []MerkleRoot) error
}
