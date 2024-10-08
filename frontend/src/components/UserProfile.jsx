import React, { useState, useEffect } from 'react';

function UserProfile({ user }) {
    const [profile, setProfile] = useState({ username: '', email: '' }); // Adjust fields as necessary
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const apiURL = process.env.REACT_APP_API_URL;

    useEffect(() => {
        const fetchUserProfile = async () => {
            try {
                const response = await fetch(`'${process.env.REACT_APP_API_URL}/user/${user.username}`);
                if (!response.ok) {
                    throw new Error('Failed to fetch profile');
                }
                const data = await response.json();
                setProfile(data);
            } catch (err) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        fetchUserProfile();
    }, [apiURL, user.username]);

    const handleInputChange = (e) => {
        const { name, value } = e.target;
        setProfile({ ...profile, [name]: value });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const response = await fetch(`'${process.env.REACT_APP_API_URL}/user/${user.id}/update`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(profile),
            });
            if (!response.ok) {
                throw new Error('Failed to update profile');
            }
            alert('Profile updated successfully');
        } catch (err) {
            setError(err.message);
        }
    };

    if (loading) return <div>Loading...</div>;
    if (error) return <div>Error: {error}</div>;

    return (
        <div>
            <h2>Your Profile</h2>
            <form onSubmit={handleSubmit}>
                <div className="form-group">
                    <label htmlFor="username">Username</label>
                    <input
                        type="text"
                        className="form-control"
                        id="username"
                        name="username"
                        value={profile.username}
                        onChange={handleInputChange}
                        disabled
                    />
                </div>
                <div className="form-group">
                    <label htmlFor="email">Email</label>
                    <input
                        type="email"
                        className="form-control"
                        id="email"
                        name="email"
                        value={profile.email}
                        onChange={handleInputChange}
                        placeholder="Enter email"
                    />
                </div>
                {/* Add more fields as needed */}
                <button type="submit" className="btn btn-primary">Update Profile</button>
            </form>
        </div>
    );
}

export default UserProfile;
