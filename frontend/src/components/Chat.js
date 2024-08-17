import React, { useEffect, useState, useCallback } from 'react';

function Chat({ user }) {
    const [messages, setMessages] = useState([]);
    const [message, setMessage] = useState('');
    const [ws, setWs] = useState(null);

    const wsURL = `${process.env.REACT_APP_API_URL.replace('http', 'ws')}/ws`;

    // Function to initialize WebSocket connection
    const connectWebSocket = useCallback(() => {
        const socket = new WebSocket(wsURL);

        socket.onopen = () => {
            console.log('Connected to WebSocket');
        };

        socket.onmessage = (event) => {
            const data = JSON.parse(event.data); // Assuming message data comes in JSON format
            setMessages(prevMessages => [...prevMessages, data]);
        };

        socket.onerror = (error) => {
            console.error('WebSocket Error:', error);
        };

        socket.onclose = (event) => {
            console.log('WebSocket connection closed', event);
            if (event.code !== 1000) { // Reconnect only if the close was abnormal
                setTimeout(() => {
                    connectWebSocket();
                }, 5000);
            }
        };

        setWs(socket);
        return () => {
            if (socket.readyState === WebSocket.OPEN) {
                socket.close();
            }
        };
    }, [wsURL]);

    // Establish WebSocket connection on component mount
    useEffect(() => {
        const cleanup = connectWebSocket();
        return cleanup; // Cleanup WebSocket on component unmount
    }, [connectWebSocket]);

    // Handler for sending messages
    const handleSubmit = (e) => {
        e.preventDefault();
        if (ws && ws.readyState === WebSocket.OPEN) {
            const messageData = { username: user.username, message }; // Send username with the message
            ws.send(JSON.stringify(messageData));
            setMessage('');
        } else {
            console.error('WebSocket is not open. Cannot send message.');
        }
    };

    return (
        <div className="chat-container">
            <h2>Chat</h2>
            <div className="messages">
                {messages
                    .filter(msg => msg.username !== user.username) // Filter out messages sent by the current user
                    .map((msg, idx) => (
                    <div key={idx} className="message">
                        <strong>{msg.username || "Anonymous"}:</strong> {msg.message}
                    </div>
                ))}
            </div>
            <form onSubmit={handleSubmit} className="message-form">
                <input
                    type="text"
                    className="form-control"
                    value={message}
                    onChange={(e) => setMessage(e.target.value)}
                    placeholder="Enter message"
                    required
                />
                <button type="submit" className="btn btn-primary">Send</button>
            </form>
        </div>
    );
}

export default Chat;
