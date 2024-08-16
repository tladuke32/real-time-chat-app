import React, { useState } from 'react';

function Login({ onLogin }) {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const apiURL = process.env.REACT_APP_API_URL;

    const handleSubmit = async (e) => {
        e.preventDefault();

        try {
            const response = await fetch(`${apiURL}/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username, password }),
            });

            if (!response.ok) {
                const errorMessage = await response.text();
                console.error('Login failed:', errorMessage);
                alert(`Login failed: ${errorMessage}`);
                return;
            }

            const data = await response.json();

            if (data.token) {
                localStorage.setItem('token', data.token);
                console.log('Login successful, token stored');
                onLogin({ username, token: data.token }); // Trigger the onLogin callback with user details
            } else {
                console.error('Login failed: Token not received');
                alert('Login failed: Token not received');
            }
        } catch (error) {
            console.error('Login request failed:', error);
            alert('Login request failed. Please try again later.');
        }
    };

    return (
        <form onSubmit={handleSubmit} className="mt-3">
            <div className="form-group">
                <label htmlFor="username">Username</label>
                <input
                    type="text"
                    className="form-control"
                    id="username"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    placeholder="Enter username"
                    required
                />
            </div>
            <div className="form-group">
                <label htmlFor="password">Password</label>
                <input
                    type="password"
                    className="form-control"
                    id="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="Password"
                    required
                />
            </div>
            <button type="submit" className="btn btn-primary">Login</button>
        </form>
    );
}

export default Login;

