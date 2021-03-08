import { query } from "../../graphql/queries"
import { Value } from "baseui/select"
import React from "react"
import { TablePicker } from "../../components/tablePicker"
import { LatestTrackActionColumn } from "../../components/common"

export interface TablePickerProps {
	/* select value */
	setValue: (v: Value) => void
	value?: Value
	/** The column to use for the label */
	labelResolver?: (value: any) => void
}

const Cartons = (props: TablePickerProps) => (
	<TablePicker
		query={query.CARTONS}
		{...props}
		queryName={"cartons"}
		columns={[
			{ label: "Cartons", value: "code" },
			{ label: "Weight", value: "weight" },
			{ label: "Products Amount", value: "productCount" },
			{
				label: "Last Track Action",
				value: "latestTrackAction",
				resolver: row => <LatestTrackActionColumn value={row.latestTrackAction} />,
			},
			{ label: "Description", value: "description" },
			{ label: "Meat Type", value: "meatType" },
			{ label: "Date Created", value: "createdAt", dateTime: true },
		]}
	/>
)

const Pallets = (props: TablePickerProps) => (
	<TablePicker
		query={query.PALLETS}
		{...props}
		queryName={"pallets"}
		columns={[
			{ label: "Pallets", value: "code" },
			{ label: "Cartons Amount", value: "cartonCount" },
			{
				label: "Last Track Action",
				value: "latestTrackAction",
				resolver: row => <LatestTrackActionColumn value={row.latestTrackAction} />,
			},
			{ label: "Description", value: "description" },
			{ label: "Date Created", value: "createdAt", dateTime: true },
		]}
	/>
)

const Containers = (props: TablePickerProps) => (
	<TablePicker
		query={query.CONTAINERS}
		{...props}
		queryName={"containers"}
		columns={[
			{ label: "Containers", value: "code" },
			{ label: "Pallets Amount", value: "palletCount" },
			{ label: "Description", value: "description" },
			{ label: "Date Created", value: "createdAt", dateTime: true },
		]}
	/>
)

const Skus = (props: TablePickerProps) => (
	<TablePicker
		query={query.SKUS}
		{...props}
		queryName={"skus"}
		columns={[
			{ label: "SKU", value: "name" },
			{ label: "Products Amount", value: "productCount" },
			{ label: "Date Created", value: "createdAt", dateTime: true },
		]}
	/>
)

const Distributors = (props: TablePickerProps) => (
	<TablePicker
		query={query.DISTRIBUTORS}
		{...props}
		queryName={"distributors"}
		columns={[
			{ label: "Distributor", value: "name" },
			{ label: "Date Created", value: "createdAt", dateTime: true },
		]}
	/>
)

const LivestockSpecifications = ({ labelResolver, ...rest }: TablePickerProps) => (
	<TablePicker
		query={query.CONTRACTS}
		{...rest}
		labelResolver={labelResolver ? labelResolver : value => `${value.name} (${value.supplierName})`}
		queryName={"contracts"}
		columns={[
			{ label: "Livestock Specification", value: "name" },
			{ label: "Supplier Name", value: "supplierName" },
			{ label: "Date Signed", value: "dateSigned", dateTime: true },
		]}
	/>
)

const Orders = (props: TablePickerProps) => (
	<TablePicker
		query={query.ORDERS}
		{...props}
		queryName={"orders"}
		columns={[
			{ label: "Order", value: "code" },
			{ label: "Products Amount", value: "productCount" },
			{ label: "Date Created", value: "createdAt", dateTime: true },
		]}
	/>
)

export const TablePickers = {
	Cartons,
	Pallets,
	Containers,
	Skus,
	Distributors,
	LivestockSpecifications,
	Orders,
}
