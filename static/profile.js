// code that is executed on every user profile
$(document).ready(function() {
  var wl = window.location;
  var newPathName = wl.pathname;
  // userID is defined in profile.html
  if (newPathName.split("/")[2] != userID) {
    newPathName = "/u/" + userID;
  }
  // if there's no mode parameter in the querystring, add it
  if (wl.search.indexOf("mode=") === -1)
    window.history.replaceState('', document.title, newPathName + "?mode=" + favouriteMode + wl.hash);
  else if (wl.pathname != newPathName)
    window.history.replaceState('', document.title, newPathName + wl.search + wl.hash);
  // when an item in the mode menu is clicked, it means we should change the mode.
  $("#mode-menu>.item").click(function(e) {
    e.preventDefault();
    if ($(this).hasClass("active"))
      return;
    var m = $(this).data("mode");
    $("[data-mode]:not(.item):not([hidden])").attr("hidden", "");
    $("[data-mode=" + m + "]:not(.item)").removeAttr("hidden");
    $("#mode-menu>.active.item").removeClass("active");
    var needsLoad = $("#scores-zone>[data-mode=" + m + "][data-loaded=0]");
    if (needsLoad.length > 0)
      initialiseScores(needsLoad, m);
    $(this).addClass("active");
    window.history.replaceState('', document.title, wl.pathname + "?mode=" + m + wl.hash);
  });
  initialiseFriends();
  // load scores page for the current favourite mode
  initialiseScores($("#scores-zone>div[data-mode=" + favouriteMode + "]"), favouriteMode);
});

function initialiseFriends() {
  var b = $("#add-friend-button");
  if (b.length == 0) return;
  api('friends/with', {id: userID}, setFriendOnResponse);
  b.click(friendClick);
}
function setFriendOnResponse(r) {
  var x = 0;
  if (r.friend) x++;
  if (r.mutual) x++;
  setFriend(x);
}
function setFriend(i) {
  var b = $("#add-friend-button");
  b.removeClass("loading green blue red");
  switch (i) {
  case 0:
    b
      .addClass("blue")
      .attr("title", "Add friend")
      .html("<i class='plus icon'></i>");
    break;
  case 1:
    b
      .addClass("green")
      .attr("title", "Remove friend")
      .html("<i class='minus icon'></i>");
    break;
  case 2:
    b
      .addClass("red")
      .attr("title", "Unmutual friend")
      .html("<i class='heart icon'></i>");
    break;
  }
  b.attr("data-friends", i > 0 ? 1 : 0)
}
function friendClick() {
  var t = $(this);
  if (t.hasClass("loading")) return;
  t.addClass("loading");
  api("friends/" + (t.attr("data-friends") == 1 ? "del" : "add"), {id: userID}, setFriendOnResponse);
}

var defaultScoreTable = $("<table class='ui table score-table' />")
  .append(
    $("<thead />").append(
      $("<tr />").append(
        $("<th>General info</th>"),
        $("<th>Score</th>")
      )
    )
  )
  .append(
    $("<tbody />")
  )
  .append(
    $("<tfoot />").append(
      $("<tr />").append(
        $("<th colspan=2 />").append(
          $("<div class='ui right floated pagination menu' />").append(
            $("<a class='disabled item load-more-button'>Load more</a>").click(loadMoreClick)
          )
        )
      )
    )
  );
