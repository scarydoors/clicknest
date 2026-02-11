import type { UiNode, UiNodeInputAttributes } from "@ory/client-fetch";

export type UiNodeInput = UiNode & {
    type: "input",
    attributes: UiNodeInputAttributes
}
