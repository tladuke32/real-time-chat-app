import React, { useState } from 'react';
import Login from './components/Login';
import Register from './components/Register';
import Chat from './components/Chat';
import Notifications from './components/Notifications'; // Make sure to import the Notifications component
import 'bootstrap/dist/css/bootstrap.min.css';

function App() {
    const [user, setUser] = useState(null); // State to manage user authentication

    const handleLogin = (userData) => {
        setUser(userData); // Assume userData contains user details on successful login
        console.log('User logged in:', userData);
    };

    const handleLogout = () => {
        setUser(null); // Clear user data to log out
        console.log('User logged out');
    };

    return (
        <div className="App container mt-3">
            <header className="mb-4">
                <h1>Welcome to Chat App</h1>
                {user && (
                    <div>
                        Logged in as: <strong>{user.username}</strong>
                        <button onClick={handleLogout} className="btn btn-danger ml-2">Logout</button>
                    </div>
                )}
            </header>
            {user ? (
                <>
                    <Notifications user={user} />
                    <Chat user={user} onLogout={handleLogout} />
                </>
            ) : (
                <div className="authentication">
                    <Register />
                    <Login onLogin={handleLogin} />
                </div>
            )}
        </div>
    );
}

export default App;

