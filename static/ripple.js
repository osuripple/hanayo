// Ripple custom JS that goes on all pages

var singlePageSnippets = {

  "/login": function() {
    $("#login-form").submit(function(e) {
      $("button").addClass("disabled");
      
      var fix = function(errorMessage) {
        $("button").removeClass("disabled");
        $(".ui.form").removeClass("loading");
        showMessage("error", errorMessage);
      };
      
      if (!/^[a-zA-Z0-9 \[\]\@\.\+-]+$/.test($("input[name='username']").val())) {
        fix("Invalid username.");
        return false;
      }
    });
  },

  "/": function() {
    $(".expand-icon").popup().click(function() {
      var addTo = $(this).closest(".raised.segment");
      if (addTo.attr("data-expanded") == "true") {
        addTo.removeAttr("data-expanded")
        addTo.children(".post-content").slideUp();
        $(this).attr("data-content", "Load post inline")
        $(this).removeClass("up").addClass("down");
      } else {
        addTo.attr("data-expanded", "true");
        $(this).removeClass("down").addClass("up");
        $(this).attr("data-content", "Reduce")
        api("blog/posts/content", {
          id: addTo.data("post-id"),
          html: "",
        }, function(data) {
          var el = $("<div class='post-content' />").hide().append(data.content)
          addTo.append(el);
          el.slideDown();
        });
      }
    });
  },

  "/settings/avatar": function() {
    // TODO
    // https://foliotek.github.io/Croppie/demo/demo.js
    $uploadCrop = $('#croppie-thing').croppie({
      enableExif: true,
      viewport: {
        width: 128,
        height: 128
      },
      boundary: {
        width: 300,
        height: 300
      }
    });
  },
};

$(document).ready(function(){
  /* semantic stuff */
  $('.message .close').on('click', closeClosestMessage);

  $('.ui.checkbox').checkbox();
  
  $('.ui.form').submit(function() {
    $(this).addClass("loading");
  });
  
  $('.ui.dropdown').dropdown();
  
  /* ripple stuff */
  var f = singlePageSnippets[window.location.pathname];
  if (typeof f === 'function') {
    f();
  }
});

var closeClosestMessage = function() {
  $(this)
    .closest('.message')
    .transition('fade');
};

var showMessage = function(type, message) {
  var newEl = $('<div class="ui ' + type + ' message hidden"><i class="close icon"></i>' + message + '</div>');
  newEl.find(".close.icon").click(closeClosestMessage);
  $("#messages-container").append(newEl);
  newEl.transition("slide down");
};

// function for all api calls
var api = function(endpoint, data, success) {
  if (typeof data == "function") {
    success = data;
    data = null;
  }
  
  var errorMessage = "An error occurred while contacting the Ripple API. Please report this to a Ripple developer.";

  $.ajax({
    dataType: "json",
    url:      "/api/v1/" + endpoint,
    data:     data,
    success:  function(data) {
      if (data.code != 200) {
        console.warn(data);
        showMessage("error", errorMessage);
      }
      success(data);
    },
    error:    function(jqXHR, textStatus, errorThrown) {
      console.warn(jqXHR, textStatus, errorThrown);
      showMessage("error", errorMessage);
    },
  });
};
