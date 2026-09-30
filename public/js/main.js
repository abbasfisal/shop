let win = $(window),
    breakpoint_lg = 1024,
    header_height;

mainMenu();
setTimeout(() => {
    header()
}, 20)

searchHeader();

textWithShowMore()
modal()
quickAddTocart()
removeByClick()
productSlider();
textTruncate()
popover()
tabs()
profile();
dropDown();
goToElement()
gotoTop()



$("input[type='text'].number").on({
    keyup: function (e) {
        if (e.keyCode == 8) {
            let value = convertedValue = $(this).val()
            let i = this.selectionStart;


            if ($(this).hasClass('currency'))
                convertedValue = numberCurrency(convertedValue);

            if ($(this).hasClass('fanumber'))
                convertedValue = convertToFaDigit(convertedValue);

            $(this).val(convertedValue)


            if (value.length == convertedValue.length)
                this.selectionStart = this.selectionEnd = i;
            else if (value.length < convertedValue.length)
                this.selectionStart = this.selectionEnd = i + 1;
            else
                this.selectionStart = this.selectionEnd = i - 1;
        }
    },

    keypress: function (e) {
        var charCode = e.keyCode;

        if (charCode > 31 && (charCode < 48 || charCode > 57))
            return false;

        let value = convertedValue = $(this).val()
        let i = this.selectionStart;


        if ($(this).hasClass('fanumber'))
            convertedValue = value = value.substr(0, i) + convertToFaDigit(e.key) + value.substr(this.selectionEnd);


        if ($(this).hasClass('currency'))
            convertedValue = numberCurrency(convertedValue);

        if ($(this).hasClass('fanumber'))
            convertedValue = convertToFaDigit(convertedValue);

        $(this).val(convertedValue);

        if (convertedValue.length > value.length)
            this.selectionStart = this.selectionEnd = i + 2;
        else
            this.selectionStart = this.selectionEnd = i + 1;

        return false;

    }
})

function numberCurrency(n) {
    n = convertToEnDigit(n)
    return n.replace(/\D/g, "").replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}


(function ($) {
    $.fn.ToCurrency = function () {
        $(this).each(function () {
            let txt;

            if ($(this).is('input')) txt = $(this).val()
            else txt = $(this).text()

            let str = numberCurrency(txt);

            if ($(this).is('input')) $(this).val(str)
            else $(this).text(str)

        })
    }
}(jQuery));

(function ($) {
    $.fn.ToFaDigit = function () {
        $(this).each(function () {
            let txt;

            if ($(this).is('input')) txt = $(this).val()
            else txt = $(this).text()

            txt = convertToFaDigit(txt);

            if ($(this).is('input')) $(this).val(txt)
            else $(this).text(txt)

        })
    }
}(jQuery));

(function ($) {
    $.fn.ToCountDown = function () {
        $(this).each(function () {
            let that = $(this),
                eventTime = that.data('date'),
                currenctTime = '1366547400000',
                leftTime = eventTime - currenctTime,
                duration = moment.duration(leftTime, 'seconds');



            setInterval(function () {
                duration = moment.duration(duration.asSeconds() - 1, 'seconds')
                let value = duration.seconds() + ' : ' + duration.minutes() + ' : ' + duration.hours()

                if (that.hasClass('fanumber'))
                    value = convertToFaDigit(value);

                that.text(value)

            }, 1000)

        })
    }
}(jQuery));

function convertToFaDigit(text) {
    var en_number = ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9'];
    var fa_number = ['۰', '۱', '۲', '۳', '۴', '۵', '۶', '۷', '۸', '۹'];

    for (i = 0; i <= 9; i++) {

        var regex = new RegExp(en_number[i], 'g');
        text = text.toString().replace(regex, fa_number[i])

    }

    return text;
}

function convertToEnDigit(text) {
    var en_number = ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9'];
    var fa_number = ['۰', '۱', '۲', '۳', '۴', '۵', '۶', '۷', '۸', '۹'];

    for (i = 0; i <= 9; i++) {

        var regex = new RegExp(fa_number[i], 'g');
        text = text.toString().replace(regex, en_number[i])

    }

    return text;
}

$('.currency').ToCurrency()
$('.fanumber').ToFaDigit()
$('.countDown').ToCountDown()

$(function () {
    $('.input-focus-changeborder').children($('input')).focus(function () {
        $(this).parent().addClass('focused')
    })

    $('.input-focus-changeborder').children($('input')).blur(function () {
        $(this).parent().removeClass("focused")
    });

    $('.focus-changeparentBorder').children($('input')).focus(function () {
        $(this).parent().addClass('focused')
    })

    $('.focus-changeparentBorder').children($('input')).blur(function () {
        $(this).parent().removeClass("focused")
    });



})

function removeByClick() {
    $('.removeByClick').click(function () { $(this).remove() })
}

function gotoTop() {
    $('.gotoTop').click(function () {
        $('html,body').animate({
            scrollTop: 0
        }, 250)
    })
}

function goToElement() {
    let el = $('.goto');

    $('.goto').click(function () {
        var contentId = $(this).data('target');

        $('html,body').animate({
            scrollTop: $('[data-id=' + contentId + ']').offset().top - header_height
        }, 250)
    })
}






function openHeaderOverlay() {
    $('.header-overlay').addClass('active')
    $('body').addClass('no-overflow')
}

function closeOverlay() {
    $('.overlay').removeClass('active')
    $('body').removeClass('no-overflow')
}


function header() {
    let position = win.scrollTop(),
        nav = $('.header-nav'),
        header = $('header'),
        nav_height = nav.height(),
        headerTop_height = $('.header-top').outerHeight(),
        currentPosition = 0,
        main = $('.main');
    header_height = header.height()
    main.css('padding-top', header_height)

    win.scroll(function () {
        currentPosition = win.scrollTop();


        if (currentPosition < 80) return;

        if (win.width() < breakpoint_lg) return;

        if (currentPosition > position) {

            nav.css('transform', 'translateY(-100%)');

            header.height(header_height - nav_height)
        }
        else {

            nav.css('transform', '')
            header.height(header_height)

        }

        position = currentPosition;
    })


    win.on('resize', function () {
        nav_height = nav.height();
        if (win.width() < breakpoint_lg) {
            header.height('');
            header_height = header.height()
        }
        else {
            header_height = $('.header-top').outerHeight() + nav_height
        }


        main.css('padding-top', header_height)



    })



    $('.header-profilebtn').click(function (e) {
        $(".header-profile").toggleClass("active");
        //e.stopPropagation();
    })

    $(document).on("click", function (e) {

        if ($(".header-profile").has(e.target).length === 0) {
            $(".header-profile").removeClass("active");
        }
    });
}

