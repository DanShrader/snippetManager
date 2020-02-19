// thanks  https://stackoverflow.com/questions/1144783/how-to-replace-all-occurrences-of-a-string-in-javascript
String.prototype.replaceAll = function (search, replacement) {
	var target = this;
	return target.replace(new RegExp(search, 'g'), replacement);
};

// Thanks
// https://stackoverflow.com/questions/3066586/get-string-in-yyyymmdd-format-from-js-date-object
Date.prototype.yyyymmdd = function () {
	var mm = this.getMonth() + 1; // getMonth() is zero-based
	var dd = this.getDate();

	return [this.getFullYear(),
		(mm > 9 ? '' : '0') + mm,
		(dd > 9 ? '' : '0') + dd
	].join('-');
};

Date.prototype.toInt = function () {
	var mm = this.getMonth() + 1; // getMonth() is zero-based
	var dd = this.getDate();
	var hh = this.getHours();
	var MM = this.getMinutes();
	var ss = this.getSeconds();

	return parseInt([this.getFullYear(),
		(mm > 9 ? '' : '0') + mm,
		(dd > 9 ? '' : '0') + dd,
		(hh > 9 ? '' : '0') + hh,
		(MM > 9 ? '' : '0') + MM,
		(ss > 9 ? '' : '0') + ss
	].join(''));
};

var app = function () {

	var settings = {
		mode: "initial"
	};

	var filesToDelete = [];
	var orginalFiles = {};

	var globalKey = APIkey;

	$.ajaxSetup({
		cache: false
	});

  url = "http://100.115.92.2:10000/api/v1/snippet"

	var urlWithKey = function (url, key) {
		return function () {
			return url;
		};
	};



    // {
    //     "Name": "new thing with a note",
    //     "Description": "",
    //     "Tags": "",
    //     "Notes": "This is a new thing that has an awesome note",
    //     "Reference": "",
    //     "Related": "",
    //     "Files": "",
    //     "Language": "",
    //     "Public": false
    // },

	var model = Backbone.Model.extend({
	  idAttribute: 'Id',
		defaults: {
			Name: "----",
			Description: "",
			Notes: "",
			updated_at: 9999999999999999,
			Language: "",
			Public: false,
			Tags: " "
		}
	});

	var GistCollection = Backbone.Collection.extend({
		model: model,
		url: url
	});


	var gists = new GistCollection();
	gists.fetch()
	window.gists = gists
	

	var ChildView = Marionette.View.extend({
		template: '#item',
		className: "list-group-item",
		tagName: 'li',
		events: {
			"click": "clicked"
		},
		initialize: function () {
			this.listenTo(this.model, 'change', this.render);
		},
		clicked: function () {
			$('.sidebar-nav .active').removeClass('active');
			settings.mode = "view"
			this.$el.addClass('active');
			gist.model = this.model;
			gist.render();
		}
	});

	var MyEmptyGistCollectionView = Marionette.View.extend({
		// template: '#item',
		className: "list-group-item",
		tagName: 'li',
		template: _.template('Loading....')
	});

	var CollectionView = Marionette.CollectionView.extend({
		className: "sidebar-nav list-group",
		tagName: 'ul',
		childView: ChildView,
		collection: gists,
		reorderOnSort: true,
		viewComparator: -'intUpdateDate',
		emptyView: MyEmptyGistCollectionView
	});

	var gistList = new CollectionView();
	gistList.render();



	var SearchView = Marionette.View.extend({
		template: '#search-template',
		tagName: 'form',
		className: 'form-inline',
		ui: {
			searchBTN: '.searchBTN',
			newGist: '.new-gist',
			searchInput: '.searchInput'
		},
		events: {
			"click @ui.searchBTN": "filter",
			"click @ui.newGist": "newGist",
			'keypress @ui.searchInput': 'preventEnter'
		},

		preventEnter: function (e) {
			if (e.which === 13) {
				e.preventDefault();
				this.filter(e)
			}
		},

		filter: function (e) {
			e.preventDefault();
			var searchValue = this.ui.searchInput.val();
			$('.language-wrapper .active').removeClass('active');

			var filter = function (child, index, collection) {
				return (child.get('Description').toLowerCase()).indexOf(searchValue.toLowerCase()) >= 0;
			};
			gistList.setFilter(filter, {
				preventRender: true
			});
			gistList.render();
		},

		newGist: function () {
			console.log('new clicked');

			gist.newView()

		}
	});

	var searchView = new SearchView();
	searchView.render();


// window.tags = tags;

	var detailView = Marionette.View.extend({
		ui: {
			edit: '.edit-gist',
			deleteGists: '.deleteGists',
			save: '.save-gist',
			cancel: '.cancel-gist',
			cancelNew: '.cancel-gist-new',
			add: '.add-gist',
			desc: 'textarea.description'
		},
		events: {
			"click @ui.edit": "editView",
			"click @ui.deleteGists": "deleteGists",
			"click @ui.save": "saveView",
			"click @ui.cancel": "readView",
			"click @ui.cancelNew": "cancelView",
			"click @ui.add": "addGist",
			"click @ui.cancel": "readView",
			"click @ui.cancelNew": "cancelView",
			"click @ui.add": "addGist"
		},

		template: '#details',

		onBeforeRender: function () {
			tmp = '#template-initial-load'
			if (settings.mode === "edit" && typeof (settings.mode) !== "undefined") {
				tmp = '#template-edit-details'
			}
			if (settings.mode === "new" && typeof (settings.mode) !== "undefined") {
				tmp = '#template-new-details'
			}
			if (settings.mode === "loading" && typeof (settings.mode) !== "undefined") {
				tmp = '#template-loading'
			}
			if (settings.mode === "view" && typeof (settings.mode) !== "undefined") {
				tmp = '#details'
			}
			this.template = tmp
		},

		editView: function () {
			settings.mode = "edit";
			this.render();
		},
		newView: function () {
			settings.mode = "new";
			this.render();
			this.addGist();
		},
		deleteGists: function () {
        console.log(this.model)
        console.log(this.model.attributes)
        this.model.destroy();
				this.$el.html('<h1>Deleted the gist</h1>');
				gistList.$el.find('li').first().click();
		},

		cancelView: function () {
			settings.mode = "view";
			gistList.$el.find('li').first().click();
		},
		readView: function () {
			settings.mode = "view";
			this.render();

			$('pre code').each(function (i, block) {
				hljs.highlightBlock(block);
			});

		},

		loadView: function () {
			settings.mode = "loading";
			this.render();
		},

		saveView: function () {
			console.log('save button');
			console.log(this.model);
			
			console.log(this.ui.desc.val());
			var desc = this.ui.desc.val();
			
			if(settings.mode == "new"){
			  this.model = new model()
			  gists.add(this.model)
			}
			
			this.model.set({"Description": desc})
			this.model.set({"Notes": desc})
			this.model.save();
			this.readView();
		},

		addGist: function () {

		},

		onRender: function () {
		// 	this.$el.find('.files').append(files.el);
		}

	});

	var gist = new detailView();
	gist.render()

	$("#list").html(gistList.el);
	$("#content").html(gist.el);
	$(".searchbox").html(searchView.el);

// 	window.gist = gist;
// 	window.gists = gists;
// 	window.files = files;
// 	window.gistList = gistList;
// 	window.tagSummary = tagSummary;
// 	window.tagViewSummary = tagViewSummary;
// 	window.filtersAndTags = filtersAndTags;
}