function initialiseScores(el, mode) {
  el.attr("data-loaded", "1");
  var best = defaultScoreTable.clone(true).addClass("orange");
  var recent = defaultScoreTable.clone(true).addClass("blue");
  best.attr("data-type", "best");
  recent.attr("data-type", "recent");
  recent.addClass("no bottom margin");
  el.append($("<div class='ui segments no bottom margin' />").append(
    $("<div class='ui segment' />").append("<h2 class='ui header'>Best scores</h2>", best),
    $("<div class='ui segment' />").append("<h2 class='ui header'>Recent scores</h2>", recent)
  ));
  loadScoresPage("best", mode);
  loadScoresPage("recent", mode);
};
function loadMoreClick() {
  var t = $(this);
  if (t.hasClass("disabled"))
    return;
  t.addClass("disabled");
  var type = t.parents("table[data-type]").data("type");
  var mode = t.parents("div[data-mode]").data("mode");
  loadScoresPage(type, mode);
}
// currentPage for each mode
var currentPage = {
  0: {best: 0, recent: 0},
  1: {best: 0, recent: 0},
  2: {best: 0, recent: 0},
  3: {best: 0, recent: 0},
};
var scoreStore = {};
function loadScoresPage(type, mode) {
  var table = $("#scores-zone div[data-mode=" + mode + "] table[data-type=" + type + "] tbody");
  var page = ++currentPage[mode][type];
  console.log("loadScoresPage with", {
    page: page,
    type: type,
    mode: mode,
  });
  api("users/scores/" + type, {
    mode: mode,
    p: page,
    l: 20,
    id: userID,
  }, function(r) {
    if (r.scores == null) {
      disableLoadMoreButton(type, mode);
      return;
    }
    r.scores.forEach(function(v){
      scoreStore[v.id] = v;
      var scoreRank = getRank(mode, v.mods, v.accuracy, v.count_300, v.count_100, v.count_50, v.count_miss);
      table.append($("<tr class='new score-row' data-scoreid='" + v.id + "' />").append(
        $(
          "<td><img src='/static/ranking-icons/" + scoreRank + ".png' alt='" + scoreRank.toUpperCase() + "'> " +
          escapeHTML(v.beatmap.song_name) + " <b>" + getScoreMods(v.mods) + "</b> <i>(" + v.accuracy.toFixed(2) + "%)</i><br />" +
          "<div class='subtitle'><time class='new timeago' datetime='" + v.time + "'>" + v.time + "</time></div></td>"
        ),
        $("<td><b>" + ppOrScore(v.pp, v.score) + "</b></td>")
      ));
    });
    $(".new.timeago").timeago().removeClass("new");
    $(".new.score-row").click(viewScoreInfo).removeClass("new");
    var enable = true;
    if (r.scores.length != 20)
      enable = false;
    disableLoadMoreButton(type, mode, enable);
  });
};
function disableLoadMoreButton(type, mode, enable) {
  var button = $("#scores-zone div[data-mode=" + mode + "] table[data-type=" + type + "] .load-more-button");
  if (enable) button.removeClass("disabled");
  else button.addClass("disabled");
}
function viewScoreInfo() {
  var scoreid = $(this).data("scoreid");
  if (!scoreid && scoreid !== 0) return;
  var s = scoreStore[scoreid];
  if (s === undefined) return;

  // data to be displayed in the table.
  var data = {
    "Points":       addCommas(s.score),
    "PP":           addCommas(s.pp),
    "Beatmap":      "<a href='/b/" + s.beatmap.beatmap_id + "'>" + escapeHTML(s.beatmap.song_name) + "</a>",
    "Accuracy":     s.accuracy + "%",
    "Max combo":    addCommas(s.max_combo) + "/" + addCommas(s.beatmap.max_combo)
                      + (s.full_combo ? " (full combo)" : ""),
    "Difficulty":   s.beatmap.difficulty2[modesShort[s.play_mode]] + " stars",
    "Mods":         getScoreMods(s.mods, true),
    "300s":         addCommas(s.count_300),
    "100s":         addCommas(s.count_100),
    "50s":          addCommas(s.count_50),
    "Gekis":        addCommas(s.count_geki),
    "Katus":        addCommas(s.count_katu),
    "Misses":       addCommas(s.count_miss),
    "Ranked?":      s.completed == 3 ? "Yes" : "No",
    "Achieved":     s.time,
    "Mode":         modes[s.play_mode],
  };

  var els = [];
  $.each(data, function(key, value) {
    els.push(
      $("<tr />").append(
        $("<td>" + key + "</td>"),
        $("<td>" + value + "</td>")
      )
    );
  });

  $("#score-data-table tr").remove();
  $("#score-data-table").append(els);
  $(".ui.modal").modal("show");
}

