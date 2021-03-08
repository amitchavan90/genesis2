import * as React from "react"
import { Input } from "baseui/input"
import { Button, ButtonOverrides } from "baseui/button"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { useStyletron } from "baseui"
import { useForm } from "react-hook-form"
import { StatefulPopover } from "baseui/popover"
import { StatefulMenu } from "baseui/menu"
import { SearchFilter } from "../types/types"
import { FilterOption, SortByOption, SortDir, FilterOptionItem, SortByOptionItem } from "../types/enums"
import { MenuOptionLabel, MenuListOverride, MenuOptionStyleOverride } from "./menuList"
import { useHistory } from "react-router-dom"
import { Label1 } from "baseui/typography"
import { ChevronLeft, ChevronRight } from "baseui/icon"

const FILTER_OPTIONS: FilterOptionItem[] = [
	{ label: "All", id: FilterOption.All },
	{ label: "Active", id: FilterOption.Active },
	{ label: "Archived", id: FilterOption.Archived },
]
const SORT_OPTION = { label: "Sort...", id: undefined }
const SORT_BY_OPTIONS = [
	{ label: "Date Created", id: SortByOption.DateCreated },
	{ label: "Date Updated", id: SortByOption.DateUpdated },
	{ label: "Alphabetical", id: SortByOption.Alphabetical },
]
const SORT_ORDER_OPTIONS = [
	{ label: "Ascending", id: SortDir.Ascending },
	{ label: "Descending", id: SortDir.Descending },
]

interface PaginationBarProps {
	offset: number
	total: number
	setSearch: (value: SearchFilter) => void
	setOffset: (value: number) => void

	limit?: number
	defaultFilter?: FilterOption
	defaultSortOption?: SortByOption
	defaultSortDir?: SortDir

	extraFilterOptions?: FilterOptionItem[]
	extraSortByOptions?: SortByOptionItem[]

	searchType?: "button" | "submit" | "reset" | undefined
}

export const PaginationBar = (props: PaginationBarProps) => {
	const history = useHistory()
	const searchArgs = new URLSearchParams(window.location.search)

	const pageArg = searchArgs.get("page")
	const searchArg = searchArgs.get("search")

	const { offset, total } = props

	const limit = props.limit || 20

	const totalPages = Math.ceil(total / limit)

	const [css] = useStyletron()
	const flex = css({
		display: "flex",
	})
	const flexOne = css({
		display: "flex",
		flex: 1,
	})
	const totalPageStyle = css({
		alignSelf: "center",
		marginRight: "15px",
	})

	const { register, handleSubmit, getValues, setValue } = useForm<{ search: string }>()
	const [filter, setFilter] = React.useState<FilterOption>(props.defaultFilter || FilterOption.Active)
	const [sortBy, setSortBy] = React.useState<SortByOption>(props.defaultSortOption || SortByOption.DateUpdated)
	const [sortDir, setSortDir] = React.useState<SortDir>(props.defaultSortDir || SortDir.Descending)
	const [searchText, setSearchText] = React.useState("")

	const [pageField, setPageField] = React.useState((offset / limit + 1).toString())

	const onSearchChange = (search: string) => {
		props.setOffset(0)
		setSearchText(search)
		props.setSearch({ search, filter, sortBy, sortDir })
	}

	React.useEffect(() => {
		onSearchChange(getValues().search)
	}, [filter, sortBy, sortDir])

	React.useEffect(() => {
		// init search bar/ offset
		props.setOffset((+(pageArg || 0) - 1 <= 0 ? 0 : +(pageArg || 0) - 1) * limit)
		props.setSearch({ search: searchArg || "", filter, sortBy, sortDir })
		setValue("search", searchArg || "")
	}, [])

	React.useEffect(() => {
		setPageField((offset / limit + 1).toString())
		history.push({ search: `page=${offset / limit + 1}&search=${searchText}` })
	}, [offset, searchText])

	const filterOptions = props.extraFilterOptions ? FILTER_OPTIONS.concat(...props.extraFilterOptions, SORT_OPTION) : FILTER_OPTIONS.concat(SORT_OPTION)

	return (
		<div className={flex}>
			<form
				className={flexOne}
				onSubmit={handleSubmit(({ search }) => {
					onSearchChange(search)
					history.push({ search: `page=1&search=${search}` })
				})}
			>
				<Input inputRef={register} name="search" placeholder="Search..." />
				<Button>
					<FontAwesomeIcon icon={["fas", "search"]} />
				</Button>
				<StatefulPopover
					showArrow
					placement="bottomRight"
					dismissOnClickOutside={false}
					content={() => (
						<StatefulMenu
							items={filterOptions.map((option, index) => {
								return {
									label: <MenuOptionLabel label={option.label} checked={filter == option.id} borderTop={index == filterOptions.length - 1} />,
									id: option.id,
								}
							})}
							onItemSelect={({ item }) => setFilter(item.id)}
							overrides={{
								List: MenuListOverride,
								Option: {
									style: MenuOptionStyleOverride,
									props: {
										getChildMenu: (item: { label: JSX.Element; id?: FilterOption }) => {
											if (item.id == undefined) {
												return (
													<SortMenu
														sortBy={sortBy}
														sortDir={sortDir}
														setSortBy={setSortBy}
														setSortDir={setSortDir}
														extraSortByOptions={props.extraSortByOptions}
													/>
												)
											}

											return null
										},
									},
								},
							}}
						/>
					)}
					accessibilityType="menu"
					overrides={{
						Body: { style: { boxShadow: "rgba(0, 0, 0, 0.13) 0px 4px 4px" } },
					}}
				>
					<Button type="button" kind="secondary">
						<FontAwesomeIcon icon={["fas", "filter"]} />
					</Button>
				</StatefulPopover>
			</form>
			<Button
				onClick={() => props.setOffset(Math.min(Math.max(offset - limit, 0), total - 1))}
				disabled={offset <= 0}
				startEnhancer={ChevronLeft}
				overrides={PaginationButtonOverrides}
			>
				Prev
			</Button>
			<Input
				value={pageField}
				onChange={e => setPageField(e.currentTarget.value)}
				onKeyDown={e => e.keyCode === 13 && props.setOffset(Math.min(Math.max((+pageField - 1) * limit, 0), total - 1))}
				overrides={{
					Root: {
						style: {
							maxWidth: "80px",
						},
					},
					Input: {
						style: ({ $isFocused }) => ({
							backgroundColor: !$isFocused ? "white" : "rgb(246, 246, 246)",
							borderColorLeft: !$isFocused ? "unset" : "black",
							borderColorRight: !$isFocused ? "unset" : "black",
							borderColorTop: !$isFocused ? "unset" : "black",
							borderColorBottom: !$isFocused ? "unset" : "black",
							textAlign: "right",
						}),
					},
					InputContainer: {
						style: ({ $isFocused }) => ({
							borderColorLeft: !$isFocused ? "white" : "black",
							borderColorRight: !$isFocused ? "white" : "black",
							borderColorTop: !$isFocused ? "white" : "black",
							borderColorBottom: !$isFocused ? "white" : "black",
						}),
					},
				}}
			/>
			<div className={totalPageStyle}>{`of ${totalPages}`}</div>
			<Button
				onClick={() => props.setOffset(Math.min(Math.max(offset + limit, 0), total - 1))}
				disabled={offset >= total - limit}
				endEnhancer={ChevronRight}
				overrides={PaginationButtonOverrides}
			>
				Next
			</Button>
		</div>
	)
}