function searchHeader() {
    new Swiper(".searchTrends-swiper", {
        spaceBetween: 8,
        freeMode: true,
        slidesPerView: "auto",
        navigation: {
            nextEl: ".swiper-btn-next",
            prevEl: ".swiper-btn-prev",
        },
    });

    var input = $('.header-search-input'),
        result = $('.header-search-result'),
        search = $('.search');

    win.on("click", function (e) {
        if (input.is(e.target) ||
            result.is(e.target) ||
            result.has(e.target).length) {
            search.addClass("active");
            openHeaderOverlay();
        }
        else {
            search.removeClass("active");
            if ($('.header-overlay').hasClass('active'))
                closeOverlay();
        }
    })


}


function mainMenu() {
    menuHover();
    showSubMenu();
    mobileMenu();

    function menuHover() {
        let li_hover = $('.mainMenu-hover'),
            li = $('.mainMenu > li');

        li.hover(function () {
            li_hover.css('width', $(this).width());
            li_hover.css('left', $(this).offset().left)
            li_hover.css('transform', 'scaleX(1)')
        }, function () {
            li_hover.css('transform', 'scaleX(0)')
        })



    }


    function showSubMenu() {
        let mainMenuList = $('.mainMenu-list');

        $('.mainMenu-li-hasSubMenu').hover(function () {
            mainMenuList.addClass('active');
            $('.search').removeClass('active')
            openHeaderOverlay()

        }, function () {
            mainMenuList.removeClass('active');
            closeOverlay()
        })

        $('.mainMenu-list-categoryTitle').hover(function () {
            let id = $(this).data('id');

            $('.mainMenu-list-categoryTitle').removeClass('hovered');
            $('.mainMenu-sublist-content').removeClass('active');


            $(this).addClass('hovered')
            $('.mainMenu-sublist-content[data-id=' + id + ']').addClass('active')
        })
    }

    function mobileMenu() {
        addMobileClass();

        win.on('resize', addMobileClass)
        function addMobileClass() {
            var winwidth = win.width();


            if (winwidth < breakpoint_lg) {
                $('.header-nav').css({ 'transform': "translateX(312px)" })
                setTimeout(function () {
                    $('.header-nav').addClass('is-mobile')
                    $('.header-nav').css({ 'transform': "" })


                }, 5)

            }
            else {
                $('.header-nav ').removeClass('is-mobile')

            }
        }

        $('.mobile-hamburgerMenubtn').click(function () {
            $('nav.is-mobile').addClass("active");

            $('.mobile-header-overlay').addClass('active')
        })

        $('.mobile-header-overlay').click(function () {
            $('nav.is-mobile').removeClass("active");
            closeOverlay()
        })

        $('.mainMenu-list-categoryTitle').click(function () {
            if ($(this).hasClass('active')) {
                $(this).removeClass('active')
                $(this).next('.mainMenu-sublist-menulist').removeClass('d-block')
            }
            else {
                $(this).addClass('active');
                $(this).next('.mainMenu-sublist-menulist').addClass('d-block')
            }

        })

        $('.mainMenu-sublist-item-main').click(function () {
            if ($(this).hasClass('active')) {
                $(this).removeClass('active')
                $(this).nextUntil('.mainMenu-sublist-item-main').removeClass('active')
            }
            else {
                $(this).addClass('active');
                $(this).nextUntil('.mainMenu-sublist-item-main').addClass('active')
            }

        })
    }
}

function mainSlider() {
    new Swiper(".mainSlider-swiper", {

        loop: true,
        slidesPerView: "auto",
        navigation: {
            nextEl: ".swiper-btn-next",
            prevEl: ".swiper-btn-prev",
        },
        autoplay: {
            delay: 7000,
          },
        pagination: {
            el: '.swiper-pagination'
        }
    });
}


function amazingCarousel() {
    new Swiper(".amazing-carousel-swiper", {


        slidesPerView: "auto",
        spaceBetween: 2,
        freeMode: true,
        navigation: {
            nextEl: ".swiper-btn-next",
            prevEl: ".swiper-btn-prev",
        },
    });
}


function digikalaRecommendationSwiper() {
    new Swiper(".digikalaRecommendation-swiper", {
        slidesPerView: "auto",
        freeMode: true,

    });
}

function popularBrandSwiper() {
    new Swiper(".popularBrands-swiper", {


        slidesPerView: "auto",
        spaceBetween: 2,
        freeMode: true,
        navigation: {
            nextEl: ".swiper-btn-next",
            prevEl: ".swiper-btn-prev",
        },
    });
}


function bestSellingswiper() {
    new Swiper(".bestSelling-swiper", {


        slidesPerView: "auto",
        spaceBetween: 20,
        freeMode: true,
        navigation: {
            nextEl: ".swiper-btn-next",
            prevEl: ".swiper-btn-prev",
        },
    });
}

function textWithShowMore() {
    let showMore = $('.textWithShowore-btn-showMore'),
        span = showMore.find('span'),
        footerdesc = $('.textWithShowore ');


    showMore.click(function () {
        span.text(span.text() == 'بستن' ? 'مشاهده بیشتر' : 'بستن');

        footerdesc.toggleClass('opened')
    })
}


function modal() {
    $('[data-modal]').click(function () {
        let id = $(this).data('modal');
        $('.modal[data-id=' + id + ']').addClass('active')
        $('.modal-overlay').addClass('active');
        $('body').addClass('no-overflow')
    });

    $('.modal-close').click(function () {
        $('.modal-overlay').removeClass('active')
        $('.modal').removeClass('active')
        $('body').removeClass('no-overflow')
    })


    $('[data-modal-mobile]').click(function () {
        let id = $(this).data('modal-mobile');
        $('.mobileModal[data-id=' + id + ']').addClass('active')
        $('.modal-overlay').addClass('active');
        $('body').addClass('no-overflow')
    });

    $('.mobile-modal-close').click(function () {
        $('.modal-overlay').removeClass('active')
        $('.mobileModal').removeClass('active')
        $('body').removeClass('no-overflow')
    })

    $('.modal-overlay').click(function () {
        $('.modal-overlay').removeClass('active')
        $('.modal').removeClass('active')
        $('.mobileModal').removeClass('active')
        $('body').removeClass('no-overflow')
    });

    $('.modal-close-currenct').click(function () {
        $(this).parents($('.modal')).removeClass('active');
    });
}



