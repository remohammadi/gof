var startTime;

var getUserID = function() {
  if (typeof(Storage) !== "undefined") {
    userID = localStorage.getItem("userID");
    if (!userID) {
      userID = Math.random().toString(36).substring(2, 15);
      localStorage.setItem("userID", userID);
    }
    return userID;
  }
  return ""
}

var getUserFP = function() {
  var client = new ClientJS();
  return client.getFingerprint();
}

var av_judge = function(genuine) {
  var endTime = new Date();
  $("#buttons").hide();
  if (genuine == referenceGenuine) {
    $("#resultPassed").show();
    $("#resultFailed").hide();
  } else {
    $("#resultPassed").hide();
    $("#resultFailed").show();
  }
  $("#userID").val(getUserID());
  $("#userFP").val(getUserFP());
  $("#userChoice").val(genuine);
  var diff = endTime.getTime() - window.startTime.getTime();
  $("#duration").val(diff);
  $("#result").show(600);
  return false;
}

$(document).ready(function(){
  window.startTime = new Date();
});
