import * as React from "react"
import { Transaction } from "../types/types"
import { useStyletron } from "baseui"
import { timeAgo } from "../helpers/time"
import { Modal, ModalBody, ModalHeader } from "baseui/modal"
import MapPreview from "./mapPreview"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { useQuery } from "@apollo/react-hooks"
import { graphql } from "../graphql"
import { Spinner } from "baseui/icon"
import { sha256 } from "js-sha256"
import { VerticalTimeline, VerticalTimelineElement } from "react-vertical-timeline-component"

interface TransactionsViewProps {
	transactions?: Transaction[]
	flipNames?: boolean
	showManifestDetails: boolean
}

const fallbackImage = require("../assets/images/default-image.png")

export const TransactionsView = (props: TransactionsViewProps) => {
	const { transactions, flipNames, showManifestDetails } = props

	// Styling
	const [css] = useStyletron()
	const dateStyle = css({
		paddingLeft: "20px !important",
	})
	const paddingStyle = css({
		padding: "0 25px",
	})

	if (!transactions) return <div className={paddingStyle}>Failed to load tracking history.</div>

	if (transactions.length == 0) return <div className={paddingStyle}>No transactions have been made.</div>

	const sortedTransactions = transactions.sort((a, b) => new Date(b.scannedAt || b.createdAt).getTime() - new Date(a.scannedAt || a.createdAt).getTime())

	return (
		<VerticalTimeline layout="1-column">
			{sortedTransactions.map((t) => (
				<VerticalTimelineElement
				key={`timeline-${t.id}`}
					className="vertical-timeline-element--work"
					contentStyle={{
						background: "#58b7e8",
						color: "#fff",
						opacity: t.action.private ? 0.7 : "unset",
						padding: 0,
					}}
					contentArrowStyle={{ borderRight: "7px solid  #58b7e8" }}
					iconStyle={{ background: "#58b7e8", color: "#fff" }}
					dateClassName={dateStyle}
					date={`${timeAgo(t.scannedAt || t.createdAt)} - ${new Date(t.scannedAt || t.createdAt).toLocaleString()}`}
				>
					<TransactionItem key={`Transaction${t.id}`} transaction={t} flipNames={flipNames} showManifestDetails={showManifestDetails} />
				</VerticalTimelineElement>
			))}
		</VerticalTimeline>
	)
}

