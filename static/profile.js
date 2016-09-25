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
    window.history.pushState('', document.title, newPathName + "?mode=" + favouriteMode + wl.hash);
  else if (wl.pathname != newPathName)
    window.history.pushState('', document.title, newPathName + wl.search + wl.hash);
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
      initialiseScores(needsLoad);
    $(this).addClass("active");
    window.history.pushState('', document.title, wl.pathname + "?mode=" + m + wl.hash);
  });
  // load scores page for the current favourite mode
  initialiseScores($("#scores-zone>[data-mode=" + favouriteMode + "]"));
});

var loadMoreClick = function() {
  var t = $(this);
  console.log(t);
  var type = t.parents("[data-type]").data("type");
  var mode = t.parents("[data-mode]").data("mode");
  console.log(type, mode);
};
var defaultScoreTable = $("<table class='ui very basic celled table score-table' />")
  .append(
    $("<thead />").append(
      $("<tr />").append(
        $("<th>Rank</th>"),
        $("<th>Song info</th>"),
        $("<th>Score/PP</th>")
      )
    )
  )
  .append(
    $("<tbody />")
  )
  .append(
    $("<tfoot />").append(
      $("<tr />").append(
        $("<th colspan=3 />").append(
          $("<div class='ui right floated pagination menu' />").append(
            $("<a class='item'>Load more</a>").click(loadMoreClick)
          )
        )
      )
    )
  );
var initialiseScores = function(el) {
  el.attr("data-loaded", "1");
  var best = defaultScoreTable.clone(true);
  var recent = defaultScoreTable.clone(true);
  best.attr("data-type", "best");
  recent.attr("data-type", "recent");
  recent.addClass("no bottom margin");
  el.append($("<div class='ui segments no bottom margin' />").append(
    $("<div class='ui segment' />").append("<h2 class='ui header'>Best scores</h2>").append(best),
    $("<div class='ui segment' />").append("<h2 class='ui header'>Recent scores</h2>").append(recent)
  ))
};
var loadScoresPage = function(type, element, page) {
  if (page === undefined) {
    page = type;
    type = undefined;
  }
};
