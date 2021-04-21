import gql from "graphql-tag"

const ROLE = gql`
	fragment RoleFragment on Role {
		id
		name
		tier
		archived
		permissions
		trackActions {
			id
		}
		createdAt
	}
`

const USER = gql`
	fragment UserFragment on User {
		id
		firstName
		lastName
		email
		verified
		archived
		role {
			...RoleFragment
		}
		affiliateOrg
		wechatID
		mobilePhone
		mobileVerified
		createdAt
	}
	${ROLE}
`

const CONSUMER = gql`
	fragment ConsumerFragment on User {
		id
		firstName
		lastName
		email
		verified
		role {
			...RoleFragment
		}
		affiliateOrg
		wechatID
		createdAt

		loyaltyPoints
	}
	${ROLE}
`

const SKU = gql`
	fragment SKUFragment on SKU {
		id
		name
		code
		description
		isBeef
		isPointBound
		isAppBound
		brand
		ingredients
		weight
		price
		purchasePoints
		ingredients
		archived
		productCount
		loyaltyPoints

		hasClones
		cloneParentID
		gif{
			id
			file_url
		}
		brandLogo {
			id
			file_url
		}
		masterPlan {
			id
			file_url
		}
		video {
			id
			file_url
		}
		photos {
			id
			file_url
		}
		urls {
			title
			content
		}
		retailLinks {
			name
			url
		}
		productInfo {
			title
			content
		}
		categories {
            name
        }
        productCategories {
            name
        }
	}
`
const SKU_LIST = gql`
	fragment SKUListFragment on SKU {
		id
		name
		code
		description
		archived
		createdAt
		productCount

		masterPlan {
			id
			file_url
		}
	}
`

const CONTAINER = gql`
	fragment ContainerFragment on Container {
		id
		code
		archived
		createdAt
		palletCount
		description
	}
`

const PALLET = gql`
	fragment PalletFragment on Pallet {
		id
		code
		archived
		createdAt
		cartonCount
		description

		container {
			id
			code
		}
	}
`
const PALLET_LIST = gql`
	fragment PalletListFragment on Pallet {
		id
		code
		archived
		createdAt
		cartonCount
		description
		latestTrackAction {
			name
			createdAt
		}
	}
`

const CARTON = gql`
	fragment CartonFragment on Carton {
		id
		code
		archived
		createdAt
		weight
		meatType
		processedAt
		description
		spreadsheetLink
		productCount

		pallet {
			id
			code
			container {
				id
				code
			}
		}
		sku {
			id
			code
		}
		order {
			id
			code
		}
		distributor {
			id
			code
		}
	}
`
const CARTON_LIST = gql`
	fragment CartonListFragment on Carton {
		id
		code
		archived
		createdAt
		weight
		meatType
		productCount
		spreadsheetLink
		latestTrackAction {
			name
			createdAt
		}
		description

		pallet {
			id
			code
		}
	}
`

const PRODUCT = gql`
	fragment ProductFragment on Product {
		id
		code
		registerID
		loyaltyPoints
		loyaltyPointsExpire
		archived
		createdAt
		latestTrackAction {
			name
			createdAt
		}
		description

		sku {
			id
			code
			name
		}
		carton {
			id
			code
			pallet {
				id
				code
				container {
					id
					code
				}
			}
		}
		order {
			id
			code
		}
		distributor {
			id
			code
		}
		contract {
			id
			code
			name
			supplierName
		}

		registered
		registeredBy {
			id
			wechatID
		}
	}
`
const PRODUCT_LIST = gql`
	fragment ProductListFragment on Product {
		id
		code
		description
		loyaltyPoints
		loyaltyPointsExpire
		archived
		createdAt
		latestTrackAction {
			name
			createdAt
		}
		transactions {
			id
			transactionHash
			carton {
				id
			}
		}

		sku {
			id
			code
			name
		}
		carton {
			id
			code
		}
		order {
			id
			code
		}
		distributor {
			id
			code
		}
		contract {
			id
			code
			name
		}
	}
`

const ORDER = gql`
	fragment OrderFragment on Order {
		id
		code
		archived
		createdAt
		productCount
		isAppBound
		sku {
			id
			code
			name
			description
			masterPlan {
				id
				file_url
			}
		}
	}
`
const ORDER_LIST = gql`
	fragment OrderListFragment on Order {
		id
		code
		archived
		createdAt
		productCount
		isAppBound
		sku {
			id
			code
		}
	}
`

const CONTRACT = gql`
	fragment ContractFragment on Contract {
		id
		code
		name
		description
		latitude
		longitude
		supplierName
		dateSigned
		archived
		createdAt
	}
`

