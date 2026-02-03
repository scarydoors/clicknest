import { handleFlowError, type ApiResponse, type FlowType, type RegistrationFlow, type SuccessfulNativeRegistration, type UpdateRegistrationFlowBody } from "@ory/client-fetch";
import { page } from "$app/state";
import { goto } from "$app/navigation";

type Flow = {
    id: string;
}

type UpdateBody<T extends Flow> =
    T extends RegistrationFlow ? UpdateRegistrationFlowBody :
    never;

type UpdateResponse<T extends Flow> = 
    T extends RegistrationFlow ? SuccessfulNativeRegistration :
    never;

type FlowStoreConfig<T extends Flow> = {
    flowType: FlowType,
    createFlow: (params: URLSearchParams) => Promise<ApiResponse<T>>,
    getFlow: (id: string) => Promise<ApiResponse<T>>,
    updateFlow: UpdateFlow<T>,
}

type UpdateFlow<T extends Flow> = (id: string, body: UpdateBody<T>) => Promise<ApiResponse<UpdateResponse<T>>>;

export class FlowStore<T extends Flow> {
    flow = $state<T>();

    private onUpdateFlow: UpdateFlow<T>

    constructor({
        flowType,
        createFlow,
        getFlow,
        updateFlow
    }: FlowStoreConfig<T>) {
        const errorHandler = handleFlowError({
            onValidationError: () => {},
            // TODO: onRestartFlow, use ory/kratos built-in flow redirect URLs?
            onRestartFlow: () => {
                console.log("restart flow requested");
            },
            onRedirect: (url, external) => {
                if (external) {
                    window.location.assign(url)
                } else {
                    goto(url)
                }
            }
        })

        this.onUpdateFlow = updateFlow

        const params = page.url.searchParams;

        const flowId = params.get('flow');
        if (flowId) {
            getFlow(flowId).then((resp) => resp.value()).then((resp) => this.handleSetFlow(resp)).catch(errorHandler)
            return
        }

        createFlow(params).then((resp) => resp.value()).then((resp) => this.handleSetFlow(resp)).catch(errorHandler)
    }

    updateFlow(body: UpdateBody<T>) {
        const errorHandler = handleFlowError({
            onValidationError: (body: T) => this.flow = body,
            // TODO: onRestartFlow, use ory/kratos built-in flow redirect URLs?
            onRestartFlow: () => {},
            onRedirect: (url, external) => {
                if (external) {
                    window.location.assign(url)
                } else {
                    goto(url)
                }
            }
        })

        if (this.flow) {
            this.onUpdateFlow(this.flow.id, body).then((resp) => resp.value()).then(() => {
                throw new Error("not implemented");
            }).catch(errorHandler)
        } else {
            throw new Error("Tried calling `updateFlow` when flow is undefined");
        }
    }

    private async handleSetFlow(flow: T) {
        this.flow = flow;

        const params = new URLSearchParams(page.url.searchParams);
        params.set('flow', flow.id)
        
        goto(`?${params.toString()}`, {
            keepFocus: true,
            noScroll: true,
        })
    }
}
