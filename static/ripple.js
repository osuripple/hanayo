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

  "/2fa_gateway": function() {
    $('#telegram-code').on('input', function() {
      if ($(this).val().length >= 8) {
        $.get("/2fa_gateway/verify", {
          token: $(this).val().trim().substr(0, 8),
        }, function(resp) {
          switch (resp) {
          case "0":
            $("#telegram-code").closest(".field").addClass("success");
            redir = redir ? redir : "/"; 
            window.location.href = redir;
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
        user: +t.data("userid")
      }, function(data) {
        t.removeClass("disabled");
        t.data("deleted", data.friend ? "0" : "1");
        t.removeClass("green red blue");
        t.addClass(data.friend ? (data.mutual ? "red" : "green") : "blue");
        t.find(".icon").removeClass("minus plus heart").
          addClass(data.friend ? (data.mutual ? "heart" : "minus") : "plus");
        t.find("span").text(data.friend ? (data.mutual ? "Mutual" : "Remove") : "Add");
      }, true);
    });
  },

  "/team": function() {
    $("#everyone").click(function() {
      $(".ui.modal").modal("show");
    });
  },

  "/register/verify": function() {
    var qu = query("u");
    setInterval(function() {
      $.getJSON(hanayoConf.banchoAPI + "/api/v1/verifiedStatus?u=" + qu, function(data) {
        if (data.result >= 0) {
          window.location.href = "/register/welcome?u=" + qu;
        }
      })
    }, 5000)
  },

  "/settings": function() {
    $("input[name='custom_badge.icon']").on("input", function() {
      $("#badge-icon").attr("class", "circular big icon " + escapeHTML($(this).val()));
    });
    $("input[name='custom_badge.name']").on("input", function() {
      $("#badge-name").html(escapeHTML($(this).val()));
    });
    $("input[name='custom_badge.show']").change(function() {
      if ($(this).is(":checked"))
        $("#custom-badge-fields").slideDown();
      else
        $("#custom-badge-fields").slideUp();        
    });
    $("form").submit(function(e) {
      e.preventDefault();
      var obj = formToObject($(this));
      var ps = 0;
      $(this).find("input[data-sv]").each(function(_, el) {
        el = $(el);
        if (el.is(":checked")) {
          ps |= el.data("sv");
        }
      });
      obj.play_style = ps;
      var f = $(this);
      api("users/self/settings", obj, function(data) {
        showMessage("success", "Your new settings have been saved.");
        f.removeClass("loading");
      }, true);
      return false;
    });
  },

  "/settings/userpage": function() {
    var lastTimeout = null;
    $("textarea[name='data']").on('input', function() {
      if (lastTimeout !== null) {
        clearTimeout(lastTimeout);
      }
      var v = $(this).val();
      lastTimeout = setTimeout(function() {
        $("#userpage-content").addClass("loading");
        $.post("/settings/userpage/parse", $("textarea[name='data']").val(), function(data) {
          var e = $("#userpage-content").removeClass("loading").html(data);
          if (typeof twemoji !== "undefined") {
            twemoji.parse(e[0]);
          }
        }, "text");
      }, 800);
    });
    $("form").submit(function(e) {
      e.preventDefault();
      var obj = formToObject($(this));
      var f = $(this);
      api("users/self/userpage", obj, function(data) {
        showMessage("success", "Your userpage has been saved.");
        f.removeClass("loading");
      }, true);
      return false;
    });
  },

  "/donate": function() {
    var sl = $("#months-slider")[0];
    noUiSlider.create(sl, {
      start: [1],
      step: 1,
      connect: [true, false],
      range: {
        min: [1],
        max: [24],
      }
    });
    var rates = {};
    var us = sl.noUiSlider;
    var doneOne = false;
    $.getJSON("/donate/rates", function(data) {
      rates = data;
      us.on('update', function() {
        var months = us.get();
        var priceEUR = Math.pow(months * 30 * 0.2, 0.70);
        // 3.22 : 700 = x : 1, where x = priceBTC and 700 is the cost of bitcoin
        // (3.22 * 1) / 700 = 3.22 / 700
        var priceBTC = priceEUR / rates.BTC;
        var priceUSD = priceEUR * rates.USD;
        $("#cost").html("<b>" + (+months).toFixed(0) + "</b> month" + (months == 1 ? "" : "s") +
          " = <b>€ " + priceEUR.toFixed(2) + "</b><br><i>($ " + priceUSD.toFixed(2) + " / BTC " + priceBTC.toFixed(6) + ")</i>");
        $("input[name='os0']").attr("value", (+months).toFixed(0) + " month" + (months == 1 ? "" : "s"));
        $("#bitcoin-amt").text(priceBTC.toFixed(6));
      });
    });
  },
  
  "/settings/avatar": function() {
    $("#file").change(function(e) {
      var f = e.target.files;
      if (f.length < 1) {
        return;
      }
      var u = window.URL.createObjectURL(f[0]);
      var i = $("#avatar-img")[0];
      i.src = u;
      i.onload = function() {
        window.URL.revokeObjectURL(this.src);
      };
    });
  },

  "/beatmaps/rank_request": function() {
    function updateRankRequestPage(data) {
      $("#queue-info").html(data.submitted + "/" + data.queue_size +
        (data.submitted_by_user > 0 ? " <i>(" + data.submitted_by_user + "/" + data.max_per_user + " by you)</i>" : ""));
      var perc = (data.submitted / data.queue_size * 100).toFixed(0);
      $("#progressbar .progress").text(perc + "%");
      $("#progressbar").progress({
        percent: perc,
      });
      if (data.can_submit)
        $("#b-form .input, #b-form .button").removeClass("disabled");
      else
        $("#b-form .input, #b-form .button").addClass("disabled");
    }
    setInterval(function() {
      api("beatmaps/rank_requests/status", {}, updateRankRequestPage);
    }, 10000);
    var re = /^https?:\/\/osu.ppy.sh\/(s|b)\/(\d+)$/gi;
    $("#b-form").submit(function(e) {
      e.preventDefault();
      var v = $("#beatmap").val().trim();
      var reData = re.exec(v);
      if (reData === null) {
        showMessage("error", "Please provide a valid link, in the form " +
          "of either https://osu.ppy.sh/s/&lt;ID&gt; or https://osu.ppy.sh/b/&lt;ID&gt;.");
        $(this).removeClass("loading");
        return false;
      }
      var postData = {};
      if (reData[1] == "s")
        postData.set_id = +reData[2];
      else
        postData.id = +reData[2];
      var t = $(this);     
      api("beatmaps/rank_requests", postData, function(data) {
        t.removeClass("loading");
        showMessage("success", "Beatmap rank request has been submitted.");
        updateRankRequestPage(data);
      }, true)
      return false;
    });
  },

  "/settings/profbackground": function() {
    $("#colorpicker").minicolors({
      inline: true,
    });
    $("#background-type").change(function() {
      $("[data-type]:not([hidden])").attr("hidden", "hidden");
      $("[data-type=" + $(this).val() + "]").removeAttr("hidden");
    });
    $("#file").change(function(e) {
      var f = e.target.files;
      if (f.length < 1) {
        return;
      }
      var u = window.URL.createObjectURL(f[0]);
      var i = document.createElement("img");
      i.src = u;
      i.onload = function() {
        window.URL.revokeObjectURL(this.src);
      };
      $("#image-background").empty().append(i);
    });
  }
};

$(document).ready(function(){
  // semantic stuff
  $('.message .close').on('click', closeClosestMessage);
  $('.ui.checkbox').checkbox();
  $('.ui.dropdown').dropdown();
  $('.ui.progress').progress();
  $('.ui.form').submit(function(e) {
    var t = $(this);
    if (t.hasClass("loading") || t.hasClass("disabled")) {
      e.preventDefault();
      return false;
    }
    t.addClass("loading");
    var f = t.attr("id");
    $("[form='" + f + "']").addClass("loading");
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
function api(endpoint, data, success, post) {
  if (typeof data == "function") {
    success = data;
    data = null;
  }
  
  var errorMessage = "An error occurred while contacting the Ripple API. Please report this to a Ripple developer.";

  $.ajax({
    method:   (post ? "POST" : "GET"),
    dataType: "json",
    url:      "/api/v1/" + endpoint,
    data:     (post ? JSON.stringify(data) : data),
    contentType: (post ? "application/json; charset=utf-8" : ""),
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

window.URL = window.URL || window.webkitURL;

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

// http://stackoverflow.com/a/901144/5328069
function query(name, url) {
    if (!url) {
      url = window.location.href;
    }
    name = name.replace(/[\[\]]/g, "\\$&");
    var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
        results = regex.exec(url);
    if (!results) return null;
    if (!results[2]) return '';
    return decodeURIComponent(results[2].replace(/\+/g, " "));
}

// Useful for forms contacting the Ripple API
function formToObject(form) {
  var inputs = form.find("input, textarea, select");
  var obj = {};
  inputs.each(function(_, el) {
    el = $(el);
    if (el.attr("name") === undefined) {
      return;
    }
    var parts = el.attr("name").split(".");
    var value;
    switch (el.attr("type")) {
    case "checkbox":
      value = el.is(":checked");
      break;
    default:
      switch (el.data("cast")) {
      case "int":
        value = +el.val();
        break;
      default:
        value = el.val();
        break;        
      }
      break;
    }
    obj = modifyObjectDynamically(obj, parts, value);
  });
  return obj;
}

// > modifyObjectDynamically({}, ["nice", "meme", "dude"], "lol")
// { nice: { meme: { dude: 'lol' } } }
function modifyObjectDynamically(obj, inds, set) {
  if (inds.length === 1) {
    obj[inds[0]] = set;
  } else if (inds.length > 1) {
    if (typeof obj[inds[0]] !== "object")
      obj[inds[0]] = {};
    obj[inds[0]] = modifyObjectDynamically(obj[inds[0]], inds.slice(1), set);
  }
  return obj;
}