const DISTRIBUTOR = gql`
	fragment DistributorFragment on Distributor {
		id
		name
		code
		archived
		createdAt
	}
`

const TRACK_ACTION = gql`
	fragment TrackActionFragment on TrackAction {
		id
		code
		name
		nameChinese
		private
		requirePhotos
		blockchain
		system
		archived
		createdAt
	}
`

const TRANSACTION = gql`
	fragment TransactionFragment on Transaction {
		id
		transactionHash
		manifestID
		manifestLineSha256
		transactionPending
		locationGeohash
		locationName
		memo
		action {
			...TrackActionFragment
		}
		createdAt
		createdBy {
			id
			email
			firstName
			lastName
			affiliateOrg
			wechatID
		}
		createdByName
		createdAt
		scannedAt
	}
	${TRACK_ACTION}
`
const TRANSACTION_VIEW = gql`
	fragment TransactionViewFragment on Transaction {
		id
		transactionHash
		manifestID
		manifest {
			id
			transactionHash
			merkleRootSha256
		}
		manifestLineJson
		manifestLineSha256
		locationGeohash
		locationName
		photos {
			cartonPhoto {
				id
				file_url
			}
			productPhoto {
				id
				file_url
			}
		}
		memo
		action {
			...TrackActionFragment
		}
		createdByName
		createdAt
		scannedAt
	}
	${TRACK_ACTION}
`

/** Basic product info w/ sku info (for consumer product view) */
const PRODUCT_VIEW = gql`
	fragment ProductViewFragment on Product {
		id
		code
		loyaltyPoints
		loyaltyPointsExpire
		archived
		registered

		sku {
			...SKUFragment
		}
		contract {
			id
			supplierName
			description
		}

		transactions {
			...TransactionViewFragment
		}
		createdAt
	}
	${SKU}
	${TRANSACTION_VIEW}
`
const TASK = gql`
	fragment TaskFragment on Task {
		id,
		title,
		code,
		description,
		loyaltyPoints,
		isTimeBound,
		isPeopleBound,
		isProductRelevant,
		finishDate,
		maximumPeople,
		bannerPhoto{
			id
			file_url
		},
		sku {
			id,
			name,
			code
		},
		subtasks {
			id,
			title,
			description
		}
		createdAt
	}
`
const TASK_LIST = gql`
	fragment TaskListFragment on Task {
			id,
			title,
			code,
			description,
			loyaltyPoints,
			isTimeBound,
			isPeopleBound,
			isProductRelevant,
			isFinal,
			finishDate,
			maximumPeople,
			createdAt
	}
`
const REFERRAL_LIST = gql`
fragment ReferralListFragment on Referral{
	id,
	code,
    referee {
        id,
        firstName,
        lastName,
		email,
        referralCode
    },
    isRedemmed,
    createdAt,
    user {
        firstName,
        lastName,
        email
    }
}`

const  USER_TASK_LIST = gql`
fragment UserTaskListFragment on UserTask {
	id,
	isComplete,
	status,
	code,
	task {
		id,
		title,
		description,
	},
	user {
		id,
		firstName,
		lastName,
		email
	}
	userSubtasks {
		id,
		isComplete,
		status
	},
	createdAt
}`

const  USER_TASK = gql`
fragment UserTaskFragment on UserTask {
	id,
	isComplete,
	status,
	code,
	task {
		id,
		title,
		description,
	},
	user {
		id,
		firstName,
		lastName,
		email
	}
	userSubtasks {
		id,
		isComplete,
		status
	},
	createdAt
}`

const  USER_PURCHASE_ACTIVITY = gql`
fragment UserPurchaseActivityFragment on  UserPurchaseActivity{
	id,
	loyaltyPoints,
	user {
		id,
		firstName,
		lastName,
		email
	},
	product {
		id
		code
	}
	createdAt
}`

export const fragment = {
	USER,
	CONSUMER,
	ROLE,

	SKU,
	SKU_LIST,

	CONTAINER,

	PALLET,
	PALLET_LIST,

	CARTON,
	CARTON_LIST,

	PRODUCT,
	PRODUCT_LIST,
	PRODUCT_VIEW,

	ORDER,
	ORDER_LIST,

	CONTRACT,

	DISTRIBUTOR,

	TRACK_ACTION,
	TRANSACTION,
	TRANSACTION_VIEW,

	TASK,
	TASK_LIST,

	REFERRAL_LIST,

	USER_TASK,
	USER_TASK_LIST,
	USER_PURCHASE_ACTIVITY
}
