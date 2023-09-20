var selectedChat = "General"

class Event {
    constructor(type, payload) {
        this.type = type;
        this.payload = payload;
    }
}
function routeEvent(event) {
    if (event.type === undefined) {
        alert("no Type is defined")
    }
    switch (event.type) {
        case "new message": {
            console.log("New message");
            break
        }
        // case "send message": {
        //     console.log("send message");
        //     break
        // }
        default: {
            alert("Unsupported event type ");
            break;
        }

    }

}
function sendEvent(EventType, payload) {
    const event = new Event(EventType, payload)
    conn.send(JSON.stringify(event))
}
function changeChatRoom() {
    var newChat = document.getElementById("chatroom");
    if (newChat != null && newChat.value != selectedChat) {
        console.log("newChat : " + newChat);
    }
    return false;
}

function sendMessage() {
    var newMessage = document.getElementById("message");
    if (newMessage != null) {
        // console.log("newMessage : "+newMessage);
        // conn.send(newMessage.value);
        sendEvent("send_message", newMessage.value)
    }
    return false;
}

function login() {
    console.log("Login Method called");
    let formData = {
        "username": document.getElementById("username").value,
        "password": document.getElementById("password").value
    }
    fetch("login", {
        method: "post",
        body: JSON.stringify(formData),
        mode: 'cors',
    }).then((response) => {
        if (response.ok) {
            console.log("Secound response");
            return response.json();
        } else {
            console.log("Unauthorized");
            throw "Unathorized"
        }
    }).then((res) => {
        // we are Authorized
        console.log("Secound res");
        webSocketConnection(res.otp)
    }).catch((e) => {
        console.log("error occurs");
        alert(e)
    });
    return false
}
function webSocketConnection(otp) {
    if (window["WebSocket"]) {
        console.log("Browser support Web socket");
        // logic to connect to web-socket
        conn = new WebSocket("ws://" + document.location.host + "/ws?otp=",otp)
        console.log("Connected to the web Socket");
        conn.onopen = function(evt){
            var text = document.getElementById("connection-header").innerHTML = "Connected to web-socket : true";
        }
        conn.onclose = function(evt){
            var text = document.getElementById("connection-header").innerHTML = "Connected to web-socket : false";
        }
        conn.onmessage = function (evt) {
            // console.log("evt : ", evt);
            const eventData = JSON.parse(evt.data);
            const event = Object.assign(new Event, eventData) // assign the event data to Event class
            routeEvent(event);
        }
    } else {
        alert("Browser does not support web socket")
    }
}
window.onload = function () {
    document.getElementById("chatroom-selection").onsubmit = changeChatRoom;
    document.getElementById("chatroom-message").onsubmit = sendMessage;
    // console.log("Login method called");
    // alert("Login")
    document.getElementById("login-form").onsubmit = login;
}