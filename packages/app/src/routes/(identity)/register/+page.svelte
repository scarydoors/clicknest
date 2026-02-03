<script lang="ts">
import { PUBLIC_KRATOS_API_URL } from "$env/static/public"
import { Configuration, FlowType, FrontendApi, type UiNodeInputAttributes } from "@ory/client-fetch";
	import { FlowStore } from "../stores/flow-store.svelte";

const frontendClient = new FrontendApi(
    new Configuration({
        headers: {
            Accept: "application/json",
        },
        credentials: 'include',
        basePath: PUBLIC_KRATOS_API_URL
    })
)

const flowStore = new FlowStore({
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

function update() {
    flowStore.updateFlow({
        method: "profile",
        traits: {
            email: "what@email.com",
            name: {
                first: "yeah",
                last: "nah",
            }
        },
        csrf_token: "uakS0RAwm3StSU/fUy8qs1hV6bW3S7tFqFcfqpoP19I="
    })
}
$inspect(flowStore.flow);
</script>

<button onclick={update}>test</button>
