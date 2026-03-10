import {
	handleContinueWith,
	handleFlowError,
	type ApiResponse,
	FlowType,
	type RegistrationFlow,
	type SuccessfulNativeRegistration,
	type UpdateRegistrationFlowBody,
	type VerificationFlow,
	type UpdateVerificationFlowBody
} from '@ory/client-fetch';
import { page } from '$app/state';
import { goto } from '$app/navigation';
import { getContext, setContext } from 'svelte';

type FlowTypeMap = {
	[FlowType.Registration]: {
		flow: RegistrationFlow;
		updateBody: UpdateRegistrationFlowBody;
		updateResponse: SuccessfulNativeRegistration;
	};
	[FlowType.Verification]: {
		flow: VerificationFlow;
		updateBody: UpdateVerificationFlowBody;
		updateResponse: VerificationFlow;
	};
};

type Flow<K extends keyof FlowTypeMap> = FlowTypeMap[K]['flow'];

type UpdateBody<K extends keyof FlowTypeMap> = FlowTypeMap[K]['updateBody'];

type UpdateResponse<K extends keyof FlowTypeMap> = FlowTypeMap[K]['updateResponse'];

type FlowStoreProps<K extends keyof FlowTypeMap> = {
	flowType: K;
	createFlow: (params: URLSearchParams) => Promise<ApiResponse<K>>;
	getFlow: (id: string) => Promise<ApiResponse<K>>;
	updateFlow: UpdateFlow<K>;
};

type UpdateFlow<K extends keyof FlowTypeMap> = (
	id: string,
	body: UpdateBody<K>
) => Promise<ApiResponse<UpdateResponse<K>>>;

export class FlowStore<K extends keyof FlowTypeMap> {
	flowType: K;
	flow = $state<Flow<K>>();

	private onUpdateFlow: UpdateFlow<K>;

	constructor({ flowType, createFlow, getFlow, updateFlow }: FlowStoreProps<K>) {
		this.flowType = flowType;

		const errorHandler = handleFlowError({
			onValidationError: () => {},
			// TODO: onRestartFlow, use ory/kratos built-in flow redirect URLs?
			onRestartFlow: () => {
				console.log('restart flow requested');
			},
			onRedirect: (url, external) => {
				if (external) {
					window.location.assign(url);
				} else {
					goto(url);
				}
			}
		});

		this.onUpdateFlow = updateFlow;

		const params = page.url.searchParams;

		const flowId = params.get('flow');
		if (flowId) {
			getFlow(flowId)
				.then((resp) => resp.value())
				.then((resp) => this.handleSetFlow(resp))
				.catch(errorHandler);
			return;
		}

		createFlow(params)
			.then((resp) => resp.value())
			.then((resp) => this.handleSetFlow(resp))
			.catch(errorHandler);
	}

	async updateFlow(body: UpdateBody<K>) {
		const errorHandler = handleFlowError({
			onValidationError: (body: Flow<K>) => (this.flow = body),
			// TODO: onRestartFlow, use ory/kratos built-in flow redirect URLs?
			onRestartFlow: () => {},
			onRedirect: (url, external) => {
				console.log('i want to redirect');
				if (external) {
					window.location.assign(url);
				} else {
					goto(url);
				}
			}
		});

		if (this.flow) {
			await this.onUpdateFlow(this.flow.id, body)
				.then((resp) => resp.value())
				.then((body) => {
					if (this.flowType == FlowType.Registration) {
						const didContinueWith = handleContinueWith(
							(body as UpdateResponse<FlowType.Registration>).continue_with,
							{
								onRedirect: (url, external) => {
									console.log('i want to redirect');
									if (external) {
										window.location.assign(url);
									} else {
										goto(url);
									}
								}
							}
						);

						// eslint-disable-next-line promise/always-return
						if (didContinueWith) {
							return;
						}

						// We did not receive a valid continue_with, but the state flow is still a success. In this case we re-initialize
						// the registration flow which will redirect the user to the default url.
						//onRedirect(registrationUrl(config), true)
					} else if (this.flowType == FlowType.Verification) {
						this.flow = body as UpdateResponse<FlowType.Verification>;
					}
				})
				.catch(errorHandler);
		} else {
			throw new Error('Tried calling `updateFlow` when flow is undefined');
		}
	}

	private async handleSetFlow(flow: Flow<K>) {
		this.flow = flow;

		const params = new URLSearchParams(page.url.searchParams);
		params.set('flow', flow.id);

		goto(`?${params.toString()}`, {
			keepFocus: true,
			noScroll: true
		});
	}
}

const SYMBOL_KEY = 'identity-flow-store';

export function setFlowStore<K extends keyof FlowTypeMap>(
	flowStoreProps: FlowStoreProps<K>
): FlowStore<K> {
	return setContext(Symbol.for(SYMBOL_KEY), new FlowStore(flowStoreProps));
}

export function getFlowStore<K extends keyof FlowTypeMap>(): FlowStore<K> {
	return getContext(Symbol.for(SYMBOL_KEY));
}
