import { getContext, setContext } from "svelte";
import { getFlowStore } from "./flow-store.svelte";
import { superForm } from "$lib/utils";
import { getDefaultValues, transformIntoNestedForm, type FormValues } from "./form.ts";
import type { UpdateRegistrationFlowBody } from "@ory/client-fetch";
import type { SuperForm } from "sveltekit-superforms";

export class FormStore {
    flowStore = getFlowStore();
    superForm: SuperForm<FormValues>

    constructor() {
        const flowStore = getFlowStore();
        this.superForm = superForm(
            getDefaultValues(this.flowStore.flow),
            {
                validators: false,
                resetForm: false,
                async onUpdate({ form }) {
                    if (!form.valid) {
                        return;
                    }

                    console.log(transformIntoNestedForm(form.data));
                    await flowStore.updateFlow(
                        transformIntoNestedForm(form.data) as unknown as UpdateRegistrationFlowBody
                    )

                    form.data = getDefaultValues(flowStore.flow)
                }
            }
        )
    }
}

const SYMBOL_KEY = "identity-form-store";

export function setFormStore(): FormStore {
    return setContext(Symbol.for(SYMBOL_KEY), new FormStore());
}

export function getFormStore(): FormStore {
    return getContext(Symbol.for(SYMBOL_KEY));
}
