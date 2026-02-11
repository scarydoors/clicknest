<script lang="ts">
import * as Card from "$lib/components/ui/card/index.js";
	import * as Field from "$lib/components/ui/field/index";
import { getFlowStore } from "../flow-store.svelte";
import { setFormStore } from "../form-store.svelte";
import Node from "./node.svelte";

const flowStore = getFlowStore()
const formStore = setFormStore()
const { enhance } = formStore.superForm

let nodes = $derived(flowStore.flow?.ui.nodes);

</script>

<!-- TODO: dynamic title, description -->
<!-- TODO: different card content, we want to sort nodes differently depending on what stage we are at. -->
<!-- TODO: method selector, I want passwordless -->
<Card.Root>
    <Card.Header>
		<Card.Title>Create an account</Card.Title>
		<Card.Description>Enter your information below to create your account</Card.Description>
    </Card.Header>
    <Card.Content>
        <form use:enhance method="POST">
            <Field.Group>
                {#each nodes as node}
                    <Node node={node}/>
                {/each}
            </Field.Group>
        </form>
    </Card.Content>
</Card.Root>
