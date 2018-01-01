(function() {
  var mapset = {};
  setData.ChildrenBeatmaps.forEach(function(diff) {
    mapset[diff.BeatmapID] = diff;
  });
  console.log(mapset);
  function loadLeaderboard(b, m) {
    var wl = window.location;
    window.history.replaceState('', document.title,
      "/b/" + b + "?mode=" + m + wl.hash);
    api("scores", {
      mode : m,
      b : b,
      p : 1,
      l : 50,
    },
    function(data) {
      console.log(data);
      var tb = $(".ui.table tbody");
      tb.find("tr").remove();
      if (data.scores == null) {
        data.scores = [];
      }
      var i = 0;
      data.scores.forEach(function(score) {
        var user = score.user;
        tb.append($("<tr />").append(
          $("<td />").text("#" + ((page - 1) * 50 + (++i))),
          $("<td />").html("<a href='/u/" + user.id +
                                 "' title='View profile'><i class='" +
                                 user.country.toLowerCase() + " flag'></i>" +
                                 escapeHTML(user.username) + "</a>"),
          $("<td />").html(addCommas(score.score)),
          $("<td />").html(modbits.string(score.mods)),
          $("<td />").text(score.accuracy.toFixed(2) + "%"),
          $("<td />").text(addCommas(score.max_combo)),
          $("<td />").html(score.pp.toFixed(2))));
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

    loadLeaderboard(bid, currentMode);
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
      loadLeaderboard(beatmapID, currentMode);
    });
})();