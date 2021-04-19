import gql from "graphql-tag"
import { fragment } from "./fragments"

const FILE_UPLOAD = gql`
	mutation fileUpload($file: Upload!) {
		fileUpload(file: $file) {
			id
			file_url
		}
	}
`
const FILE_UPLOAD_MULTIPLE = gql`
	mutation fileUploadMultiple($files: [Upload!]!) {
		fileUploadMultiple(files: $files) {
			id
			file_url
		}
	}
`

const CHANGE_PASSWORD = gql`
	mutation changePassword($oldPassword: String!, $password: String!) {
		changePassword(oldPassword: $oldPassword, password: $password)
	}
`

const CHANGE_DETAILS = gql`
	mutation changeDetails($input: UpdateUser!) {
		changeDetails(input: $input) {
			...UserFragment
		}
	}
	${fragment.USER}
`

const FORGOT_PASSWORD = gql`
	mutation forgotPassword($email: String!, $viaSMS: Boolean) {
		forgotPassword(email: $email, viaSMS: $viaSMS)
	}
`
const RESET_PASSWORD = gql`
	mutation resetPassword($token: String!, $password: String!, $email: NullString) {
		resetPassword(token: $token, password: $password, email: $email)
	}
`
const RESEND_EMAIL_VERIFICATION = gql`
	mutation resendEmailVerification($email: String!) {
		resendEmailVerification(email: $email)
	}
`

const UPDATE_USER = gql`
	mutation userUpdate($id: ID!, $input: UpdateUser!) {
		userUpdate(id: $id, input: $input) {
			...UserFragment
		}
	}
	${fragment.USER}
`
const CREATE_USER = gql`
	mutation userCreate($input: UpdateUser!) {
		userCreate(input: $input) {
			...UserFragment
		}
	}
	${fragment.USER}
`
const ARCHIVE_USER = gql`
	mutation userArchive($id: ID!) {
		userArchive(id: $id) {
			...UserFragment
		}
	}
	${fragment.USER}
`
const UNARCHIVE_USER = gql`
	mutation userUnarchive($id: ID!) {
		userUnarchive(id: $id) {
			...UserFragment
		}
	}
	${fragment.USER}
`

// Role
const UPDATE_ROLE = gql`
	mutation roleUpdate($id: ID!, $input: UpdateRole!) {
		roleUpdate(id: $id, input: $input) {
			...RoleFragment
		}
	}
	${fragment.ROLE}
`
const CREATE_ROLE = gql`
	mutation roleCreate($input: UpdateRole!) {
		roleCreate(input: $input) {
			...RoleFragment
		}
	}
	${fragment.ROLE}
`
const ARCHIVE_ROLE = gql`
	mutation roleArchive($id: ID!) {
		roleArchive(id: $id) {
			...RoleFragment
		}
	}
	${fragment.ROLE}
`
const UNARCHIVE_ROLE = gql`
	mutation roleUnarchive($id: ID!) {
		roleUnarchive(id: $id) {
			...RoleFragment
		}
	}
	${fragment.ROLE}
`

// SKU
const UPDATE_SKU = gql`
	mutation skuUpdate($id: ID!, $input: UpdateSKU!) {
		skuUpdate(id: $id, input: $input) {
			...SKUFragment
		}
	}
	${fragment.SKU}
`
const CREATE_SKU = gql`
	mutation skuCreate($input: UpdateSKU!) {
		skuCreate(input: $input) {
			...SKUFragment
		}
	}
	${fragment.SKU}
`
const ARCHIVE_SKU = gql`
	mutation skuArchive($id: ID!) {
		skuArchive(id: $id) {
			...SKUFragment
		}
	}
	${fragment.SKU}
`
const UNARCHIVE_SKU = gql`
	mutation skuUnarchive($id: ID!) {
		skuUnarchive(id: $id) {
			...SKUFragment
		}
	}
	${fragment.SKU}
`
const BATCH_ACTION_SKU = gql`
	mutation skuBatchAction($ids: [ID!]!, $action: Action!, $value: BatchActionInput) {
		skuBatchAction(ids: $ids, action: $action, value: $value)
	}
`

