import * as React from "react"
import { RouteComponentProps, useHistory } from "react-router-dom"
import { useStyletron } from "baseui"
import { useQuery, useMutation } from "@apollo/react-hooks"
import { graphql } from "../../graphql"
import { SKU, Blob, SKUContent, Settings } from "../../types/types"
import { useForm } from "react-hook-form"
import { FormControl, FormControlOverrides } from "baseui/form-control"
import { Input } from "baseui/input"
import { Textarea } from "baseui/textarea"
import { H1 } from "baseui/typography"
import { Button } from "baseui/button"
import { Notification } from "baseui/notification"
import { Spaced } from "../../components/spaced"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { LoadingSimple } from "../../components/loading"
import { ErrorNotification } from "../../components/errorBox"
import { Spread } from "../../components/spread"
import { ImageUpload, ImageUploadHandler, ImageUploadMultipleHandler } from "../../components/imageUpload"
import { VideoUpload, VideoUploadHandler } from "../../components/videoUpload"
import { IconName } from "@fortawesome/fontawesome-svg-core"
import { Tabs, Tab } from "baseui/tabs"
import { paddingZero, ButtonMarginLeftOverride } from "../../themeOverrides"
import { CenteredPage } from "../../components/common"
import { ItemList } from "../../components/itemList"
import { FilterOption } from "../../types/enums"
import { ActionItemSet } from "../../types/actions"
import { Modal, ModalButton, ModalFooter, ModalHeader } from "baseui/modal"
import { SKUViewComponent } from "../consumer/landingPage"
import { invalidateListQueries } from "../../apollo"
import { SKUCloneTree } from "../../components/skuCloneTree"
import { Checkbox } from "baseui/checkbox"
import { promiseTimeout, TIMED_OUT } from "../../helpers/timeout"

type FormData = {
	name: string
	code: string
	urls: SKUContent[]
	productInfo: SKUContent[]
	loyaltyPoints: number
}

