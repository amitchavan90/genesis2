import gql from "graphql-tag"
import { fragment } from "./fragments"

const ME = gql`
	{
		me {
			...UserFragment
		}
	}
	${fragment.USER}
`

const USER = gql`
	query user($email: String, $wechatID: String) {
		user(email: $email, wechatID: $wechatID) {
			...UserFragment
		}
	}
	${fragment.USER}
`
const CONSUMER = gql`
	query user($email: String, $wechatID: String) {
		user(email: $email, wechatID: $wechatID) {
			...ConsumerFragment
		}
	}
	${fragment.CONSUMER}
`
const USERS = gql`
	query users($search: SearchFilter!, $limit: Int!, $offset: Int!) {
		users(search: $search, limit: $limit, offset: $offset) {
			users {
				...UserFragment
			}
			total
		}
	}
	${fragment.USER}
`
const CONSUMERS = gql`
	query consumers($search: SearchFilter!, $limit: Int!, $offset: Int!) {
		consumers(search: $search, limit: $limit, offset: $offset) {
			consumers {
				...ConsumerFragment
			}
			total
		}
	}
	${fragment.CONSUMER}
`

const ROLE = gql`
	query role($name: String!) {
		role(name: $name) {
			...RoleFragment
		}
	}
	${fragment.ROLE}
`
const ROLES = gql`
	query roles($search: SearchFilter!, $limit: Int!, $offset: Int!) {
		roles(search: $search, limit: $limit, offset: $offset, excludeSuper: false) {
			roles {
				...RoleFragment
			}
			total
		}
	}
	${fragment.ROLE}
`
const ROLES_LIMITED = gql`
	query roles($search: SearchFilter!, $limit: Int!, $offset: Int!) {
		roles(search: $search, limit: $limit, offset: $offset, excludeSuper: true) {
			roles {
				...RoleFragment
			}
			total
		}
	}
	${fragment.ROLE}
`

const SKU = gql`
	query sku($code: String!) {
		sku(code: $code) {
			...SKUFragment
		}
	}
	${fragment.SKU}
`

/** Get SKU by ID for consumer SKU view (no permission check)*/
const SKU_VIEW = gql`
	query skuByID($id: ID!) {
		skuByID(id: $id) {
			...SKUFragment
		}
	}
	${fragment.SKU}
`

const SKUS = gql`
	query skus($search: SearchFilter!, $limit: Int!, $offset: Int!) {
		skus(search: $search, limit: $limit, offset: $offset) {
			skus {
				...SKUListFragment
			}
			total
		}
	}
	${fragment.SKU_LIST}
`

const SKU_CLONE_TREE = gql`
	query skuCloneTree($id: ID!) {
		skuCloneTree(id: $id) {
			sku {
				id
				code
				name
			}
			depth
		}
	}
`

const CONTAINER = gql`
	query container($code: String!) {
		container(code: $code) {
			...ContainerFragment
		}
	}
	${fragment.CONTAINER}
`

const CONTAINERS = gql`
	query containers($search: SearchFilter!, $limit: Int!, $offset: Int!) {
		containers(search: $search, limit: $limit, offset: $offset) {
			containers {
				...ContainerFragment
			}
			total
		}
	}
	${fragment.CONTAINER}
`
const CONTAINERS_BASIC = gql`
	query containers($search: SearchFilter!, $limit: Int!, $offset: Int!) {
		containers(search: $search, limit: $limit, offset: $offset) {
			containers {
				id
				code
			}
			total
		}
	}
`

const PALLET = gql`
	query pallet($code: String!) {
		pallet(code: $code) {
			...PalletFragment
		}
	}
	${fragment.PALLET}
`

const PALLETS = gql`
	query pallets($search: SearchFilter!, $limit: Int!, $offset: Int!, $containerID: String, $trackActionID: String) {
		pallets(search: $search, limit: $limit, offset: $offset, containerID: $containerID, trackActionID: $trackActionID) {
			pallets {
				...PalletListFragment
			}
			total
		}
	}
	${fragment.PALLET_LIST}
`

const PALLETS_BASIC = gql`
	query pallets($search: SearchFilter!, $limit: Int!, $offset: Int!, $containerID: String, $trackActionID: String) {
		pallets(search: $search, limit: $limit, offset: $offset, containerID: $containerID, trackActionID: $trackActionID) {
			pallets {
				id
				code
			}
			total
		}
	}
`

const CARTON = gql`
	query carton($code: String!) {
		carton(code: $code) {
			...CartonFragment
		}
	}
	${fragment.CARTON}
`

