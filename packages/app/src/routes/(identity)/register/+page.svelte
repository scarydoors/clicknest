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

      async function submitRegistration() {
    const res = await fetch(
      'http://localhost:4433/self-service/registration?flow=101d13eb-e987-49f8-b593-8e4a1c4d9873',
      {
        method: 'POST',
        credentials: 'include', // REQUIRED for Kratos browser flows
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          method: 'profile',
          csrf_token: 'zGiQgDDuDmATrUJqZ4VyTiA33UoRtThopAV+Fe81RfA=',
          traits: {
            email: 'test.user@example.com',
            name: {
              first: 'Test',
              last: 'User'
            }
          }
        })
      }
    );

    const data = await res.json();
    console.log('Kratos response:', data);
}
</script>

<button on:click={submitRegistration}>test</button>
