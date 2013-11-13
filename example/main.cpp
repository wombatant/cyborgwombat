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
#ifdef CYBORGBEAR_BOOST_ENABLED
#include <boost/archive/text_oarchive.hpp>
#include <boost/archive/text_iarchive.hpp>
#endif
#include "Model.hpp"

using namespace std;
using namespace models;

void testBoost(string json) {
#ifdef CYBORGBEAR_BOOST_ENABLED
	Model1 mod1;
	Model1 mod2;

	mod1.fromJson(json);

	stringstream out;
	{
		boost::archive::text_oarchive oa(out);
		oa << mod1;
	}
	while (out.good())
		cout << out.get();
	cout << endl;
	//boost::archive::text_iarchive ia(out);
	//ia >> mod2;
	//cout << "Boost Test: " << (mod2.toJson().compare(json) == 0 ? "Pass" : "Fail") << endl;
#endif
}

int main() {
	Model1 mod;
	mod.fromJson("{\"Field1\": \"Ni!\", \"Field2\": 1, \"Field3\": [4, 2], \"Field4\": [[\"Narf!\", \"Narf!\"], [\"Narf!\", \"Narf!\"]], \"Field5\": {\"Narf\": \"Ni\"}}");
	mod.Field1 = "Narf!";
	cout << mod.toJson(cyborgbear::Readable) << endl;

	testBoost(mod.toJson());
	return 0;
}
