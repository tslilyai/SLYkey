#include <cstring>

typedef struct request {
  int order; // Maybe useful later to optimize lookup
  int type; // Update or Create
  string email;
  string public_key;
  string signature; 
} request;

typedef struct block {
    long public_key;
    long signature;
} block;
