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


    /*------------------------------------ Carousel */
    $(".carousel-md").owlCarousel({
        rtl: true,
        margin: 10,
        nav: true,
        navText: [
            '<svg style="width="30px" height="30px" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path opacity="0.1" d="M3 12C3 4.5885 4.5885 3 12 3C19.4115 3 21 4.5885 21 12C21 19.4115 19.4115 21 12 21C4.5885 21 3 19.4115 3 12Z" fill="#323232"></path><path d="M11 15L13.7158 12.2842V12.2842C13.8728 12.1272 13.8728 11.8728 13.7158 11.7158V11.7158L11 9" stroke="#323232" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"></path><path d="M3 12C3 4.5885 4.5885 3 12 3C19.4115 3 21 4.5885 21 12C21 19.4115 19.4115 21 12 21C4.5885 21 3 19.4115 3 12Z" stroke="#323232" stroke-width="2"></path></svg>',
            '<svg width="30px" height="30px" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path opacity="0.1" d="M3 12C3 4.5885 4.5885 3 12 3C19.4115 3 21 4.5885 21 12C21 19.4115 19.4115 21 12 21C4.5885 21 3 19.4115 3 12Z" fill="#323232"/><path d="M3 12C3 4.5885 4.5885 3 12 3C19.4115 3 21 4.5885 21 12C21 19.4115 19.4115 21 12 21C4.5885 21 3 19.4115 3 12Z" stroke="#323232" stroke-width="2"/><path d="M13 9L10.2625 11.7375V11.7375C10.1175 11.8825 10.1175 12.1175 10.2625 12.2625V12.2625L13 15" stroke="#323232" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>',
          ],
        dots: true,
        responsiveClass: true,
        responsive: {
          0: {
            items: 2,
            slideBy: 1,
          },
          480: {
            items: 2,
            slideBy: 1,
          },
          576: {
            items: 3,
            slideBy: 1,
          },
          768: {
            items: 3,
            slideBy: 2,
          },
          992: {
            items: 4,
            slideBy: 2,
          },
          1200: {
            items: 4,
            slideBy: 3,
          },
          1400: {
            items: 5,
            slideBy: 4,
          },
        },
    });
    /*------------------------------------ Carousel */

    /*------------------------------------ Shopping */
    $("#btn-checkout-contact-location").click(function () {
        $(".checkout-address").addClass("show");
        $(".checkout-contact-content").addClass("hidden");
    });
    $("#cancel-change-address-btn").click(function () {
        $(".checkout-address").removeClass("show");
        $(".checkout-contact-content").removeClass("hidden");
    });
    /*------------------------------------ Shopping */
    
  
    /*---------------------------------- Nice Select */
    if ($(".custom-select-ui").length) {
        $(".custom-select-ui select").niceSelect();
    }
    /*---------------------------------- Nice Select */ 