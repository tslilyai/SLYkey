#include <cstring>

using namespace std;

int RegisterPublicKey(string key);
int RevokeAndUpdate(string new_key, string signature);
int GetPublicKey(string email);
