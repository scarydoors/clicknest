<script lang="ts">
import Input from "$lib/components/ui/input/input.svelte";
import { isUiNodeInputAttributes, type UiNode } from "@ory/client-fetch";
import * as Field from "$lib/components/ui/field/index";
import { getFormStore } from "../stores/form-store.svelte";

const { node }: { node: UiNode } = $props();

const formStore = getFormStore()
const { form } = formStore.superForm

let attr = $derived(node.attributes);
</script>
{#if isUiNodeInputAttributes(attr)}
<Field.Field>
    <Field.Label>{node.meta.label?.text ?? attr.name}</Field.Label>
    <Input id={attr.name} type={attr.type} required={attr.required} bind:value={$form[attr.name]} />
</Field.Field>
{/if}

