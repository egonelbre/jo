/* page/items */
var ノpageノitems = {};
(function(Ɛ){
	
	/* page/items/registry.js */
	var Views = {};
	
	function Register(type, component){
		Views[type] = component;
	}
	
	Ɛ.Register = Register;
	
	function ComponentFor(item){
		return Views[item.type];
	}
	
	Ɛ.ComponentFor = ComponentFor;
})(ノpageノitems);

/* page */
var ノpage = {};
(function(Ɛ){
	
	/* page/page.js */
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
	
	Ɛ.Page = Page;
	/* page/view.js */
	var items = ノpageノitems;
	
	
	var View = React.createClass({
		render: function(){
			var page = this.props.page;
			
			var components = page.story.map(function(item, i){
				var element = items.ComponentFor(item);
				return React.createElement(element, {
					key: i,
					item: item
				});
			});
			
			return React.DOM.section({}, 
				React.DOM.h1({}, page.title),
				components
			);
		}
	});
	
	Ɛ.View = View;
})(ノpage);

/* page/items/basic */
var ノpageノitemsノbasic = {};
(function(Ɛ){
	
	/* page/items/basic/paragraph.js */
	var items = ノpageノitems;
	
	
	var Paragraph = React.createClass({
		render: function(){
			var item = this.props.item;
			return React.DOM.p({}, item.text);
		}
	});
	
	items.Register("paragraph", Paragraph);
	
	Ɛ.Paragraph = Paragraph;
})(ノpageノitemsノbasic);

/*  */
var ノ = {};
(function(Ɛ){
	
	/* Application.js */
	var page = ノpage;
	var basic = ノpageノitemsノbasic;
	
	
	var MainPage = new page.Page();
	
	var Application = React.createClass({
		render: function(){
			var MainPage = this.props.page;
			return React.createElement(page.View, {page: MainPage});
		}
	});
	
	React.render(
		React.createElement(Application, {page: MainPage}),
		document.getElementById("application")
	);
})(ノ);

