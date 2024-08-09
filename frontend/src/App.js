import React from 'react';
import Login from './components/Login';
import Register from './components/Register';
import Chat from './components/Chat';

function App() {
    return (
        <div className="App">
            <Register />
            <Login />
            <Chat />
        </div>
    );
}

export default App;
