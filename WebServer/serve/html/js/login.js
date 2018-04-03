var URL = "http://172.20.0.5:44456/";
var username = ""

$(document).ready(function(){
	$('#userName').on('click', () => $('#userName').val(''));
});

// Logs the user into the webserver.
function loginRequest() {
	var userName = $('#userName').val();
	$.ajax({
		type: 'POST',
		url: "LOGIN/",
		data: {username: userName},
		success: function(data, status){
			setCookie("dayTradingUsername", $('#userName').val(), 10)
			// Redirect user to actions page
			window.location.replace("http://172.20.0.5:44456/actions.html");
		},
		error:function(){
			alert("Error occured while logging in!")
		}
	});
};

function setCookie(cname, cvalue, exdays) {
    var d = new Date();
    d.setTime(d.getTime() + (exdays*24*60*60*1000));
    var expires = "expires="+ d.toUTCString();
    document.cookie = cname + "=" + cvalue + ";" + expires + ";path=/";
}
