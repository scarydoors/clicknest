<script lang="ts">
import { PUBLIC_KRATOS_API_URL } from "$env/static/public"
import { Configuration, FrontendApi, UpdateRegistrationFlowBodyFromJSON, type RegistrationFlow } from "@ory/client-fetch";
import { onMount } from "svelte";

const frontendClient = new FrontendApi(
    new Configuration({
        headers: {
            Accept: "application/json",
        },
        basePath: PUBLIC_KRATOS_API_URL
    })
)

async function createRegisterFlow() {
    const registerFlow = await frontendClient.createBrowserRegistrationFlow();

    return registerFlow
}


// TODO: create 
onMount(() => {
    createRegisterFlow().then((flow) => {
        console.log(flow);
    })
})
</script>