var APIkey = localStorage.getItem("gistyAPIKey") || "";
var themeColor = localStorage.getItem("gistyTheme") || "";


// Thanks
// https://stackoverflow.com/questions/7846980/how-do-i-switch-my-css-stylesheet-using-jquery
var themeDark = function(){
   console.warn('Dark')
   $('link[href="/static/bootstrap.light.min.css"]').attr('href','bootstrap.dark.min.css');
   $('link[href="/static/light.css"]').attr('href','dark.css');
   $('link[href="/static/https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.12.0/styles/github.min.css"]').attr('href','https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.12.0/styles/railscasts.min.css');
   localStorage.setItem("gistyTheme", "dark");
}
var themeLight = function(){
   console.warn('Light')
   $('link[href="/static/bootstrap.dark.min.css"]').attr('href','bootstrap.light.min.css');
   $('link[href="/static/dark.css"]').attr('href','light.css');
   $('link[href="/static/https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.12.0/styles/railscasts.min.css"]').attr('href','https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.12.0/styles/github.min.css');
   localStorage.setItem("gistyTheme", "light");
   
}


$(document).ready(function() {
  $('#dark').click(function (){
    themeDark()
  });
  
  $('#light').click(function (){
    themeLight()
  });
  
  if (1==2) {
  	// Moved to below
  	$('.sidebar-wrapper, .language-wrapper').hide()
  	$("#page-content-wrapper").css('margin-left','50px');
  } else {
    if(themeColor === "light"){
      themeLight()
    } else {
      localStorage.setItem("gistyTheme", "dark");
    }
    $('.sidebar-wrapper, .language-wrapper').show()
    $("#page-content-wrapper").css('margin-left','550px');
  	app();
  }
  
});