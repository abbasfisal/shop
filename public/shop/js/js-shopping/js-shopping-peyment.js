var preloader = $(".page-loader");
$(window).on("load", function() {
    var preloaderFadeOutTime = 500;

    function hidePreloader() {
        preloader.fadeOut(preloaderFadeOutTime);
    }

    hidePreloader();
});


  /*------------------------------- Sidebar Sticky */
  if ($(".container .sticky-sidebar").length) {
      $(".container .sticky-sidebar").theiaStickySidebar();
  }
  /*------------------------------- Sidebar Sticky */