const TransactionItem = (props: { transaction: Transaction; flipNames?: boolean; showManifestDetails: boolean }) => {
	const { transaction, flipNames, showManifestDetails } = props
	const nameChinese = transaction.action.nameChinese || transaction.action.name

	const [showDetailsModal, setShowDetailsModal] = React.useState<boolean>()
	const [showPreviewModal, setShowPreviewModal] = React.useState<boolean>()
	const [previewModalImage, setPreviewModalImage] = React.useState<string>()

	const [css, theme] = useStyletron()
	const recordContainerStyle = css({
		display: "flex",
		flexDirection: "column",
		cursor: "pointer",
	})
	const flexColumnStyle = css({
		display: "flex",
		flexDirection: "column",
		paddingTop: "20px",
		paddingLeft: "20px",
		width: "65%",
	})
	const flexStyle = css({
		display: "flex",
		justifyContent: "space-between",
	})
	const locationStyle = css({
		paddingLeft: "20px",
		paddingRight: "20px",
		fontSize: "0.8rem",
		lineHeight: "0.8rem",
		marginTop: "8px",
	})
	const nameContainerStyle = css({
		display: "flex",
		flexDirection: "column",
		marginBottom: "8px",
	})
	const nameStyle = css({
		fontSize: "1.15rem",
		lineHeight: "1rem",
		fontWeight: "bold",
	})
	const nameSubStyle = css({
		fontSize: flipNames ? "0.8rem" : "0.9rem",
		fontWeight: flipNames ? "unset" : "bold",
		lineHeight: "1rem",
		color: "#FFFFFF8f",
	})
	const createdByNameStyle = css({
		fontSize: "1rem",
		lineHeight: "1rem",
	})
	const imageBox = css({
		width: "100%",
		display: "block",
		objectFit: "contain",
		cursor: "pointer",
		marginRight: "15px",
	})
	const photos = css({
		display: "flex",
		maxWidth: "500px",
		marginTop: "8px",
	})
	const photoDescriptionStyle = css({
		fontSize: "0.8rem",
		lineHeight: "0.8rem",
		color: "#0000008f",
		marginTop: "5px",
		marginBottom: "5px",
	})
	const imagePreviewStyle = css({
		width: "100%",
		cursor: "pointer",
	})
	const mapContainerStyle = css({
		width: "200px",
		marginLeft: "10px",
	})
	const showDetailsButtonStyle = css({
		padding: "10px",
		position: "absolute",
		bottom: "0px",
		right: "2px",
	})
	const detailsMapContainerStyle = css({
		width: "100%",
	})
	const detailNameSubStyle = css({
		fontSize: flipNames ? "0.8rem" : "0.9rem",
		fontWeight: flipNames ? "unset" : "bold",
		lineHeight: "1rem",
		color: "#0000008f",
		marginLeft: "10px",
	})
	const detailDateStyle = css({
		color: "#0000008f",
		padding: "10px 0",
	})

	return (
		<>
			<div className={recordContainerStyle} onClick={() => setShowDetailsModal(true)}>
				<div className={flexStyle}>
					<div className={flexColumnStyle}>
						{/* Action */}
						<div className={nameContainerStyle}>
							<div className={nameStyle}>{flipNames ? nameChinese : transaction.action.name}</div>
							{transaction.action.nameChinese && <div className={nameSubStyle}>{flipNames ? transaction.action.name : nameChinese}</div>}
						</div>

						{/* Created by */}
						<div className={nameContainerStyle}>
							<div className={createdByNameStyle}>{transaction.createdByName || "Unknown Organisation"}</div>
							{transaction.createdBy && <div className={createdByNameStyle}>{transaction.createdBy.email || transaction.createdBy.wechatID}</div>}
						</div>
					</div>

					{/* Map */}
					{transaction.locationGeohash && (
						<div className={mapContainerStyle}>
							<MapPreview locationName={transaction.locationName} locationGeohash={transaction.locationGeohash} roundedCorner largeMarker zoom={4} />
						</div>
					)}
				</div>

				{/* Location */}
				<div className={locationStyle}>{transaction.locationName}</div>

				{/* Details Buttons */}
				<div className={showDetailsButtonStyle}>
					<FontAwesomeIcon icon={["fas", "chevron-down"]} />
				</div>

				{/* Photos of carton and product */}
				{transaction.photos && (
					<div className={photos}>
						{transaction.photos.cartonPhoto && (
							<div style={{ marginRight: "10px" }}>
								<div className={photoDescriptionStyle}>photo of carton:</div>
								<img
									className={imageBox}
									src={transaction.photos.cartonPhoto.file_url || fallbackImage}
									alt="carton image"
									onClick={() => {
										if (!transaction.photos.cartonPhoto?.file_url) return
										setPreviewModalImage(transaction.photos.cartonPhoto?.file_url)
										setShowPreviewModal(true)
									}}
								/>
							</div>
						)}
						{transaction.photos.productPhoto && (
							<div>
								<div className={photoDescriptionStyle}>photo of product:</div>
								<img
									className={imageBox}
									src={transaction.photos.productPhoto.file_url || fallbackImage}
									alt="product image"
									onClick={() => {
										if (!transaction.photos.productPhoto?.file_url) return
										setPreviewModalImage(transaction.photos.productPhoto?.file_url)
										setShowPreviewModal(true)
									}}
								/>
							</div>
						)}
					</div>
				)}
			</div>

			{/* Details Modal */}
			<Modal
				isOpen={showDetailsModal}
				onClose={() => setShowDetailsModal(false)}
				overrides={{
					Root: { style: { zIndex: 100 } },
				}}
			>
				<ModalHeader>
					{/* Action */}
					<span>
						<span className={nameStyle}>{flipNames ? nameChinese : transaction.action.name}</span>
						{transaction.action.nameChinese && <span className={detailNameSubStyle}>({flipNames ? transaction.action.name : nameChinese})</span>}
					</span>
				</ModalHeader>
				<ModalBody>
					{/* Created by */}
					<div className={nameContainerStyle}>
						<div className={createdByNameStyle}>{transaction.createdByName || "Unknown Organisation"}</div>
						{transaction.createdBy && <div className={createdByNameStyle}>{transaction.createdBy.email || transaction.createdBy.wechatID}</div>}
					</div>

					{/* Date */}
					<div className={detailDateStyle}>{`${timeAgo(transaction.scannedAt || transaction.createdAt)} - ${new Date(
						transaction.scannedAt || transaction.createdAt,
					).toLocaleString()}`}</div>

					{/* Location */}
					<div>{transaction.locationName}</div>

					{/* Map */}
					{transaction.locationGeohash && (
						<div className={detailsMapContainerStyle}>
							<MapPreview locationName={transaction.locationName} locationGeohash={transaction.locationGeohash} />
						</div>
					)}

					{/* Manifest Details */}
					{showManifestDetails && <TransactionManifestDetails transaction={transaction} />}
				</ModalBody>
			</Modal>

			{/* Image Modal */}
			<Modal
				isOpen={showPreviewModal}
				onClose={() => setShowPreviewModal(false)}
				overrides={{
					Root: { style: { zIndex: 100 } },
					Dialog: {
						style: {
							width: "unset",
							backgroundColor: "unset",
						},
					},
				}}
			>
				<img className={imagePreviewStyle} src={previewModalImage} />
			</Modal>
		</>
	)
}

