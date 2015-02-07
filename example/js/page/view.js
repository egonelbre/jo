import (
	"page/items"
)

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

export View;