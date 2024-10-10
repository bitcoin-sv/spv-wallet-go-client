package fixtures

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/bitcoin-sv/spv-wallet-go-client/models"
)

// simulate a storage of merkle roots that exists on a client side that is using SyncMerkleRoots method
type DB struct {
	MerkleRoots []models.MerkleRoot
}

func (db *DB) SaveMerkleRoots(syncedMerkleRoots []models.MerkleRoot) error {
	db.MerkleRoots = append(db.MerkleRoots, syncedMerkleRoots...)
	time.Sleep(5 * time.Millisecond)
	return nil
}

func (db *DB) GetLastMerkleRoot() string {
	if len(db.MerkleRoots) == 0 {
		return ""
	}
	return db.MerkleRoots[len(db.MerkleRoots)-1].MerkleRoot
}

// CreateRepository creates a simulated repository a client passes to SyncMerkleRoots()
func CreateRepository(merkleRoots []models.MerkleRoot) *DB {
	return &DB{
		MerkleRoots: merkleRoots,
	}
}

func sendJSONResponse(data interface{}, w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(*w).Encode(data); err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
	}
}

func MockMerkleRootsAPIResponseNormal() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/v1/merkleroots" && r.Method == http.MethodGet:
			lastEvaluatedKey := r.URL.Query().Get("lastEvaluatedKey")
			sendJSONResponse(MockedMerkleRootsAPIResponseFn(lastEvaluatedKey), &w)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))

	return server
}

func MockMerkleRootsAPIResponseDelayed() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/v1/merkleroots" && r.Method == http.MethodGet:
			lastEvaluatedKey := r.URL.Query().Get("lastEvaluatedKey")
			// it is to limit the result up to 3 merkle roots per request to ensure
			// that the sync merkleroots will loop more than once and hit the timeout
			all := MockedMerkleRootsAPIResponseFn(lastEvaluatedKey)
			if len(all.Content) > 3 {
				all.Content = all.Content[:3]
			}

			all.Page.Size = len(all.Content)

			if len(all.Content) > 0 {
				all.Page.LastEvaluatedKey = all.Content[len(all.Content)-1].MerkleRoot
			} else {
				all.Page.LastEvaluatedKey = ""
			}

			time.Sleep(50 * time.Millisecond)
			sendJSONResponse(all, &w)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))

	return server
}

func MockMerkleRootsAPIResponseStale() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/v1/merkleroots" && r.Method == http.MethodGet:
			staleLastEvaluatedKeyResponse := models.ExclusiveStartKeyPage[[]models.MerkleRoot]{
				Content: []models.MerkleRoot{
					{
						MerkleRoot:  "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b",
						BlockHeight: 0,
					},
					{
						MerkleRoot:  "0e3e2357e806b6cdb1f70b54c3a3a17b6714ee1f0e68bebb44a74b1efd512098",
						BlockHeight: 1,
					},
					{
						MerkleRoot:  "9b0fc92260312ce44e74ef369f5c66bbb85848f2eddd5a7a1cde251e54ccfdd5",
						BlockHeight: 2,
					},
				},
				Page: models.ExclusiveStartKeyPageInfo{
					LastEvaluatedKey: "9b0fc92260312ce44e74ef369f5c66bbb85848f2eddd5a7a1cde251e54ccfdd5",
					Size:             3,
					TotalElements:    len(MockedSPVWalletData),
				},
			}
			sendJSONResponse(staleLastEvaluatedKeyResponse, &w)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))

	return server
}
