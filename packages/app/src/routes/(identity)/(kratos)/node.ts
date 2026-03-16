import { isUiNodeInputAttributes, type UiNode, type UiNodeInputAttributes } from '@ory/client-fetch';

export type UiNodeInput = UiNode & {
	type: 'input';
	attributes: UiNodeInputAttributes;
};

export function isUiNodeInput(node: UiNode): node is UiNodeInput {
    return isUiNodeInputAttributes(node.attributes);
}
