import React, { useState } from 'react';

function Register() {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const apiURL = process.env.REACT_APP_API_URL; // Ensure this is correctly set in Dockerfile

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const response = await fetch(`${apiURL}/register`, { // Use the API URL and endpoint
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username, password }),
            });
            if (response.ok) {
                console.log('Registration successful');
                // Handle successful registration logic
            } else {
                console.error('Registration failed', await response.text());
            }
        } catch (error) {
            console.error('Error:', error);
        }
    };

    return (
        <form onSubmit={handleSubmit}>
            <input
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                placeholder="Username"
                required
            />
            <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Password"
                required
            />
            <button type="submit">Register</button>
        </form>
    );
}

export default Register;
