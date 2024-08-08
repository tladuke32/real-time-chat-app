import React, { useEffect, useState } from 'react';

function Chat() {
    const [messages, setMessages] = useState([]);
    const [message, setMessage] = useState('');

    useEffect(() => {
        const ws = new WebSocket('ws://localhost:8080/ws');
        ws.onmessage = (event) => {
            setMessages((prevMessages) => [...prevMessages, event.data]);
        };
        return () => ws.close();
    }, []);

    const handleSubmit = (e) => {
        e.preventDefault();
        const ws = new WebSocket('ws://localhost:8080/ws');
        ws.onopen = () => {
            ws.send(message);
            setMessage('');
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
