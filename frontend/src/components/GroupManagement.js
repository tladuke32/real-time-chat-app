import React, { useState, useEffect } from 'react';

function GroupManagement({ user }) {
    const [groups, setGroups] = useState([]);
    const [selectedGroup, setSelectedGroup] = useState(null);
    const [newGroupName, setNewGroupName] = useState('');
    const [groupMessage, setGroupMessage] = useState('');
    const [groupMessages, setGroupMessages] = useState([]);

    useEffect(() => {
        fetchGroups();
    }, []);

    const fetchGroups = async () => {
        try {
            const response = await fetch(`${process.env.REACT_APP_API_URL}/groups`);
            const data = await response.json();
            setGroups(data);
        } catch (error) {
            console.error('Error fetching groups:', error);
        }
    };

    const createGroup = async () => {
        try {
            const response = await fetch(`${process.env.REACT_APP_API_URL}/groups/create`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name: newGroupName }),
            });
            if (response.ok) {
                fetchGroups();
                setNewGroupName('');
            } else {
                console.error('Failed to create group');
            }
        } catch (error) {
            console.error('Error creating group:', error);
        }
    };

    const joinGroup = async (groupId) => {
        try {
            const response = await fetch(`${process.env.REACT_APP_API_URL}/groups/add_member`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ group_id: groupId, user_id: user.id }),
            });
            if (response.ok) {
                setSelectedGroup(groupId);
                FetchGroupMessages(groupId);
            } else {
                console.error('Failed to join group');
            }
        } catch (error) {
            console.error('Error joining group:', error);
        }
    };

    const FetchGroupMessages = async (groupId) => {
        try {
            const response = await fetch(`${process.env.REACT_APP_API_URL}/groups/messages?group_id=${groupId}`);
            const data = await response.json();
            setGroupMessages(data);
        } catch (error) {
            console.error('Error fetching group messages:', error);
        }
    };

    const sendMessageToGroup = async () => {
        try {
            const response = await fetch(`${process.env.REACT_APP_API_URL}/groups/send_message`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ group_id: selectedGroup, user_id: user.id, content: groupMessage }),
            });
            if (response.ok) {
                fetchGroupMessages(selectedGroup);
                setGroupMessage('');
            } else {
                console.error('Failed to send message to group');
            }
        } catch (error) {
            console.error('Error sending message to group:', error);
        }
    };

    return (
        <div className="group-management">
            <h2>Group Management</h2>

            <div className="create-group">
                <h3>Create Group</h3>
                <input
                    type="text"
                    className="form-control"
                    value={newGroupName}
                    onChange={(e) => setNewGroupName(e.target.value)}
                    placeholder="Group Name"
                />
                <button onClick={createGroup} className="btn btn-success mt-2">Create Group</button>
            </div>

            <div className="group-list mt-4">
                <h3>Available Groups</h3>
                <ul className="list-group">
                    {groups.map((group) => (
                        <li key={group.id} className="list-group-item">
                            {group.name}
                            <button onClick={() => joinGroup(group.id)} className="btn btn-primary ml-2">Join Group</button>
                        </li>
                    ))}
                </ul>
            </div>

            {selectedGroup && (
                <div className="group-chat mt-4">
                    <h3>Group Chat</h3>
                    <div className="messages">
                        {groupMessages.map((msg, idx) => (
                            <div key={idx} className="message">
                                <strong>{msg.username || "Anonymous"}:</strong> {msg.content}
                            </div>
                        ))}
                    </div>
                    <form onSubmit={(e) => { e.preventDefault(); sendMessageToGroup(); }} className="message-form">
                        <input
                            type="text"
                            className="form-control"
                            value={groupMessage}
                            onChange={(e) => setGroupMessage(e.target.value)}
                            placeholder="Enter message"
                            required
                        />
                        <button type="submit" className="btn btn-primary mt-2">Send</button>
                    </form>
                </div>
            )}
        </div>
    );
}

export default GroupManagement;
