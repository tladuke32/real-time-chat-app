import React, { useState } from 'react';
import { BrowserRouter as Router, Route, Routes, Link } from 'react-router-dom';
import Login from './components/Login';
import Register from './components/Register';
import Chat from './components/Chat';
import GroupManagement from './components/GroupManagement';
import UserProfile from './components/UserProfile';
import 'bootstrap/dist/css/bootstrap.min.css';

function App() {
    const [user, setUser] = useState(null); // State to manage user authentication

    const handleLogin = (userData) => {
        setUser(userData); // Set user data on successful login
        console.log('User logged in:', userData);
    };

    const handleLogout = () => {
        setUser(null); // Clear user data to log out
        console.log('User logged out');
    };

    return (
        <Router>
            <div className="App container mt-3">
                <header className="mb-4">
                    <h1>Welcome to Chat That</h1>
                    {user && (
                        <div>
                            Logged in as: <strong>{user.username}</strong>
                            <button onClick={handleLogout} className="btn btn-danger ml-2">Logout</button>
                            <Link to="/chat" className="btn btn-primary ml-2">Chat</Link>
                            <Link to="/groups" className="btn btn-secondary ml-2">Groups</Link>
                            <Link to="/profile" className="btn btn-info ml-2">Profile</Link>
                        </div>
                    )}
                </header>
                {user ? (
                    <Routes>
                        <Route path="/chat" element={<Chat user={user} />} />
                        <Route path="/groups" element={<GroupManagement user={user} />} />
                        <Route path="/profile" element={<UserProfile user={user} />} />
                        <Route path="/" element={<Chat user={user} />} />
                    </Routes>
                ) : (
                    <div className="authentication">
                        <Register />
                        <Login onLogin={handleLogin} />
                    </div>
                )}
            </div>
        </Router>
    );
}

export default App;
