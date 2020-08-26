window.addEventListener('message', function(event) {
  if (event.data.rpcId !== "0") {
    return;
  }
  if (event.data.error) {
    console.log("ERROR: " + event.data.error);
  } else {
    const el = document.getElementById("offer-iframe");
    el.setAttribute("src", event.data.uri);
  }
});
document.addEventListener('DOMContentLoaded', function() {
  const template = window.location.protocol.replace('http', 'ws') +
    "//$API_HOST/.sandstorm-token/$API_TOKEN/socket";
  window.parent.postMessage({renderTemplate: {
    rpcId: "0",
    template: template,
    clipboardButton: 'left'
  }}, "*");
})
