import { isUiNodeInputAttributes, type UiNode } from "@ory/client-fetch"

export type FormValues = Record<string, string | boolean | number | undefined>

export function getDefaultValues(flow?: {
  active?: string
  ui: { nodes: UiNode[] }
}): FormValues {
    return flow?.ui.nodes.reduce<FormValues>((form, node) => {
        const attrs = node.attributes;
        if (isUiNodeInputAttributes(attrs)) {
            form[attrs.name] = attrs.value ?? "";
        }

        return form;
    }, {}) ?? {};
}
