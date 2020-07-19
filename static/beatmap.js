(function() {
  var mapset = {};
  setData.ChildrenBeatmaps.forEach(function(diff) {
    mapset[diff.BeatmapID] = diff;
  });
  function loadLeaderboard(b, m, r) {
    var wl = window.location;
    window.history.replaceState('', document.title,
      "/b/" + b + "?mode=" + m + "&relax=" + r + wl.hash);
    api("scores?sort=score,desc&sort=id,asc", {
      mode : m,
      b : b,
      p : 1,
      l : 50,
      relax: r,
    },
    function(data) {
      var tb = $(".ui.table tbody");
      tb.find("tr").remove();
      if (data.scores == null) {
        data.scores = [];
      }
      var i = 0;
      data.scores.sort(function(a, b) { return b.score - a.score; });
      data.scores.forEach(function(score) {
        var user = score.user;
        tb.append($("<tr />").append(
          $("<td data-sort-value=" + (++i) + " />")
            .text("#" + ((page - 1) * 50 + i)),
          $("<td />").html("<a href='/u/" + user.id +
                                 "' title='View profile'><i class='" +
                                 user.country.toLowerCase() + " flag'></i>" +
                                 escapeHTML(user.username) + "</a>"),
          $("<td data-sort-value=" + score.score + " />")
            .html(addCommas(score.score)),
          $("<td />").html(getScoreMods(score.mods, true)),
          $("<td data-sort-value=" + score.accuracy + " />")
            .text(score.accuracy.toFixed(2) + "%"),
          $("<td data-sort-value=" + score.max_combo + " />")
            .text(addCommas(score.max_combo)),
          $("<td data-sort-value=" + score.pp + " />")
            .html(score.pp.toFixed(2))));
      });
    });
  }
  function changeDifficulty(bid) {
    // load info
    var diff = mapset[bid];

    // column 2
    $("#cs").html(diff.CS);
    $("#hp").html(diff.HP);
    $("#od").html(diff.OD);
    $("#passcount").html(addCommas(diff.Passcount));
    $("#playcount").html(addCommas(diff.Playcount));

    // column 3
    $("#ar").html(diff.AR);
    $("#stars").html(diff.DifficultyRating.toFixed(2));
    $("#length").html(timeFormat(diff.TotalLength));
    $("#drainLength").html(timeFormat(diff.HitLength));
    $("#bpm").html(diff.BPM);

    // hide mode for non-std maps
    if (diff.Mode != 0) {
      // Non-std! Force right game mode
      currentMode = diff.Mode;
      $("#mode-menu").hide();
    } else {
      // Std, all modes supported
      // Choose fav if no mode qs param was provided
      if (currentMode === null) {
        currentMode = favMode;
      }
      $("#mode-menu").show();
      $("#relax-menu").show();
    }

    // hide classic/relax switcher for mania only-maps
    if (diff.Mode == 3) {
      currentRelax = 0;
      $("#relax-menu").hide();
    } else if (currentRelax === null) {
      // Chose fav if no relax qs param was provided
      currentRelax = favRelax;
    }

    // update mode menu
    $("#mode-menu .active.item").removeClass("active");
    $("#mode-" + currentMode).addClass("active");

    // brico meiser
    $("#relax-menu>[data-relax=" + favRelax + "]").addClass("active");
    loadLeaderboard(bid, currentMode, currentRelax);
  }
  window.loadLeaderboard = loadLeaderboard;
  window.changeDifficulty = changeDifficulty;
  changeDifficulty(beatmapID);
  // loadLeaderboard(beatmapID, currentMode);
  $("#diff-menu .item")
    .click(function(e) {
      e.preventDefault();
      $(this).addClass("active");
      beatmapID = $(this).data("bid");
      changeDifficulty(beatmapID);
    });
  $("#mode-menu .item")
    .click(function(e) {
      e.preventDefault();
      $("#mode-menu .active.item").removeClass("active");
      $(this).addClass("active");
      currentMode = $(this).data("mode");
      if (currentMode == 3) {
        $("#relax-menu>[data-relax=1]").addClass("disabled").removeClass("active");
        $("#relax-menu>[data-relax=0]").addClass("active");
      } else {
        $("#relax-menu>[data-relax=1]").removeClass("disabled");
      }
      loadLeaderboard(beatmapID, currentMode, currentRelax);
    });
  $("#relax-menu .item")
    .click(function(e) {
      e.preventDefault();
      if ($(this).hasClass("disabled")) {
        return;
      }
      $("#relax-menu .active.item").removeClass("active");
      $(this).addClass("active");
      currentRelax = $(this).data("relax");
      loadLeaderboard(beatmapID, currentMode, currentRelax);
    })
  $("table.sortable").tablesort();
})();