var eklpsAdviser = {

    pageBody: "",

    horseId: 0,

    ssid: '',

    request : {
        Age : 0,
        Sex : '',
        Spec: '',
        Distances: [],
        Filters: []
    },

    sendReq : function()
    {
        $("#main").append("<label>Запрос к сервису скачек</label>");
        var url = 'http://194.135.95.211:8011/getraces?data='+JSON.stringify(eklpsAdviser.request);
        //var url = 'http://127.0.0.1:8011/getraces?data='+JSON.stringify(eklpsAdviser.request);
        eklpsAdviser.requestData(url);
    },

    requestData: function(request) {
        var req = new XMLHttpRequest();

        req.open("GET", request, true);

        req.onreadystatechange = function() {
            if (req.readyState == 4) {
                // innerText does not let the attacker inject HTML elements.
                var races = {};
                try {
                    races = JSON.parse(req.responseText);
                }
                catch (e) {
                    $("#main").html("<label>Ошибка на сервере скачек:<br><span style='color: #f38245;'> HTTP Status: "+req.statusText+"; Response: "+req.responseText+"</span></label>");
                    return false;
                }
                eklpsAdviser.showData(races);
            }
        };

//        req.onload = this.showData_.bind(this);
        req.send(null);
    },

    showData: function (races) {

        if(races == null) {
            $("#main").html("<label>Подходящих скачек не найдено</label>");
            return true;
        }

        $("#main").html("<label>Подходящие скачки</label>");
        var list = '<table style="width: 520px;"><tr><td style="width: 130px; font-weight: bold;">Дата</td>' +
            '<td style=" font-weight: bold;">Скачка</td><td style="width: 30px; font-weight: bold;">Ограничение</td><td style="width: 50px; font-weight: bold;"></td></tr>';
        for(var i=0; i<races.length; i++) {
            list += '<tr>';
            list += '<td style="width: 130px;">'+races[i].RaceDate+'</td><td>'+races[i].RaceTitle+'</td>';
            list += '<td style="width: 30px;">'+races[i].RaceClass+'</td>';
            list += '<td style="width: 50px;"><span class="fakelink submitter" raceId="'+races[i].RaceId+'">Записать</span></td>';
            list += '</tr>';
        }
        list += '</table>';
        $("#main").append(list);

        eklpsAdviser.raceFinalSend();

        return true;
    },

    showData_: function (e) {
        alert(e.target);
        //var kittens = e.target.responseXML.querySelectorAll('photo');
        //for (var i = 0; i < kittens.length; i++) {
        //    var img = document.createElement('img');
        //    img.src = this.constructKittenURL_(kittens[i]);
        //    img.setAttribute('alt', kittens[i].getAttribute('title'));
        //    document.body.appendChild(img);
        //}
    },

    run : function()
    {
        $("#main").html("<label>Откройте, пожалуйста, список отделения</label>");

        chrome.tabs.query({active: true}, function(tabs){
            if(/http:\/\/(www\.)?eklps\.com\/eclipse\/stable/.test(tabs[0].url)) {

                //var matches = /http:\/\/(www\.)?eklps\.com\/eclipse\/stable/.exec(tabs[0].url);
                //if(matches[2] != undefined && matches[2] != 'undefined') {
                //    eklpsAdviser.horseId = matches[2];
                //}
                //else {
                //    return false;
                //}

                $("#main").html("<label>Поиск подходящих скачек...</label>");
                var requestBody = eklpsAdviser.prepareRequest();
            }
        });

    },

    refresh : function() {
        $("#main").html("Данные обновляются");
        eklpsAdviser.setFilters();
        eklpsAdviser.sendReq();
    },

    prepareRequest : function() {
        var request = '';

        $("#main").append("<label>Подготовка запроса...</label>");
        chrome.extension.onMessage.addListener(function(req){
            if(req.src == "EKLPS_HORSECARD") {
                eklpsAdviser.pageBody = req.text;
                res = eklpsAdviser.parseBody();
                if(res) {
                    eklpsAdviser.sendReq();
                }
                eklpsAdviser.ssid = req.ssid;
            }
        });

        chrome.tabs.executeScript(null, {file: "js/content_script.js"}, function(result) {

        });

        return request;
    },

    setFilters : function() {
        var filterSet = [];
        $(".class-selector").each(function() {
            if($(this).hasClass("btn-down")) {
                filterSet.push($(this).attr("data"));
            }
        });
        eklpsAdviser.request.Filters = filterSet;
    },

    parseBody : function()
    {
        var domBody = $(eklpsAdviser.pageBody);
        eklpsAdviser.horseId = parseInt(domBody.find("#cboxLoadedContent strong").eq(0).html());
        eklpsAdviser.request.Age = parseInt(domBody.find("#cboxLoadedContent tr").eq(1).find('td').eq(1).html());
        eklpsAdviser.request.Sex = domBody.find("#infoblocktbl tr").eq(3).find('td').eq(1).html()[0];

        //Обе таблицы
        for(var t=3; t<=4; t++) {
            var distContainer = domBody.find("#infoblockheader").eq(t).parent();

            distContainer.find("#charter").each(function (obj) {

                //Нам не нужны дистанции без классов
                if ((/-+/.test(distContainer.find("#charter").eq(obj).children().eq(1).html())
                    && /-+/.test(distContainer.find("#charter").eq(obj).children().eq(2).html()))
                    || /[^0-9]+/.test(distContainer.find("#charter").eq(obj).children().eq(0).html())
                ) {
                    return;
                }

                eklpsAdviser.request.Distances.push({
                    Distance: distContainer.find("#charter").eq(obj).children().eq(0).html(),
                    Fl: (/-+/.test(distContainer.find("#charter").eq(obj).children().eq(1).html()) ? '' : distContainer.find("#charter").eq(obj).children().eq(1).html()),
                    Sc: (/-+/.test(distContainer.find("#charter").eq(obj).children().eq(2).html()) ? '' : distContainer.find("#charter").eq(obj).children().eq(2).html())
                });
            });
        }
        $("#main").append("<label>Вытащил данные о лошади</label>");
        return true;
    },

    raceFinalSend : function() {
        $(".submitter").bind('click', function() {
            $("#racecode").val($(this).attr('raceId'));
            $("#redirector").submit();
        });
    },

    _raceFinalSend : function() {
        $(".submitter").bind('click', function() {

            $.cookie('PHPSESSID', eklpsAdviser.ssid);

            $.post( "http://www.eklps.com/enter_racerfinal.php", {
                horse_enter_id: eklpsAdviser.horseId,
                racecode: $(this).attr('raceId'),
                register: "Записать лошадь на скачку"
            })
                .done(function( data ) {
                    alert( data );
                })
                .fail(function() {
                    alert( "Что-то пошло не так..." );
                });
        });
    }

};

document.addEventListener('DOMContentLoaded', function () {

    $(".hint-block").hide();

    eklpsAdviser.run();

    $(".class-selector").bind('click', function(e) {
        $(this).toggleClass("btn-down");
    });

    $(".hint").bind('click', function(){
        $(".hint-block").toggle();
    })

    $("#refreshAction").bind('click', function() {
        eklpsAdviser.refresh();
    });

});