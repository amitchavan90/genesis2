import * as React from "react"
import { useMutation } from "@apollo/react-hooks"
import { useStyletron } from "baseui"
import { ExecutionResult, MutationFunctionOptions } from "@apollo/react-common"

import { FormControl } from "baseui/form-control"
import { FileUploader } from "baseui/file-uploader"
import { Block } from "baseui/block"
import { mutation } from "../graphql/mutations"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { Blob } from "../types/types"
import { Button } from "baseui/button"
import { Modal } from "baseui/modal"

interface ImageUploadOutput {
	fileUpload: Blob
}
interface ImageUploadMultipleOutput {
	fileUploadMultiple: Blob[]
}

interface ImageUploadProps {
	name: string
	imageURL: string
	disabled?: boolean
	setImageUploader: (image: ImageUploadHandler | undefined) => void
	label?: string
	buttonLabel?: string
	caption?: string
	maxFileSize?: number

	// Used for profile editing (allows sync between profile preview and editing modal)
	imageRemoved?: boolean
	file?: File
	setFile?: (file?: File) => void

	previewHeight?: string
	previewWidth?: string

	clearable: boolean
}

export interface ImageUploadHandler {
	removeImage: boolean
	upload?: (options?: MutationFunctionOptions<ImageUploadOutput, { file: File }>) => Promise<void | ExecutionResult<ImageUploadOutput>>
	setUploadError?: React.Dispatch<React.SetStateAction<string>>
}

const Single: React.FunctionComponent<ImageUploadProps> = props => {
	const [uploadError, setUploadError] = React.useState<string>("")
	const [removed, setRemoved] = React.useState<boolean>(props.imageRemoved === true)
	const [file, setFile] = React.useState<File | undefined>(props.file)
	const [uploadImage, mutUploadImage] = useMutation<ImageUploadOutput, { file: File | undefined }>(mutation.FILE_UPLOAD, {
		variables: { file },
	})

	const previewImage = !removed ? (file ? URL.createObjectURL(file) : props.imageURL) : ""

	const [css, theme] = useStyletron()
	const imagePreview = css({
		display: "flex",
		height: props.previewHeight || "100px",
		width: props.previewWidth || "unset",
	})
	const imageStyle = css({
		borderRadius: "5px",
		height: "200px",
		objectFit: "contain",
		cursor: "pointer",
	})
	const imagePreviewStyle = css({
		width: "100%",
		cursor: "pointer",
	})
	const removeButton = css({
		position: "relative",
		top: "-10px",
		right: "10px",
		backgroundColor: "grey",
		width: "20px",
		height: "20px",
		borderRadius: "10px",
		":hover": {
			backgroundColor: "#d63916",
		},
		transition: "0.2s",
	})
	const removeButtonX = css({
		color: "white",
		height: "20px",
		padding: "0 5px",
	})

	const [showPreviewModal, setShowPreviewModal] = React.useState<boolean>()

	React.useEffect(() => {
		// Remove blob attachment after uploading
		if (!mutUploadImage.data) return
		props.setImageUploader({
			removeImage: false,
			upload: undefined,
			setUploadError: undefined,
		})
	}, [mutUploadImage.data])

	return (
		<FormControl label={props.label} disabled={props.disabled} caption={props.caption} error={uploadError} positive="">
			<Block as="div" display="flex">
				{!previewImage && (
					<FileUploader
						multiple={false}
						maxSize={props.maxFileSize || 5e6}
						disabled={props.disabled}
						onDrop={(acceptedFiles, rejectedFiles) => {
							if (acceptedFiles.length == 0) {
								if (rejectedFiles[0].size >= (props.maxFileSize || 5e6)) {
									setUploadError("File size too large")
								} else {
									setUploadError("Invalid file")
								}

								return
							}

							setUploadError("")
							setRemoved(false)
							setFile(acceptedFiles[0])
							if (props.setFile) props.setFile(acceptedFiles[0])

							props.setImageUploader({
								removeImage: false,
								upload: uploadImage,
								setUploadError: setUploadError,
							})
						}}
						accept=".jpg, .jpeg, .png, .svg, .gif, .bmp"
						progressMessage={mutUploadImage.loading ? "Uploading..." : ""}
						overrides={{
							Root: {
								style: {
									width: "unset",
									boxShadow: `0px 6px 12px #00000027`,
								},
							},

							ContentMessage: {
								style: { display: "none" },
							},
							ContentSeparator: {
								style: { display: "none" },
							},
							FileDragAndDrop: {
								style: {
									paddingLeft: 0,
									paddingRight: 0,
									paddingTop: 0,
									paddingBottom: 0,
									borderStyle: "unset",
									outline: "unset",
								},
							},
							ButtonComponent: { component: BrowseFileButton },
						}}
					/>
				)}

				{previewImage && (
					<div className={imagePreview}>
						<img className={imageStyle} src={previewImage} onClick={() => setShowPreviewModal(true)} />
						{props.clearable && (
							<div
								className={removeButton}
								onClick={() => {
									setUploadError("")
									setFile(undefined)
									if (props.setFile) props.setFile(undefined)
									setRemoved(true)
									props.setImageUploader({
										removeImage: true,
										upload: undefined,
										setUploadError: undefined,
									})
								}}
							>
								<FontAwesomeIcon icon={["fal", "times"]} className={removeButtonX} />
							</div>
						)}
						<Modal
							isOpen={showPreviewModal}
							onClose={() => setShowPreviewModal(false)}
							overrides={{
								Dialog: {
									style: {
										width: "unset",
										backgroundColor: "unset",
									},
								},
							}}
						>
							<img className={imagePreviewStyle} src={previewImage} />
						</Modal>
					</div>
				)}
			</Block>
		</FormControl>
	)
}
Single.defaultProps = { label: "Avatar" }