// helper functions copied from user.js in old-frontend
function getScoreMods(m, noplus) {
  // TODO: refactoring
  // 1. use $.foreach and an object instead of ugly copypasted ifs
  // 2. r = [], r.append("NF"), r.join(", ")
	var r = '';
	var hasNightcore = false;
  var hasPerfect = false;
	if (m & NoFail)
		r += 'NF, ';
	if (m & Easy)
		r += 'EZ, ';
	if (m & NoVideo)
		r += 'NV, ';
	if (m & Hidden)
		r += 'HD, ';
	if (m & HardRock)
		r += 'HR, ';
	if (m & Perfect) {
		r += 'PF, ';
    hasPerfect = true;
	}
	if (!hasPerfect && (m & SuddenDeath))
		r += 'SD, ';
	if (m & Nightcore) {
		r += 'NC, ';
    hasNightcore = true;
	}
	if (!hasNightcore && (m & DoubleTime))
		r += 'DT, ';
	if (m & Relax)
		r += 'RX, ';
	if (m & HalfTime)
		r += 'HT, ';
	if (m & Flashlight)
		r += 'FL, ';
	if (m & Autoplay)
		r += 'AP, ';
	if (m & SpunOut)
		r += 'SO, ';
	if (m & Relax2)
		r += 'AP, ';
	if (m & Key4)
		r += '4K, ';
	if (m & Key5)
		r += '5K, ';
	if (m & Key6)
		r += '6K, ';
	if (m & Key7)
		r += '7K, ';
	if (m & Key8)
		r += '8K, ';
	if (m & FadeIn)
		r += 'FD, ';
	if (m & Random)
		r += 'RD, ';
	if (m & LastMod)
		r += 'CN, ';
	if (m & Key9)
		r += '9K, ';
	if (m & Key10)
		r += '10K, ';
	if (m & Key1)
		r += '1K, ';
	if (m & Key3)
		r += '3K, ';
	if (m & Key2)
		r += '2K, ';
	if (r.length > 0) {
		return (noplus ? "" : "+ ") + r.slice(0, -2);
	} else {
		return (noplus ? 'None' : '');
	}
}

var None = 0;
var NoFail = 1;
var Easy = 2;
var NoVideo = 4;
var Hidden = 8;
var HardRock = 16;
var SuddenDeath = 32;
var DoubleTime = 64;
var Relax = 128;
var HalfTime = 256;
var Nightcore = 512;
var Flashlight = 1024;
var Autoplay = 2048;
var SpunOut = 4096;
var Relax2 = 8192;
var Perfect = 16384;
var Key4 = 32768;
var Key5 = 65536;
var Key6 = 131072;
var Key7 = 262144;
var Key8 = 524288;
var keyMod = 1015808;
var FadeIn = 1048576;
var Random = 2097152;
var LastMod = 4194304;
var Key9 = 16777216;
var Key10 = 33554432;
var Key1 = 67108864;
var Key3 = 134217728;
var Key2 = 268435456;

function getRank(gameMode, mods, acc, c300, c100, c50, cmiss) {
	var total = c300+c100+c50+cmiss;

	var hdfl = (mods & (Hidden | Flashlight | FadeIn)) > 0;

	var ss = hdfl ? "sshd" : "ss";
	var s = hdfl ? "shd" : "s";

	switch(gameMode) {
		case 0:
		case 1: // TODO: check taiko's
			var ratio300 = c300 / total;
			var ratio50 = c50 / total;

			if (ratio300 == 1)
				return ss;

			if (ratio300 > 0.9 && ratio50 <= 0.01 && cmiss == 0)
				return s;

			if ((ratio300 > 0.8 && cmiss == 0) || (ratio300 > 0.9))
				return "a";

			if ((ratio300 > 0.7 && cmiss == 0) || (ratio300 > 0.8))
				return "b";

			if (ratio300 > 0.6)
				return "c";

			return "d";

		case 2:
			if (acc == 100)
				return ss;

			if (acc > 98)
				return s;

			if (acc > 94)
				return "a";

			if (acc > 90)
				return "b";

			if (acc > 85)
				return "c";

			return "d";

		case 3:
			if (acc == 100)
				return ss;

			if (acc > 95)
				return s;

			if (acc > 90)
				return "a";

			if (acc > 80)
				return "b";

			if (acc > 70)
				return "c";

			return "d";
	}
}

function ppOrScore(pp, score) {
  if (pp != 0)
    return addCommas(pp.toFixed(2)) + "pp";
  return addCommas(score);
}

function beatmapLink(type, id) {
  if (type == "s")
    return "<a href='/s/" + id + "'>" + id + '</a>';
  return "<a href='/b/" + id + "'>" + id + '</a>';  
}
