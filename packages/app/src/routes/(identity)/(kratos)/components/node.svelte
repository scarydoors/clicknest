<script lang="ts">
import { isUiNodeInputAttributes, type UiNode } from "@ory/client-fetch";
import NodeInput from "./node-input.svelte";

type NodeProps = { node: UiNode };

let { node }: NodeProps = $props();

</script>
{#if isUiNodeInputAttributes(node.attributes)}
    <NodeInput node={node} />
{/if}
