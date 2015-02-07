
var Views = {};

function Register(type, component){
	Views[type] = component;
}

export Register;

function ComponentFor(item){
	return Views[item.type];
}

export ComponentFor;