const SKUEdit = (props: RouteComponentProps<{ code: string }>) => {
	const code = props.match.params.code
	const isNewSKU = code === "new"

	const sQuery = new URLSearchParams(props.location.search)
	const cloneFrom = isNewSKU ? sQuery.get("clone") : null

	const history = useHistory()

	const [activeKey, setActiveKey] = React.useState(props.location.hash || "#details")
	const [showPreviewModal, setShowPreviewModal] = React.useState<boolean>()

	// Get SKU
	const [sku, setSKU] = React.useState<SKU>()
	const { data, loading, error } = useQuery<{ sku: SKU }>(graphql.query.SKU, {
		variables: { code: isNewSKU ? cloneFrom : code },
		fetchPolicy: isNewSKU && !cloneFrom ? "cache-only" : undefined, // prevent query if new
	})
	const settingsQuery = useQuery<{ settings: Settings }>(graphql.query.SETTINGS)

	// Mutations
	const [updateSKU, mutUpdateSKU] = useMutation(isNewSKU ? graphql.mutation.CREATE_SKU : graphql.mutation.UPDATE_SKU)
	const [archiveSKU, mutArchiveSKU] = useMutation<{ skuArchive: SKU }>(graphql.mutation.ARCHIVE_SKU)
	const [unarchiveSKU, mutUnarchiveSKU] = useMutation<{ skuUnarchive: SKU }>(graphql.mutation.UNARCHIVE_SKU)

	const toggleArchive = () => {
		if (!sku) return

		if (sku.archived) {
			unarchiveSKU({
				variables: { id: sku.id },
				update: (cache: any) => invalidateListQueries(cache, "skus"),
			})
		} else {
			archiveSKU({
				variables: { id: sku.id },
				update: (cache: any) => invalidateListQueries(cache, "skus"),
			})
		}
	}

	// Form submission
	const [timedOut, setTimedOut] = React.useState(false)
	const [showSuccessModal, setShowSuccessModal] = React.useState(false)
	const [description, setDescription] = React.useState("")
	const [isBeef, setIsBeef] = React.useState<boolean>()

	const [masterPlanFile, setMasterPlanFile] = React.useState<File>()
	const [masterPlanURL, setMasterPlanURL] = React.useState<string>()
	const [uploadMasterPlan, setUploadMasterPlan] = React.useState<ImageUploadHandler>()

	const [videoFile, setVideoFile] = React.useState<File>()
	const [videoURL, setVideoURL] = React.useState<string>()
	const [uploadVideo, setUploadVideo] = React.useState<VideoUploadHandler>()

	const [photos, setPhotos] = React.useState<Blob[] | undefined>()
	const [uploadPhotos, setUploadPhotos] = React.useState<ImageUploadMultipleHandler | undefined>()

	const [urlCount, setURLCount] = React.useState(isNewSKU && !cloneFrom ? 0 : 4) // count must start a max for the initial load to work
	const [productInfoCount, setProductInfoCount] = React.useState(isNewSKU && !cloneFrom ? 0 : 25)

	const { register, setValue, handleSubmit, errors, getValues } = useForm<FormData>()

	const onSubmit = handleSubmit(async ({ name, code, urls, productInfo, loyaltyPoints }) => {
		setTimedOut(false)

		// Upload Master Plan
		let masterPlanBlobID: string | null = null
		if (uploadMasterPlan) {
			if (uploadMasterPlan.upload && uploadMasterPlan.setUploadError) {
				const response = await uploadMasterPlan.upload()

				if (!response || !response.data) {
					uploadMasterPlan.setUploadError("Upload Failed")
					return
				}

				if (response.data.fileUpload) {
					setMasterPlanURL(response.data.fileUpload.file_url)
					masterPlanBlobID = response.data.fileUpload.id
				}
			} else if (uploadMasterPlan.removeImage) {
				masterPlanBlobID = "-"
			}
		} else if (cloneFrom && sku) {
			masterPlanBlobID = sku.masterPlan.id
		}

		// Upload Video
		let videoBlobID: string | null = null
		if (uploadVideo) {
			if (uploadVideo.upload && uploadVideo.setUploadError) {
				const response = await uploadVideo.upload()

				if (!response || !response.data) {
					uploadVideo.setUploadError("Upload Failed")
					return
				}

				if (response.data.fileUpload) {
					setVideoURL(response.data.fileUpload.file_url)
					videoBlobID = response.data.fileUpload.id
				}
			} else if (uploadVideo.removeFile) {
				videoBlobID = "-"
			}
		} else if (cloneFrom && sku) {
			videoBlobID = sku.video.id
		}

		// Upload Photos
		let photoBlobIDs: string[] | null = null
		if (uploadPhotos) {
			if (uploadPhotos.upload && uploadPhotos.setUploadError) {
				const response = await uploadPhotos.upload()

				if (!response || !response.data) {
					uploadPhotos.setUploadError("Upload Failed")
					return
				}

				const blobs = [...uploadPhotos.currentFiles, ...response.data.fileUploadMultiple]
				setPhotos(blobs)
				photoBlobIDs = blobs.map(p => p.id)
				if (photoBlobIDs.length == 0) photoBlobIDs = [""] // no photos
			}
		} else if (cloneFrom && sku) {
			photoBlobIDs = sku.photos.map(p => p.id)
		}

		const input = {
			name,
			code,
			description,
			isBeef,
			masterPlanBlobID,
			videoBlobID,
			photoBlobIDs,
			urls,
			productInfo,
			loyaltyPoints: loyaltyPoints || 0,
			cloneParentID: isNewSKU && cloneFrom ? sku?.id : undefined,
		}

		if (isNewSKU) {
			updateSKU({
				variables: { input },
				update: (cache: any) => invalidateListQueries(cache, "skus"),
			})
		} else if (sku) {
			promiseTimeout(updateSKU({ variables: { id: sku.id, input } })).catch(reason => {
				if (reason !== TIMED_OUT) return
				setTimedOut(true)
			})
		}
	})

	// On load sku
	React.useEffect(() => {
		if (!data || !data.sku) return
		setSKU(data.sku)
	}, [data, loading, error])
	React.useEffect(() => {
		if (activeKey != "#details") return
		if (!sku) return

		setValue("name", sku.name)
		setValue("code", sku.code)
		setValue("loyaltyPoints", sku.loyaltyPoints)
		setDescription(sku.description)
		setIsBeef(sku.isBeef)
		if (sku.masterPlan) setMasterPlanURL(sku.masterPlan.file_url)
		if (sku.video) setVideoURL(sku.video.file_url)
		setPhotos(sku.photos)

		setURLCount(sku.urls.length)
		sku.urls.forEach((url, index) => {
			setValue(`urls[${index}].title`, url.title)
			setValue(`urls[${index}].content`, url.content)
		})

		setProductInfoCount(sku.productInfo.length)
		sku.productInfo.forEach((info, index) => {
			setValue(`productInfo[${index}].title`, info.title)
			setValue(`productInfo[${index}].content`, info.content)
		})
	}, [sku, activeKey])

	// On mutation (update/create sku)
	React.useEffect(() => {
		if (!mutUpdateSKU.data) return

		if (isNewSKU) {
			if (!mutUpdateSKU.data.skuCreate) return
			setShowSuccessModal(true)
			history.replace(`/portal/sku/${mutUpdateSKU.data.skuCreate.code}`)
			return
		}

		if (!mutUpdateSKU.data.skuUpdate) return

		setSKU(mutUpdateSKU.data.skuUpdate)
		setShowSuccessModal(true)
		setTimedOut(false)
	}, [mutUpdateSKU])

	React.useEffect(() => {
		if (!mutArchiveSKU.data?.skuArchive) return
		setSKU(mutArchiveSKU.data.skuArchive)
	}, [mutArchiveSKU])
	React.useEffect(() => {
		if (!mutUnarchiveSKU.data?.skuUnarchive) return
		setSKU(mutUnarchiveSKU.data.skuUnarchive)
	}, [mutUnarchiveSKU])

	// Preview SKU
	const getPreviewSKU = () => {
		const formValues = getValues({ nest: true })

		const newPhotos = uploadPhotos?.newFiles
			? uploadPhotos.newFiles.map((f: File) => {
					return {
						file_url: URL.createObjectURL(f),
					} as Blob
			  })
			: []

		return {
			code: formValues.code,
			name: formValues.name,
			description: description,
			productInfo: formValues.productInfo || [],
			urls: formValues.urls || [],
			masterPlan: masterPlanFile ? { file_url: URL.createObjectURL(masterPlanFile) } : sku?.masterPlan,
			video: videoFile ? { file_url: URL.createObjectURL(videoFile) } : sku?.video,
			photos: sku && sku.photos.length > 0 ? sku.photos.concat(newPhotos) : newPhotos,
		} as SKU
	}

	// Styling
	const [css, theme] = useStyletron()
	const marginLeftStyle = css({
		marginLeft: "10px",
	})
	const breakLineStyle = css({
		border: "1px solid rgba(0, 0, 0, 0.1)",
		marginBottom: "10px",
	})
	const plusCircleMargin = css({
		marginRight: "10px",
	})

	const plusCircle = (
		<div className={plusCircleMargin}>
			<FontAwesomeIcon icon={["fal", "plus-circle"]} size="lg" />
		</div>
	)

	const breakLine = <div className={breakLineStyle} />

	const editForm = (
		<form onSubmit={onSubmit}>
			{isNewSKU && cloneFrom && <div>Cloned from {cloneFrom}</div>}

			{mutUpdateSKU.error && <ErrorNotification message={mutUpdateSKU.error.message} />}

			{!isNewSKU && (
				<FormControl label="Code" error={errors.code ? errors.code.message : ""} positive="" disabled>
					<Input name="code" inputRef={register} disabled />
				</FormControl>
			)}

			<FormControl label="Name" error={errors.name ? errors.name.message : ""} positive="">
				<Input name="name" inputRef={register({ required: "Required" })} />
			</FormControl>

			<FormControl label="Description" error="" positive="">
				<Textarea
					name="description"
					value={description}
					onChange={e => setDescription(e.currentTarget.value)}
					overrides={{
						Input: {
							style: {
								resize: "vertical",
								height: "170px",
							},
						},
					}}
				/>
			</FormControl>

			<FormControl caption="(will show basic scan page on product registration if not beef)">
				<Checkbox checked={isBeef} onChange={e => setIsBeef(e.currentTarget.checked)}>
					Is Beef
				</Checkbox>
			</FormControl>

			{breakLine}

			<ImageUpload.Single
				// client requested to rename it to Hero Image
				label="Hero Image"
				name="masterPlan"
				imageURL={!masterPlanURL || uploadMasterPlan?.removeImage ? "" : masterPlanURL}
				setImageUploader={imageUploader => setUploadMasterPlan(imageUploader)}
				imageRemoved={uploadMasterPlan?.removeImage}
				file={masterPlanFile}
				setFile={(file?: File) => setMasterPlanFile(file)}
				previewHeight="200px"
				caption="Please select a jpg/png file smaller than 10MB"
				maxFileSize={1e7}
				clearable
			/>

			{breakLine}

			<VideoUpload
				label="Video"
				name="video"
				videoURL={!videoURL || uploadVideo?.removeFile ? "" : videoURL}
				setVideoUploader={videoUploader => setUploadVideo(videoUploader)}
				imageRemoved={uploadVideo?.removeFile}
				file={videoFile}
				setFile={(file?: File) => setVideoFile(file)}
				caption="Please select a mp4 file smaller than 30MB"
				maxFileSize={3e7}
			/>

			{breakLine}

			<ImageUpload.Multiple
				label="Additional images"
				name="photos"
				files={photos}
				setImageUploader={imageUploader => setUploadPhotos(imageUploader)}
				maxFiles={9}
				maxWidth="850px"
			/>

			{breakLine}

			<FormControl label={`URLs (${urlCount}/4)`} error="" positive="">
				<div>
					{Array.from({ length: urlCount }).map((_, index) => (
						<InputPair
							key={`sku_url_${index}`}
							prefix="urls"
							index={index}
							icon="link"
							label2="URL"
							titleInputRef={register({ required: "Required" })}
							contentInputRef={register({ required: "Required" })}
							titleError={errors.urls && errors.urls[index] && errors.urls[index].title?.message}
							contentError={errors.urls && errors.urls[index] && errors.urls[index].content?.message}
							onDelete={async () => {
								const urls = getValues({ nest: true }).urls
								const newURLs = [...urls.slice(0, index), ...urls.slice(index + 1)]
								urls.forEach((_, index) => {
									setValue(`urls[${index}].title`, index >= newURLs.length ? "" : newURLs[index].title)
									setValue(`urls[${index}].content`, index >= newURLs.length ? "" : newURLs[index].content)
								})
								setURLCount(urlCount - 1)
							}}
						/>
					))}
					{urlCount < 4 && (
						<Button type="button" kind="secondary" onClick={() => setURLCount(urlCount + 1)}>
							{plusCircle} Add URL
						</Button>
					)}
				</div>
			</FormControl>

			{breakLine}

			<FormControl label={`Product Information (${productInfoCount}/25)`} error="" positive="">
				<div>
					{Array.from({ length: productInfoCount }).map((_, index) => (
						<InputPair
							key={`sku_info_${index}`}
							prefix="productInfo"
							index={index}
							icon="info-circle"
							titleInputRef={register({ required: "Required" })}
							contentInputRef={register({ required: "Required" })}
							titleError={errors.productInfo && errors.productInfo[index] && errors.productInfo[index].title?.message}
							contentError={errors.productInfo && errors.productInfo[index] && errors.productInfo[index].content?.message}
							onDelete={async () => {
								const productInfo = getValues({ nest: true }).productInfo
								const newInfo = [...productInfo.slice(0, index), ...productInfo.slice(index + 1)]
								productInfo.forEach((_, index) => {
									setValue(`productInfo[${index}].title`, index >= newInfo.length ? "" : newInfo[index].title)
									setValue(`productInfo[${index}].content`, index >= newInfo.length ? "" : newInfo[index].content)
								})
								setProductInfoCount(productInfoCount - 1)
							}}
						/>
					))}
					{productInfoCount < 25 && (
						<Button type="button" kind="secondary" onClick={() => setProductInfoCount(productInfoCount + 1)}>
							{plusCircle} Add Item
						</Button>
					)}
				</div>
			</FormControl>

			{breakLine}

			<FormControl label="Loyality Points" error={errors.loyaltyPoints ? errors.loyaltyPoints.message : ""} positive="">
				<Input name="loyaltyPoints" type="number" inputRef={register} />
			</FormControl>

			{breakLine}

			<Spread>
				<Button type="button" kind="secondary" onClick={() => history.push("/portal/skus")}>
					Cancel
				</Button>

				<Spaced>
					{sku !== undefined && !isNewSKU && (
						<Button
							type="button"
							kind="secondary"
							isLoading={mutArchiveSKU.loading || mutUnarchiveSKU.loading}
							onClick={toggleArchive}
							startEnhancer={<FontAwesomeIcon icon={["fas", sku.archived ? "undo" : "archive"]} size="lg" />}
						>
							{sku.archived ? "Unarchive" : "Archive"}
						</Button>
					)}

					<Button isLoading={mutUpdateSKU.loading && !timedOut} startEnhancer={<FontAwesomeIcon icon={["fas", "save"]} size="lg" />}>
						{isNewSKU ? "Create SKU" : timedOut ? "Timed out... Try again" : "Save"}
					</Button>
				</Spaced>
			</Spread>
		</form>
	)

	if (!isNewSKU && !sku) {
		return <LoadingSimple />
	}

	if (settingsQuery.loading) {
		return <LoadingSimple />
	}
	if (settingsQuery.error) {
		return <p>Error: {settingsQuery.error.message}</p>
	}
	if (!settingsQuery.data) {
		return <p>No settings returned</p>
	}

	return (
		<CenteredPage>
			<Spread>
				<Spaced>
					<FontAwesomeIcon icon={["fal", "barcode-alt"]} size="3x" />
					<H1>{isNewSKU ? "New SKU" : code}</H1>
				</Spaced>

				{!isNewSKU && sku ? (
					<div>
						<Button kind="secondary" type="button" onClick={() => window.open(`/portal/sku/new?clone=${sku.code}`, "_self")}>
							<FontAwesomeIcon icon={["far", "clone"]} />
							<div className={marginLeftStyle}>Clone SKU</div>
						</Button>
						<Button kind="secondary" type="button" onClick={() => setShowPreviewModal(true)} overrides={ButtonMarginLeftOverride}>
							<FontAwesomeIcon icon={["far", "mobile-android"]} color="#276EF1" />
							<div className={marginLeftStyle}>Preview</div>
						</Button>
					</div>
				) : (
					<div />
				)}
			</Spread>

			{isNewSKU ? (
				editForm
			) : (
				<Tabs
					onChange={({ activeKey }) => {
						setActiveKey(activeKey.toString())
						history.push(`/portal/sku/${code}${activeKey}`)
					}}
					activeKey={activeKey}
					overrides={{
						TabContent: {
							style: { ...paddingZero },
						},
						TabBar: {
							style: { ...paddingZero },
						},
					}}
				>
					<Tab
						key="#details"
						title={
							<Spaced>
								<FontAwesomeIcon icon={["fal", "pencil"]} />
								<div>Details</div>
							</Spaced>
						}
					>
						{editForm}
					</Tab>

					{sku && (sku.cloneParentID || sku.hasClones) && (
						<Tab
							key="#cloneTree"
							title={
								<Spaced>
									<FontAwesomeIcon icon={["fal", "folder-tree"]} />
									<div>Clone Tree</div>
								</Spaced>
							}
						>
							<SKUCloneTree skuID={sku.id} />
						</Tab>
					)}

					{!isNewSKU && sku && (
						<Tab
							key="#products"
							title={
								<Spaced>
									<FontAwesomeIcon icon={["fal", "steak"]} />
									<div>Products</div>
								</Spaced>
							}
						>
							<ItemList
								skuID={sku.id}
								itemName="product"
								query={graphql.query.PRODUCTS}
								batchActionMutation={graphql.mutation.BATCH_ACTION_PRODUCT}
								extraFilterOptions={[
									{ label: "Not in Carton", id: FilterOption.ProductWithoutCarton },
									{ label: "Not in Order", id: FilterOption.ProductWithoutOrder },
								]}
								itemLinks={["order", "contract", "distributor", "carton"]}
								actions={ActionItemSet.Products}
								showQRCodesToggle
							/>
						</Tab>
					)}
				</Tabs>
			)}

			<Modal isOpen={showPreviewModal} onClose={() => setShowPreviewModal(false)}>
				<SKUViewComponent sku={getPreviewSKU()} />
			</Modal>

			{/* Success Modal */}
			<Modal onClose={() => setShowSuccessModal(false)} isOpen={showSuccessModal}>
				<ModalHeader>
					<span>
						<FontAwesomeIcon icon={["fas", "check"]} />
						<span style={{ marginLeft: "10px" }}>SKU Updated Successfully</span>
					</span>
				</ModalHeader>
				<ModalFooter>
					<ModalButton onClick={() => setShowSuccessModal(false)}>OK</ModalButton>
				</ModalFooter>
			</Modal>
		</CenteredPage>
	)
}

