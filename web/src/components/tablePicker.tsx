import { Value } from "baseui/select"
import { SearchFilter } from "../types/types"
import { FilterOptionItem, SortByOptionItem, FilterOption } from "../types/enums"
import { ItemListColumn } from "./itemList"
import React from "react"
import { PaginationBar } from "./paginationBar"
import { useQuery } from "@apollo/react-hooks"
import { TableBuilder, TableBuilderColumn } from "baseui/table-semantic"
import { Checkbox } from "baseui/checkbox"
import { timeAgo } from "../helpers/time"
import { Caption1 } from "baseui/typography"

export interface TablePickerProps {
	value?: Value
	setValue: (value: Value) => void
	filter?: SearchFilter
	extraFilterOptions?: FilterOptionItem[]
	extraSortByOptions?: SortByOptionItem[]
	query: any
	identifier?: string
	/** The column to use for the label */
	labelResolver?: (value: any) => void
	queryName: string
	columns?: ItemListColumn[]
}
export const TablePicker = (props: TablePickerProps) => {
	const [search, setSearch] = React.useState<SearchFilter>({ filter: FilterOption.Active })
	const identifier = props.identifier || "code"

	// Pagination
	const [offset, setOffset] = React.useState(0)
	const [total, setTotal] = React.useState(0)
	const [items, setItems] = React.useState<any[]>([])

	const pagination = (
		<PaginationBar
			searchType={"button"}
			offset={offset}
			total={total}
			limit={20}
			defaultFilter={props.filter && props.filter.filter}
			defaultSortOption={props.filter && props.filter.sortBy}
			defaultSortDir={props.filter && props.filter.sortDir}
			setSearch={(value: SearchFilter) => setSearch(value)}
			setOffset={(value: number) => setOffset(value)}
			extraFilterOptions={props.extraFilterOptions}
			extraSortByOptions={props.extraSortByOptions}
		/>
	)

	// query
	const { data, loading, error } = useQuery(props.query, {
		variables: {
			search,
			limit: 20,
			offset: offset,
		},
	})

	React.useEffect(() => {
		if (!data || !data[props.queryName]) return
		setItems(data[props.queryName][props.queryName])
		setTotal(data[props.queryName].total)
	}, [data, loading, error])

	function toggle(event: any) {
		const { name, checked } = event.currentTarget
		setItems((items) => {
			return items.map((row) => {
				if (String(row[identifier]) === name) {
					props.setValue([{ id: row["id"], label: props.labelResolver ? props.labelResolver(row) : row[identifier] }])
				}
				return {
					...row,
					selected: String(row[identifier]) === name ? checked : false,
				}
			})
		})
	}

	return (
		<div>
			{pagination}
			{items.length < 1 ? (
				<Caption1>{`No ${props.queryName} in the database`}</Caption1>
			) : (
				<TableBuilder
					data={items}
					overrides={{
						Root: {
							style: {
								height: "60vh",
							},
						},
					}}
				>
					<TableBuilderColumn
						overrides={{
							TableHeadCell: { style: { width: "1%" } },
							TableBodyCell: { style: { width: "1%" } },
						}}
					>
						{(row) => {
							return <Checkbox name={row[identifier]} checked={row.selected} onChange={toggle} />
						}}
					</TableBuilderColumn>
					{props.columns &&
						props.columns.map((c) => {
							return (
								<TableBuilderColumn
									key={`column-${c.label}`}
									id={c.label}
									header={c.label}
									sortable={c.filterable}
									overrides={{
										TableBodyCell: {
											style: {
												verticalAlign: "center",
												color: "inherit",
											},
										},
									}}
								>
									{(item) => {
										// Custom resolver
										if (c.resolver !== undefined) return c.resolver(item)
										if (!c.value || !item[c.value]) return <></>

										// Get cell value
										let value = item[c.value]
										if (c.subValues !== undefined) c.subValues.forEach((s) => (value = value[s]))
										if (!value && c.subValuesAlt !== undefined) {
											value = item[c.value]
											c.subValuesAlt.forEach((s) => (value = value[s]))
										}

										// Date/Time?
										if (c.dateTime) {
											return (
												<div>
													<div>{timeAgo(value)}</div>
													<div>{`(${new Date(value).toLocaleString()})`}</div>
												</div>
											)
										}
										// String
										return <div>{value.toString()}</div>
									}}
								</TableBuilderColumn>
							)
						})}
				</TableBuilder>
			)}
		</div>
	)
}
