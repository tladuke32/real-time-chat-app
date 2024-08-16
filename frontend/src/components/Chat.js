import React, { useEffect, useState, useCallback } from 'react';

function Chat({ username }) { // Assuming username is passed as a prop to the component
    const [messages, setMessages] = useState([]);
    const [message, setMessage] = useState('');
    const wsURL = process.env.REACT_APP_API_URL.replace('http', 'ws') + '/ws';
    const [ws, setWs] = useState(null);

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

        socket.onclose = () => {
            console.log('WebSocket connection closed');
            // Automatically try to reconnect on unexpected closure
            setTimeout(() => {
                connectWebSocket();
            }, 5000);
        };

        setWs(socket);
        return () => {
            if (socket.readyState === WebSocket.OPEN) {
                socket.close();
            }
        };
    }, [wsURL]);

    // Establish WebSocket connection on component mount and clean up on unmount
    useEffect(() => {
        const cleanup = connectWebSocket();
        return cleanup;
    }, [connectWebSocket]);

    // Handler for sending messages
    const handleSubmit = (e) => {
        e.preventDefault();
        if (ws && ws.readyState === WebSocket.OPEN) {
            const messageToSend = JSON.stringify({ username: username || "Anonymous", message });
            ws.send(messageToSend);
            setMessage('');
        } else {
            console.error('WebSocket is not open. Cannot send message.');
        }
    };

    return (
        <div className="chat-container">
            <h2>Chat</h2>
            <div className="messages">
                {messages.map((msg, idx) => (
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
