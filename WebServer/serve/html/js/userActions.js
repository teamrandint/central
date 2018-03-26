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
	$('#usernameText').text(userName);
});

// Handles the rendering of the page according to the command selection.

function commandChanged(){
    var selection = $('#commands').val();

    // Display the corresponding input section on the page.
    switch (selection) {
        case "ADD":
            break
        case "QUOTE":
            break
        case "BUY":
            break
        case "COMMIT_BUY":
            break
        case "CANCEL_BUY":
            break
        case "SELL":
            break
        case "COMMIT_SELL":
            break
        case "CANCEL_SELL":
            break
        case "SET_BUY_AMOUNT":
            break
        case "CANCEL_SET_BUY":
            break
        case "SET_BUY_TRIGGER":
            break
        case "SET_SELL_AMOUNT":
            break
        case "SET_SELL_TRIGGER":
            break
        case "CANCEL_SET_SELL":
            break
        case "DUMPLOG":
            break
        case "DISPLAY_SUMMARY":
            break
    }
}

function submitRequest() {
    var command = $('commands').val();

}

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