<script lang="ts">
	import { getFormStore } from '../form-store.svelte';
	import type { UiNodeInput } from '../node';
	import { useNodeAttributes } from '../use-node-attributes.svelte';
	import * as Field from '$lib/components/ui/field/index';
	import { Input } from '$lib/components/ui/input';
	import { getNodeLabel } from '@ory/client-fetch';

	const { node }: { node: UiNodeInput } = $props();

	const formStore = getFormStore();
	const { form } = formStore.superForm;

	let attr = useNodeAttributes(node.attributes);

	let label = $derived(getNodeLabel(node));
</script>

<Field.Field>
	<Field.Label>{label?.text}</Field.Label>
	<Input id={attr.name} type={attr.type} required={attr.required} bind:value={$form[attr.name]} />
</Field.Field>
