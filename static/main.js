var showingAnswer = false;

var av_onpause = function(v) {
  console.log(v);
  if ((2 * v.currentTime) > v.duration) {
    $("#manipulatedBtn").removeAttr("disabled");
    $("#genuineBtn").removeAttr("disabled");
  }
}

var av_judge = function(genuine) {
  if (genuine == referenceGenuine) {
    $("#resultPassed").show();
    $("#resultFailed").hide();
  } else {
    $("#resultPassed").hide();
    $("#resultFailed").show();
  }
  return false;
}
