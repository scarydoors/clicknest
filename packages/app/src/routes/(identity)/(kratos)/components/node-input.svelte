<script lang="ts">
	import { UiNodeInputAttributesTypeEnum } from "@ory/client-fetch";
import type { UiNodeInput } from "../node";
	import NodeHiddenInputRenderer from "./node-hidden-input-renderer.svelte";
import NodeInputButtonRenderer from "./node-input-button-renderer.svelte";
import NodeInputRenderer from "./node-input-renderer.svelte";

const { node }: { node: UiNodeInput } = $props();

const buttonProps: UiNodeInputAttributesTypeEnum[] = [UiNodeInputAttributesTypeEnum.Button, UiNodeInputAttributesTypeEnum.Submit];

</script>
{#if buttonProps.includes(node.attributes.type)}
    <NodeInputButtonRenderer node={node} />
{:else if node.attributes.type === "hidden"}
    <NodeHiddenInputRenderer node={node} />
{:else}
    <NodeInputRenderer node={node} />
{/if}

