function formatTime(t) {
	if (~~(t / 3600) > 0) {
		var hour = ~~(t / 3600)
		var min  = ~~((t % 3600) / 60)
		var sec  = t % 60
		return hour + (min > 9 ? ":" : ":0") + min + (sec > 9 ? ":" : ":0") + sec;
	} else {
		var min = ~~(t / 60)
		var sec = t%60
		return min + (sec > 9 ? ":" : ":0") + sec;
	}
}

var nui = 0;
var forest = 0;
var atc = 0
var btc = 0
$(document).ready(function() {
	getFromDB();
	$("#nui").append(nui);
	$("#forest").append(forest);	
	$("#atc").append(atc);
	$("#btc").append(btc);
	setInterval("getFromDB()", 10000);
	setInterval("countDown()", 1000);
});

function countDown() {
	if (nui >= 1) {
		nui--;
		$("#nui").html(formatTime(nui));
		$("#nuied").css("display","none");
	}
	if (nui == 0) {
		$("#nuied").css("display","inline");
	}
	if (forest >= 1) {
		forest--;
		$("#forest").html(formatTime(forest));
		$("#forested").css("display","none");
	}
	if (forest == 0) {
		$("#forested").css("display","inline");
	}
	if (atc >= 1) {
		atc--;
		$("#atc").html(formatTime(atc));
		$("#atced").css("display","none");
	}
	if (atc == 0) {
		$("#atced").css("display","inline");
	}
	if (btc >= 1) {
		btc--;
		$("#btc").html(formatTime(btc));
		$("#btced").css("display","none");
	}
	if (btc == 0) {
		$("#btced").css("display","inline");
	}
}

function getFromDB() {
	$.post("/nuitimer", function(data, status) {
			nui = data;
		});
	$.post("/foresttimer", function(data, status) {
			forest = data;
		});
	$.post("/atctimer", function(data, status) {
			atc = data;
		});
	$.post("/btctimer", function(data, status) {
			btc = data;
		});
}
