import { getConfig } from "./config";
import { monkeypatchPost } from "./utils";

interface EventRequest {
    domain: string;
    kind: EventKind;
    url: string;
    timestamp: Date;
    data?: EventData;
};

const baseUrl = 'http://localhost:6969';

type Flat = string 
| number 
| boolean 
| null 
| undefined
| Date;
export type EventData = Record<string, Flat>
export type EventKind = 'pageview' | (string & {})

export function track(kind: EventKind, data?: EventData) {
    const { domain } = getConfig()
    const body: EventRequest = {
        url: window.location.href,
        domain: domain,
        kind,
        timestamp: new Date(),
        data
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
