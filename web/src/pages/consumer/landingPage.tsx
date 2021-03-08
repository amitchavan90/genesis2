import * as React from "react"
import { RouteComponentProps } from "react-router-dom"
import { useStyletron } from "baseui"
import { useQuery } from "@apollo/react-hooks"
import { graphql } from "../../graphql"
import { Product, SKUContent, SKU, Contract, Transaction } from "../../types/types"
import { LoadingSimple } from "../../components/loading"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { TransactionsView } from "../../components/transactionsView"
import { LightBlue } from "../../themeOverrides"
import Slider from "react-slick"

const SteakAndGrill = React.lazy(() => import("../../components/consumer/steakAndGrill"))
const SteakCloseCongrats = React.lazy(() => import("../../components/consumer/steakCloseCongrats"))

const ProductLandingPage = (props: RouteComponentProps) => {
	const searchArgs = new URLSearchParams(props.location.search)
	const id = searchArgs.get("productID")
	const registered = searchArgs.get("registered")
	const pointsGiven = searchArgs.get("pointsGiven")

	// Get Product info
	const [product, setProduct] = React.useState<Product>()
	const { data, loading, error } = useQuery<{ productByID: Product }>(graphql.query.PRODUCT_VIEW, { variables: { id } })

	React.useEffect(() => {
		if (!data || !data.productByID) return
		setProduct(data.productByID)
	}, [data, loading, error])

	// Styling
	const [css, theme] = useStyletron()
	const container = css({
		overflow: "hidden",

		fontFamily: "Spartan Book Classified",
		fontSize: "16px",
		color: theme.colors.primary,
	})
	const page = css({
		display: "flex",
		flexDirection: "column",
		minHeight: "100vh",
		maxWidth: "400px",
		margin: "16px auto",
	})
	const errorStyle = css({
		...theme.typography.font1250,
		height: "100vh",
		width: "100%",
		display: "flex",
		alignItems: "center",
		justifyContent: "center",
	})

	// Loading
	if (loading || (!product && data)) return <LoadingSimple />

	// Error
	if (!product) {
		return (
			<div className={page}>
				<div className={errorStyle}>Invalid UUID</div>
			</div>
		)
	}

	if (registered || product.registered)
		return (
			<div className={container}>
				<SteakCloseCongrats pointsGiven={pointsGiven} />
				<div className={page}>
					<SKUViewComponent sku={product.sku} contract={product.contract} transactions={product.transactions} />
				</div>
			</div>
		)

	return (
		<div className={container}>
			<div className={page}>
				<SteakAndGrill />

				<SKUViewComponent sku={product.sku} contract={product.contract} transactions={product.transactions} />
			</div>
		</div>
	)
}

