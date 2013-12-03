/*
 * Copyright 2013 gtalent2@gmail.com
 * 
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * 
 *   http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
#include <sstream>
#include <iostream>
#include <fstream>
#ifdef CYBORGBEAR_BOOST_ENABLED
#include <boost/archive/text_oarchive.hpp>
#include <boost/archive/text_iarchive.hpp>
#endif
#include "Model.hpp"

using namespace std;
using namespace models;

void testJson(Model1 &orig) {
	Model1 copy;

	copy.fromJson(orig.toJson());

	cout << "JSON Test:  " << (copy.toJson().compare(orig.toJson()) == 0 ? "Pass" : "Fail") << endl;
}

void testBoost(Model1 &orig) {
#ifdef CYBORGBEAR_BOOST_ENABLED
	Model1 copy;

	copy.fromBoostBinary(orig.toBoostBinary());

	cout << "Boost Test: " << (copy.toJson().compare(orig.toJson()) == 0 ? "Pass" : "Fail") << endl;
#endif
}

int main() {
	Model1 mod;
	mod.fromJson("{\"Field1\": \"Ni!\", \"Field2\": \"asdf\", \"Field3\": [4, 2], \"Field4\": [[\"Narf!\", \"Narf!\"], [\"Narf!\", \"Narf!\"]], \"Field5\": {\"Narf\": \"Ni\"}}");
	mod.Field1 = "Narf!";
	cout << mod.toJson(cyborgbear::Readable) << endl;

	testJson(mod);
	testBoost(mod);
	return 0;
}
