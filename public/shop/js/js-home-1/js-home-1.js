(function($) {
    "use strict";
    var Mihanjs = {};
  
    /*----------------------------------- Preloader */
    var preloader = $(".page-loader");
    $(window).on("load", function() {
        var preloaderFadeOutTime = 500;
  
        function hidePreloader() {
            preloader.fadeOut(preloaderFadeOutTime);
        }
  
        hidePreloader();
    });
    /*----------------------------------- Preloader */
  
  
    /*-------------------------------- Sticky Header */
    Mihanjs.StickyHeader = function() {
        $(window).scroll(function() {
            if ($(this).scrollTop() > 55) {
                $('.page-header').addClass('fixed');
                $('.page-header .top-page-header').slideUp(200);
            } else {
                $('.page-header').removeClass('fixed');
                $('.page-header .top-page-header').slideDown(200);
            }
        });
        var lastScrollTop = 0;
        window.addEventListener('scroll', function() {
            var scrollTop = window.pageYOffset || document.documentElement.scrollTop;
            if (scrollTop > lastScrollTop) {
                $('.page-header .nav-wrapper').addClass('hidden--nav-wrapper');
            } else {
                $('.page-header .nav-wrapper').removeClass('hidden--nav-wrapper');
            }
            lastScrollTop = scrollTop;
        });
    }
    /*-------------------------------- Sticky Header */
  
  
    /*-------------------------------- Category List */
    Mihanjs.CategoryList = function() {
        $('.category-list>ul>li:first-child').addClass('active');
        $('.category-list>ul>li').on('mouseenter', function() {
            $(this).addClass('active').siblings().removeClass('active');
        });
    }
    /*-------------------------------- Category List */
  
  
    /*----------------------------- ResponsiveHeader */
    Mihanjs.ResponsiveHeader = function() {
        $('.header-responsive .btn-toggle-side-navigation').on('click', function(event) {
            $(this).siblings('.side-navigation').addClass('toggle');
            $(this).siblings('.overlay-side-navigation').fadeIn(200);
        });
        $('.header-responsive .side-navigation ul li > a').on('click', function(event) {
            event.preventDefault();
            $(this).toggleClass('active');
            $(this).siblings('ul').slideToggle(100);
        });
        $('.header-responsive .overlay-side-navigation').on('click', function(event) {
            $(this).siblings('.side-navigation').removeClass('toggle');
            $(this).fadeOut(200);
        });
    }
    /*----------------------------- ResponsiveHeader */
  
  
    /*-------------------------------- Search Result */
    Mihanjs.SearchResult = function() {
        $('.search-box form input').on('click', function() {
            $(this).parents('.search-box').addClass('show-result').find('.search-result').fadeIn();
        })
        $('.search-box form input').keyup(function() {
            $(this).parents('.search-box').addClass('show-result').find('.search-result').fadeIn();
            $(this).parents('.search-box').find('.search-result-list').fadeIn();
            if ($(this).val().length == 0) {
                // Hide the element
                $(this).parents('.search-box').removeClass('show-result').find('.search-result').fadeOut(100);
                $(this).parents('.search-box').find('.search-result-list').fadeOut();
            } else {
                // Otherwise show it
                $(this).parents('.search-box').find('.search-result').fadeIn(100);
                $(this).parents('.search-box').find('.search-result-list').fadeIn();
            }
        });
        $(document).click(function(e) {
            if ($(e.target).is('.search-box *')) return;
            $('.search-result').hide();
        });
    }
    /*-------------------------------- Search Result */
  
  
 
    
  
    /*------------------------------------ BackToTop */
    Mihanjs.BackToTop = function() {
        $(".back-to-top .back-btn").click(function() {
            $("body,html").animate({
                scrollTop: 0
            }, 700);
            return false;
        });
    }
    /*------------------------------------ BackToTop */
  
    $(window).on("load", function() {});
    $(document).ready(function() {
        Mihanjs.StickyHeader(),
        Mihanjs.CategoryList(), 
        Mihanjs.ResponsiveHeader(), 
        Mihanjs.SearchResult(),
        Mihanjs.BackToTop();
       
    });
  })(jQuery);


    $('.pointer').click(function(){
        $("html, body").animate({ scrollTop: 0 }, 600);
        return false;
    });


    /*------------------------------- Sidebar Sticky */
    if ($(".container .sticky-sidebar").length) {
        $(".container .sticky-sidebar").theiaStickySidebar();
    }
     /*------------------ Slider Profile Order Steps */
    $(".carousel-lg").owlCarousel({
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
            items: 4,
            slideBy: 2,
          },
          992: {
            items: 4,
            slideBy: 2,
          },
          1200: {
            items: 5,
            slideBy: 3,
          },
          1400: {
            items: 6,
            slideBy: 4,
          },
        },
    });
    /*------------------ Slider Profile Order Steps */

   /*---------------------------------- Nice Select */
    if ($(".custom-select-ui").length) {
        $(".custom-select-ui select").niceSelect();
      }
   /*---------------------------------- Nice Select */ 

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
    $("#suggestion-slider").owlCarousel({
        rtl: true,
        items: 1,
        autoplay: true,
        autoplayTimeout: 5000,
        loop: true,
        dots: true,
        onInitialized: startProgressBar,
        onTranslate: resetProgressBar,
        onTranslated: startProgressBar,
    });

    function startProgressBar() {
        // apply keyframe animation
        $(".slide-progress").css({
        width: "100%",
        transition: "width 5000ms",
        });
    }

    function resetProgressBar() {
        $(".slide-progress").css({
        width: 0,
        transition: "width 0s",
        });
    }
    $(".category-slider").owlCarousel({
      rtl: true,
      margin: 10,
      nav: true,
      navText: [
        '<svg width="30px" height="30px" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path opacity="0.1" d="M3 12C3 4.5885 4.5885 3 12 3C19.4115 3 21 4.5885 21 12C21 19.4115 19.4115 21 12 21C4.5885 21 3 19.4115 3 12Z" fill="#323232"></path><path d="M11 15L13.7158 12.2842V12.2842C13.8728 12.1272 13.8728 11.8728 13.7158 11.7158V11.7158L11 9" stroke="#323232" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"></path><path d="M3 12C3 4.5885 4.5885 3 12 3C19.4115 3 21 4.5885 21 12C21 19.4115 19.4115 21 12 21C4.5885 21 3 19.4115 3 12Z" stroke="#323232" stroke-width="2"></path></svg>',
        '<svg width="30px" height="30px" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path opacity="0.1" d="M3 12C3 4.5885 4.5885 3 12 3C19.4115 3 21 4.5885 21 12C21 19.4115 19.4115 21 12 21C4.5885 21 3 19.4115 3 12Z" fill="#323232"/><path d="M3 12C3 4.5885 4.5885 3 12 3C19.4115 3 21 4.5885 21 12C21 19.4115 19.4115 21 12 21C4.5885 21 3 19.4115 3 12Z" stroke="#323232" stroke-width="2"/><path d="M13 9L10.2625 11.7375V11.7375C10.1175 11.8825 10.1175 12.1175 10.2625 12.2625V12.2625L13 15" stroke="#323232" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>',
      ],
      dots: false,
      responsiveClass: true,
      responsive: {
        0: {
          items: 2,
          slideBy: 1,
        },
        576: {
          items: 3,
          slideBy: 2,
        },
        768: {
          items: 4,
          slideBy: 2,
        },
        992: {
          items: 6,
          slideBy: 3,
        },
        1400: {
          items: 8,
          slideBy: 4,
        },
      },
    });