import React, { useEffect, useState } from 'react';

function Chat() {
    const [messages, setMessages] = useState([]);
    const [message, setMessage] = useState('');
    const wsURL = process.env.REACT_APP_API_URL.replace('http', 'ws') + '/ws';

    useEffect(() => {
        const ws = new WebSocket(wsURL);

        ws.onmessage = (event) => {
            setMessages((prevMessages) => [...prevMessages, event.data]);
        };

        ws.onerror = (error) => {
            console.error('WebSocket Error:', error);
        };

        return () => {
            ws.close();
        };
    }, [wsURL]);

    const handleSubmit = (e) => {
        e.preventDefault();
        const ws = new WebSocket(wsURL);

        ws.onopen = () => {
            ws.send(message);
            setMessage('');
        };

        ws.onerror = (error) => {
            console.error('WebSocket Error:', error);
        };
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
