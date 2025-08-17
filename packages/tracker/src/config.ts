export interface Config {
    domain: string;
}

let cfg: Config | undefined;
export function setConfig(config: Config) {
    if (cfg) {
        throw new Error('Tracker is already initialized.');
    }
    
    cfg = config
}

export function getConfig(): Config {
    if (!cfg) {
        throw new Error('Tracker is not initialized. Use provided init() function.');
    }

    return cfg
}
