package trieutil

import (
	"errors"

	"github.com/prysmaticlabs/prysm/shared/hashutil"
)

// MerkleTrie implements a sparse, general purpose Merkle trie to be used
// across ETH2.0 Phase 0 functionality.
type MerkleTrie struct {
	branches [][][32]byte
}

// GenerateTrieFromItems constructs a Merkle trie from a sequence of byte slices.
func GenerateTrieFromItems(items [][]byte, depth int) (*MerkleTrie, error) {
	if len(items) == 0 {
		return nil, errors.New("no items provided to generate Merkle trie")
	}
	leaves := make([][32]byte, len(items))
	emptyNodes := generateEmptyNodes(depth)
	// We then construct the leaves of the trie by hashing every
	// value in the items slice.
	for i, val := range items {
		leaves[i] = hashutil.Hash(val)
	}
	// Append the leaves to the branches.
	branches := [][][32]byte{leaves}
	for i := 0; i < depth-1; i++ {
		if len(branches[i])%2 == 1 {
			branches[i] = append(branches[i], emptyNodes[i])
		}
		// We append the layer that results from hashing the trie's current layer.
		branches = append(branches, hashLayer(branches[i]))
	}
	// Reverse the branches so as to have the root in the 0th layer.
	for i, j := 0, len(branches)-1; i < j; i, j = i+1, j-1 {
		branches[i], branches[j] = branches[j], branches[i]
	}
	return &MerkleTrie{branches}, nil
}

// VerifyMerkleProof verifies a Merkle branch against a root of a trie.
func (m *MerkleTrie) VerifyMerkleProof(item []byte, merkleIndex int, proof [][32]byte) bool {
	node := hashutil.Hash(item)
	branchIndices := BranchIndices(merkleIndex, len(proof))
	for i := 0; i < len(proof); i++ {
		if branchIndices[i]%2 == 0 {
			node = parentHash(node, proof[i])
		} else {
			node = parentHash(proof[i], node)
		}
	}
	return m.Root() == node
}

// BranchIndices returns the indices of all ancestors for a node with up to the root
// given the node's index by utilizing the depth of the trie.
func BranchIndices(merkleIndex int, depth int) []int {
	indices := make([]int, depth)
	idx := merkleIndex
	indices[0] = idx
	for i := 1; i < depth; i++ {
		idx /= 2
		indices[i] = idx
	}
	return indices
}

// Root of the Merkle trie.
func (m *MerkleTrie) Root() [32]byte {
	return m.branches[0][0]
}

// MerkleProof obtains a Merkle proof for an item at a given
// index in the Merkle trie up to the root of the trie.
func (m *MerkleTrie) MerkleProof(merkleIndex int) ([][32]byte, error) {
	lastLevel := m.branches[len(m.branches)-1]
	if merkleIndex < 0 || merkleIndex >= len(lastLevel) || lastLevel[merkleIndex] == [32]byte{} {
		return nil, errors.New("merkle index out of range in trie")
	}
	branchIndices := BranchIndices(merkleIndex, len(m.branches))
	// We create a list of proof indices, which do not include the root so the length
	// of our proof will be the length of the branch indices - 1.
	proofIndices := make([]int, len(branchIndices)-1)
	for i := 0; i < len(proofIndices); i++ {
		// We fetch the sibling by flipping the rightmost bit.
		proofIndices[i] = branchIndices[i] ^ 1
	}
	proof := make([][32]byte, len(proofIndices))
	for j := 0; j < len(proofIndices); j++ {
		// We fetch the layer that corresponds to the proof element index
		// in our Merkle trie's branches. Since the length of proof indices
		// is the len(tree)-1, this will ignore the root.
		layer := m.branches[len(m.branches)-1-j]
		proof[j] = layer[proofIndices[j]]
	}
	return proof, nil
}

// parentHash takes a left and right node and hashes their concatenation.
func parentHash(left [32]byte, right [32]byte) [32]byte {
	res := hashutil.Hash(append(left[:], right[:]...))
	return res
}

// hashLayer computes the layer on top of another one by hashing left and right
// nodes to compute the nodes in the trie above.
func hashLayer(layer [][32]byte) [][32]byte {
	chunks := partition(layer)
	topLayer := [][32]byte{}
	for i := 0; i < len(chunks); i++ {
		topLayer = append(topLayer, parentHash(chunks[i][0], chunks[i][1]))
	}
	return topLayer
}

// generateEmptyNodes creates a trie of empty nodes up a path given a trie depth.
// This is necessary given the Merkle trie is a balanced trie and empty nodes serve
// as padding along the way if an odd number of leaves are originally provided.
func generateEmptyNodes(depth int) [][32]byte {
	nodes := make([][32]byte, depth)
	for i := 0; i < depth; i++ {
		nodes[i] = parentHash([32]byte{}, [32]byte{})
	}
	return nodes
}

// partition a slice into chunks of size two.
// Example: [1, 2, 3, 4] -> [[1, 2], [3, 4]]
func partition(layer [][32]byte) [][][32]byte {
	chunks := [][][32]byte{}
	size := 2
	for i := 0; i < len(layer); i += size {
		chunks = append(chunks, layer[i:i+size])
	}
	return chunks
}