import React, { Component } from "react";
import "./App.css";
import { connect, sendMsg } from "./api";
import ChatInput from './components/ChatInput/ChatInput';

class App extends Component{
  constructor(props){
    super(props);
    this.showNotification = this.showNotification.bind(this);
  }

  componentDidMount() {
    if (!("Notification" in window)) {
      console.log("This browser does not support desktop notification");
    } else {
      Notification.requestPermission();
    }
    connect((msg) => {
      this.showNotification(msg.data)
    });
  }
  
  showNotification(data) {
    var options = {
      body: data,
    };
   var notification = new Notification("Desktop Notifier", options);
  }

  send(event) {
    if (event.keyCode === 13) {
      sendMsg(event.target.value);
      event.target.value = "";
    }
  }
  render() {
    return (
      <div className="App">
        <ChatInput send={this.send} />
      </div>
    );
  }

}

export default App;