// Container
const UPDATE_CONTAINER = gql`
	mutation containerUpdate($id: ID!, $input: UpdateContainer!) {
		containerUpdate(id: $id, input: $input) {
			...ContainerFragment
		}
	}
	${fragment.CONTAINER}
`
const CREATE_CONTAINER = gql`
	mutation containerCreate($input: CreateContainer!) {
		containerCreate(input: $input)
	}
`
const ARCHIVE_CONTAINER = gql`
	mutation containerArchive($id: ID!) {
		containerArchive(id: $id) {
			...ContainerFragment
		}
	}
	${fragment.CONTAINER}
`
const UNARCHIVE_CONTAINER = gql`
	mutation containerUnarchive($id: ID!) {
		containerUnarchive(id: $id) {
			...ContainerFragment
		}
	}
	${fragment.CONTAINER}
`
const BATCH_ACTION_CONTAINER = gql`
	mutation containerBatchAction($ids: [ID!]!, $action: Action!, $value: BatchActionInput) {
		containerBatchAction(ids: $ids, action: $action, value: $value)
	}
`

// Pallet
const UPDATE_PALLET = gql`
	mutation palletUpdate($id: ID!, $input: UpdatePallet!) {
		palletUpdate(id: $id, input: $input) {
			...PalletFragment
		}
	}
	${fragment.PALLET}
`
const CREATE_PALLET = gql`
	mutation palletCreate($input: CreatePallet!) {
		palletCreate(input: $input)
	}
`
const ARCHIVE_PALLET = gql`
	mutation palletArchive($id: ID!) {
		palletArchive(id: $id) {
			...PalletFragment
		}
	}
	${fragment.PALLET}
`
const UNARCHIVE_PALLET = gql`
	mutation palletUnarchive($id: ID!) {
		palletUnarchive(id: $id) {
			...PalletFragment
		}
	}
	${fragment.PALLET}
`
const BATCH_ACTION_PALLET = gql`
	mutation palletBatchAction($ids: [ID!]!, $action: Action!, $value: BatchActionInput) {
		palletBatchAction(ids: $ids, action: $action, value: $value)
	}
`

// Carton
const UPDATE_CARTON = gql`
	mutation cartonUpdate($id: ID!, $input: UpdateCarton!) {
		cartonUpdate(id: $id, input: $input) {
			...CartonFragment
		}
	}
	${fragment.CARTON}
`
const CREATE_CARTON = gql`
	mutation cartonCreate($input: CreateCarton!) {
		cartonCreate(input: $input)
	}
`
const ARCHIVE_CARTON = gql`
	mutation cartonArchive($id: ID!) {
		cartonArchive(id: $id) {
			...CartonFragment
		}
	}
	${fragment.CARTON}
`
const UNARCHIVE_CARTON = gql`
	mutation cartonUnarchive($id: ID!) {
		cartonUnarchive(id: $id) {
			...CartonFragment
		}
	}
	${fragment.CARTON}
`
const BATCH_ACTION_CARTON = gql`
	mutation cartonBatchAction($ids: [ID!]!, $action: Action!, $value: BatchActionInput) {
		cartonBatchAction(ids: $ids, action: $action, value: $value)
	}
`

// Product
const UPDATE_PRODUCT = gql`
	mutation productUpdate($id: ID!, $input: UpdateProduct!) {
		productUpdate(id: $id, input: $input) {
			...ProductFragment
		}
	}
	${fragment.PRODUCT}
`
const CREATE_PRODUCT = gql`
	mutation productCreate($input: UpdateProduct!) {
		productCreate(input: $input) {
			...ProductFragment
		}
	}
	${fragment.PRODUCT}
`
const ARCHIVE_PRODUCT = gql`
	mutation productArchive($id: ID!) {
		productArchive(id: $id) {
			...ProductFragment
		}
	}
	${fragment.PRODUCT}
`
const UNARCHIVE_PRODUCT = gql`
	mutation productUnarchive($id: ID!) {
		productUnarchive(id: $id) {
			...ProductFragment
		}
	}
	${fragment.PRODUCT}
`
const BATCH_ACTION_PRODUCT = gql`
	mutation productBatchAction($ids: [ID!]!, $action: Action!, $value: BatchActionInput) {
		productBatchAction(ids: $ids, action: $action, value: $value)
	}
`

