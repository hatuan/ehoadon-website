function google_maps_58af1ee65f2cb() {
    var latlng = new google.maps.LatLng(21.004560, 105.804807);
    var myOptions = {
        zoom : 17,
        center : latlng,
        mapTypeId : google.maps.MapTypeId.ROADMAP,
        zoomControl : true,
        mapTypeControl : false,
        streetViewControl : false,
        scrollwheel : false
    };
    var displayID = document.getElementById("google-map-area-58af1ee65f2cb");
    if(displayID != undefined) {
        var marker = new google.maps.Marker({
            position : latlng,
            map :  new google.maps.Map(displayID, myOptions)
        });
    }
}
 
jQuery(document).ready(function($) {
    google_maps_58af1ee65f2cb();

    //show #mfn-rev-slider if in home
    if(window.location.pathname == "/" || window.location.pathname == "/index" || window.location.pathname == "/index.html") {
        $("#mfn-rev-slider").show();
    } else {
        $("#mfn-rev-slider").hide();
    }

    //Validation RegisterForm
    $("form[name='registerform']").validate({
        rules: {
            Description: "required",
            VatNumber: "required",
            CompanyAddress: "required",
            AddressTransition:"required",
            Telephone: "required",
            Email: {
                required: true,
                email: true
            },
            Website: {
                required: false,
                url: true
            },
            ContactName:"required",
            Mobile: "required",
            CaptchaInput: "required",
        },
        messages: {
            Description: "Tên công ty không được để trống",
            VatNumber: "Mã số thuế không được để trống",
            CompanyAddress: "Địa chỉ công ty không được để trống",
            AddressTransition: "Địa chỉ làm việc công ty không được để trống",
            Telephone: "Điện thoại không được để trống",
            Email: {
                required: "Email không được để trống",
                email: "Không đúng định dạng Email"
            },
            Website: {
                url: "Không đúng định dạng URL, cập nhật theo định dạng http://..."
            },
            ContactName: "Người liên lạc không được để trống",
            Mobile: "Mobile không được để trống",
            CaptchaInput: "Mã xác nhận không được để trống"
        },
        submitHandler: function(form) {
          form.submit();
        }
    });
});

var tpj = jQuery;
var revapi1;
tpj(document).ready(function() {
    if (tpj("#rev_slider_1_1").revolution == undefined) {
        revslider_showDoubleJqueryError("#rev_slider_1_1");
    } else {
        revapi1 = tpj("#rev_slider_1_1").show().revolution({
            sliderType: "hero",
            sliderLayout: "auto",
            dottedOverlay: "none",
            delay: 9000,
            navigation: {},
            visibilityLevels: [1240, 1024, 778, 480],
            gridwidth: 1240,
            gridheight: 760,
            lazyType: "none",
            parallax: {
                type: "scroll",
                origo: "enterpoint",
                speed: 400,
                levels: [5, 10, 15, 20, 25, 30, 35, 40, 45, 46, 47, 48, 49, 50, 51, 55],
            },
            shadow: 0,
            spinner: "spinner2",
            autoHeight: "off",
            disableProgressBar: "on",
            hideThumbsOnMobile: "off",
            hideSliderAtLimit: 0,
            hideCaptionAtLimit: 0,
            hideAllCaptionAtLilmit: 0,
            debugMode: false,
            fallbacks: {
                simplifyAll: "off",
                disableFocusListener: false,
            }
        });
    }
});

jQuery(window).load(function() {
    var retina = window.devicePixelRatio > 1 ? true : false;
    if (retina) {
        var retinaEl = jQuery("#logo img.logo-main");
        var retinaLogoW = retinaEl.width();
        var retinaLogoH = retinaEl.height();
        retinaEl.attr("src", "content/betheme/images/logo-retina.png").width(retinaLogoW).height(retinaLogoH);
        var stickyEl = jQuery("#logo img.logo-sticky");
        var stickyLogoW = stickyEl.width();
        var stickyLogoH = stickyEl.height();
        stickyEl.attr("src", "content/betheme/images/logo-retina.png").width(stickyLogoW).height(stickyLogoH);
        var mobileEl = jQuery("#logo img.logo-mobile");
        var mobileLogoW = mobileEl.width();
        var mobileLogoH = mobileEl.height();
        mobileEl.attr("src", "content/betheme/images/logo-retina.png").width(mobileLogoW).height(mobileLogoH);
        var mobileStickyEl = jQuery("#logo img.logo-mobile-sticky");
        var mobileStickyLogoW = mobileStickyEl.width();
        var mobileStickyLogoH = mobileStickyEl.height();
        mobileStickyEl.attr("src", "content/betheme/images/logo-retina.png").width(mobileStickyLogoW).height(mobileStickyLogoH);
    }
});

function base64ToArrayBuffer(base64) {
    var binary_string =  window.atob(base64);
    var len = binary_string.length;
    var bytes = new Uint8Array( len );
    for (var i = 0; i < len; i++)        {
        bytes[i] = binary_string.charCodeAt(i);
    }
    return bytes.buffer;
}

function arrayBufferToBase64(buffer) {
    var binary = '';
    var len = buffer.length;
    for (var i = 0; i < len; i++) {
        binary += String.fromCharCode(buffer[i]);
    }
    return window.btoa( binary );
}

function searchInvoice() {
    $("#input-err").text("");
    jQuery.ajax("/invoices",
        {
            dataType: 'json',
            data: {
                "s": jQuery("#SearchID").val(),
                "CaptchaInput": jQuery("#CaptchaInput").val(),
                "CaptchaID" : jQuery("#CaptchaID").val(),
            }
        }
    ).done(function(response, status) {
        download("data:application/pdf;base64,"+response, jQuery("#SearchID").val() + ".pdf", "application/pdf");
        jQuery("#CaptchaInput").val("");
    }).fail(function( jqXHR, textStatus ) {
        //TODO : Thong bao sai captcha, chuoi xac thuc sai (file khong co ...)
        var response =  JSON.parse(jqXHR.responseText);
        if(response.data) {
            recaptcha();
            $("#CaptchaInput").val("");
            if(response.data == "FileNotFound")
                $("#input-err").text("Không tìm thấy file - Chuỗi xác thực hóa đơn bị sai !");
        }
    });
}

/* ---------------------------------------------------------------------------
* Header Search
* --------------------------------------------------------------------------- */
function recaptcha() {
    jQuery.ajax("/captcha",  {
        async: false
    })
    .done(function(response, status) {
            $("#CaptchaImg").attr('src', "");
            setTimeout(function(){
                $("#CaptchaImg").attr('src', "/captcha/" + response.captchaid +".png");
                $("#CaptchaID").val(response.captchaid);
        }, 0);
    })
    .fail(function(response, status) {
        //TODO : Thong bao sai captcha, chuoi xac thuc sai (file khong co ...)
    });
}

/* ---------------------------------------------------------------------------
* Header Search
* --------------------------------------------------------------------------- */
jQuery("#search_button_new").click(function(e) {
    e.preventDefault();

    var faceOff =  $("#Top_bar .search_wrapper").css('display')
    if(faceOff == "none") {
        recaptcha();
        jQuery('#Top_bar .search_wrapper').fadeToggle();
    } else {
        jQuery('#Top_bar .search_wrapper').fadeOut();
    }
});