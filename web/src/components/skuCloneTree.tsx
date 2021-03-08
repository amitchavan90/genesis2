import * as React from "react"
import { TreeView, toggleIsExpanded, TreeNode, TreeLabelProps, TreeLabel } from "baseui/tree-view"
import { SKUClone } from "../types/types"
import { graphql } from "../graphql"
import { useQuery } from "@apollo/react-hooks"
import { LoadingSimple } from "./loading"
import { useStyletron } from "baseui"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { useHistory } from "react-router-dom"

export const SKUCloneTree = (props: { skuID: String }) => {
	const history = useHistory()

	const [tree, setTree] = React.useState<TreeNode[]>([])
	const { data, loading, error } = useQuery<{ skuCloneTree: SKUClone[] }>(graphql.query.SKU_CLONE_TREE, {
		variables: { id: props.skuID },
	})

	React.useEffect(() => {
		if (!data || !data.skuCloneTree || data.skuCloneTree.length == 0) return
		setTree(getTreeNodes(data.skuCloneTree))
	}, [data, loading, error])

	const getTreeNodes = (cloneTree: SKUClone[]) => {
		const nodes: TreeNode[] = []
		const added = cloneTree.map(c => false)
		for (let i = 0; i < cloneTree.length; i++) {
			if (added[i]) continue
			nodes.push({
				id: cloneTree[i].sku.id,
				label: cloneTree[i].sku.code,
				info: cloneTree[i].sku,
				isExpanded: true,
				children: getTreeNodeChildren(cloneTree, added, cloneTree[i].depth + 1, i + 1),
			})
			added[i] = true
		}

		return nodes
	}
	const getTreeNodeChildren = (cloneTree: SKUClone[], added: boolean[], depth: number, offset?: number) => {
		const children: TreeNode[] = []
		for (let i = offset || 0; i < cloneTree.length; i++) {
			if (cloneTree[i].depth < depth) return children
			if (added[i] || cloneTree[i].depth != depth) continue
			children.push({
				id: cloneTree[i].sku.id,
				label: cloneTree[i].sku.code,
				info: cloneTree[i].sku,
				isExpanded: true,
				children: getTreeNodeChildren(cloneTree, added, cloneTree[i].depth + 1, i + 1),
			})
			added[i] = true
		}
		return children
	}

	const [css] = useStyletron()
	const containerStyle = css({
		paddingTop: "10px",
	})
	const labelStyle = css({
		display: "flex",
		cursor: "pointer",
	})
	const infoStyle = css({
		color: "grey",
		marginLeft: "20px",
	})
	const cardStyle = css({
		display: "flex",
		flexDirection: "column",
		marginLeft: "6px",
	})
	const nameStyle = css({
		fontSize: "1.15rem",
		lineHeight: "1rem",
	})
	const boldStyle = css({
		fontWeight: "bold",
	})
	const nameSubStyle = css({
		fontSize: "0.8rem",
		lineHeight: "1rem",
		color: "grey",
		marginTop: "3px",
	})

	if (loading) return <LoadingSimple />
	if (tree.length == 0) return <></>

	const CustomLabel = (node: TreeNode) => {
		const currentSKU = node.info && node.info.id == props.skuID
		return (
			<div className={labelStyle} onClick={() => history.push(`/portal/sku/${node.label}#cloneTree`)}>
				<FontAwesomeIcon icon={["fal", "barcode-alt"]} size="2x" />
				<div className={cardStyle + (currentSKU ? ` ${boldStyle}` : "")}>
					<div className={nameStyle}>{node.info.name}</div>
					<div className={nameSubStyle}>{node.label}</div>
				</div>
				{currentSKU && <div className={infoStyle}>{"<- current"}</div>}
			</div>
		)
	}
	const CustomTreeLabel = (props: TreeLabelProps) => <TreeLabel {...props} label={CustomLabel} />

	return (
		<div className={containerStyle}>
			<TreeView
				data={tree}
				onToggle={node => setTree(prevData => toggleIsExpanded(prevData, node))}
				overrides={{
					TreeLabel: {
						component: CustomTreeLabel,
					},
				}}
			/>
		</div>
	)
}