const TransactionManifestDetails = (props: { transaction: Transaction }) => {
	const { transaction } = props

	const hexToBytes = (hex: string) => {
		for (var bytes = [], c = 0; c < hex.length; c += 2) bytes.push(parseInt(hex.substr(c, 2), 16))
		return bytes
	}

	const settingsQuery = useQuery<{ settings: { etherscanHost: string; consumerHost: string } }>(graphql.query.SETTINGS)

	const [verifying, setVerifying] = React.useState(false)
	const [lineVerified, setLineVerified] = React.useState<boolean>()
	const [manifestVerified, setManifestVerified] = React.useState<boolean>()

	const verify = () => {
		if (
			!transaction.manifestLineJson ||
			!transaction.manifestLineSha256 ||
			!transaction.manifest?.transactionHash ||
			!transaction.manifest?.merkleRootSha256 ||
			!settingsQuery.data?.settings?.consumerHost
		)
			return

		setVerifying(true)

		// verify line
		setLineVerified(sha256(transaction.manifestLineJson) === transaction.manifestLineSha256)

		// verify manifest
		fetch(`${settingsQuery.data.settings.consumerHost}/api/manifest/tx/${transaction.manifest?.transactionHash}`).then((r) => {
			r.text().then((text) => {
				const lines = text.split("\n")
				let data: number[] = []
				lines.forEach((l) => data.push(...hexToBytes(l.substr(0, 64))))

				if (!transaction.manifest?.merkleRootSha256) return

				setManifestVerified(sha256(data) === transaction.manifest.merkleRootSha256)
			})
		})
	}

	// Styling
	const [css, theme] = useStyletron()
	const containerStyle = css({
		fontFamily: "monospace",
		padding: "6px",
		backgroundColor: "#f0f0f0",
		boxShadow: theme.lighting.shadow400,
		marginTop: "6px",
	})
	const labelStyle = css({
		fontWeight: "bold",
		whiteSpace: "nowrap",
		marginRight: "6px",
	})
	const groupStyle = css({
		marginBottom: "10px",
	})
	const jsonStyle = css({
		whiteSpace: "nowrap",
		padding: "2px",
		overflowX: "auto",
		marginRight: "5px",
	})
	const flexStyle = css({
		display: "flex",
	})
	const hashStyle = css({
		wordBreak: "break-word",
		padding: "2px",
		backgroundColor: "#f0f0f0",
		textOverflow: "ellipsis",
	})
	const buttonStyle = css({
		padding: "6px",
		alignSelf: "center",
		textAlign: "center",
		boxShadow: theme.lighting.shadow400,
		backgroundColor: "#58b7e8",
		color: "white",
		":hover": {
			backgroundColor: "#85d5ff",
		},
		":visited": {
			color: "white",
		},
		textDecoration: "none",
		cursor: "pointer",
		whiteSpace: "nowrap",
		marginLeft: "auto",
		minWidth: "78px",
	})
	const blockchainLinksStyle = css({
		display: "flex",
		flexDirection: "column",
		marginLeft: "auto",
	})
	const blockLinkStyle = css({
		width: "calc(100% - 12px)",
		textAlign: "left",
		marginBottom: "6px",
	})
	const resultStyle = css({
		padding: "6px",
		textAlign: "center",
		borderRadius: "3px",
		color: "white",
		backgroundColor: lineVerified ? theme.colors.borderPositive : theme.colors.borderNegative,
		marginTop: manifestVerified ? "unset" : "6px",
	})
	const result2Style = css({
		padding: "6px",
		textAlign: "center",
		borderRadius: "3px",
		color: "white",
		backgroundColor: manifestVerified ? theme.colors.borderPositive : theme.colors.borderNegative,
		marginTop: "6px",
	})

	return (
		<div className={containerStyle}>
			<div className={groupStyle}>
				<div className={labelStyle}>data json:</div>
				<div className={flexStyle}>
					<div className={jsonStyle}>{transaction.manifestLineJson}</div>
					<a
						className={buttonStyle}
						href={
							transaction.manifestLineJson && transaction.manifestLineSha256
								? `data:text/json;charset=utf-8,${encodeURIComponent(transaction.manifestLineJson)}`
								: undefined
						}
						download={`${transaction.manifestLineSha256}.json`}
					>
						download
					</a>
				</div>
			</div>

			<div className={groupStyle}>
				<div className={labelStyle}>data hash:</div>
				<div className={flexStyle}>
					<div className={hashStyle}>{transaction.manifestLineSha256}</div>
					<div className={buttonStyle} onClick={() => transaction.manifestLineSha256 && navigator.clipboard.writeText(transaction.manifestLineSha256)}>
						copy
					</div>
				</div>
			</div>

			<div className={groupStyle}>
				<div className={labelStyle}>manifest merkle root:</div>
				<div className={flexStyle}>
					<div className={hashStyle}>{transaction.manifest?.merkleRootSha256}</div>
					<div
						className={buttonStyle}
						onClick={() => transaction.manifest?.merkleRootSha256 && navigator.clipboard.writeText(transaction.manifest.merkleRootSha256)}
					>
						copy
					</div>
				</div>
			</div>

			<div className={groupStyle}>
				<div className={labelStyle}>blockchain:</div>
				<div className={flexStyle}>
					<div className={hashStyle}>{transaction.manifest?.transactionHash}</div>
					<div className={blockchainLinksStyle}>
						<a
							className={buttonStyle + " " + blockLinkStyle}
							href={
								settingsQuery.data?.settings?.etherscanHost && transaction.manifest
									? `${settingsQuery.data.settings.etherscanHost}/tx/${transaction.manifest?.transactionHash}`
									: undefined
							}
							target="_blank"
						>
							<FontAwesomeIcon icon={["fas", "external-link"]} style={{ marginRight: "5px" }} />
							etherscan
						</a>
						<a
							className={buttonStyle + " " + blockLinkStyle}
							href={
								settingsQuery.data?.settings?.consumerHost && transaction.manifest
									? `${settingsQuery.data.settings.consumerHost}/api/manifest/tx/${transaction.manifest?.transactionHash}`
									: undefined
							}
							target="_blank"
						>
							<FontAwesomeIcon icon={["fas", "external-link"]} style={{ marginRight: "5px" }} />
							manifest
						</a>
					</div>
				</div>
			</div>

			{manifestVerified === undefined && (
				<div className={buttonStyle} onClick={verifying ? undefined : verify}>
					{verifying ? <Spinner /> : "verify"}
				</div>
			)}
			{lineVerified !== undefined && <div className={resultStyle}>{lineVerified === true ? "LINE: PASS" : "LINE: FAIL"}</div>}
			{manifestVerified !== undefined && <div className={result2Style}>{manifestVerified === true ? "MANIFEST: PASS" : "MANIFEST: FAIL"}</div>}
		</div>
	)
}
