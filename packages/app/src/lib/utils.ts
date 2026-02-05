import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";
import { superForm as originalSuperForm } from "sveltekit-superforms"

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChild<T> = T extends { child?: any } ? Omit<T, "child"> : T;
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChildren<T> = T extends { children?: any } ? Omit<T, "children"> : T;
export type WithoutChildrenOrChild<T> = WithoutChildren<WithoutChild<T>>;
export type WithElementRef<T, U extends HTMLElement = HTMLElement> = T & { ref?: U | null };

export function superForm(
    form: Parameters<typeof originalSuperForm>[0],
    formOptions: Parameters<typeof originalSuperForm>[1],
): ReturnType<typeof originalSuperForm> {
    return originalSuperForm(form, {
        SPA: true,
        ...formOptions
    })
}

