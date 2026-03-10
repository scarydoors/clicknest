import type { UiNodeInputAttributes } from '@ory/client-fetch';
import { getFormStore } from './form-store.svelte';

export function useNodeAttributes(attributes: UiNodeInputAttributes) {
	const formStore = getFormStore();

	const attrs = $derived({
		name: attributes.name,
		disabled: attributes.disabled || formStore.superForm.submitting,
		required: attributes.required,
		autocomplete: attributes.autocomplete,
		maxLength: attributes.maxlength,
		type: attributes.type
	});

	return attrs;
}