// Order
const UPDATE_ORDER = gql`
	mutation orderUpdate($id: ID!, $input: UpdateOrder!) {
		orderUpdate(id: $id, input: $input) {
			...OrderFragment
		}
	}
	${fragment.ORDER}
`
const CREATE_ORDER = gql`
	mutation orderCreate($input: CreateOrder!) {
		orderCreate(input: $input) {
			...OrderFragment
		}
	}
	${fragment.ORDER}
`
const ARCHIVE_ORDER = gql`
	mutation orderArchive($id: ID!) {
		orderArchive(id: $id) {
			...OrderFragment
		}
	}
	${fragment.ORDER}
`
const UNARCHIVE_ORDER = gql`
	mutation orderUnarchive($id: ID!) {
		orderUnarchive(id: $id) {
			...OrderFragment
		}
	}
	${fragment.ORDER}
`
const BATCH_ACTION_ORDER = gql`
	mutation orderBatchAction($ids: [ID!]!, $action: Action!, $value: BatchActionInput) {
		orderBatchAction(ids: $ids, action: $action, value: $value)
	}
`

// Contract
const UPDATE_CONTRACT = gql`
	mutation contractUpdate($id: ID!, $input: UpdateContract!) {
		contractUpdate(id: $id, input: $input) {
			...ContractFragment
		}
	}
	${fragment.CONTRACT}
`
const CREATE_CONTRACT = gql`
	mutation contractCreate($input: UpdateContract!) {
		contractCreate(input: $input) {
			...ContractFragment
		}
	}
	${fragment.CONTRACT}
`
const ARCHIVE_CONTRACT = gql`
	mutation contractArchive($id: ID!) {
		contractArchive(id: $id) {
			...ContractFragment
		}
	}
	${fragment.CONTRACT}
`
const UNARCHIVE_CONTRACT = gql`
	mutation contractUnarchive($id: ID!) {
		contractUnarchive(id: $id) {
			...ContractFragment
		}
	}
	${fragment.CONTRACT}
`

// Distributor
const UPDATE_DISTRIBUTOR = gql`
	mutation distributorUpdate($id: ID!, $input: UpdateDistributor!) {
		distributorUpdate(id: $id, input: $input) {
			...DistributorFragment
		}
	}
	${fragment.DISTRIBUTOR}
`
const CREATE_DISTRIBUTOR = gql`
	mutation distributorCreate($input: UpdateDistributor!) {
		distributorCreate(input: $input) {
			...DistributorFragment
		}
	}
	${fragment.DISTRIBUTOR}
`
const ARCHIVE_DISTRIBUTOR = gql`
	mutation distributorArchive($id: ID!) {
		distributorArchive(id: $id) {
			...DistributorFragment
		}
	}
	${fragment.DISTRIBUTOR}
`
const UNARCHIVE_DISTRIBUTOR = gql`
	mutation distributorUnarchive($id: ID!) {
		distributorUnarchive(id: $id) {
			...DistributorFragment
		}
	}
	${fragment.DISTRIBUTOR}
`

// TrackAction
const UPDATE_TRACK_ACTION = gql`
	mutation trackActionUpdate($id: ID!, $input: UpdateTrackAction!) {
		trackActionUpdate(id: $id, input: $input) {
			...TrackActionFragment
		}
	}
	${fragment.TRACK_ACTION}
`
const CREATE_TRACK_ACTION = gql`
	mutation trackActionCreate($input: UpdateTrackAction!) {
		trackActionCreate(input: $input) {
			...TrackActionFragment
		}
	}
	${fragment.TRACK_ACTION}
`
const ARCHIVE_TRACK_ACTION = gql`
	mutation trackActionArchive($id: ID!) {
		trackActionArchive(id: $id) {
			...TrackActionFragment
		}
	}
	${fragment.TRACK_ACTION}
`
const UNARCHIVE_TRACK_ACTION = gql`
	mutation trackActionUnarchive($id: ID!) {
		trackActionUnarchive(id: $id) {
			...TrackActionFragment
		}
	}
	${fragment.TRACK_ACTION}
`

