<script lang="ts">
import { PUBLIC_KRATOS_API_URL } from "$env/static/public"
import { Configuration, FlowType, FrontendApi } from "@ory/client-fetch";
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

const flowStore = new FlowStore(
    FlowType.Registration,
    (params) => frontendClient.createBrowserRegistrationFlowRaw({
        returnTo: params.get("returnTo") ?? undefined,
        loginChallenge: params.get("loginChallenge") ?? undefined,
        afterVerificationReturnTo: params.get("afterVerificationReturnTo") ?? undefined,
        organization: params.get("organization") ?? undefined,
    }),
    (id) => frontendClient.getRegistrationFlowRaw({ id })
)

$inspect(flowStore.flow);
</script>

