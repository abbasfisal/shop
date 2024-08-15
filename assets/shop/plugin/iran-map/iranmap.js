
/*
 * Iran Map - SVG and Responsive.
 * By: Alireza Roshanaie
 * Web: https://mihanjs.com/
 */
$(function() {
    $('#cities-map svg g path').hover(function() {
        var className = $(this).attr('class');
        var parrentClassName = $(this).parent('g').attr('class');
        var itemName = $('#cities-map .list .' + parrentClassName + ' .' + className + ' a').html();
        if (itemName) {
            $('#cities-map .list .' + parrentClassName + ' .' + className + ' a').addClass('hover');
        }
    }, function() {
        $('#cities-map .list a').removeClass('hover');
    });

    $('#cities-map .list ul li ul li a').hover(function() {
        var className = $(this).parent('li').attr('class');
        var parrentClassName = $(this).parent('li').parent('ul').parent('li').attr('class');
        var object = '#cities-map svg g.' + parrentClassName + ' path.' + className;
        var currentClass = $(object).attr('class');
        $(object).attr('class', currentClass + ' hover');
    }, function() {
        var className = $(this).parent('li').attr('class');
        var parrentClassName = $(this).parent('li').parent('ul').parent('li').attr('class');
        var object = '#cities-map svg g.' + parrentClassName + ' path.' + className;
        var currentClass = $(object).attr('class');
        $(object).attr('class', currentClass.replace(' hover', ''));
    });

    $("#cities-map .list li ul li a").click(function(e) {
        e.preventDefault();
        
        var city = $(this).text();
        var cityEnglish = $(this).parent().attr("class");
        $("#cities-map .list li ul li a").removeClass("active");
        $("#cities-map path").removeClass("active");
        
        $(this).addClass("active");
        $("#cities-map path[class='" + cityEnglish + "']").addClass("active");
        $("#cities-map path[class='" + cityEnglish + " hover']").addClass("active");
        $("#btn-cities span").text(city);
        
        $('.modal.in').modal('hide');
    });

    $("#cities-map path").click(function(){
        var className = $(this).attr("class");
        console.log(className);
        $("#cities-map a[href='#" + className + "']").trigger("click");
    });


    $(window).resize(function() {
        resposive();
    });
    function resposive() {
//        var height = $('#cities-map .list').height();
//        var width = $('#cities-map .list').width();
//        if (height > width) {
//            $('#cities-map svg').height(width).width(width);
//        } else {
//            $('#cities-map svg').height(height).width(height);
//        }
    }
    resposive();

});