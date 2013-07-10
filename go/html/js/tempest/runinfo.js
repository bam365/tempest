define(
[
	"dojo/dom",
	"dojo/dom-construct",
	"dojo/request",
	"dojo/domReady!"
], function(dom, domCons, request) {

	var tagRunStart = dom.byId("txtRunStart");
	var tagRunDuration = dom.byId("txtRunDuration");


	return {
		updateInfo: function() {
			request("/ajax/runinfo", {
				handleAs: "json"
			}).then(function(runInfo) {
				tagRunStart.innerHTML = runInfo.runstart;
				tagRunDuration.innerHTML = runInfo.duration;
			},
			function(error) {
				//TODO: Should probably do something
				//here...
			});
		}
	};
});
