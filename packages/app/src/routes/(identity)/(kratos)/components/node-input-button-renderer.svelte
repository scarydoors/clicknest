<script lang="ts">
import Button from '$lib/components/ui/button/button.svelte';
import { getNodeLabel } from '@ory/client-fetch';
import type { UiNodeInput } from '../node';
import { getFormStore } from '../form-store.svelte';

type NodeInputButtonProps = { node: UiNodeInput };

let { node }: NodeInputButtonProps = $props();

let label = $derived(getNodeLabel(node));

const formStore = getFormStore();
const { form } = formStore.superForm;
</script>

<Button
	type={node.attributes.type == "submit" ? "submit" : "button"}
	onclick={() => {
		$form[node.attributes.name] = node.attributes.value;
        if (node.attributes.name == "screen") {
            $form.method = "profile";
        }
	}}>{label?.text}</Button
>
