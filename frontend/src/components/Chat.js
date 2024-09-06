import React, { useEffect, useState, useRef } from 'react';
import '../styles/chat.css'; // Assuming CSS is set for styling

function Chat({ user }) {
    const [messages, setMessages] = useState([]);
    const [message, setMessage] = useState('');
    const wsRef = useRef(null);
    const wsURL = `${process.env.REACT_APP_API_URL.replace(/^http/, 'ws')}/ws`;

    // Fetch chat history when the component mounts
    const fetchChatHistory = async () => {
        try {
            const response = await fetch(`${process.env.REACT_APP_API_URL}/chat-history`);
            const history = await response.json();
            setMessages(history);
        } catch (error) {
            console.error('Error fetching chat history:', error);
        }
    };

    // Function to initialize WebSocket connection
    const connectWebSocket = () => {
        if (wsRef.current && (wsRef.current.readyState === WebSocket.OPEN || wsRef.current.readyState === WebSocket.CONNECTING)) {
            console.log('WebSocket is already connected or connecting.');
            return;
        }

        const socket = new WebSocket(wsURL);

        socket.onopen = () => {
            console.log('Connected to WebSocket');
        };

        socket.onmessage = (event) => {
            const data = JSON.parse(event.data);
            console.log('Message received from WebSocket:', data); // Log the received message

            // Check for duplicates based on timestamp or unique ID if available
            setMessages(prevMessages => {
                if (prevMessages.some(msg => msg.id === data.id)) {
                    console.log('Duplicate message detected, ignoring:', data);
                    return prevMessages;
                }
                return [...prevMessages, data];
            });
        };

        socket.onerror = (error) => {
            console.error('WebSocket Error:', error);
        };

        socket.onclose = (event) => {
            console.log('WebSocket connection closed', event);
            wsRef.current = null;
            if (event.code !== 1000) {
                setTimeout(() => {
                    connectWebSocket();
                }, 5000);
            }
        };

        wsRef.current = socket;
    };

    // Establish WebSocket connection on component mount
    useEffect(() => {
        console.log('Connecting WebSocket on component mount.');
        fetchChatHistory();
        connectWebSocket();

        // Cleanup WebSocket on component unmount
        return () => {
            if (wsRef.current) {
                console.log('Cleaning up WebSocket connection on component unmount.');
                wsRef.current.close();
            }
        };
    }, []); // Empty dependency array ensures this runs only once on mount

    // Handler for sending messages
    const handleSubmit = (e) => {
        e.preventDefault();
        if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
            const messageData = { username: user.username, message };
            console.log('Sending message:', messageData);
            wsRef.current.send(JSON.stringify(messageData));

            // Add the message to the local state immediately
            setMessages(prevMessages => [...prevMessages, { ...messageData, local: true }]);
            setMessage(''); // Clear the input field
        } else {
            console.error('WebSocket is not open. Cannot send message.');
        }
    };

    return (
        <div className="chat-container">
            <h2>Chat</h2>
            <div className="messages">
                {messages.map((msg) => (
                    <div
                        key={msg.id}
                        className={`message ${msg.username === user.username ? 'sent' : 'received'}`}
                    >
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
