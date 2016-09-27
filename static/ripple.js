// Ripple custom JS that goes on all pages

// this object contains tiny snippets that were deemed too small to be worth
// their own file.
var singlePageSnippets = {

  "/": function() {
    $(".expand-icon").popup().click(function() {
      var addTo = $(this).closest(".segment");
      if (addTo.attr("data-expanded") == "true") {
        addTo.removeAttr("data-expanded")
        var ch = addTo.children(".post-content");
        ch.slideUp(function() {
          ch.remove();
        });
        $(this).attr("data-content", "Expand")
        $(this).removeClass("up").addClass("down");
      } else {
        addTo.attr("data-expanded", "true");
        $(this).removeClass("down").addClass("up");
        $(this).attr("data-content", "Collapse");
        api("blog/posts/content", {
          id: addTo.data("post-id"),
          html: "",
        }, function(data) {
          var el = $("<div class='post-content' />").hide().append(data.content);
          addTo.append(el);
          el.slideDown();
        });
      }
    });
  },

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

  "/2fa_gateway": function() {
    $('#telegram-code').on('input', function() {
      if ($(this).val().length >= 8) {
        $.get("/2fa_gateway/verify", {
          token: $(this).val().trim().substr(0, 8),
        }, function(resp) {
          switch (resp) {
          case "0":
            $("#telegram-code").closest(".field").addClass("success");
            window.location.href = "/";
            break;
          case "1":
            $("#telegram-code").closest(".field").addClass("error");
            break;
          }
        });
      } else {
        $("#telegram-code").closest(".field").removeClass("error");
      }
    });
  },

  "/leaderboard": function() {
    var wl = window.location;
    if (window.location.search.indexOf("mode=") === -1)
      window.history.pushState('', document.title, wl.pathname + "?mode=" + favouriteMode + wl.hash);
  },
};

$(document).ready(function(){
  // semantic stuff
  $('.message .close').on('click', closeClosestMessage);

  $('.ui.checkbox').checkbox();
  
  $('.ui.form').submit(function() {
    $(this).addClass("loading");
  });
  
  $('.ui.dropdown').dropdown();
  
  $('.ui.progress').progress();
  
  // ripple stuff
  var f = singlePageSnippets[window.location.pathname];
  if (typeof f === 'function')
    f();

  // setup user search
  $("#user-search").search({
    onSelect: function(val) {
      window.location.href = val.url;
      return false;
    },
    apiSettings: {
      url: "/api/v1/users/lookup?name={query}",
      onResponse: function(resp) {
        var r = {
          results: [],
        };
        $.each(resp.users, function(index, item) {
          r.results.push({
            title: item.username,
            url  : "/u/" + item.id,
            image: hanayoConf.avatars + "/" + item.id,
          });
        });
        return r;
      },
    },
  });
  $("#user-search-input").keypress(function (e) {
    if (e.which == 13) {
      window.location.pathname = "/u/" + $(this).val();
    }
  });
  
  // setup timeago
  $.timeago.settings.allowFuture = true;
  $("time.timeago").timeago();
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

var modes = {
  0: "osu! standard",
  1: "Taiko",
  2: "Catch the Beat",
  3: "osu!mania",
};
var modesShort = {
  0: "std",
  1: "taiko",
  2: "ctb",
  3: "mania",
};

var entityMap = {
  "&": "&amp;",
  "<": "&lt;",
  ">": "&gt;",
  '"': '&quot;',
  "'": '&#39;',
  "/": '&#x2F;',
};
function escapeHTML(str) {
  return String(str).replace(/[&<>"'\/]/g, function(s) {
    return entityMap[s];
  });
}