const PaginationButtonOverrides: ButtonOverrides = {
	BaseButton: {
		style: {
			backgroundColor: "white",
			color: "black",
			":hover": {
				backgroundColor: "#F6F6F6",
			},
		},
	},
}

interface SortMenuProps {
	sortBy: SortByOption
	sortDir: SortDir
	setSortBy: (value: SortByOption) => void
	setSortDir: (value: SortDir) => void

	extraSortByOptions?: { label: string; id: SortByOption }[]
}
const SortMenu = (props: SortMenuProps) => {
	const { sortBy, sortDir, setSortBy, setSortDir } = props
	const [css, theme] = useStyletron()
	const back = css({
		backgroundColor: "white",
	})
	const menuSubTitle = css({
		...theme.typography.font200,
		color: "#0078D4",
		fontSize: "14px",
		padding: "10px 16px 2px",
	})
	const borderTopStyle = css({
		borderTop: "solid rgba(0, 0, 0, 0.2) 1px",
	})

	const sortByOptions = props.extraSortByOptions ? SORT_BY_OPTIONS.concat(...props.extraSortByOptions) : SORT_BY_OPTIONS

	return (
		<div className={back}>
			<div className={menuSubTitle}>Sort by</div>
			<StatefulMenu
				items={sortByOptions.map(option => {
					return {
						label: <MenuOptionLabel label={option.label} checked={sortBy == option.id} />,
						id: option.id,
					}
				})}
				onItemSelect={({ item }) => setSortBy(item.id)}
				overrides={{
					List: MenuListOverride,
					Option: { style: MenuOptionStyleOverride },
				}}
			/>
			<div className={menuSubTitle + " " + borderTopStyle}>Sort order</div>
			<StatefulMenu
				items={SORT_ORDER_OPTIONS.map(option => {
					return {
						label: <MenuOptionLabel label={option.label} checked={sortDir == option.id} />,
						id: option.id,
					}
				})}
				onItemSelect={({ item }) => setSortDir(item.id)}
				overrides={{
					List: MenuListOverride,
					Option: { style: MenuOptionStyleOverride },
				}}
			/>
		</div>
	)
}
