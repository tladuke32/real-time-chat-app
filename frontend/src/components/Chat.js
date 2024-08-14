import React, { useEffect, useState } from 'react';

function Chat() {
    const [messages, setMessages] = useState([]);
    const [message, setMessage] = useState('');
    const [ws, setWs] = useState(null);
    const [reconnect, setReconnect] = useState(false);
    const wsURL = process.env.REACT_APP_API_URL.replace('http', 'ws') + '/ws';

    useEffect(() => {
        const connectWebSocket = () => {
            const socket = new WebSocket(wsURL);
            setWs(socket);

            socket.onopen = () => {
                console.log('Connected to WebSocket');
                setReconnect(false); // Reset reconnect flag on successful connection
            };

            socket.onmessage = (event) => {
                setMessages((prevMessages) => [...prevMessages, event.data]);
            };

            socket.onerror = (error) => {
                console.error('WebSocket Error:', error);
            };

            socket.onclose = (event) => {
                console.log('WebSocket connection closed', event.reason);
                if (!reconnect) {
                    setReconnect(true);
                }
            };

            // Clean up WebSocket connection on component unmount
            return () => {
                if (ws && ws.readyState === WebSocket.OPEN) {
                    ws.close();
                }
            };
        };

        const socketCleanup = connectWebSocket();

        // Reconnect logic
        const reconnectInterval = setInterval(() => {
            if (reconnect) {
                socketCleanup();
            }
        }, 5000); // Attempt to reconnect every 5 seconds

        return () => {
            clearInterval(reconnectInterval);
            socketCleanup();
        };
    }, [wsURL, reconnect]);

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