const CARTONS = gql`
	query cartons($search: SearchFilter!, $limit: Int!, $offset: Int!, $palletID: String, $trackActionID: String) {
		cartons(search: $search, limit: $limit, offset: $offset, palletID: $palletID, trackActionID: $trackActionID) {
			cartons {
				...CartonListFragment
			}
			total
		}
	}
	${fragment.CARTON_LIST}
`

const CARTONS_BASIC = gql`
	query cartons($search: SearchFilter!, $limit: Int!, $offset: Int!, $palletID: String, $trackActionID: String) {
		cartons(search: $search, limit: $limit, offset: $offset, palletID: $palletID, trackActionID: $trackActionID) {
			cartons {
				id
				code
			}
			total
		}
	}
`

const CARTON_TRANSACTIONS = gql`
	query carton($code: String!) {
		carton(code: $code) {
			id
			transactions {
				...TransactionViewFragment
				createdBy {
					id
					email
					wechatID
				}
			}
		}
	}
	${fragment.TRANSACTION_VIEW}
`

const PRODUCT = gql`
	query product($code: String!) {
		product(code: $code) {
			...ProductFragment
		}
	}
	${fragment.PRODUCT}
`

/** Get product by ID w/ sku info for consumer product view (no permission check) */
const PRODUCT_VIEW = gql`
	query productByID($id: String!) {
		productByID(id: $id) {
			...ProductViewFragment
		}
	}
	${fragment.PRODUCT_VIEW}
`

const PRODUCTS = gql`
	query products(
		$search: SearchFilter!
		$limit: Int!
		$offset: Int!
		$cartonID: String
		$orderID: String
		$skuID: String
		$distributorID: String
		$contractID: String
		$trackActionID: String
	) {
		products(
			search: $search
			limit: $limit
			offset: $offset
			cartonID: $cartonID
			orderID: $orderID
			skuID: $skuID
			distributorID: $distributorID
			contractID: $contractID
			trackActionID: $trackActionID
		) {
			products {
				...ProductListFragment
			}
			total
		}
	}
	${fragment.PRODUCT_LIST}
`

const PRODUCT_TRANSACTIONS = gql`
	query productByID($id: String!) {
		productByID(id: $id) {
			transactions {
				...TransactionFragment
			}
		}
	}
	${fragment.TRANSACTION}
`
const PRODUCT_TRANSACTIONS_VIEW = gql`
	query productByID($id: String!) {
		productByID(id: $id) {
			transactions {
				...TransactionViewFragment
				createdBy {
					id
					email
					wechatID
				}
			}
		}
	}
	${fragment.TRANSACTION_VIEW}
`

const ORDER = gql`
	query order($code: String!) {
		order(code: $code) {
			...OrderFragment
		}
	}
	${fragment.ORDER}
`

const ORDERS = gql`
	query orders($search: SearchFilter!, $limit: Int!, $offset: Int!) {
		orders(search: $search, limit: $limit, offset: $offset) {
			orders {
				...OrderListFragment
			}
			total
		}
	}
	${fragment.ORDER_LIST}
`

const ORDERS_BASIC = gql`
	query orders($search: SearchFilter!, $limit: Int!, $offset: Int!) {
		orders(search: $search, limit: $limit, offset: $offset) {
			orders {
				id
				code
			}
			total
		}
	}
`

const TRACK_ACTION = gql`
	query trackAction($code: String!) {
		trackAction(code: $code) {
			...TrackActionFragment
		}
	}
	${fragment.TRACK_ACTION}
`

const TRACK_ACTIONS = gql`
	query trackActions($search: SearchFilter!, $limit: Int!, $offset: Int!) {
		trackActions(search: $search, limit: $limit, offset: $offset) {
			trackActions {
				...TrackActionFragment
			}
			total
		}
	}
	${fragment.TRACK_ACTION}
`

const CONTRACT = gql`
	query contract($code: String!) {
		contract(code: $code) {
			...ContractFragment
		}
	}
	${fragment.CONTRACT}
`
const CONTRACTS = gql`
	query contracts($search: SearchFilter!, $limit: Int!, $offset: Int!) {
		contracts(search: $search, limit: $limit, offset: $offset) {
			contracts {
				...ContractFragment
			}
			total
		}
	}
	${fragment.CONTRACT}
`
const CONTRACTS_BASIC = gql`
	query contracts($search: SearchFilter!, $limit: Int!, $offset: Int!) {
		contracts(search: $search, limit: $limit, offset: $offset) {
			contracts {
				id
				name
			}
			total
		}
	}
`

