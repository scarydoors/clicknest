import type { UiNodeInputAttributes } from "@ory/client-fetch";
import { getFormStore } from "./form-store.svelte";

export function useNodeAttributes(attributes: UiNodeInputAttributes) {
    const formStore = getFormStore();

    return $derived({
        disabled: attributes.disabled || formStore.superForm.submitting,
        autocomplete: attributes.autocomplete,
        maxLength: attributes.maxlength,
        type: attributes.type,
    });
}
