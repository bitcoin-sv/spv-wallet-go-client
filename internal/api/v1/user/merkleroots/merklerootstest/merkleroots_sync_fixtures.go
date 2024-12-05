package merklerootstest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/jarcoal/httpmock"
)

// DB simulates a storage of Merkle roots for testing.
type DB struct {
	MerkleRoots []models.MerkleRoot
}

// SaveMerkleRoots appends synced Merkle roots to the simulated storage.
func (db *DB) SaveMerkleRoots(syncedMerkleRoots []models.MerkleRoot) error {
	db.MerkleRoots = append(db.MerkleRoots, syncedMerkleRoots...)
	return nil
}

// GetLastMerkleRoot retrieves the last Merkle root from storage.
func (db *DB) GetLastMerkleRoot() string {
	if len(db.MerkleRoots) == 0 {
		return ""
	}
	return db.MerkleRoots[len(db.MerkleRoots)-1].MerkleRoot
}

// CreateRepository initializes a simulated repository with the provided Merkle roots.
func CreateRepository(merkleRoots []models.MerkleRoot) *DB {
	return &DB{
		MerkleRoots: merkleRoots,
	}
}

// sendJSONResponse sends a JSON response from the mock server.
func sendJSONResponse(data any, w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(*w).Encode(data); err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
	}
}

// MockMerkleRootsAPIResponseNormal creates a mock server with normal API responses.
func MockMerkleRootsAPIResponseNormal() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/merkleroots" && r.Method == http.MethodGet {
			lastEvaluatedKey := r.URL.Query().Get("lastEvaluatedKey")
			sendJSONResponse(MockedMerkleRootsAPIResponseFn(lastEvaluatedKey), &w)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	return server
}

// MockMerkleRootsAPIResponseDelayed creates a mock server with delayed API responses.
func MockMerkleRootsAPIResponseDelayed() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/merkleroots" && r.Method == http.MethodGet {
			lastEvaluatedKey := r.URL.Query().Get("lastEvaluatedKey")
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
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	return server
}

// MockMerkleRootsAPIResponseStale creates a mock server with a stale last evaluated key.
func MockMerkleRootsAPIResponseStale() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/merkleroots" && r.Method == http.MethodGet {
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
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	return server
}

// MockedSPVWalletData is mocked  merkle roots data on spv-wallet side
var MockedSPVWalletData = []models.MerkleRoot{
	{
		BlockHeight: 0,
		MerkleRoot:  "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b",
	},
	{
		BlockHeight: 1,
		MerkleRoot:  "0e3e2357e806b6cdb1f70b54c3a3a17b6714ee1f0e68bebb44a74b1efd512098",
	},
	{
		BlockHeight: 2,
		MerkleRoot:  "9b0fc92260312ce44e74ef369f5c66bbb85848f2eddd5a7a1cde251e54ccfdd5",
	},
	{
		BlockHeight: 3,
		MerkleRoot:  "999e1c837c76a1b7fbb7e57baf87b309960f5ffefbf2a9b95dd890602272f644",
	},
	{
		BlockHeight: 4,
		MerkleRoot:  "df2b060fa2e5e9c8ed5eaf6a45c13753ec8c63282b2688322eba40cd98ea067a",
	},
	{
		BlockHeight: 5,
		MerkleRoot:  "63522845d294ee9b0188ae5cac91bf389a0c3723f084ca1025e7d9cdfe481ce1",
	},
	{
		BlockHeight: 6,
		MerkleRoot:  "20251a76e64e920e58291a30d4b212939aae976baca40e70818ceaa596fb9d37",
	},
	{
		BlockHeight: 7,
		MerkleRoot:  "8aa673bc752f2851fd645d6a0a92917e967083007d9c1684f9423b100540673f",
	},
	{
		BlockHeight: 8,
		MerkleRoot:  "a6f7f1c0dad0f2eb6b13c4f33de664b1b0e9f22efad5994a6d5b6086d85e85e3",
	},
	{
		BlockHeight: 9,
		MerkleRoot:  "0437cd7f8525ceed2324359c2d0ba26006d92d856a9c20fa0241106ee5a597c9",
	},
	{
		BlockHeight: 10,
		MerkleRoot:  "d3ad39fa52a89997ac7381c95eeffeaf40b66af7a57e9eba144be0a175a12b11",
	},
	{
		BlockHeight: 11,
		MerkleRoot:  "f8325d8f7fa5d658ea143629288d0530d2710dc9193ddc067439de803c37066e",
	},
	{
		BlockHeight: 12,
		MerkleRoot:  "3b96bb7e197ef276b85131afd4a09c059cc368133a26ca04ebffb0ab4f75c8b8",
	},
	{
		BlockHeight: 13,
		MerkleRoot:  "9962d5c704ec27243364cbe9d384808feeac1c15c35ac790dffd1e929829b271",
	},
	{
		BlockHeight: 14,
		MerkleRoot:  "e1afd89295b68bc5247fe0ca2885dd4b8818d7ce430faa615067d7bab8640156",
	},
}

