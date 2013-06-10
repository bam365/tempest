define([
        "dojo/dom",
        "dojo/dom-construct",
        "dojo/request",
        "dojo/domReady!"
], function(dom, domCons, request) {

        var readingsTag = dom.byId("readings");

        var addReading = function(name, rdg, ul) {
                var lnode = domCons.create("div", {
                        className: "sensorItem",
                }, ul);
                domCons.create("span", {
                        className: "sensorName",
                        innerHTML: name + ": "
                }, lnode);
                var rdgtxt;
                if (!(rdg["err"] === "")) {
                        rdgtxt = "ERROR: " + rdg["err"];
                } else {
                        rdgtxt = rdg["data"].toString();
                }
                domCons.create("span", {
                        className: "sensorReading",
                        innerHTML: rdgtxt
                }, lnode);
        };
                        
        var setReadings = function(rdgs) {
                domCons.empty("readings");
                var ul = domCons.create("ul", {
                        className: "readingsList"
                }, readingsTag);

                for (var rdg in rdgs) {
                        addReading(rdg, rdgs[rdg], ul);
                }
        };


        return {
                updateReadings: function() {
                        request("/ajax/readings", {
                                handleAs: "json"
                        }).then(function(rdgs) {
                                        setReadings(rdgs)
                        },
                        function(error) {
                                //TODO: Should probably do something
                                //here...
                        });
                }
        };
});






