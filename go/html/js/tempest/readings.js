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

	var readingsTag = dom.byId("boxReadings");
	var indicators = {};

	var rdgGroupClass = "readingGroup";
	var rdgTxtClass = "sensorReading";
	var errClass = "readingTxtError";

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
			className: rdgGroupClass
		}, ul);
		domCons.create("h3", {
			className: "sensorName",
			innerHTML: name 
		}, lnode);
		var gnode = domCons.create("div", {
			id: gaugeID(name)
		}, lnode);
		domCons.create("h4", {
			id: readingTextID(name),
			className: rdgTxtClass,
			innerHTML: "17"
		}, lnode);
		makeGauge(gnode, name, sinfo);
	};

	var setReadings = function(rdgs) {
		for (var rdg in rdgs) {
			var val = rdgs[rdg].data;
			var err = rdgs[rdg].err;
			var rdgtxt = dom.byId(readingTextID(rdg));
			var rdghdr = dom.byId(readingGroupID(rdg));
			if (err !== "") {
				rdgtxt.innerHTML = err;
				rdgtxt.className = errClass;
				dom.byId(readingGroupID(rdg)).className = errClass; 
			} else {
				rdgtxt.className = rdgTxtClass;
				rdghdr.className = rdgGroupClass;
				indicators[rdg].update(val);
				rdgtxt.innerHTML = val.toString();
			}

		}
	};

	var calcTickIntervals = function(ginfo) {
		var range = ginfo.range.hi - ginfo.range.lo;
		var roughMajor = range / 15;
		var roundedMajor = Math.ceil(roughMajor / 5) * 5;
		return {
			major: roundedMajor,
			minor: roundedMajor / 3.0
		};
	};


	var makeAlertIndicators = function(ginfo) {
		ret = [];
		if (ginfo.alert.hasOwnProperty("lo") && ginfo.alert.lo > ginfo.range.lo) {
			ret.push(new dojox.gauges.AnalogLineIndicator({
				value: ginfo.alert.lo,
				color: "#CC0000",
				width: 2,
				title: 'Low alert',
				noChange: true,
				hideValue: true
			}));
		}
		if (ginfo.alert.hasOwnProperty("hi") && ginfo.alert.hi < ginfo.range.hi) {
			ret.push(new dojox.gauges.AnalogLineIndicator({
				value: ginfo.alert.hi,
				color: "#CC0000",
				width: 2,
				title: 'Low alert',
				noChange: true,
				hideValue: true
			}));
		}
		return ret;
	};


	var makeGauge = function(parent, name, ginfo) {
		tickIntervals = calcTickIntervals(ginfo);
		var g = new dojox.gauges.AnalogGauge({
			background: [255, 255, 255, 0],
			id: gaugeID(name),
			width: 350,
			height: 175,
			cy: 160,
			radius: 125,
			min: ginfo.range.lo,
			max: ginfo.range.hi,
			ranges: [ {low: ginfo.range.lo, high: ginfo.range.hi } ],
			majorTicks: {
				offset: 125,
				interval: tickIntervals.major,
				length: 5,
				color: 'gray'
			},
			minorTicks: {
				offset: 125,
				interval: tickIntervals.minor,
				length: 3,
				color: 'gray'
			},
			indicators: makeAlertIndicators(ginfo),
			hideValue: true
		}, parent);
                g.startup();
                var arrow = new dojox.gauges.AnalogArrowIndicator({
                        value: 17, 
                        width: 3,
                        title: 'Reading',
                        noChange: false,
                        hideValue: true
                });
		g.addIndicator(arrow);
                indicators[name] = arrow;
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






