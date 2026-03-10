<script lang="ts">
	import { PUBLIC_KRATOS_API_URL } from '$env/static/public';
	import { Configuration, FlowType, FrontendApi } from '@ory/client-fetch';
	import { setFlowStore } from '../(kratos)/flow-store.svelte';
	import Form from '../(kratos)/components/form.svelte';

	const frontendClient = new FrontendApi(
		new Configuration({
			headers: {
				Accept: 'application/json'
			},
			credentials: 'include',
			basePath: PUBLIC_KRATOS_API_URL
		})
	);

	const flowStore = setFlowStore({
		flowType: FlowType.Verification,
		createFlow: (params) => frontendClient.createBrowserVerificationFlowRaw({}),
		getFlow: (id) => frontendClient.getVerificationFlowRaw({ id }),
		updateFlow: (id, body) =>
			frontendClient.updateVerificationFlowRaw({
				flow: id,
				updateVerificationFlowBody: body
			})
	});
</script>

{#if flowStore.flow}
	<Form></Form>
{/if}
