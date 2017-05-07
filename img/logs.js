function formatTime(t) {
	// If t is greater than or equal to one hour
	if (~~(t / 3600) > 0) {
		// ~~ is essentially a shorthand for Math.floor()
		// this is, essentially, time / 3600 (time in hours)
		// rounded to the nearest whole number
		var hour = ~~(t / 3600)
		var min  = ~~((t % 3600) / 60)
		var sec  = t % 60
		// Append a 0 to the : if minutes / seconds are less than 9
		// that way it displays as 10:09 rather than 10:9
		return hour + (min > 9 ? ":" : ":0") + min + (sec > 9 ? ":" : ":0") + sec;
	} else {
		var min = ~~(t / 60)
		var sec = t%60
		return min + (sec > 9 ? ":" : ":0") + sec;
	}
}

// I would like to use something like a struct here,
// maybe an object or something of the sort.
// Can't be bothered right now
var nui = 0;
var forest = 0;
var atc = 0
var btc = 0
var brazier = 0;
// When the page loads
$(document).ready(function() {
	getTimers();
	// Sync timers with db every 10 seconds to make sure 
	// values haven't deviated
	setInterval("getTimers()", 10000);
	// Decrement timer every 1 second
	setInterval("countDown()", 1000);
});

function countDown() {
	// Tons of if statements like this rub me the wrong way
	if (nui >= 1) {
		nui--;
		$("#nui").html(formatTime(nui));
		$("#nuied").css("display","none");
	}
	if (nui == 0) {
		$("#nui").html(formatTime(nui));
		$("#nuied").css("display","inline");
	}
	if (forest >= 1) {
		forest--;
		$("#forest").html(formatTime(forest));
		$("#forested").css("display","none");
	}
	if (forest == 0) {
		$("#forest").html(formatTime(forest));
		$("#forested").css("display","inline");
	}
	if (atc >= 1) {
		atc--;
		$("#atc").html(formatTime(atc));
		$("#atced").css("display","none");
	}
	if (atc == 0) {
		$("#atc").html(formatTime(atc));
		$("#atced").css("display","inline");
	}
	if (btc >= 1) {
		btc--;
		$("#btc").html(formatTime(btc));
		$("#btced").css("display","none");
	}
	if (btc == 0) {
		$("#btc").html(formatTime(btc));
		$("#btced").css("display","inline");
	}
	if (brazier >= 1) {
		brazier--;
		$("#lit").html("Lit");
		$("#bcd").html(formatTime(brazier));
		$("#bcd").css("display","block");
		$("#bri").css("display","none");
	}
	if (brazier == 0) {
		$("#lit").html("Not lit!");
		$("#bri").css("display","block");
		$("#bcd").css("display","none");
	}
}

// There's a better way to do this
function getTimers() {
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
	$.post("/braziertimer", function(data, status) {
		brazier = data;
	});
}

