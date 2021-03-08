import { useQuery } from "@apollo/react-hooks"
import { Checkbox } from "baseui/checkbox"
import { Select, Value } from "baseui/select"
import { TableBuilder, TableBuilderColumn } from "baseui/table-semantic"
import * as React from "react"
import { SmallItemLink } from "../components/smallItemLink"
import { graphql } from "../graphql"
import { timeAgo } from "../helpers/time"
import { TablePickers } from "../pages/admin/tablePickers"
import { GetItemIcon } from "../themeOverrides"
import { FilterOption, FilterOptionItem, SortByOptionItem } from "../types/enums"
import { SearchFilter, SKU } from "../types/types"
import { ItemListColumn } from "./itemList"
import { SKUItemPreview } from "./itemPreview"
import { PaginationBar } from "./paginationBar"

interface ItemSelectListProps {
	itemName: string
	value?: Value
	setValue: (value: Value) => void
	limit?: number
	filter?: SearchFilter
	query: any
	/** Name of list returned from query (defaults to itemName + s) */
	queryName?: string
	link?: string
	/** Default tab to open to (eg: products) */
	hash?: string
	/** The main value to id the item (default: "code") */
	identifier?: string
	/** The column to use for the label */
	labelResolver?: (value: any) => void
	clearable?: boolean
	disableSearch?: boolean
	error?: boolean
	/** if true renders a table instead of a selectbox*/
	useTable?: boolean
}
export const ItemSelectList = (props: ItemSelectListProps) => {
	const { itemName, value, setValue, limit, filter, query, hash, clearable, disableSearch, labelResolver } = props
	const queryName = props.queryName || `${itemName}s`
	const link = props.link || `/portal/${itemName}`
	const icon = GetItemIcon(itemName, true)
	const identifier = props.identifier || "code"

	const [search, setSearch] = React.useState<SearchFilter>(filter || { filter: FilterOption.Active })
	const [items, setItems] = React.useState<{ id: string; label: string }[]>([])

	const { data, loading, error } = useQuery(query, {
		variables: {
			search,
			limit: limit || 20,
			offset: 0,
		},
	})
	React.useEffect(() => {
		if (!data || !data[queryName]) return
		setItems(data[queryName][queryName].map((c: any) => ({ id: c.id, label: labelResolver ? labelResolver(c) : c[identifier] })))
	}, [data, loading, error])

	React.useEffect(() => {
		if (!disableSearch) setSearch({ search: value && value.length > 0 ? (value[0].label as string) : "" })
	}, [value])

	// returns a tablePicker based on query name
	const getTablePicker = (q: string) => {
		switch (q) {
			case "containers":
				return <TablePickers.Containers setValue={setValue} value={value} labelResolver={labelResolver} />
			case "pallets":
				return <TablePickers.Pallets setValue={setValue} value={value} labelResolver={labelResolver} />
			case "cartons":
				return <TablePickers.Cartons setValue={setValue} value={value} labelResolver={labelResolver} />
			case "skus":
				return <TablePickers.Skus setValue={setValue} value={value} labelResolver={labelResolver} />
			case "distributors":
				return <TablePickers.Distributors setValue={setValue} value={value} labelResolver={labelResolver} />
			case "contracts": // note: this is also "livestock specifications"
				return <TablePickers.LivestockSpecifications setValue={setValue} value={value} labelResolver={labelResolver} />
			case "orders":
				return <TablePickers.Orders setValue={setValue} value={value} labelResolver={labelResolver} />
		}
		return <div>Something went wrong</div>
	}

	return props.useTable ? (
		// table picker for "set carton" or "set pallet" or "set conainer" actions
		getTablePicker(queryName)
	) : (
		<Select
			openOnClick
			autoFocus
			clearable={clearable}
			type={disableSearch ? undefined : "search"}
			options={items}
			labelKey="label"
			valueKey="id"
			value={value}
			onChange={({ value }) => setValue(value)}
			onInputChange={(e) => !disableSearch && setSearch({ ...filter, search: e.currentTarget.value })}
			getValueLabel={({ option }) => <SmallItemLink code={option.label as string} link={link} hash={hash} icon={icon.icon} iconLight={icon.light} disabled />}
			getOptionLabel={({ option }) =>
				option ? <SmallItemLink code={option.label as string} link={link} hash={hash} icon={icon.icon} iconLight={icon.light} disabled /> : <div />
			}
			overrides={{
				Input: {
					style: {
						display: value && value.length > 0 ? "none" : "inline-block",
					},
				},
			}}
			error={props.error}
		/>
	)
}

interface SKUSelectListProps {
	value?: Value
	setValue: (value: Value) => void
	limit?: number
	filter?: SearchFilter
	clearable?: boolean
	disableSearch?: boolean
	error?: boolean
}
export const SKUSelectList = (props: SKUSelectListProps) => {
	const { value, setValue, limit, filter, clearable, disableSearch } = props

	const [search, setSearch] = React.useState<SearchFilter>(filter || {})
	const [items, setItems] = React.useState<{ id: string; label: string }[]>([])

	const { data, loading, error } = useQuery(graphql.query.SKUS, {
		variables: {
			search,
			limit: limit || 20,
			offset: 0,
		},
	})
	React.useEffect(() => {
		if (!data || !data.skus) return
		setItems(data.skus.skus.map((c: SKU) => ({ id: c.id, label: c.name, sku: c })))
	}, [data, loading, error])

	React.useEffect(() => {
		if (!disableSearch) setSearch({ search: value && value.length > 0 ? (value[0].label as string) : "" })
	}, [value])

	return (
		<Select
			clearable={clearable}
			type={disableSearch ? undefined : "search"}
			options={items}
			labelKey="label"
			valueKey="id"
			value={value}
			onChange={({ value }) => setValue(value)}
			onInputChange={(e) => !disableSearch && setSearch({ ...filter, search: e.currentTarget.value })}
			getValueLabel={({ option }) => (option && option.sku ? <SKUItemPreview sku={option.sku} /> : <div />)}
			getOptionLabel={({ option }) => (option && option.sku ? <SKUItemPreview sku={option.sku} /> : <div />)}
			overrides={{
				ValueContainer: { style: { minHeight: "105px" } },
				Input: {
					style: {
						display: value && value.length > 0 ? "none" : "inline-block",
					},
				},
			}}
			error={props.error}
		/>
	)
}
