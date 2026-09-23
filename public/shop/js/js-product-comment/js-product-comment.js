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


    /*---------------------------------- Advantage Input */

    var inputs = $('#advantage-input, #disadvantage-input');
    var inputChangeCallback = function () {
        var self = $(this);
        if (self.val().trim().length > 0) {
            self.siblings('.js-icon-form-add').show();
        } else {
            self.siblings('.js-icon-form-add').hide();
        }
    };
    inputs.each(function () {
        inputChangeCallback.bind(this)();
        $(this).on('change keyup', inputChangeCallback.bind(this));
    });
    $("#advantages").delegate(".js-icon-form-add", 'click', function (e) {
        var parent = $('.js-advantages-list');
        if (parent.find(".js-advantage-item").length >= 5) {
            return;
        }
        var advantageInput = $('#advantage-input');
        if (advantageInput.val().trim().length > 0) {
            parent.append(
                '<div class="ui-dynamic-label ui-dynamic-label--positive js-advantage-item">\n' +
                advantageInput.val() +
                '<button type="button" class="ui-dynamic-label-remove js-icon-form-remove"></button>\n' +
                '<input type="hidden" name="comment[advantages][]" value="' + advantageInput
                    .val() + '">\n' +
                '</div>');
            advantageInput.val('').change();
            advantageInput.focus();
        }
    }).delegate(".js-icon-form-remove", 'click', function (e) {
        $(this).parent('.js-advantage-item').remove();
    });
    $("#disadvantages").delegate(".js-icon-form-add", 'click', function (e) {
        var parent = $('.js-disadvantages-list');
        if (parent.find(".js-disadvantage-item").length >= 5) {
            return;
        }
        var disadvantageInput = $('#disadvantage-input');
        if (disadvantageInput.val().trim().length > 0) {
            parent.append(
                '<div class="ui-dynamic-label ui-dynamic-label--negative js-disadvantage-item">\n' +
                disadvantageInput.val() +
                '<button type="button" class="ui-dynamic-label-remove js-icon-form-remove"></button>\n' +
                '<input type="hidden" name="comment[disadvantages][]" value="' +
                disadvantageInput.val() + '">\n' +
                '</div>');
            disadvantageInput.val('').change();
            disadvantageInput.focus();
        }
    }).delegate(".js-icon-form-remove", 'click', function (e) {
        $(this).parent('.js-disadvantage-item').remove();
    });
    /*---------------------------------- Advantage Input */