#include "merkle_tree.hh"

// Returns true if all of the requests in the block are valid.
int ValidateBlock(block *b);

// Compute proof of work on a block
int ComputeProofOfWork(block *b);

// Broadcasts a block to nodes once proof of work has been computed.
int BroadcastBlock(block *b);