const DISTRIBUTOR = gql`
	query distributor($code: String!) {
		distributor(code: $code) {
			...DistributorFragment
		}
	}
	${fragment.DISTRIBUTOR}
`
const DISTRIBUTORS = gql`
	query distributors($search: SearchFilter!, $limit: Int!, $offset: Int!) {
		distributors(search: $search, limit: $limit, offset: $offset) {
			distributors {
				...DistributorFragment
			}
			total
		}
	}
	${fragment.DISTRIBUTOR}
`
const DISTRIBUTORS_BASIC = gql`
	query distributors($search: SearchFilter!, $limit: Int!, $offset: Int!) {
		distributors(search: $search, limit: $limit, offset: $offset) {
			distributors {
				id
				code
			}
			total
		}
	}
`

const VERIFY_RESET_TOKEN = gql`
	query verifyResetToken($token: String!, $email: NullString) {
		verifyResetToken(token: $token, email: $email)
	}
`

const GET_LOYALTY_ACTIVITY = gql`
	query getLoyaltyActivity($userID: ID!) {
		getLoyaltyActivity(userID: $userID) {
			id
			amount
			bonus
			message
			transactionHash
			createdAt
			product {
				id
				code
			}
		}
	}
`

const ALL_TRANSACTIONS = gql`
	query transactions($search: SearchFilter!, $limit: Int!, $offset: Int!, $productID: String, $cartonID: String, $trackActionID: String) {
		transactions(search: $search, limit: $limit, offset: $offset, productID: $productID, cartonID: $cartonID, trackActionID: $trackActionID) {
			transactions {
				...TransactionFragment
				carton {
					id
					code
				}
				product {
					id
					code
				}
			}
			total
		}
	}
	${fragment.TRANSACTION}
`

const USER_ACTIVITY = gql`
	query userActivities($search: SearchFilter!, $limit: Int!, $offset: Int!, $userID: String) {
		userActivities(search: $search, limit: $limit, offset: $offset, userID: $userID) {
			userActivities {
				id
				user {
					id
					firstName
					lastName
					email
					affiliateOrg
					role {
						id
						name
					}
				}
				action
				objectID
				objectCode
				objectType
				createdAt
			}
			total
		}
	}
`
const SETTINGS = gql`
	{
		settings {
			consumerHost
			adminHost
			etherscanHost
			fieldappVersion
			smartContractAddress
		}
	}
`
const FIELDAPP_VERSION = gql`
	{
		settings {
			fieldappVersion
		}
	}
`
const ETH_ACCOUNT_INFO = gql`
	{
		ethereumAccountAddress
		ethereumAccountBalance
	}
`
const PENDING_TRANSACTIONS_COUNT = gql`
	{
		pendingTransactionsCount
	}
`

const TICKER_INFO = gql`
	{
		getTickerInfo {
			lastTick
			tickInterval
		}
	}
`
//task 
const TASK = gql`
	query task($code: String!) {
		task(code: $code) {
			...TaskFragment
		}
	}
	${fragment.TASK}
`
const TASKS = gql`
	query tasks($search: SearchFilter!, $limit: Int!, $offset: Int!) {
		tasks(search: $search, limit: $limit, offset: $offset) {
			tasks {
				...TaskListFragment
			}
			total
		}
	}
	${fragment.TASK_LIST}
`
export const query = {
	ME,

	USER,
	USERS,
	CONSUMER,
	CONSUMERS,

	ROLE,
	ROLES,
	ROLES_LIMITED,

	SKU,
	SKUS,
	SKU_VIEW,
	SKU_CLONE_TREE,

	CONTAINER,
	CONTAINERS,
	CONTAINERS_BASIC,

	PALLET,
	PALLETS,
	PALLETS_BASIC,

	CARTON,
	CARTONS,
	CARTONS_BASIC,
	CARTON_TRANSACTIONS,

	PRODUCT,
	PRODUCTS,
	PRODUCT_VIEW,
	PRODUCT_TRANSACTIONS,
	PRODUCT_TRANSACTIONS_VIEW,

	ORDER,
	ORDERS,
	ORDERS_BASIC,

	TRACK_ACTION,
	TRACK_ACTIONS,

	ALL_TRANSACTIONS,

	CONTRACT,
	CONTRACTS,
	CONTRACTS_BASIC,

	DISTRIBUTOR,
	DISTRIBUTORS,
	DISTRIBUTORS_BASIC,

	VERIFY_RESET_TOKEN,

	GET_LOYALTY_ACTIVITY,
	USER_ACTIVITY,

	SETTINGS,
	ETH_ACCOUNT_INFO,
	PENDING_TRANSACTIONS_COUNT,
	FIELDAPP_VERSION,
	TICKER_INFO,

	TASKS,
	TASK,
}
