var username;
var URL = "http://172.20.0.5:44456/";
var pendingSells = [];
var pendingBuys = [];
var triggers = [];
var submitRequest;

$(document).ready(function(){
	userName = getCookie("dayTradingUsername");
	// Check if user has a valid login
	checkLogin(userName);

    // TODO: pull account information from database

	// Main logic of the program here!
	$('#usernameText').text(userName);
	// Hide inputs by default
	$('#inputTwoTitle').hide();
	$('#textInputTwo').hide();

    // bind 
    $('#submitButton').on('click', submitRequest);
    $('#textInputOne').on('click', () => $('#textInputOne').val(''));
    $('#textInputTwo').on('click', () => $('#textInputTwo').val(''));
});

function checkLogin(userName) {
    if (userName === "") {
       // User isn't logged in! Redirect back to login!
       window.location.replace("http://172.20.0.5:44456/");
    }
}

function getCookie(cname) {
    var name = cname + "=",
    decodedCookie = decodeURIComponent(document.cookie),
    ca = decodedCookie.split(';');
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

// Handles the rendering of the page according to the command selection.
function commandChanged(){
    var selection = $('#commands').val();

    // Display the corresponding input section on the page.
    switch (selection) {
        case "ADD":
        	toggleDisplayInputs(1, "Amount:", "");
            break;
        case "QUOTE":
        	toggleDisplayInputs(1, "Stock:", "");
            break;
        case "BUY":
        	toggleDisplayInputs(2, "Amount:", "Stock:");
            break;
        case "COMMIT_BUY":
        	toggleDisplayInputs(0, "", "");
            break;
        case "CANCEL_BUY":
        	toggleDisplayInputs(0, "", "");
            break;
        case "SELL":
        	toggleDisplayInputs(2, "Amount:", "Stock:");
            break;
        case "COMMIT_SELL":
        	toggleDisplayInputs(0, "", "");
            break;
        case "CANCEL_SELL":
        	toggleDisplayInputs(0, "", "");
            break;
        case "SET_BUY_AMOUNT":
        	toggleDisplayInputs(2, "Amount:", "Stock:");
            break;
        case "CANCEL_SET_BUY":
        	toggleDisplayInputs(1, "Stock:", "");
            break;
        case "SET_BUY_TRIGGER":
        	toggleDisplayInputs(2, "Amount:", "Stock:");
            break;
        case "SET_SELL_AMOUNT":
        	toggleDisplayInputs(2, "Amount:", "Stock:");
            break;
        case "SET_SELL_TRIGGER":
        	toggleDisplayInputs(2, "Amount:", "Stock:");
            break;
        case "CANCEL_SET_SELL":
        	toggleDisplayInputs(1, "Stock:", "");
            break;
        case "DUMPLOG":
        	toggleDisplayInputs(1, "Filename:", "");
            break;
        case "DISPLAY_SUMMARY":
        	toggleDisplayInputs(0, "", "");
            break;
    }
}

function toggleDisplayInputs(numVisible, displayTextOne, displayTextTwo){
	switch (numVisible) {
		case 0:
			$('#inputOneTitle').hide();
			$('#textInputOne').hide();
	    	$('#inputTwoTitle').hide();
			$('#textInputTwo').hide();
			break;
		case 1:
			if (displayTextOne.localeCompare("Amount:") === 0) {
				$('#inputOneTitle').show();
				$('#textInputOne').show();
				$('#inputOneTitle').text(displayTextOne);
		    	$('#inputTwoTitle').hide();
				$('#textInputTwo').hide();
			} else {
				$('#inputOneTitle').hide();
				$('#textInputOne').hide();
		    	$('#inputTwoTitle').show();
		    	$('#inputTwoTitle').text(displayTextOne);
				$('#textInputTwo').show();
			}
			break;
		case 2:
			$('#inputOneTitle').show();
			$('#textInputOne').show();
			$('#inputOneTitle').text(displayTextOne);
	    	$('#inputTwoTitle').show();
			$('#textInputTwo').show();
			$('#inputTwoTitle').text(displayTextTwo);
			break
	}
	// Clear inputs of values...doesn't work
	$('#textInputOne').val('')
	$('#textInputTwo').val('')
	// Clear status area 
	$('#resultsDiv').text('')
}

function submitRequest() {
    var command = $('#commands').val(),
    params = {username: userName};
    submitRequest = {command: $('#commands').val(),
					 amount: 0,
					 stock: ''};

    if ($('#textInputOne').is(":visible")) {
    	params.amount = $('#textInputOne').val();
    	submitRequest.amount = $('#textInputOne').val();
    }
    if ($('#textInputTwo').is(":visible")) {
    	if (command.localeCompare('DUMPLOG') === 0) {
    		params.filename = $('#textInputTwo').val();
    	} else {
    		params.stock = $('#textInputTwo').val();
    		submitRequest.stock = $('#textInputTwo').val();
    	}
    }
    $.ajax({
    	url: URL + command + "/",
    	data: params,
    	type: 'POST',
    	success: function(data, status) {
            // Clear values from inputs
            $('#textInputOne').val('');
            $('#textInputTwo').val('');
    		displaySuccess(data);
    	},
    	error: function(jqXHR, textStatus, errorThrown) {
    		// Display error message to user.
    		var err = jqXHR.responseText;
    		$('#resultsDiv').text('Error occured: ' + err);
    	}
    });
}

function displaySuccess(data) {
	var successMsg = submitRequest.command.toLowerCase().replace(/\b\w/g, l => l.toUpperCase());

	if (submitRequest.amount !== 0) {
		successMsg += " Amount: $" + submitRequest.amount;
	}
	if (submitRequest.stock.localeCompare('') !== 0) {
		successMsg += " Stock: " + submitRequest.stock;
	}
	successMsg += " successful.";

	// Display different message for successful quotes
	if (submitRequest.command.localeCompare('QUOTE') === 0) {
		successMsg = "Stock: " + submitRequest.stock + " - $" + data;
	}

	// Handle dumplog
	if (submitRequest.command.localeCompare('DUMPLOG') === 0) {
		// DO something with the data
	}

	$('#resultsDiv').text(successMsg);
    // TODO: Open save prompt for dumplog
}