export const SKUViewComponent = (props: { sku: SKU; contract?: Contract; transactions?: Transaction[] }) => {
	const { sku, contract, transactions } = props

	// Styling
	const [css, theme] = useStyletron()
	const containerStyle = css({
		padding: "0 26px",
		fontFamily: "Spartan Book Classified",
		fontSize: "16px",
		color: theme.colors.primary,
	})
	const titleStyle = css({
		fontFamily: "Elephant Black",
		fontSize: "27px",
		color: LightBlue,
		marginTop: "23px",
		marginBottom: "23px",
		maxWidth: "320px",
		textTransform: "uppercase",
	})
	const fieldLabel = css({
		fontFamily: "Elephant Black",
		fontSize: "18px",
		marginTop: "20px",
		marginBottom: "7px",
		textTransform: "uppercase",
	})
	const productInfo = css({
		width: "100%",
		marginTop: "10px",
	})
	const sectionStyle = css({
		width: "100%",
		textAlign: "center",
	})
	const headerImageContainer = css({
		marginTop: "18px",
		marginBottom: "10px",
		width: "100%",
		textAlign: "center",
	})
	const imageStyle = css({
		width: "100%",
		maxWidth: "800px",
	})
	const videoStyle = css({
		width: "100%",
		maxWidth: "800px",
		margin: sku && sku.urls.length == 0 ? "10px auto 0" : "0 auto",
	})
	const negativePaddingStyle = css({
		width: "calc(100% + 52px)",
		marginLeft: "-26px", // ignore page padding
	})
	const carouselImageStyle = css({
		//width: "100%",
	})
	const carouselArrowLeftStyle = css({
		left: "8px",
		position: "absolute",
		zIndex: 2,
		color: LightBlue,
	})
	const carouselArrowRightStyle = css({
		right: "8px",
		position: "absolute",
		zIndex: 2,
		color: LightBlue,
	})

	return (
		<div className={containerStyle}>
			<h2 className={titleStyle} style={{ marginBottom: "0" }}>
				About your Steak
			</h2>

			{sku.masterPlan && (
				<div className={headerImageContainer}>
					<img src={sku.masterPlan.file_url} className={imageStyle} />
				</div>
			)}

			<h4 className={fieldLabel}>About the Steak</h4>
			<div style={{ marginBottom: "10px" }}>{sku.name}</div>
			<div>{sku.description}</div>

			<h4 className={fieldLabel}>About the Farm</h4>
			{contract ? (
				<>
					<div style={{ marginBottom: "10px" }}>{contract.supplierName}</div>
					<div>{contract.description}</div>
				</>
			) : (
				<div>Livestock Specification info here.</div>
			)}

			{sku.productInfo && sku.productInfo.length > 0 && (
				<>
					<div className={fieldLabel}>Product Information</div>
					<div className={productInfo}>
						{sku.productInfo.map((info, index) => (
							<ProductInfoRow info={info} key={`product-info-row-${index}`} />
						))}
					</div>
				</>
			)}

			<div className={fieldLabel}>Tracking History</div>
			{transactions ? (
				<div className={negativePaddingStyle}>
					<TransactionsView transactions={transactions} flipNames showManifestDetails />
				</div>
			) : (
				<div>No transactions have been made.</div>
			)}

			{(sku.video || sku.urls.length > 0) && (
				<>
					<div className={fieldLabel}>More Information</div>
					{sku.urls.length > 0 && (
						<div className={sectionStyle}>
							{sku.urls.map((url, index) => (
								<ProductURL url={url} key={`product-url-${index}`} />
							))}
						</div>
					)}
					{sku.video && <video src={sku.video.file_url + "#t=0.001"} className={videoStyle} controls />}
				</>
			)}

			{sku.photos && sku.photos.length > 0 && (
				<div className={negativePaddingStyle}>
					<Slider infinite slidesToShow={1}>
						{sku.photos.map((photo, index) => (
							<img className={carouselImageStyle} src={photo.file_url} key={`product-photo-${index}`} />
						))}
					</Slider>
				</div>
			)}
		</div>
	)
}

const ProductInfoRow = (props: { info: SKUContent }) => {
	const { info } = props

	const [css, theme] = useStyletron()
	const container = css({
		display: "flex",
		width: "100%",
		justifyContent: "space-between",
		borderBottom: `1px solid grey`,
		marginBottom: "10px",
	})
	const leftSide = css({
		width: "50%",
	})
	const rightSide = css({
		textAlign: "right",
		width: "50%",
	})

	return (
		<div className={container}>
			<div className={leftSide}>{info.title}</div>
			<div className={rightSide}>{info.content}</div>
		</div>
	)
}
const ProductURL = (props: { url: SKUContent }) => {
	const { url } = props

	const [css, theme] = useStyletron()
	const container = css({
		fontWeight: "normal",
		display: "flex",
		width: "100%",
		marginBottom: "5px",
		textDecoration: "unset",
		alignItems: "center",
		padding: "5px",
		color: LightBlue,
		":hover": {
			backgroundColor: LightBlue,
			color: "white",
		},
	})
	const iconStyle = css({
		marginRight: "8px",
	})

	return (
		<a className={container} href={url.content.indexOf("://") == -1 ? `http://${url.content}` : url.content}>
			<FontAwesomeIcon icon={["far", "globe"]} className={iconStyle} />
			{url.title}
		</a>
	)
}

export default ProductLandingPage
