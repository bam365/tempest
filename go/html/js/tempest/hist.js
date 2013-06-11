
define([
        "dojo/dom",
        "dojo/dom-construct",
        "dojo/request",
        "dojox/charting/Chart",
        "dojox/charting/themes/Claro",
        "dojox/charting/widget/Legend",
        "dojox/charting/plot2d/Lines",
        "dojox/charting/axis2d/Default",
        "dojo/domReady!"
], function(dom, domCons, request, Chart, theme, Legend) {

        var chartNode = dom.byId("chartNode");

        var addChart = function(series) {
                var histChart = new Chart("chartNode");
                histChart.setTheme(theme);
                histChart.addPlot("default", {
                        type: "Lines",
                        hAxis: "x",
                        vAxis: "y",
                });
                histChart.addAxis("x");
                histChart.addAxis("y", { 
                        vertical: true,
                        min: 0,
                        //max: 200,
                        fixUpper: "major",
                        fixLower: "major"
                });
                for (var s in series) {
                        histChart.addSeries(s, series[s]);
                }
                histChart.render();
                //Must do this after chart is rendered...
                var legend = new Legend({
                        chart: histChart,
                        horizontal: false
                }, "legendNode");
                        
        };
                        

        return {
                updateChart: function() {
                        request.post("/ajax/hist", {
                                handleAs: "json",
                                data: "{ \"interval\": 60 }"
                        }).then(function(hist) {
                                        addChart(hist);
                        },
                        function(error) {
                                //TODO: Should probably do something
                                //here...
                        });
                }
        };
});