// LastMockedMerkleRoot returns last merkleroot value from MockedSPVWalletData
func LastMockedMerkleRoot() models.MerkleRoot {
	return MockedSPVWalletData[len(MockedSPVWalletData)-1]
}

func MockedMerkleRootsAPIResponseFn(lastMerkleRoot string) models.ExclusiveStartKeyPage[[]models.MerkleRoot] {
	// If no lastMerkleRoot is provided, return the full dataset
	if lastMerkleRoot == "" {
		return models.ExclusiveStartKeyPage[[]models.MerkleRoot]{
			Content: MockedSPVWalletData,
			Page: models.ExclusiveStartKeyPageInfo{
				LastEvaluatedKey: MockedSPVWalletData[len(MockedSPVWalletData)-1].MerkleRoot, // Last Merkle root as key
				TotalElements:    len(MockedSPVWalletData),
				Size:             len(MockedSPVWalletData),
			},
		}
	}

	// Find the index of the lastMerkleRoot
	lastMerkleRootIdx := slices.IndexFunc(MockedSPVWalletData, func(mr models.MerkleRoot) bool {
		return mr.MerkleRoot == lastMerkleRoot
	})

	// If lastMerkleRoot is not found, return an empty response (or handle as error if desired)
	if lastMerkleRootIdx == -1 {
		return models.ExclusiveStartKeyPage[[]models.MerkleRoot]{
			Content: []models.MerkleRoot{},
			Page: models.ExclusiveStartKeyPageInfo{
				LastEvaluatedKey: "",
				TotalElements:    len(MockedSPVWalletData),
				Size:             0,
			},
		}
	}

	// If lastMerkleRoot is the highest in the server database, return no new content
	if lastMerkleRootIdx >= len(MockedSPVWalletData)-1 {
		return models.ExclusiveStartKeyPage[[]models.MerkleRoot]{
			Content: []models.MerkleRoot{},
			Page: models.ExclusiveStartKeyPageInfo{
				LastEvaluatedKey: "",
				TotalElements:    len(MockedSPVWalletData),
				Size:             0,
			},
		}
	}

	// Return all Merkle roots after the given lastMerkleRoot
	content := MockedSPVWalletData[lastMerkleRootIdx+1:]

	// Set the LastEvaluatedKey to the last Merkle root in the current page, or "" if it's the final one
	lastEvaluatedKey := ""
	if len(content) > 0 && content[len(content)-1].MerkleRoot != MockedSPVWalletData[len(MockedSPVWalletData)-1].MerkleRoot {
		lastEvaluatedKey = content[len(content)-1].MerkleRoot
	}

	return models.ExclusiveStartKeyPage[[]models.MerkleRoot]{
		Content: content,
		Page: models.ExclusiveStartKeyPageInfo{
			LastEvaluatedKey: lastEvaluatedKey,
			TotalElements:    len(MockedSPVWalletData),
			Size:             len(content),
		},
	}
}

