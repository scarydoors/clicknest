import { getContext, setContext } from "svelte";
import { getFlowStore } from "./flow-store.svelte";
import { superForm } from "$lib/utils";
import { getDefaultValues, transformIntoNestedForm, type FormValues } from "../utils/form";
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
                onUpdate({ form }) {
                    if (!form.valid) {
                        return;
                    }

                    console.log(transformIntoNestedForm(form.data));

                    flowStore.updateFlow(
                        transformIntoNestedForm(form.data) as unknown as UpdateRegistrationFlowBody
                    )
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
