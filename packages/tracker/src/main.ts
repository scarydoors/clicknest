export interface EventRequest {
    domain: string;
    kind: string;
    url: string;
    timestamp: Date;
};

const baseUrl = 'http://localhost:6969';
function track() {
    const body: EventRequest = {
        url: window.location.href,
        domain: "localhost",
        kind: "pageview",
        timestamp: new Date(),
    };
    fetch(`${baseUrl}/api/event`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(body)
    });
}

track();
