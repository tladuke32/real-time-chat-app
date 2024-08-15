import React, { useEffect, useState } from 'react';

function Notifications() {
    const [notifications, setNotifications] = useState([]);
    const wsURL = process.env.REACT_APP_API_URL.replace('http', 'ws') + '/ws';

    useEffect(() => {
        const ws = new WebSocket(wsURL);

        ws.onmessage = (event) => {
            const data = JSON.parse(event.data);
            if (data.type === 'notification') {
                setNotifications(notifs => [...notifs, data.message]);
            }
        };

        return () => {
            ws.close();
        };
    }, []);

    return (
        <div>
            {notifications.map((notif, index) => (
                <div key={index}>{notif}</div>
            ))}
        </div>
    );
}

export default Notifications;