interface MultipleImageUploadProps {
	name: string
	files?: Blob[]
	disabled?: boolean
	setImageUploader: (image: ImageUploadMultipleHandler | undefined) => void
	label?: string
	buttonLabel?: string
	caption?: string
	maxFiles?: number
	maxFileSize?: number
	maxWidth?: string
}

export interface ImageUploadMultipleHandler {
	newFiles?: File[]
	currentFiles: Blob[]
	upload?: (options?: MutationFunctionOptions<ImageUploadMultipleOutput, { files: File[] }>) => Promise<void | ExecutionResult<ImageUploadMultipleOutput>>
	setUploadError?: React.Dispatch<React.SetStateAction<string>>
}

const Multiple: React.FunctionComponent<MultipleImageUploadProps> = props => {
	const [uploadError, setUploadError] = React.useState<string>("")
	const [files, setFiles] = React.useState<File[]>([])
	const [currentFiles, setCurrentFiles] = React.useState<Blob[]>([])
	const [uploadImages, mutUploadImages] = useMutation<ImageUploadMultipleOutput, { files: File[] }>(mutation.FILE_UPLOAD_MULTIPLE, {
		variables: { files },
	})

	const [css, theme] = useStyletron()
	const imagePreviews = css({
		display: "grid",
		gridTemplateColumns: "30% 30% 30%",
		gridColumnGap: "5%",
		gridRowGap: "10px",
		width: "100%",
		maxWidth: props.maxWidth,
		marginBottom: "10px",
	})

	React.useEffect(() => {
		if (props.files) {
			setCurrentFiles(props.files)
			setFiles([])
		}
	}, [props.files])

	React.useEffect(() => {
		// Remove blob attachment after uploading
		if (!mutUploadImages.data) return
		props.setImageUploader({
			newFiles: undefined,
			currentFiles: currentFiles,
			upload: undefined,
			setUploadError: undefined,
		})
	}, [mutUploadImages.data])

	return (
		<>
			<FormControl
				label={props.label + (props.maxFiles ? ` (${files.length + currentFiles.length}/${props.maxFiles})` : "")}
				disabled={props.disabled}
				caption={props.caption}
				error={uploadError}
				positive=""
			>
				<Block as="div" display="flex">
					<FileUploader
						multiple={true}
						maxSize={props.maxFileSize || 5e6}
						disabled={props.disabled}
						onDrop={(acceptedFiles, rejectedFiles) => {
							if (acceptedFiles.length == 0) {
								if (rejectedFiles[0].size >= (props.maxFileSize || 5e6)) {
									setUploadError("File size too large")
								} else {
									setUploadError("Invalid file")
								}

								return
							}
							if (props.maxFiles && files.length + acceptedFiles.length > props.maxFiles) {
								setUploadError(`Max images: ${props.maxFiles}`)
								return
							}
							setUploadError("")

							const f = [...files, ...acceptedFiles]
							setFiles(f)

							props.setImageUploader({
								newFiles: f,
								currentFiles: currentFiles,
								upload: uploadImages,
								setUploadError: setUploadError,
							})
						}}
						accept=".jpg, .jpeg, .png, .svg, .gif, .bmp"
						progressMessage={mutUploadImages.loading ? "Uploading..." : ""}
						overrides={{
							Root: {
								style: {
									width: "100%",
									boxShadow: `0px 6px 12px #00000027`,
									maxWidth: props.maxWidth,
								},
							},

							ContentMessage: {
								style: { display: "none" },
							},
							ContentSeparator: {
								style: { display: "none" },
							},
							FileDragAndDrop: {
								style: {
									paddingLeft: 0,
									paddingRight: 0,
									paddingTop: 0,
									paddingBottom: 0,
									borderStyle: "unset",
									outline: "unset",
								},
							},
							ButtonComponent: { component: BrowseFilesButton },
						}}
					/>
				</Block>
			</FormControl>
			{(files.length > 0 || currentFiles.length > 0) && (
				<div className={imagePreviews}>
					{currentFiles.map((file: Blob, index: number) => (
						<ImagePreview
							key={`currentFile-${index}-${file.id}`}
							src={file.file_url}
							onRemove={() => {
								const cf = [...currentFiles.slice(0, index), ...currentFiles.slice(index + 1)]
								setCurrentFiles(cf)
								props.setImageUploader({
									currentFiles: cf,
									upload: uploadImages,
									setUploadError: setUploadError,
								})
							}}
						/>
					))}
					{files.map((file: File, index: number) => (
						<ImagePreview
							key={`file-${index}-${file.name}`}
							src={URL.createObjectURL(file)}
							onRemove={() => {
								const f = [...files.slice(0, index), ...files.slice(index + 1)]
								setFiles(f)

								props.setImageUploader({
									newFiles: f,
									currentFiles: currentFiles,
									upload: uploadImages,
									setUploadError: setUploadError,
								})
							}}
						/>
					))}
				</div>
			)}
		</>
	)
}

