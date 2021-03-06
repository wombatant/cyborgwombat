/*
 * Copyright 2013-2014 gtalent2@gmail.com
 * 
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */
//Generated Code
#include "Model.hpp"

namespace wombat {
namespace models {
Model1::Model1() {
	this->Field1 = "";
	this->Field2 = "";
}

bool Model1::operator==(const Model1 &o) const {
	if (Field1 != o.Field1) return false;
	if (Field2 != o.Field2) return false;
	if (Field3 != o.Field3) return false;
	if (Field4 != o.Field4) return false;
	if (Field5 != o.Field5) return false;

	return true;
}

bool Model1::operator!=(const Model1 &o) const {
	if (Field1 != o.Field1) return true;
	if (Field2 != o.Field2) return true;
	if (Field3 != o.Field3) return true;
	if (Field4 != o.Field4) return true;
	if (Field5 != o.Field5) return true;

	return false;
}

}
}

