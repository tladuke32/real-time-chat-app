import React, { useEffect, useState, useCallback } from 'react';

function Chat() {
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
            setMessages(prevMessages => [...prevMessages, event.data]);
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
            socket.close();
        };
    }, [wsURL]);

    // Establish WebSocket connection on component mount and clean up on unmount
    useEffect(() => {
        const cleanup = connectWebSocket();

        return () => {
            cleanup();
        };
    }, [connectWebSocket]);

    // Handler for sending messages
    const handleSubmit = (e) => {
        e.preventDefault();

        if (ws && ws.readyState === WebSocket.OPEN) {
            ws.send(message);
            setMessage('');
        } else {
            console.error('WebSocket is not open. Cannot send message.');
        }
    };

    return (
        <div>
            <div>
                {messages.map((msg, idx) => (
                    <div key={idx}>{msg}</div>
                ))}
            </div>
            <form onSubmit={handleSubmit}>
                <input
                    type="text"
                    value={message}
                    onChange={(e) => setMessage(e.target.value)}
                    placeholder="Enter message"
                />
                <button type="submit">Send</button>
            </form>
        </div>
    );
}

export default Chat;
