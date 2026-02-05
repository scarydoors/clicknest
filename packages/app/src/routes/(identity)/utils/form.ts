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

export interface NestedFormValues {
    [key: string]: string | boolean | number | undefined | NestedFormValues
}

// TODO: fix types, I don't really want to use `as NestedFormValues`
// TODO: write tests
// TODO: handle edgecase:
//     transformIntoNestedForm({
//         'what.name': 55,
//         what: 56,
//     })
//     produces `{ what: { name: 55 } }`, shouldn't fail silently
export function transformIntoNestedForm(form: FormValues): NestedFormValues {
    return Object.entries(form).reduce<NestedFormValues>((transformedForm, [key, value]) => {
        const keys = key.split('.');
        
        let obj = transformedForm;
        for (const [idx, key] of keys.entries()) {
            const isLastKey = idx == keys.length - 1;
            if (!(key in obj)) {
                obj[key] = isLastKey ? value : {};
            }

            obj = obj[key] as NestedFormValues;
        }

        return transformedForm;
    }, {})
}
