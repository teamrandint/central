var username;
var URL = "http://172.20.0.5:44456/";
var pendingSells = [];
var pendingBuys = [];
var triggers = [];

$(document).ready(function(){
	userName = getCookie("dayTradingUsername");
	// Check if user has a valid login
	checkCookie(userName);
	// Main logic of the program here!
	$('#usernameText').text(userName)
});

function checkCookie(userName) {
    if (userName === "") {
       // User isn't logged in! Redirect back to login!
       window.location.replace("http://172.20.0.5:44456/");
    }
}

function getCookie(cname) {
    var name = cname + "=";
    var decodedCookie = decodeURIComponent(document.cookie);
    var ca = decodedCookie.split(';');
    for(var i = 0; i <ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }
    return "";
}