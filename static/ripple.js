// Ripple custom JS that goes on all pages

var singlePageSnippets = {
  "/login": function() {
    $("#login-form").submit(function(e) {
      $("button").addClass("disabled");
      
      var fix = function(errorMessage) {
        $("button").removeClass("disabled");
        $(".ui.form").removeClass("loading");
        var newEl = $('<div class="ui error message hidden"><i class="close icon"></i>' + errorMessage + '</div>');
        newEl.find(".close.icon").click(closeClosestMessage);
        $("#messages-container").append(newEl);
        newEl.transition("slide down");
      };
      
      if (!/^[a-zA-Z0-9 \[\]\@\.\+-]+$/.test($("input[name='username']").val())) {
        fix("Invalid username.");
        return false;
      }
    });
  }
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
