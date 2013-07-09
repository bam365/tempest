define(
[
	"dojo/dom",
	"dojo/dom-construct",
	"dojo/request",
	"dojox/gauges/AnalogGauge",
	"dojox/gauges/AnalogArrowIndicator",
	"dojox/gauges/AnalogLineIndicator",
	"dojo/domReady!"
], function(dom, domCons, request) {

	var readingsTag = dom.byId("readings");
	var indicators = {};


	var gaugeID = function(name) {
		return name + "-gauge";
	};

	var readingGroupID = function(name) {
		return name + "-grp";
	};

	var readingTextID = function(name) {
		return name + "-txt";
	};


	var addSensor = function(name, sinfo, ul) {
		var lnode = domCons.create("div", {
			id: readingGroupID(name),
			className: "readingGroup"
		}, ul);
		domCons.create("h3", {
			className: "sensorName",
			innerHTML: name + ": "
		}, lnode);
		var gnode = domCons.create("div", {
			id: gaugeID(name)
		}, lnode);
		domCons.create("h3", {
			id: readingTextID(name),
			className: "sensorReading",
			innerHTML: "17"
		}, lnode);
		makeGauge(gnode, name, sinfo);
	};

	var setReadings = function(rdgs) {
		for (var rdg in rdgs) {
			var val = rdgs[rdg].data;
			if (indicators.hasOwnProperty(rdg)) {
				indicators[rdg].update(val);
			}
	                //TODO: Make sure this node actually exists
	                var rdgtxt = dom.byId(readingTextID(rdg));
	                rdgtxt.innerHTML = val.toString();
	                //TODO: Check error value
	        }
	};


	var makeGauge = function(parent, name, ginfo) {
		indicators[name] = new dojox.gauges.AnalogArrowIndicator({
			value: 17, 
			width: 3,
			title: 'Reading',
			noChange: true,
			hideValue: true
		});
		var alertIndicators = [ 
			new dojox.gauges.AnalogLineIndicator({
				value: ginfo.alert.lo,
				color: "#CC0000",
				width: 2,
				title: 'Low alert',
				noChange: true,
				hideValue: true
			}),
			new dojox.gauges.AnalogLineIndicator({
				value: ginfo.alert.hi,
				color: "#CC0000",
				width: 2,
				title: 'High alert',
				noChange: true,
				hideValue: true
			})
		];
		var g = new dojox.gauges.AnalogGauge({
			background: [255, 255, 255, 0],
			id: gaugeID(name),
			width: 300,
			height: 150,
			cy: 140,
			radius: 125,
			min: ginfo.range.lo,
			max: ginfo.range.hi,
			ranges: [ {low: ginfo.range.lo, high: ginfo.range.hi } ],
			majorTicks: {
				offset: 125,
				interval: 10,
				length: 5,
				color: 'gray'
			},
			indicators: alertIndicators,
			hideValue: true
		}, parent);
		g.addIndicator(indicators[name]);
		g.startup();
		return g;
	};


	return {
		setupReadings: function() {
			request("/ajax/sensors", {
				handleAs: "json"
			}).then(function(sensors) {
				for (var s in sensors) {
					addSensor(s, sensors[s], readingsTag);
				}
				updateReadings();
			},
			function(error) {
                        //TODO: Should probably do something
                        //here...
                });
		},

		updateReadings: function() {
			request("/ajax/readings", {
				handleAs: "json"
			}).then(function(rdgs) {
				setReadings(rdgs);
			},
			function(error) {
                        //TODO: Should probably do something
                        //here...
                });
		}
	};
});