function searchPage() {
    priceRangeSlider()

    $('.search-topCategoryList-item-showMore').click(function () {
        $(this).remove()
    });


    $('.searchFilter-prop-title').click(function () {
        $(this).parents('.search-filter-aside-filter').find('.searchFilter-props ').toggleClass('opened')
    })


    $('.searchFilter-props-searchinput').keyup(function () {
        filterSearch($(this))
    })

    $('.icon-searchFilter-searchinput').click(function () {
        let input = $(this).parent().find('.searchFilter-props-searchinput');
        input.val('')
        filterSearch(input)
    })

    function filterSearch(input) {
        let value = input.val();

        let contentList = input.parents('.searchFilter-props').children("div[data-fa]")

        if (value.length > 0) {
            input.addClass('active');
            contentList.filter(function () {

                let dataEn = this.getAttribute('data-en');
                let dataFa = this.getAttribute('data-fa');

                if (dataEn.toLowerCase().indexOf(value) > -1 || dataFa.indexOf(value) > -1)
                    $(this).removeClass("d-none")
                else $(this).addClass("d-none")


            })
        }
        else {
            input.removeClass('active');
            contentList.filter(function () {
                $(this).removeClass("d-none")
            })
        }

    }


    function priceRangeSlider() {
        let slider = document.querySelector('.rangePrice-slider'),
            min = 0,
            max = 23780000,
            inputMin = document.querySelector('#pricerange-inputmin'),
            inputMax = document.querySelector('#pricerange-inputmax');

        noUiSlider.create(slider, {
            start: [min, max],
            direction: 'rtl',
            connect: true,
            range: {
                'min': min,
                'max': max
            }
        });

        slider.noUiSlider.on('update', function (values, handle) {

            var value = values[handle];

            if (handle) {
                inputMax.value = convertToFaDigit(numberCurrency(Math.round(value)))
            }
            else {

                inputMin.value = convertToFaDigit(numberCurrency(Math.round(value)))
            }
        })
        inputMax.addEventListener('keyup', function () {

            let val = convertToEnDigit(this.value).replace(/\,/g, '');
            slider.noUiSlider.set([null, val])
        })
        inputMin.addEventListener('keyup', function () {

            let val = convertToEnDigit(this.value).replace(/\,/g, '');
            slider.noUiSlider.set([val, null])
        })
    }

    $('.search-filter-aside').theiaStickySidebar({
        'additionalMarginTop': 180,
        'additionalMarginBottom': 20,
        'updateSidebarHeight': true
    })
}
productPriceAdvantage()
function productPriceAdvantage() {
 
    $('.productPriceAdvantage-content').each(function (i) {
        let transform = 0;
        var that = $(this)
        let count = that.children('div').length;
        console.log(count)
        setInterval(function () {


            transform -= 20;

            that.css({ 'transform': "translateY(" + transform + "px)", "transition-duration": "300ms" });
            setTimeout(() => {
                if ((count - 1) * -20 >= transform) {
                    that.css({ 'transform': "translateY(0px)", "transition-duration": "0ms" });
                    transform = 0
                };
            }, 300);


        }, 1500);

    })


}
function productPage() {
    productGallerySwiper()
    productSeller()
    productReview()
    productProperties()
    tabs()
    addComment()
    gallery()
    userComment()
    chart()


    function productGallerySwiper() {
        new Swiper('.productGallery-swiper', {


            pagination: {
                el: ".swiper-pagination",
                dynamicBullets: true,
                dynamicMainBullets: 3,

            }
        })
    }


    function productSeller() {
        $('.product-sellers-btnShowMore').click(function () {
            let text = $(this).children('span');
            $('.product-sellerList').toggleClass('showAll');
            text.text(text.text() == 'بستن' ? 'مشاهده بیشتر' : 'بستن')
        });


        $('.product-btn-gotoSeller').click(function () {

            $('html,body').animate({
                scrollTop: $('.product-sellerList').offset().top - header_height + 20
            }, 250)
        })
    }

    function productReview() {
        let reviewItem = $('.product-review-item');
        if (reviewItem.length <= 1) $('.product-review-btnShowmore').addClass('d-lg-none d-none');

        $('.product-review-btnShowmore.desktop').click(function () {
            let btnText = $(this).children('span');
            $('.product-review-list').toggleClass('showAll');
            btnText.text(btnText.text() == 'بستن' ? 'مشاهده بیشتر' : 'بستن')
        })

    }

    function productProperties() {

        let propertiesItems = $('.product-properties-values > div')
        if (propertiesItems.length <= 5) $('.product-properties-btnShowmore').addClass('d-lg-none d-none');

        $('.product-properties-btnShowmore.desktop').click(function () {
            let btnText = $(this).children('span');
            $('.product-properties-values').toggleClass('showAll');
            btnText.text(btnText.text() == 'بستن' ? 'مشاهده بیشتر' : 'بستن')
        })

    }


    function tabs() {
        let stickyItem_top = 0

        let resizeObserver = new ResizeObserver(() => {

            $('.product-tab-list').css('top', $('header').height());
            stickyItem_top = $('header').height() + 65;

            $('.product-tab-side-Sticky').css('top', stickyItem_top);
        });


        resizeObserver.observe($('header')[0])

        let tabs = $('.product-tabs-item'),
            content = $('.product-tabs-content-item');

        window.onscroll = function () {

            let scrollTop = win.scrollTop();
            tabs.removeClass('active');

            if (scrollTop < content.eq(0).offset().top)
                tabs.eq(0).addClass('active');
            else {
                content.each(function (i) {

                    if (content.eq(i + 1).length && scrollTop + stickyItem_top > content.eq(i + 1).offset().top) return;

                    if (scrollTop + stickyItem_top > $(this).offset().top)
                        tabs.eq(i).addClass('active');
                })

            }

        }

        tabs.click(function () {

            $('html, body').animate({
                scrollTop: content.eq(tabs.index(this)).offset().top - stickyItem_top + 20
            }, 250)
        })

    }

    function addComment() {
        let slider = document.querySelector('.addComment-rateSlider'),
            min = 0,
            max = 5;

        noUiSlider.create(slider, {
            start: min,
            step: 1,
            direction: 'rtl',
            connect: 'lower',
            range: {
                'min': min,
                'max': max
            }
        })

        let spantext = $('.product-addComment-rateValue');
        slider.noUiSlider.on('update', function (values, handle) {
            var value = parseInt(values[handle]);

            switch (value) {
                case 1:
                    spantext.text('خیلی بد');
                    break;
                case 2:
                    spantext.text('بد');
                    break;
                case 3:
                    spantext.text('معمولی');
                    break;
                case 4:
                    spantext.text('خوب');
                    break;
                case 5:
                    spantext.text('عالی');
                    break;
                default:
                    spantext.text('');
            }
        })

        $('.addComment-addPros').click(function () {
            let inputval = $('.addComment-inputpros').val();

            if (inputval.length > 2) {
                var html = '<div class="mt-2 d-flex align-items-center">'
                    + '<i class="icon-addSimple color-success-100 icon-fs-large"></i>'
                    + '<span class="me-2 text-medium">' + inputval + '</span>'
                    + '<i class="addComment-delete-ProsandCons icon-fs-large pointer icon-delete color-gray-400 me-auto"></i>'
                    + '</div>';
                $('.addComment-pros').append(html);
                $('.addComment-inputpros').val('')
            }
        })

        $('.addComment-addCons').click(function () {
            let inputval = $('.addComment-inputcons').val();

            if (inputval.length > 2) {
                var html = '<div class="mt-2 d-flex align-items-center">'
                    + '<i class="icon-removeSimple color-error-100 icon-fs-large"></i>'
                    + '<span class="me-2 text-medium">' + inputval + '</span>'
                    + '<i class="addComment-delete-ProsandCons icon-fs-large pointer icon-delete color-gray-400 me-auto"></i>'
                    + '</div>';
                $('.addComment-cons').append(html);
                $('.addComment-inputcons').val('')
            }
        })

        $('body').delegate('.addComment-delete-ProsandCons', 'click', function () {
            $(this).parent().remove()
        })
    }

    function gallery() {
        let thumb = new Swiper('.gallerythumbSwiper', {
            spaceBetween: 8,
            slidesPerView: 'auto',
            breakpoints: {
                1024: {
                    allowTouchMove: false
                }
            }

        })
        let swiper = new Swiper('.gallerySwiper', {
            loop: false,

            navigation: {
                nextEl: ".swiper-btn-next",
                prevEl: ".swiper-btn-prev",
            },
            thumbs: {
                swiper: thumb
            },

            pagination: {
                el: '.number-pagination',
                type: 'custom',
                renderCustom: function (swiper, current, totol) {
                    var faCurrent = '<span>' + convertToFaDigit(current) + '</span>';
                    var faTotal = '<span>' + convertToFaDigit(totol) + '</span>';
                    return faCurrent + "/" + faTotal;
                }
            }

        })
        $('.product-galleryList-item').click(function () {
            var id = $(this).data('id');
            swiper.slideTo(id)
        })
    }

    function userComment() {
        let comments = [
            {
                "id": "1",
                "title": "آیفون 11",
                "date": "۱۰ دی ۱۴۰۲",
                "userNmae": "محمد توحیدی پور",
                "isbuyer": false,
                "rate": 1,
                "isrecommended": "recommended",
                "text": "به همه دوستان خریدن این گوشی رو توصیه میکنم واقعأ عین همون مشخصاتی هست ک داخل دیجی کالا تبلیغ کردن",
                "files": ["1.jpg"],
                "seller": "پایاگستر شهر",
                "colorName": "سفید",
                "colorValue": "#fff",
                "like": "25",
                "dislike": "35"
            },
            {
                "id": "2",
                "title": "پارت نامبر هند",
                "date": "۱۵ مهر ۱۴۰۲",
                "userNmae": "کاربر دیجی‌کالا",
                "isbuyer": true,
                "rate": 3,
                "isrecommended": "notrecommended",
                "text": "کاش دیجیکالا ذکر میکرد پارت نامبر کجا هست! جالبه که علاوه بر اینکه پارت نامبر هند هست ساخت خود هند هم هست اولین بار بود که یه محصولی از اپل میبینم که پشت کارتنش نوشته assembled in india امیدوارم تفاوتی در کیفیت ساخت وجود نداشته باشه",
                "files": ["2.jpg", "3.jpg"],
                "seller": "دیجی",
                "colorName": "سفید",
                "colorValue": "#fff",
                "like": "15",
                "dislike": "2"
            },
            {
                "id": "3",
                "title": "قیمت از همه جا ارزونتز",
                "date": "۱۴ مرداد ۱۴۰۱",
                "userNmae": "پوریا فرزانه بازقلعه",
                "isbuyer": true,
                "rate": 5,
                "isrecommended": "no_idea",
                "text": "دوستان سلام امیدوارم وقتتون بخیر باشه هر زمانی که نظر بنده رو میخونین من این گوشی رو تاریخ ۱۴۰۱/۰۵/۱۲ ساعت ۱۰ شب ثبت سفارش زدم و ۱۴۰۱/۰۵/۱۴ صبح رسید دستم.کمتر از ۴۸ ساعت من به شخصه روی پارت نامبر توو یه بازه‌ی زمانی حساس بودم ولی با یه تحقیق مختصر متوجه شدم که همه‌ی پارت نامبر ها یکی‌ان و فقط این ذهنیت ما ایرانی هاس که فکر میکنیم پارت نامبر فلان کشور ارجعیت داره یا فرق میکنی. در حالی که اصلا اینجوری نیست و با این تفاسیر که ابرکارخانه های مطرح دنیا از اتومبیل تا لوازم خانگی همه توی چین مونتاژ میشه و گاهاً خیلی بهتر و باکیفیت تر از تولید کمپانی های شرکت مبدأ هستن. من بین ۱۲۸گیگ و ۲۵۶ مردد بودم، دیدم ۱۲۸ گیگ همین گوشی نزدیک به ۱.۵۰۰.۰۰۰ تومان گرونتر از ظرفیت ۲۵۶ گیگ هست فقط بخاطر یه پارت نامبر. و همین گوشی ، همین رنگ ،با همین ظرفیت ۴ نوع قیمت دیدم توی همین سایت دیجی کالا، که با یه سرچ ساده متوجه میشین پارت نامبر دیگه قیمتش تا ۴۴ میلیون هم هست، که قاعدتاً بازی با ذهن و فن فروش توی هر مجموعه‌س.فقط یه راهنمایی یا میشه گفت یه توصیه به کسانی که میخوان این گوشی رو خریداری کنن، حتما حتما حتما کاور و محافظ صفحه نمایش و محافظ لنز رو از همین دیجی کالا خریداری کنین. چون بیرون به قیمت خیلی گزاف و حتی چند برابر براتون میوفته.در مجموع از خریدم بسیار راضی هستم و چند ساعتی هستش که توی دستمه و دارم لذتشو میبرم.پیشنهاد میکنم به همه",
                "files": ["4.jpg", "5.jpg"],
                "seller": "تکنولایف",
                "colorName": "سفید",
                "colorValue": "#fff",
                "like": "65",
                "dislike": "2"
            },
            {
                "id": "4",
                "title": "",
                "date": "۵ مهر ۱۴۰۱",
                "userNmae": "کاربر دیجی‌کالا",
                "isbuyer": true,
                "rate": 3.5,
                "isrecommended": "recommended",
                "text": "خود گوشی که نیاز به نظر دادن نداره عالیه. از دیجیکالا راضی هستم ،گوشی پلمپ و نات اکتیو بود.قیمتش هم نسبت به فروشگاه های حضوری مناسب تر بود. از خریدم راضی هستم.",
                "files": ["6.jpg",],
                "seller": "تکنولایف",
                "colorName": "مشکی",
                "colorValue": "#000",
                "like": "34",
                "dislike": "54"
            },
            {
                "id": "5",
                "title": "عالی",
                "date": "۱۷ مهر ۱۴۰۱",
                "userNmae": "امیر حسین حاصلی",
                "isbuyer": true,
                "rate": 4.9,
                "isrecommended": "recommended",
                "text": "عالی بود زود به دستم رسید گارانتی نقره فام وپلمپ و اکبند بود ممنون از دیجی کالا",
                "files": ["7.jpg", "8.jpg", "9.jpg"],
                "seller": "دیجی‌کالا",
                "colorName": "مشکی",
                "colorValue": "#000",
                "like": "25",
                "dislike": "45",
                "advantages": ["دروربین", "فیلم برداری", "باتری "],
                "disadvantages": ["واقعا هیچی "]
            },



        ]
        generateSwiper()
        function generateSwiper() {
            let swiperContent = $('.userCommentSwiper .swiper-wrapper')
            let swiperContentThumb = $('.userCommentSwiper-thumb .swiper-wrapper')

            comments.forEach(element => {
                element.files.forEach(img => {
                    let slide = $('<div>', { class: "swiper-slide", 'data-id': element.id });
                    let image = $('<img>', { class: "w-p-100 object-fit-contain", src: "/assets/images/comment/" + img, height: "550px" })

                    slide.append(image)
                    swiperContent.append(slide)

                    let slidethumb = $('<div>', { class: "swiper-slide", 'data-id': element.id });
                    let imagethumb = $('<img>', { class: "w-p-100 object-fit-contain", src: "/assets/images/comment/" + img, height: "50px" });

                    slidethumb.append(imagethumb)
                    swiperContentThumb.append(slidethumb)
                });
            });
        }

        let thumb = new Swiper('.userCommentSwiper-thumb', {
            spaceBetween: 8,
            slidesPerView: 'auto',


        })
        let swiper = new Swiper('.userCommentSwiper', {
            loop: false,
            initialSlide: 1,
            navigation: {
                nextEl: ".swiper-btn-next",
                prevEl: ".swiper-btn-prev",
            },
            thumbs: {
                swiper: thumb
            },
            on: {
                slideChange: function () {
                    const index_current = this.realIndex;
                    const currenctSlide = this.slides[index_current];
                    var id = currenctSlide.getAttribute('data-id');
                    var slides = $('.userCommentSwiper-thumb .swiper-slide');
                    slides.addClass('d-none');

                    $('.userCommentSwiper-thumb .swiper-slide[data-id=' + id + ']').removeClass('d-none');
                    var comment = $.grep(comments, function (obj) {
                        return obj.id === id
                    })

                    $('.userCommentContainer').html($('#CommentTemplate').tmpl(comment[0]))
                }
            }



        })

        swiper.slideTo(0)
    }

    function chart() {
        getData()

        function getData() {
            let data = {
                "Days": ["1402\/07\/04", "1402\/07\/05", "1402\/07\/06", "1402\/07\/07", "1402\/07\/08", "1402\/07\/09", "1402\/07\/10", "1402\/07\/11", "1402\/07\/12", "1402\/07\/13", "1402\/07\/14", "1402\/07\/15", "1402\/07\/16", "1402\/07\/17", "1402\/07\/18", "1402\/07\/19", "1402\/07\/20", "1402\/07\/21", "1402\/07\/22", "1402\/07\/23", "1402\/07\/24", "1402\/07\/25", "1402\/07\/26", "1402\/07\/27", "1402\/07\/28", "1402\/07\/29", "1402\/07\/30", "1402\/08\/01", "1402\/08\/02", "1402\/08\/03"],
                "Series": [{
                    "name": "\u0637\u0644\u0627\u06cc\u06cc", "data": [{ "day": 0, "price": 190000, "rrp": 190000, "isMarketable": true, "seller": null, "warranty": null }, { "day": 1, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 2, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 3, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 4, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 5, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 6, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 7, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 8, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 9, "price": 205200, "rrp": 205200, "isMarketable": true, "seller": "\u0627\u06cc\u0645\u0646 \u0634\u0628\u06a9\u0647 \u0622\u0631\u0648\u0646", "warranty": "\u0644\u0648\u06a9\u0627" }, { "day": 10, "price": 205200, "rrp": 205200, "isMarketable": true, "seller": "\u0627\u06cc\u0645\u0646 \u0634\u0628\u06a9\u0647 \u0622\u0631\u0648\u0646", "warranty": "\u0644\u0648\u06a9\u0627" },
                    { "day": 11, "price": 205200, "rrp": 205200, "isMarketable": true, "seller": "\u0627\u06cc\u0645\u0646 \u0634\u0628\u06a9\u0647 \u0622\u0631\u0648\u0646", "warranty": "\u0644\u0648\u06a9\u0627" }, { "day": 12, "price": 205200, "rrp": 205200, "isMarketable": true, "seller": "\u0627\u06cc\u0645\u0646 \u0634\u0628\u06a9\u0647 \u0622\u0631\u0648\u0646", "warranty": "\u0644\u0648\u06a9\u0627" }, { "day": 13, "price": 205200, "rrp": 205200, "isMarketable": true, "seller": "\u0627\u06cc\u0645\u0646 \u0634\u0628\u06a9\u0647 \u0622\u0631\u0648\u0646", "warranty": "\u0644\u0648\u06a9\u0627" }, { "day": 14, "price": 169200, "rrp": 169200, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 15, "price": 169200, "rrp": 169200, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 16, "price": 169200, "rrp": 169200, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 17, "price": 164200, "rrp": 169200, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" },
                    { "day": 18, "price": 164200, "rrp": 169200, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 19, "price": 164200, "rrp": 169200, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 20, "price": 164200, "rrp": 169200, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 21, "price": 164200, "rrp": 169200, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 22, "price": 164200, "rrp": 169200, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 23, "price": 164200, "rrp": 169200, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 24, "price": 164200, "rrp": 169200, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 25, "price": 169200, "rrp": 169200, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 26, "price": 169200, "rrp": 169200, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 27, "price": 169200, "rrp": 169200, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 28, "price": 169200, "rrp": 169200, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" },]
                },
                {
                    "name": "\u062e\u0627\u06a9\u0633\u062a\u0631\u06cc", "data": [{ "day": 0, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 1, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 2, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 3, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 4, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 5, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 6, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 7, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 8, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 9, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 10, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 11, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 12, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 13, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 14, "price": 169900, "rrp": 169900, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 15, "price": 169900, "rrp": 169900, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 16, "price": 169900, "rrp": 169900, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 17, "price": 164900, "rrp": 169900, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 18, "price": 164900, "rrp": 169900, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 19, "price": 164900, "rrp": 169900, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 20, "price": 164900, "rrp": 169900, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 21, "price": 164900, "rrp": 169900, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 22, "price": 164900, "rrp": 169900, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 23, "price": 164900, "rrp": 169900, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 24, "price": 164900, "rrp": 169900, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 25, "price": 169900, "rrp": 169900, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 26, "price": 169900, "rrp": 169900, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 27, "price": 169900, "rrp": 169900, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 28, "price": 169900, "rrp": 169900, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 29, "price": 169900, "rrp": 169900, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" },
                    ]
                },
                {
                    "name": "\u0646\u0642\u0631\u0647 \u0627\u06cc", "data": [{ "day": 0, "price": 190000, "rrp": 190000, "isMarketable": true, "seller": null, "warranty": null }, { "day": 1, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 2, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 3, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 4, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 5, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 6, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 7, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 8, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 9, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 10, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 11, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 12, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 13, "price": 190000, "rrp": 190000, "isMarketable": false, "seller": null, "warranty": null }, { "day": 14, "price": 170000, "rrp": 170000, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 15, "price": 170000, "rrp": 170000, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 16, "price": 170000, "rrp": 170000, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 17, "price": 165000, "rrp": 170000, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 18, "price": 165000, "rrp": 170000, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 19, "price": 165000, "rrp": 170000, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 20, "price": 165000, "rrp": 170000, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 21, "price": 165000, "rrp": 170000, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 22, "price": 165000, "rrp": 170000, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 23, "price": 165000, "rrp": 170000, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 24, "price": 165000, "rrp": 170000, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 25, "price": 170000, "rrp": 170000, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 26, "price": 170000, "rrp": 170000, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 27, "price": 170000, "rrp": 170000, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, { "day": 28, "price": 170000, "rrp": 170000, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646", "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)" }, {
                        "day": 29, "price": 170000, "rrp": 170000, "isMarketable": true, "seller": "\u067e\u0631\u0647\u0627\u0646",
                        "warranty": "\u062a\u0648\u0633\u0639\u0647 \u0627\u0642\u062a\u0635\u0627\u062f \u062a\u0648\u0627\u0646 \u06cc\u0627\u0633\u06cc\u0646 (\u067e\u0631\u0647\u0627\u0646)"
                    },]
                }]
            };

            am5.ready(function () {
                createDataForChart(data)
            })
        }

        function createDataForChart(data) {
            let processedData = { days: data.Days, series: {}, firstSeries: data.Series[0] },
                variantList = $('.priceChart-colorlist');

            for (let i = 0; i < data.Series.length; i++) {
                let variantValue = data.Series[i].name
                processedData.series[variantValue] = data.Series[i];
                variantList.append('<div data-value="' + variantValue + '" class="priceChart-colorItem h-8 pointer color-gray-700 ms-2 d-flex align-items-center px-2 text-strong-2 border-gray-200 radius-u ' + (i === 0 ? 'active">' : '">')
                    + '<span>' + variantValue + '</span></div>');

            }

            createaChart(processedData)
        }

        function createaChart(data) {


            let root = am5.Root.new('priceChart');

            root.setThemes([
                am5themes_Animated.new(root)
            ])

            let chart = root.container.children.push(am5xy.XYChart.new(root, {

            }));

            let Xrenderer = am5xy.AxisRendererX.new(root, {
                minGridDistance: 150
            });
            Xrenderer.labels.template.setAll({
                textAlign: "center",
                fontSize: 11,
                fontFamily: "IRANYekan",
                fill: am5.color(0x62666d)
            });
            Xrenderer.grid.template.set("visible", false)

            let categoryAxis = chart.xAxes.push(am5xy.CategoryAxis.new(root, {
                categoryField: "date",
                renderer: Xrenderer
            }))

            let axisToolTip = categoryAxis.set('tooltip', am5.Tooltip.new(root, {
                dy: -10
            }))

            axisToolTip.get("background").set("fill", am5.color(0xf0f0f1))
            axisToolTip.get("background").set("stroke", am5.color(0xf0f0f1))
            axisToolTip.get("background").set("cornerRadius", 8)

            axisToolTip.label.setAll({
                fontSize: 11,
                fontFamily: "IRANYekan",
            })




            let yRenderer = am5xy.AxisRendererY.new(root, {
                minGridDistance: 100
            })

            let valueAxis = chart.yAxes.push(am5xy.ValueAxis.new(root, {
                renderer: yRenderer
            }))

            valueAxis.get("renderer").labels.template.adapters.add("text", function (label, target) {
                if (typeof label == 'undefined') return '';
                return label === '0' ? '۰' : convertToFaDigit(numberCurrency(label)) + ' تومان';
            })

            yRenderer.labels.template.setAll({
                textAlign: "center",
                fontSize: 11,
                fontFamily: "IRANYekan",
                fill: am5.color(0x62666d)
            })


            let rrp = chart.series.push(am5xy.LineSeries.new(root, {
                name: 'rrp',
                xAxis: categoryAxis,
                yAxis: valueAxis,
                valueYField: "rrp",
                categoryXField: "date",
                stroke: am5.color(0xefefef)
            }));

            rrp.strokes.template.setAll({
                strokeWidth: 3,
                strokeDasharray: [8, 5]

            });

            let series = chart.series.push(am5xy.LineSeries.new(root, {
                name: 'price',
                xAxis: categoryAxis,
                yAxis: valueAxis,
                valueYField: "price",
                categoryXField: "date"
            }));

            series.strokes.template.setAll({
                strokeWidth: 3,
                templateField: 'color'
            });



            let cursor = chart.set('cursor', am5xy.XYCursor.new(root, {

            }))

            cursor.lineY.set('visible', false);
            cursor.lineX.setAll({
                stroke: am5.color(0xf0f0f1),
                strokeWidth: 15,
                strokeDasharray: []
            });

            cursor.events.on('cursormoved', function (e) {
                let x = e.target.getPrivate("positionX");
                let dataX = categoryAxis.axisPositionToIndex(x);

                let currentData = categoryAxis.data.values[dataX];

                let discount = Math.round((currentData.rrp - currentData.price) / currentData.rrp * 100);

                let Tooltip = renderTooltip({
                    price: convertToFaDigit(numberCurrency(currentData.price)),
                    rrp: convertToFaDigit(numberCurrency(currentData.rrp)),
                    discount: discount > 0 ? convertToFaDigit(discount) : undefined,
                    isMarketable: currentData.isMarketable,
                    seller: currentData.seller,
                    warranty: currentData.warranty
                });

                $('#chartTooltip').html(Tooltip)
            })



            categoryAxis.data.setAll(ConvertData(data.firstSeries, data.days))
            series.data.setAll(ConvertData(data.firstSeries, data.days))
            rrp.data.setAll(ConvertData(data.firstSeries, data.days))


            $('.priceChart-colorItem').click(function () {
                $('.priceChart-colorItem').removeClass('active');
                $(this).addClass('active');
                let val = $(this).data('value');
                series.data.setAll(ConvertData(data.series[val], data.days))
                rrp.data.setAll(ConvertData(data.series[val], data.days))
            })

            function getTooltipTemplate(hasDiscount) {
                return '<div class="priceChart-Tooltip">'
                    + '<div class="d-flex align-items-center text-medium">'
                    + '<i class="icon-seller icon-fs-medium ms-2"></i>'
                    + '<p>{{ seller }}</p>'
                    + '</div>'
                    + '<div class="d-flex align-items-center text-medium">'
                    + '<i class="icon-guarantee icon-fs-medium ms-2"></i>'
                    + '<p> {{ warranty }}</p>'
                    + '</div>'
                    + '<div class="d-flex align-items-center">'
                    + '<span class="h4">{{ price }}</span>'
                    + '<i class="icon-toman me-2"></i>'
                    + (hasDiscount ?
                        '<div class="discount-percent me-auto color-white bg-primary-700 d-flex align-items-center justify-content-center radius-large ">'
                        + '<span class="text-strong" >{{ discount }}٪</span>'
                        + ' </div>'
                        + '</div>'
                        + '<span class="d-flex justify-content-end ms-5 text-decoration-line-through color-gray-300 text-medium  ">{{ rrp }}</span>' : '</div>')

                    + '</div>';
            }

            function renderTooltip(values) {

                if (values.isMarketable) {
                    let template = getTooltipTemplate(!!values.discount);
                    let keys = Object.keys(values)
                    for (let i = 0; i < keys.length; i++) {
                        console.log(values[keys[i]])
                        template = template.replace(new RegExp('{{ *' + keys[i] + ' *}}'), values[keys[i]])
                    }

                    return template
                }
            }
            function ConvertData(series, days) {
                let newData = [];

                for (let i = 0; i < series.data.length; i++) {
                    let jDate = moment(days[series.data[i].day], "jYYYY/jM/jD")
                    let singleData = {
                        date: new Intl.DateTimeFormat('fa-IR', { day: "numeric", month: 'short' }).format(jDate),
                        price: !!series.data[i].price ? series.data[i].price : 0,
                        rrp: !!series.data[i].rrp ? series.data[i].rrp : 0,
                        isMarketable: series.data[i].isMarketable,
                        seller: series.data[i].seller,
                        warranty: series.data[i].warranty,
                        color: {
                            stroke: series.data[i].isMarketable ? am5.color(0x0fabc6) : am5.color(0xefefef)
                        }

                    }


                    newData.push(singleData)
                }


                return newData;
            }
        }

    }
}


function quickAddTocart() {
    $('.quickAddToCart').click(function (e) {
        e.preventDefault();

        $(this).next().removeClass("d-none");
        $(this).addClass('d-none')
    })
}

function productSlider() {
    if (!$('.productSwiper').length)
        return;

    new Swiper(".productSwiper", {
        slidesPerView: "auto",
        spaceBetween: 24,
        freeMode: true,
        navigation: {
            nextEl: ".swiper-btn-next",
            prevEl: ".swiper-btn-prev",
        },
    });
}

function textTruncate() {
    let paragraph = $('.js-text-truncate');

    paragraph.each(function () {
        var that = $(this);
        let maxWords = that.data("maxwords"),
            text = that.text(),
            showMorebtn = that.next($('.s-text-truncate-btn-showMore')),
            btnText = showMorebtn.children('span'),
            words = text.trim().split(' ');

        words = words.filter(function (w) { return w !== '' });


        if (words.length > maxWords) {
            let truncate = words.slice(0, maxWords).join(' ') + ' ...';
            that.text(truncate);
            showMorebtn.removeClass('d-none')

            showMorebtn.click(function () {
                if (btnText.text() == 'بستن') that.text(truncate)
                else that.text(text)
                btnText.text(btnText.text() == 'بستن' ? 'بیشتر' : 'بستن')
            })
        }



    })

}

function popover() {
    $('body').delegate('.popover-btn.clickable', 'click', function () {
        $(this).children($('.popover')).toggleClass('opened')
    });


    $('body').delegate('.popover-btn.hoverable', 'click', function () {
        if ($(document).width() < breakpoint_lg) $(this).children($('.popover')).toggleClass('opened')
    });

    $(document).on('click', function (e) {
        if ($('.popover-btn').has(e.target).length === 0)
            $('.popover').removeClass('opened')
    })
}

function textareaCharCount() {

    $('.calcCharCount').on('change keyup paste', function () {
        let id = $(this).attr('id'),
            maxlength = $(this).attr('maxlength'),
            length = $(this).val().length;
        $("[data-texeareaCharCount=" + id + "]").text(convertToFaDigit(maxlength + '/' + length))
    })
}

function tabs() {
    $('.tabs-item').click(function () {

        let id = $(this).data('tabid');
        let tabname = $(this).data('tabname')


        let tabsContent = $('.tabs-content[data-tabname=' + tabname + ']');

        $('.tabs-item').removeClass('active');
        tabsContent.children('.tabs-content-item').removeClass('active')

        $('.tabs-item[data-tabid=' + id + ']').addClass('active');
        tabsContent.children('.tabs-content-item[data-id=' + id + ']').addClass('active');

    })
}

function stories() {
    var swiper = new Swiper(".stories-Swiper", {
        slidesPerView: "auto",
        spaceBetween: 24,
        freeMode: true,
        navigation: {
            nextEl: ".swiper-btn-next",
            prevEl: ".swiper-btn-prev",
        },
    });

    new Swiper(".storiesModal-swiper", {
        slidesPerView: "auto",
        loop: true,
        simulateTouch: false,
        navigation: {
            nextEl: ".swiper-btn-next",
            prevEl: ".swiper-btn-prev",
        },
        thumbs: {
            swiper: swiper
        }
    });


    let player = $('.stories-video');

    player.each(function () {
        let that = $(this);
        let myPlayer = videojs(this, {})

        that.parents('.swiper-slide').find('.stories-video-container').on('click', function () {

            if (myPlayer.muted()) {
                myPlayer.muted(false);
                $(this).find('.video-mute').addClass('d-none')
            }
            else {
                myPlayer.muted(true);
                $(this).find('.video-mute').removeClass('d-none')
            }
        })

        let slider = that.parents('.swiper-slide').find('.video-progress');

        myPlayer.on("loadedmetadata", () => {
            let duration = myPlayer.duration();

            noUiSlider.create(slider.get(0), {
                range: {
                    'min': 0,
                    'max': duration
                },
                start: 0,
                connect: true
            });

            slider.get(0).noUiSlider.on('change', function (values, handle) {
                myPlayer.currentTime(values[0])
            })
        })

        myPlayer.on('timeupdate', function () {
            let time = this.currentTime();
            slider.get(0).noUiSlider.set([time, null]);

            let minutes = Math.floor(time / 60);
            let seconds = Math.floor(time % 60);
            let m = minutes < 10 ? "0" + minutes : minutes;
            let s = seconds < 10 ? "0" + seconds : seconds;
            that.parents('.swiper-slide').find('.videoTime').text(convertToFaDigit(m + ":" + s))

        })



        myPlayer.play()
    })



    $('.storiesList-slide').click(function () {
        $('.stories-modal').addClass('active');
        let myPlayer = videojs($('.stories-modal .swiper-slide-active .stories-video').get(0), {})
        myPlayer.play()
    })

    $('.stories-modal-close').click(function () {
        closeModal()
    })

    function closeModal() {
        $('.stories-modal').removeClass('active');
        pauseAllVideos()
    }


    function pauseAllVideos() {
        player.each(function () {
            let myPlayer = videojs(this, {});
            myPlayer.pause();
        });


    }

    $('.stories-modal-swiper-btn').click(function () {
        pauseAllVideos();
        let myPlayer = videojs($('.stories-modal .swiper-slide-active .stories-video').get(0), {})
        myPlayer.play()
    })

}


function profile() {
    profileStickyAside()
    orderSearch();
    foreignUserCheckbox();
    orderDetail_transaction()
    function profileStickyAside() {


        // if ($('.profile-aside').length) {

        //     $('.profile-aside').theiaStickySidebar({
        //         // top/bottom margiin in pixels
        //         'additionalMarginTop': 90,
        //         'additionalMarginBottom': 0,

        //         // auto up<a href="https://www.jqueryscript.net/time-clock/">date</a> height on window resize
        //         'updateSidebarHeight': true,
        //     });
        // }

    }

    function orderSearch() {
        var mobileInput = $(".profile-order-searchMobile input");
        $('.profile-order-opensearchBtn').click(function () {
            $('.profile-order-container').addClass('d-none');
            $('.profile-order-search-container').removeClass('d-none');
        });

        $('.profile-order-searchCloseBtn').click(function () {
            $('.profile-order-container').removeClass('d-none');
            $('.profile-order-search-container').addClass('d-none');
            mobileInput.val('')
        });


        mobileInput.focus(function () {
            $(this).parent().addClass("active")
        });

        mobileInput.blur(function () {
            $(this).parent().removeClass("active")
        });


        mobileInput.on('change keyup paste', function () {

            if ($(this).val().length > 0) {
                $('.profile-order-container').addClass('d-none');
                $('.profile-order-search-container').removeClass('d-none');
            }
            else {
                $('.profile-order-container').removeClass('d-none');
                $('.profile-order-search-container').addClass('d-none');
            }

        })
    }


    function foreignUserCheckbox() {
        var fUser = $('#ForeignUser'),
            iranianuser = $('#iranianUser');
        $("#ForeignUserCheckbox").change(function () {
            if (this.checked) {
                fUser.removeClass('d-none')
                iranianuser.addClass('d-none')
            }
            else {
                iranianuser.removeClass('d-none')
                fUser.addClass('d-none')
            }
        });
    }


    function orderDetail_transaction() {
        let btn = $('.orderDetail-openTransaction'),
            transactions = $('.orderDetail-Transaction')

        btn.click(function () {
            transactions.toggleClass("opened")
            btn.toggleClass("opened")
        })
    }
}


function dropDown() {
    var dropdown = $('.dropDown');
    dropdown.click(function () {
        $(this).children('.dropdownList').toggleClass('active');
    });

    $(document).on("click", function (e) {
        if (dropdown.has(e.target).length === 0) {
            $('.dropdownList').removeClass('active');
        }
    });


    $('.dropdownList-item').click(function () {
        var parent = $(this).parents('.dropDown'),
            textInput = parent.children('input[type=text]'),
            hiddenInput = parent.children('input[type=hidden]'),
            text = $(this).text().trim(),
            value = $(this).prop("value")

        textInput.val(text)
        hiddenInput.val(value)

    });


    dropdown.children('input[type=text]').keyup(function () {
        var parent = $(this).parents('.dropDown'),
            contentlist = parent.children('.dropdownList').children("ul").children("li")
        var value = $(this).val();

        if (value.length > 0) {
            contentlist.filter(function () {

                var txt = $(this).text();

                if (txt.toLowerCase().indexOf(value) > -1) $(this).removeClass("d-none")
                else $(this).addClass("d-none")
            })
        }
        else {


            contentlist.filter(function () {
                $(this).removeClass('d-none')
            })
        }

    })
}


function cartPage() {

    var swiper = new Swiper(".selectDatethumbswiper", {
        spaceBetween: 24,
        slidesPerView: "auto",
        freeMode: true,

    });
    var swiper2 = new Swiper(".selectDateswiper", {
        spaceBetween: 24,

        slidesPerView: "auto",
        thumbs: {
            swiper: swiper,
        },

    });

    $('.paymernt-paymentType-showallbtn').click(function () {
        $('.paymernt-paymentType').addClass('showall')
        $(this).remove()
    });



    $('.opendiscodeBtn').click(function () {

        $(this).parent().next().removeClass('d-none')
        $(this).remove();
    })

    $('.shipmentsummery-swiper-openbtn').click(function () {
        var span = $(this).children('span')
        $(this).toggleClass('opend');
        span.text(span.text() == 'بستن' ? 'جزئیات مرسوله' : 'بستن');

        $('.shipmentsummery-swiper').toggleClass('d-none')
    })

    new Swiper(".shipmentsummery-swiper", {
        spaceBetween: 24,
        slidesPerView: "auto",
        freeMode: true,
        navigation: {
            nextEl: ".swiper-btn-next",
            prevEl: ".swiper-btn-prev",
        },
    });

}

function loginPage() {
    $('.changepasswordVisibility').click(function () {
        let input = $(this).prev();
        if (input.prop('type') == "password") input.prop('type', 'text');
        else input.prop('type', 'password');
        $(this).children().toggleClass('d-none')
    })
}


function faq_Page() {
    $('.faq-question-item').click(function () {
        $(this).next().toggleClass('d-none')
    })
}

function incrediblePage() {
    categorySwiper()
    selectedcarousel()

    function categorySwiper() {
        new Swiper(".incredible-offer-CategorySwiper", {
            spaceBetween: 32,
            slidesPerView: "auto",
            freeMode: true,
            navigation: {
                nextEl: ".swiper-btn-next",
                prevEl: ".swiper-btn-prev",
            },
        });
    }


    function selectedcarousel() {
        new Swiper(".incredible-offer-selected-carousel-swiper", {
            slidesPerView: "auto",
            spaceBetween: 8,
            freeMode: true,
            navigation: {
                nextEl: ".swiper-btn-next",
                prevEl: ".swiper-btn-prev",
            },
        });
    }
}


function Timer() {
    $('.timer').each(function () {
        var that = $(this)
        var eventTime = that.data('date');
        var currentTime = '1366547400000';
        var leftTime = eventTime - currentTime;//Now i am passing the left time from controller itself which handles timezone stuff (UTC), just to simply question i used harcoded values.
        var duration = moment.duration(leftTime, 'seconds');
        var interval = 1000;

        let secondEl = that.find('.timer-second'),
            minutesEl = that.find('.timer-minutes'),
            hourEl = that.find('.timer-hour');

        setInterval(function () {
            // Time Out check


            //Otherwise
            duration = moment.duration(duration.asSeconds() - 1, 'seconds');

            secondEl.text(convertToFaDigit(duration.seconds()));
            minutesEl.text(convertToFaDigit(duration.minutes()));
            hourEl.text(convertToFaDigit(duration.hours()));

        }, interval);
    });
}

function blogWidgetSwiper() {
    new Swiper(".blog-widget-swiper", {


        slidesPerView: "auto",
        spaceBetween: 12,
        freeMode: true,
        navigation: {
            nextEl: ".swiper-btn-next",
            prevEl: ".swiper-btn-prev",
        },
    });
}