const DEPLOY_SMART_CONTRACT = gql`
	mutation deploySmartContract {
		deploySmartContract {
			consumerHost
			adminHost
			etherscanHost
			fieldappVersion
			smartContractAddress
		}
	}
`
const FLUSH_PENDING_TRANSACTIONS = gql`
	mutation flushPendingTransactions {
		flushPendingTransactions
	}
`
// TASK
const UPDATE_TASK = gql`
	mutation taskUpdate($id: ID!, $input: UpdateTask!) {
		taskUpdate(id: $id, input: $input) {
			...TaskFragment
		}
	}
	${fragment.TASK}
`
const CREATE_TASK = gql`
	mutation taskCreate($input: UpdateTask!) {
		taskCreate(input: $input) {
			...TaskFragment
		}
	}
	${fragment.TASK}
`
const BATCH_ACTION_TASK = gql`
	mutation taskBatchAction($ids: [ID!]!, $action: Action!, $value: BatchActionInput) {
		taskBatchAction(ids: $ids, action: $action, value: $value)
	}
`
const USER_TASK_APPROVE = gql`mutation userTaskApprove($id: ID!) {
	userTaskApprove(id: $id) {
		...UserTaskFragment
	}
}
${fragment.USER_TASK}`

export const mutation = {
	FILE_UPLOAD,
	FILE_UPLOAD_MULTIPLE,

	CHANGE_PASSWORD,
	CHANGE_DETAILS,

	FORGOT_PASSWORD,
	RESET_PASSWORD,
	RESEND_EMAIL_VERIFICATION,

	UPDATE_USER,
	CREATE_USER,
	ARCHIVE_USER,
	UNARCHIVE_USER,

	UPDATE_ROLE,
	CREATE_ROLE,
	ARCHIVE_ROLE,
	UNARCHIVE_ROLE,

	UPDATE_SKU,
	CREATE_SKU,
	ARCHIVE_SKU,
	UNARCHIVE_SKU,
	BATCH_ACTION_SKU,

	UPDATE_CONTAINER,
	CREATE_CONTAINER,
	ARCHIVE_CONTAINER,
	UNARCHIVE_CONTAINER,
	BATCH_ACTION_CONTAINER,

	UPDATE_PALLET,
	CREATE_PALLET,
	ARCHIVE_PALLET,
	UNARCHIVE_PALLET,
	BATCH_ACTION_PALLET,

	UPDATE_CARTON,
	CREATE_CARTON,
	ARCHIVE_CARTON,
	UNARCHIVE_CARTON,
	BATCH_ACTION_CARTON,

	UPDATE_PRODUCT,
	CREATE_PRODUCT,
	ARCHIVE_PRODUCT,
	UNARCHIVE_PRODUCT,
	BATCH_ACTION_PRODUCT,

	UPDATE_ORDER,
	CREATE_ORDER,
	ARCHIVE_ORDER,
	UNARCHIVE_ORDER,
	BATCH_ACTION_ORDER,

	UPDATE_CONTRACT,
	CREATE_CONTRACT,
	ARCHIVE_CONTRACT,
	UNARCHIVE_CONTRACT,

	UPDATE_DISTRIBUTOR,
	CREATE_DISTRIBUTOR,
	ARCHIVE_DISTRIBUTOR,
	UNARCHIVE_DISTRIBUTOR,

	UPDATE_TRACK_ACTION,
	CREATE_TRACK_ACTION,
	ARCHIVE_TRACK_ACTION,
	UNARCHIVE_TRACK_ACTION,

	DEPLOY_SMART_CONTRACT,
	FLUSH_PENDING_TRANSACTIONS,

	CREATE_TASK,
	UPDATE_TASK,
	BATCH_ACTION_TASK,

	USER_TASK_APPROVE
}
