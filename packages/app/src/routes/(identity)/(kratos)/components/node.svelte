<script lang="ts">
import { type UiNode } from '@ory/client-fetch';
import NodeInput from './node-input.svelte';
import { isUiNodeInput } from '../node';

type NodeProps = { node: UiNode };

let { node }: NodeProps = $props();
</script>

{#if isUiNodeInput(node)}
	<NodeInput {node} />
{/if}
