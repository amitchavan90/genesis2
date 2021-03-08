import { IconName } from "@fortawesome/fontawesome-svg-core"

export enum Action {
	Archive = "Archive",
	Unarchive = "Unarchive",

	SetSKU = "SetSKU",
	SetOrder = "SetOrder",
	SetDistributor = "SetDistributor",
	SetContract = "SetContract",
	SetCarton = "SetCarton",
	SetPallet = "SetPallet",
	SetContainer = "SetContainer",

	DetachFromSKU = "DetachFromSKU",
	DetachFromOrder = "DetachFromOrder",
	DetachFromDistributor = "DetachFromDistributor",
	DetachFromContract = "DetachFromContract",
	DetachFromCarton = "DetachFromCarton",
	DetachFromPallet = "DetachFromPallet",
	DetachFromContainer = "DetachFromContainer",

	SetBonusLoyaltyPoints = "SetBonusLoyaltyPoints",
	InheritCartonHistory = "InheritCartonHistory",
}

export interface ActionItem {
	label: string
	action: Action
	/** The itemName of the value you're setting (Required for Set actions) */
	itemName?: string
	identifier?: string
	icon?: IconName
	iconLight?: boolean
}

const Archive: ActionItem[] = [
	{ label: "Archive", action: Action.Archive, icon: "archive" },
	{ label: "Unarchive", action: Action.Unarchive, icon: "undo" },
]

const Products: ActionItem[] = [
	{ label: "Detach from Carton", action: Action.DetachFromCarton, itemName: "carton" },
	{ label: "Set Carton", action: Action.SetCarton, itemName: "carton" },
	{ label: "Inherit Carton History", action: Action.InheritCartonHistory, itemName: "carton" },
	{ label: "Remove SKU", action: Action.DetachFromSKU, itemName: "sku" },
	{ label: "Set SKU", action: Action.SetSKU, itemName: "sku" },
	{ label: "Remove Distributor", action: Action.DetachFromDistributor, itemName: "distributor" },
	{ label: "Set Distributor", action: Action.SetDistributor, itemName: "distributor" },
	{ label: "Remove Livestock Specification", action: Action.DetachFromContract, itemName: "contract" },
	{ label: "Set Livestock Specification", action: Action.SetContract, itemName: "contract", identifier: "name" },
	{ label: "Remove from Order", action: Action.DetachFromOrder, itemName: "order" },
	{ label: "Set Order", action: Action.SetOrder, itemName: "order" },
	...Archive,
	{ label: "Set Bonus Loyalty Points", action: Action.SetBonusLoyaltyPoints, icon: "star" },
]

const Cartons: ActionItem[] = [
	{ label: "Detach from Pallet", action: Action.DetachFromPallet, itemName: "pallet" },
	{ label: "Set Pallet", action: Action.SetPallet, itemName: "pallet" },
	{ label: "Remove Distributor", action: Action.DetachFromDistributor, itemName: "distributor" },
	{ label: "Set Distributor", action: Action.SetDistributor, itemName: "distributor" },
	{ label: "Remove Livestock Specification", action: Action.DetachFromContract, itemName: "contract" },
	{ label: "Set Livestock Specification", action: Action.SetContract, itemName: "contract", identifier: "name" },
	...Archive,
	{ label: "Set Bonus Loyalty Points", action: Action.SetBonusLoyaltyPoints, icon: "star" },
]

const Pallets: ActionItem[] = [
	{ label: "Detach from Container", action: Action.DetachFromContainer, itemName: "container" },
	{ label: "Set Container", action: Action.SetContainer, itemName: "container" },
	...Archive,
	{ label: "Set Bonus Loyalty Points", action: Action.SetBonusLoyaltyPoints, icon: "star" },
]

const Containers: ActionItem[] = [...Archive, { label: "Set Bonus Loyalty Points", action: Action.SetBonusLoyaltyPoints, icon: "star" }]

export const ActionItemSet = {
	Archive,
	Products,
	Cartons,
	Pallets,
	Containers,
}
