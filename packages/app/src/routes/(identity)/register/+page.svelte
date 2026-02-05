<script lang="ts">
import { PUBLIC_KRATOS_API_URL } from "$env/static/public"
import { Configuration, FlowType, FrontendApi } from "@ory/client-fetch";
import { setFlowStore } from "../stores/flow-store.svelte";
import Form from "../components/form.svelte";

const frontendClient = new FrontendApi(
    new Configuration({
        headers: {
            Accept: "application/json",
        },
        credentials: 'include',
        basePath: PUBLIC_KRATOS_API_URL
    })
)

const flowStore = setFlowStore({
    flowType: FlowType.Registration,
    createFlow: (params) => frontendClient.createBrowserRegistrationFlowRaw({
        returnTo: params.get("returnTo") ?? undefined,
        loginChallenge: params.get("loginChallenge") ?? undefined,
        afterVerificationReturnTo: params.get("afterVerificationReturnTo") ?? undefined,
        organization: params.get("organization") ?? undefined,
    }),
    getFlow: (id) => frontendClient.getRegistrationFlowRaw({ id }),
    updateFlow: (id, body) => frontendClient.updateRegistrationFlowRaw({
        flow: id,
        updateRegistrationFlowBody: body
    })
})
</script>

{#if flowStore.flow}
    <Form>
    </Form>
{/if}
