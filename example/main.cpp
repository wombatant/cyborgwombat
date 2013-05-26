#include <iostream>
#include "Model.hpp"

using namespace std;
using namespace models;

int main() {
	Model1 mod;
	mod.load("{\"Field1\": \"Ni!\", Field2: 5, Field3: 42, Field4: [[\"Narf!\", \"Narf!\"]]}");
	cout << mod.write() << endl;
	return 0;
}
