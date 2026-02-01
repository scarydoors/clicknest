import { handleFlowError, type ApiResponse, type FlowType } from "@ory/client-fetch";
import { page } from "$app/state";
import { goto } from "$app/navigation";

type Flow = {
    id: string;
}

export class FlowStore<T extends Flow> {
    flow = $state<T>();

    constructor(
        flowType: FlowType,
        createFlow: (params: URLSearchParams) => Promise<ApiResponse<T>>,
        getFlow: (id: string) => Promise<ApiResponse<T>>,
    ) {
        const errorHandler = handleFlowError({
            onValidationError: () => {},
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

        const params = page.url.searchParams;

        const flowId = params.get('flow');
        if (flowId) {
            getFlow(flowId).then((resp) => resp.value()).then(this.handleSetFlow).catch(errorHandler)
            return
        }

        createFlow(params).then((resp) => resp.value()).then(this.handleSetFlow).catch(errorHandler)
    }

    private handleSetFlow = async (flow: T) => {
        this.flow = flow;

        const params = new URLSearchParams(page.url.searchParams);
        params.set('flow', flow.id)
        
        goto(`?${params.toString()}`, {
            keepFocus: true,
            noScroll: true,
        })
    }
}
