var MainView = React.createClass({

  getInitialState: function() {
    return {
      temps: {}
    };
  },

  componentWillMount: function() {
      var wsuri = "ws://c27f1925-8697-4d56-b83f-d8afd6a7192a.pub.cloud.scaleway.com:8080";
      this.sock = new WebSocket(wsuri);

  },

  componentDidMount: function() {
      var _this = this;
      this.sock.onmessage = function(e) {
          var temps = _this.state.temps;
          var parts = e.data.split(":");
          temps[parts[0]] = parts[1];
          _this.setState({temps: temps});
          console.log("message received: " + e.data);
      }

      this.sock.onopen = function() {
          console.log("connected");
      }

      this.sock.onclose = function(e) {
          console.log("connection closed (" + e.code + ")");
      }

    // this.chatRoom.bind('new_message', function(message){
    //   this.setState({messages: this.state.messages.concat(message)})
    //
    //   $("#message-list").scrollTop($("#message-list")[0].scrollHeight);
    //
    // }, this);
  },

  render: function() {
      var _this = this;
      var body = Object.keys(this.state.temps).map(function(k) {
            var val = _this.state.temps[k];
            return (
                <div class="temp">
                    {k} = {val}
                </div>
            );
      })

    return (
        <div class="temps">
            {body}
        </div>
    );
  }

});

React.render(<MainView />, document.getElementById("app"));