interface InputPairProps {
	index: number
	prefix: string
	icon: IconName
	titleInputRef: any
	contentInputRef: any
	titleError?: string
	contentError?: string
	onDelete: () => void
	label1?: string
	label2?: string
}

export const InputPair = (props: InputPairProps) => {
	const { prefix, index } = props

	const [css, theme] = useStyletron()
	const containerStyle = css({
		display: "flex",
		alignItems: "center",
		padding: "8px 0 8px 10px",
		backgroundColor: index % 2 == 1 ? "rgba(0, 0, 0, 0.015)" : "unset",
	})
	const iconStyle = css({
		marginRight: "10px",
	})
	const half = css({
		display: "flex",
		alignItems: "center",
		width: "50%",
	})

	return (
		<div className={containerStyle}>
			<div className={half}>
				<div className={iconStyle}>
					<FontAwesomeIcon icon={["fas", props.icon]} size="lg" />
				</div>
				<FormControl label={props.label1 || "Name"} overrides={InputParFormControlOverrides} error="" positive="">
					<Input name={`${prefix}[${index}].title`} inputRef={props.titleInputRef} error={props.titleError !== undefined} />
				</FormControl>
			</div>
			<div className={half}>
				<FormControl label={props.label2 || "Value"} overrides={InputParFormControlOverrides} error="" positive="">
					<Input name={`${prefix}[${index}].content`} inputRef={props.contentInputRef} error={props.contentError !== undefined} />
				</FormControl>
				<Button
					type="button"
					kind="minimal"
					onClick={props.onDelete}
					overrides={{
						BaseButton: {
							style: {
								color: "grey",
								":hover": {
									color: "#d63916",
								},
							},
						},
					}}
				>
					<FontAwesomeIcon icon={["fas", "trash"]} size="lg" />
				</Button>
			</div>
		</div>
	)
}

const InputParFormControlOverrides: FormControlOverrides = {
	Label: {
		style: {
			marginTop: 0,
			marginBottom: 0,
			marginRight: "10px",
			width: "unset",
		},
	},
	ControlContainer: {
		style: {
			marginBottom: 0,
			marginRight: "20px",
			width: "100%",
		},
	},
}

export default SKUEdit
