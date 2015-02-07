
function Page(){
	this.title = "Hello";
	this.story = [
		{type: "paragraph", text: "alpha"},
		{type: "paragraph", text: "beta"},
		{type: "paragraph", text: "gamma"},
		{type: "paragraph", text: "delta"}
	];
}
Page.prototype = {};

export Page;