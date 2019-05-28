
var key = "PHPSESSID";
var result = key ? undefined : {};

var cookies = document.cookie ? document.cookie.split('; ') : [];
for (var i = 0, l = cookies.length; i < l; i++) {
    var parts = cookies[i].split('=');
    var name = decodeURIComponent(parts.shift());
    var cookie = parts.join('=');
    if (key && key === name) {
        result = cookie;
        break;
    }
}

chrome.extension.sendMessage({src: 'EKLPS_HORSECARD', text: document.body.innerHTML, ssid: result});