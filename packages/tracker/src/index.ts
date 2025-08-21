import { setupAutoTrack, track } from "./track";
import { Config, setConfig } from "./config";

function init(config: Config) {
    setConfig(config)

    setupAutoTrack()
    track("pageview")
}

export { init, Config }
