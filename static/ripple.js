// Ripple custom JS that goes on all pages

$(document).ready(function(){
  $('.message .close')
    .on('click', function() {
      $(this)
        .closest('.message')
        .transition('fade')
      ;
    });

  $('.ui.checkbox')
    .checkbox();
});