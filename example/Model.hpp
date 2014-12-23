//Generated Code

#ifndef MODEL_HPP
#define MODEL_HPP

#include "../../json_read.hpp"
#include "../../json_write.hpp"



namespace wombat {
namespace models {

class Model1 {

	public:

		Model1();

		bool operator==(const Model1&) const;

		bool operator!=(const Model1&) const;
		string Field1;
		string Field2;
		std::vector< int > Field3;
		std::vector< std::vector< string > > Field4;
		std::map< string, string > Field5;
};

inline Error toJson(Model1 model, json_t *jo) {
	Error err = Error::Ok;
	err |= writeVal(jo, "Field1", model.Field1);
	err |= writeVal(jo, "Field2", model.Field2);
	err |= writeVal(jo, "Field3", model.Field3);
	err |= writeVal(jo, "Field4", model.Field4);
	err |= writeVal(jo, "Field5", model.Field5);
	return err;
}

inline Error fromJson(Model1 *model, json_t *jo) {
	Error err = Error::Ok;
	err |= readVal(jo, "Field1", &model->Field1);
	err |= readVal(jo, "Field2", &model->Field2);
	err |= readVal(jo, "Field3", &model->Field3);
	err |= readVal(jo, "Field4", &model->Field4);
	err |= readVal(jo, "Field5", &model->Field5);
	return err;
}

}
}


#endif