func FirstMerkleRootsPage() *queries.MerkleRootPage {
	return &queries.MerkleRootPage{
		Content: []models.MerkleRoot{
			{
				BlockHeight: 0,
				MerkleRoot:  "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b",
			},
			{
				BlockHeight: 1,
				MerkleRoot:  "0e3e2357e806b6cdb1f70b54c3a3a17b6714ee1f0e68bebb44a74b1efd512098",
			},
			{
				BlockHeight: 2,
				MerkleRoot:  "9b0fc92260312ce44e74ef369f5c66bbb85848f2eddd5a7a1cde251e54ccfdd5",
			},
		},
		Page: models.ExclusiveStartKeyPageInfo{
			OrderByField:     spvwallettest.Ptr("blockHeight"),
			SortDirection:    spvwallettest.Ptr("asc"),
			TotalElements:    9,
			Size:             3,
			LastEvaluatedKey: "e4774f7a-eb99-4cac-956e-634d2aeccc93",
		},
	}
}

func SecondMerkleRootsPage() *queries.MerkleRootPage {
	return &queries.MerkleRootPage{
		Content: []models.MerkleRoot{
			{
				BlockHeight: 3,
				MerkleRoot:  "999e1c837c76a1b7fbb7e57baf87b309960f5ffefbf2a9b95dd890602272f644",
			},
			{
				BlockHeight: 4,
				MerkleRoot:  "df2b060fa2e5e9c8ed5eaf6a45c13753ec8c63282b2688322eba40cd98ea067a",
			},
			{
				BlockHeight: 5,
				MerkleRoot:  "63522845d294ee9b0188ae5cac91bf389a0c3723f084ca1025e7d9cdfe481ce1",
			},
		},
		Page: models.ExclusiveStartKeyPageInfo{
			OrderByField:     spvwallettest.Ptr("blockHeight"),
			SortDirection:    spvwallettest.Ptr("asc"),
			TotalElements:    9,
			Size:             3,
			LastEvaluatedKey: "6bad63f5-8f2e-4756-aca9-cc9cb4a001c6",
		},
	}
}

func ThirdMerkleRootsPage() *queries.MerkleRootPage {
	return &queries.MerkleRootPage{
		Content: []models.MerkleRoot{
			{
				BlockHeight: 6,
				MerkleRoot:  "20251a76e64e920e58291a30d4b212939aae976baca40e70818ceaa596fb9d37",
			},
			{
				BlockHeight: 7,
				MerkleRoot:  "8aa673bc752f2851fd645d6a0a92917e967083007d9c1684f9423b100540673f",
			},
			{
				BlockHeight: 8,
				MerkleRoot:  "a6f7f1c0dad0f2eb6b13c4f33de664b1b0e9f22efad5994a6d5b6086d85e85e3",
			},
		},
		Page: models.ExclusiveStartKeyPageInfo{
			OrderByField:     spvwallettest.Ptr("blockHeight"),
			SortDirection:    spvwallettest.Ptr("asc"),
			TotalElements:    9,
			Size:             3,
			LastEvaluatedKey: "09232c7e-ecf7-4e33-8feb-a32170c6e7b6",
		},
	}
}

func ResponderWithThreeMerkleRootPagesSuccess(t *testing.T) httpmock.Responder {
	pages := map[int]*queries.MerkleRootPage{
		0: FirstMerkleRootsPage(),
		1: SecondMerkleRootsPage(),
		2: ThirdMerkleRootsPage(),
	}

	var num int
	return func(r *http.Request) (*http.Response, error) {
		defer func() { num++ }()

		if num < len(pages) {
			res, err := httpmock.NewJsonResponse(http.StatusPartialContent, pages[num])
			if err != nil {
				t.Fatalf("test helper - failed to generate new json response: %s", err)
			}
			return res, nil
		}

		res, err := httpmock.NewJsonResponse(http.StatusOK, queries.MerkleRootPage{})
		if err != nil {
			t.Fatalf("test helper - failed to generate new json response: %s", err)
		}
		return res, nil
	}
}
