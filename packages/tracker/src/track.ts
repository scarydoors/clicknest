import { getConfig } from "./config";
import { monkeypatchPost } from "./utils";

interface EventRequest {
    domain: string;
    kind: EventKind;
    url: string;
    timestamp: Date;
};

const baseUrl = 'http://localhost:6969';

export type EventKind = 'pageview' | (string & {})
export function track(kind: EventKind) {
    const { domain } = getConfig()
    const body: EventRequest = {
        url: window.location.href,
        domain: domain,
        kind,
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


export function setupAutoTrack() {
    let lastPath = window.location.pathname; 
    function maybeTrack() {
        if (lastPath !== window.location.pathname) {
            track("pageview")
        }
        lastPath = window.location.pathname
    }

    monkeypatchPost(window.history, "pushState", maybeTrack)
    monkeypatchPost(window.history, "replaceState", maybeTrack)

    window.addEventListener("popstate", maybeTrack)
}
