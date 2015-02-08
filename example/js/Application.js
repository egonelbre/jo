import (
	"page"
	"page/items/basic"
)

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