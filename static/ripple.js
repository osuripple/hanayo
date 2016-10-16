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
    page = page === 0 ? 1 : page;
    
    function loadLeaderboard() {
      var wl = window.location;      
      window.history.replaceState('', document.title, wl.pathname + "?mode=" + favouriteMode + "&p=" + page + wl.hash);
      api("leaderboard", {
        mode: favouriteMode,
        p: page,
        l: 50,
      }, function(data) {
        var tb = $(".ui.table tbody");
        tb.find("tr").remove();
        if (!data.users)
          disableSimplepagButtons(true);
        var i = 0;
        data.users.forEach(function(v) {
          tb.append(
            $("<tr />").append(
              $("<td />").text("#" + ((page-1) * 50 + (++i))),
              $("<td />").html("<a href='/u/" + v.id + "' title='View profile'><i class='" + 
                v.country.toLowerCase() + " flag'></i>" + escapeHTML(v.username) + "</a>"),
              $("<td />").html(scoreOrPP(v.chosen_mode.ranked_score, v.chosen_mode.pp)),
              $("<td />").text(v.chosen_mode.accuracy.toFixed(2) + "%"),
              // bonus points if you get the undertale joke
              $("<td />").html(addCommas(v.chosen_mode.playcount) +
                " <i title='Why, LOVE, of course!'>(lv. " + v.chosen_mode.level.toFixed(0) + ")</i>")
            )
          );
        });
        disableSimplepagButtons(data.users.length < 50);
      });
    }
    function scoreOrPP(s, pp) {
      if (pp === 0)
        return "<b>" + addCommas(s) + "</b>";
      return "<b>" + addCommas(pp) + "pp</b> (" + addCommas(s) + ")"
    }

    loadLeaderboard();
    setupSimplepag(loadLeaderboard);
    $("#mode-menu .item").click(function(e) {
      e.preventDefault();
      $("#mode-menu .active.item").removeClass("active");
      $(this).addClass("active");
      favouriteMode = $(this).data("mode");
      page = 1;
      loadLeaderboard();
    });
  },
  "/friends": function() {
    $(".smalltext.button").click(function() {
      var t = $(this);
      var delAdd = t.data("deleted") === "1" ? "add" : "del";
      console.log(delAdd);
      t.addClass("disabled");
      api("friends/" + delAdd, {
        id: t.data("userid")
      }, function(data) {
        t.removeClass("disabled");
        t.data("deleted", data.friend ? "0" : "1");
        t.removeClass("green red blue");
        t.addClass(data.friend ? (data.mutual ? "red" : "green") : "blue");
        t.find(".icon").removeClass("minus plus heart").
          addClass(data.friend ? (data.mutual ? "heart" : "minus") : "plus");
        t.find("span").text(data.friend ? (data.mutual ? "Mutual" : "Remove") : "Add");
      });
    });
  },
};

$(document).ready(function(){
  // semantic stuff
  $('.message .close').on('click', closeClosestMessage);
  $('.ui.checkbox').checkbox();
  $('.ui.dropdown').dropdown();
  $('.ui.progress').progress();
  $('.ui.form').submit(function() {
    $(this).addClass("loading");
  });

  // emojis!
  if (typeof twemoji !== "undefined") {
    $(".twemoji").each(function(k, v) {
      twemoji.parse(v);
    });
  }
  
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

function closeClosestMessage() {
  $(this)
    .closest('.message')
    .transition('fade');
};

 function showMessage(type, message) {
  var newEl = $('<div class="ui ' + type + ' message hidden"><i class="close icon"></i>' + message + '</div>');
  newEl.find(".close.icon").click(closeClosestMessage);
  $("#messages-container").append(newEl);
  newEl.transition("slide down");
};

// function for all api calls
 function api(endpoint, data, success) {
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

function setupSimplepag(callback) {
  var el = $(".simplepag");
  el.find(".left.floated .item").click(function() {
    if ($(this).hasClass("disabled"))
      return false;
    page--;
    callback();
  });
  el.find(".right.floated .item").click(function() {
    if ($(this).hasClass("disabled"))
      return false;
    page++;
    callback();
  });
}
function disableSimplepagButtons(right) {
  var el = $(".simplepag");
  
  if (page <= 1)
    el.find(".left.floated .item").addClass("disabled");
  else
    el.find(".left.floated .item").removeClass("disabled");

  if (right)
    el.find(".right.floated .item").addClass("disabled");
  else
    el.find(".right.floated .item").removeClass("disabled");
}

// thank mr stackoverflow
function addCommas(nStr) {
	nStr += '';
	x = nStr.split('.');
	x1 = x[0];
	x2 = x.length > 1 ? '.' + x[1] : '';
	var rgx = /(\d+)(\d{3})/;
	while (rgx.test(x1)) {
		x1 = x1.replace(rgx, '$1' + ',' + '$2');
	}
	return x1 + x2;
}