const ImagePreview = (props: { src: string; onRemove: () => void }) => {
	const [css, theme] = useStyletron()
	const imagePreview = css({
		display: "flex",
		width: "100%",
	})
	const imageStyle = css({
		borderRadius: "5px",
		width: "100%",
		objectFit: "contain",
		cursor: "pointer",
	})
	const removeButton = css({
		position: "relative",
		top: "-10px",
		right: "10px",
		backgroundColor: "grey",
		width: "20px",
		height: "20px",
		borderRadius: "10px",
		":hover": {
			backgroundColor: "#d63916",
		},
		transition: "0.2s",
	})
	const removeButtonX = css({
		color: "white",
		height: "20px",
		padding: "0 5px",
	})

	const [showPreviewModal, setShowPreviewModal] = React.useState<boolean>()

	return (
		<div className={imagePreview}>
			<img className={imageStyle} src={props.src} onClick={() => setShowPreviewModal(true)} />
			<div className={removeButton} onClick={props.onRemove}>
				<FontAwesomeIcon icon={["fal", "times"]} className={removeButtonX} />
			</div>
			<Modal
				isOpen={showPreviewModal}
				onClose={() => setShowPreviewModal(false)}
				overrides={{
					Dialog: {
						style: {
							width: "unset",
							backgroundColor: "unset",
						},
					},
				}}
			>
				<img className={imageStyle} src={props.src} />
			</Modal>
		</div>
	)
}

const BrowseFileButton = ({ children, ...rest }: any) => {
	return (
		<Button type="button" {...rest} overrides={{ BaseButton: { style: { width: "100%" } } }}>
			Select image
		</Button>
	)
}
const BrowseFilesButton = ({ children, ...rest }: any) => {
	return (
		<Button type="button" {...rest} overrides={{ BaseButton: { style: { width: "100%" } } }}>
			Select images
		</Button>
	)
}

export const ImageUpload = {
	Single,
	Multiple,
}
