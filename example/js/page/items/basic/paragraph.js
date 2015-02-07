import (
	"page/items"
)

var Paragraph = React.createClass({
	render: function(){
		var item = this.props.item;
		return React.DOM.p({}, item.text);
	}
});

items.Register("paragraph", Paragraph);

export Paragraph;