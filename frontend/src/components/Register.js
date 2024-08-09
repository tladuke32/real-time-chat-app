import React, { useState } from 'react';

function Register() {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const apiURL = process.env.REACT_APP_API_URL;

    fetch(`${apiURL}/api/http://backend:8080`)
        .then(response => response.json())
        .then(data => console.log(data))

    const handleSubmit = async (e) => {
        e.preventDefault();
        const response = await fetch('/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password }),
        });
        const data = await response.json();
        // Handle registration logic
    };

    return (
        <form onSubmit={handleSubmit}>
            <input type="text" value={username} onChange={(e) => setUsername(e.target.value)} placeholder="Username" />
            <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} placeholder="Password" />
            <button type="submit">Register</button>
        </form>
    );
}

export default Register;
