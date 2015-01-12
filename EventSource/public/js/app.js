$(function() {
  var addMessage, isBlank, stream, updateUsers, username;
  isBlank = function(string) {
    return (string == null) || string.trim() === "";
  };
  while (isBlank(username)) {
    username = prompt("What's your name? ");
    if (!isBlank(username)) {
      $("#chat-msg").attr("placeholder", "" + username + ": Message Here");
      $(".username").text(username);
    }
  }

  $("#chat-form").on("submit", function(e) {
    $.post("/messages", {
      msg: $("#chat-msg").val(),
      name : username
    });
    $("#chat-msg").val("")
    $("#chat-msg").focus();
    return false
  });

  addMessage = function(message) {
    var el, text;
    test = "<strong>" + message.name + "</strong ";
    text += message.msg;
    el = $("<li data-name='" + message.name + "'>").html(text)
    if (message.name === username) {
      el.css({
        "background-color": "#ba8a9d"
      });
    }
    if (message.name === "Admin") {
      el.css({
        "background-color": "#92b4ba"
      });
    }
    $("#chat").append(el);
  }

  updateUsers = function(users) {
    var el, user, _i, _len, _results;
    $("#users li").not(".nav-header").remove()
    _results = [];
    for (_i = 0, _len = users.length; _i < _len; _i++) {
      user = users[_i];
      el = $("<li>").html("<a>" + user+ "</a>");
      if (user === username) {
        el.addClass("active");
      }
      _results.push($("#users").append(el));
    }
  };

  stream = new EventSource("/stream");
  stream.onopen = function() {
    $.post("/users", {
      user: username
    });
    $(".chat-unavailable").hide();
    $("#chat-area").show();
    $("#chat-msg").focus();
  };
  stream.onmessage = function(e) {
    payload = JSON.parse(e.data);
    if (payload.type === "message") {
      addMessage(payload.data);
    }
    if (payload.type === "users") {
      updateUsers(payload.data)
    }
  };
  stream.onerror = function(e) {
    $(".chat-unavailable").show()
    $("#chat-area").hide();
  };

  window.onbeforeunload = function() {
    $.ajax({
      url: "users?user=" + username,
      type: "DELETE",
    });
    stream.close();
  }
})
