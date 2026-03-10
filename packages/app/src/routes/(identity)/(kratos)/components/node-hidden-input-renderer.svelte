<script lang="ts">
import { Input } from '$lib/components/ui/input';
import { getFormStore } from '../form-store.svelte';
import type { UiNodeInput } from '../node';
import { useNodeAttributes } from '../use-node-attributes.svelte';

const { node }: { node: UiNodeInput } = $props();

const formStore = getFormStore();
const { form } = formStore.superForm;

let attr = useNodeAttributes(node.attributes);
</script>

<Input id={attr.name} type={attr.type} required={attr.required} bind:value={$form[attr.name]} />
