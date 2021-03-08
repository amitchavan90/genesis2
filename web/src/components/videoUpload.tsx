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

interface VideoUploadOutput {
	fileUpload: Blob
}

interface VideoUploadProps {
	name: string
	videoURL: string
	disabled?: boolean
	setVideoUploader: (image: VideoUploadHandler | undefined) => void
	label?: string
	buttonLabel?: string
	caption?: string
	maxFileSize?: number

	// Used for profile editing (allows sync between profile preview and editing modal)
	imageRemoved?: boolean
	file?: File
	setFile?: (file?: File) => void
}

export interface VideoUploadHandler {
	removeFile: boolean
	upload?: (options?: MutationFunctionOptions<VideoUploadOutput, { file: File }>) => Promise<void | ExecutionResult<VideoUploadOutput>>
	setUploadError?: React.Dispatch<React.SetStateAction<string>>
}

export const VideoUpload: React.FunctionComponent<VideoUploadProps> = props => {
	const [uploadError, setUploadError] = React.useState<string>("")
	const [removed, setRemoved] = React.useState<boolean>(props.imageRemoved === true)
	const [file, setFile] = React.useState<File | undefined>(props.file)
	const [uploadVideo, mutUploadVideo] = useMutation<VideoUploadOutput, { file: File | undefined }>(mutation.FILE_UPLOAD, {
		variables: { file },
	})

	const previewVideo = !removed ? (file ? URL.createObjectURL(file) : props.videoURL) : ""

	const [css, theme] = useStyletron()
	const videoPreview = css({
		display: "flex",
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

	React.useEffect(() => {
		// Remove blob attachment after uploading
		if (!mutUploadVideo.data) return
		props.setVideoUploader({
			removeFile: false,
			upload: undefined,
			setUploadError: undefined,
		})
	}, [mutUploadVideo.data])

	return (
		<FormControl label={props.label || "Video"} disabled={props.disabled} caption={props.caption} error={uploadError} positive="">
			<Block as="div" display="flex">
				{!previewVideo && (
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

							props.setVideoUploader({
								removeFile: false,
								upload: uploadVideo,
								setUploadError: setUploadError,
							})
						}}
						accept=".mp4"
						progressMessage={mutUploadVideo.loading ? "Uploading..." : ""}
					/>
				)}

				{previewVideo && (
					<div className={videoPreview}>
						<video src={previewVideo} height="200px" controls />
						<div
							className={removeButton}
							onClick={() => {
								setUploadError("")
								setFile(undefined)
								if (props.setFile) props.setFile(undefined)
								setRemoved(true)
								props.setVideoUploader({
									removeFile: true,
									upload: undefined,
									setUploadError: undefined,
								})
							}}
						>
							<FontAwesomeIcon icon={["fal", "times"]} className={removeButtonX} />
						</div>
					</div>
				)}
			</Block>
		</FormControl>
	)
}
