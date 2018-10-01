new Vue({
    el: '#app',

    data: {
        ws: null, // websocket连接 
        newMsg: '', // 要发送的消息 
        chatContent: '', // 聊天框里显示的消息 
        email: null, // email地址用来获取avatar 
        username: null, // 用户名 
        joined: false // 当email和username填入input时为真
    },

    created: function() {
        var self = this;
        // 实例化一个websocket连接
        this.ws = new WebSocket('ws://' + window.location.host + '/ws');
        // 监听message事件，当从服务端获取到消息的时候，显示在屏幕里
        this.ws.addEventListener('message', function(e) {
            var msg = JSON.parse(e.data);
            self.chatContent += '<div class="chip">'
                    + '<img src="' + self.gravatarURL(msg.email) + '">' // 头像 
                    + msg.username
                + '</div>'
                + emojione.toImage(msg.message) + '<br/>'; // 解析emoji表情

            var element = document.getElementById('chat-messages');
            element.scrollTop = element.scrollHeight; // 消息流自动到底部 
        });
    },

    methods: {
        send: function () {
            // 如果消息input不为空，发送出去。
            if (this.newMsg != '') {
                this.ws.send(
                    JSON.stringify({
                        email: this.email,
                        username: this.username,
                        message: $('<p>').html(this.newMsg).text() // 发送的内容 
                    }
                ));
                this.newMsg = ''; // 发送出去以后，重置为空 
            }
        },

        join: function () {
            if (!this.email) {
                Materialize.toast('你必须输入电子邮件地址！', 2000);
                return
            }
            if (!this.username) {
                Materialize.toast('你必须输入用户名', 2000);
                return
            }
            this.email = $('<p>').html(this.email).text();
            this.username = $('<p>').html(this.username).text();
            this.joined = true;
        },

        // 获取头像
        gravatarURL: function(email) {
            return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(email);
        }
    }
});