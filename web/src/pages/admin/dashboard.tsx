import * as React from "react"
import { useStyletron } from "baseui"
import { CenteredPage } from "../../components/common"
import { H1 } from "baseui/typography"
import { Spaced } from "../../components/spaced"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { useQuery, useMutation } from "@apollo/react-hooks"
import { Spinner } from "baseui/spinner"
import { Card } from "baseui/card"
import { graphql } from "../../graphql"
import { Spread } from "../../components/spread"
import { Settings, SearchFilter, TickerInfo } from "../../types/types"
import { Button } from "baseui/button"
import { ErrorNotification } from "../../components/errorBox"
import { Notification } from "baseui/notification"
import { FilterOption } from "../../types/enums"
import TimeAgo from "react-timeago"

const CardMargin = { Root: { style: { marginTop: "10px" } } }

const Dashboard = () => {
	const ethQuery = useQuery<{
		ethereumAccountBalance: string
		ethereumAccountAddress: string
	}>(graphql.query.ETH_ACCOUNT_INFO)
	const settingsQuery = useQuery<{ settings: Settings }>(graphql.query.SETTINGS)

	const [deployContract, mutDeployContract] = useMutation<{ deploySmartContract: Settings }>(graphql.mutation.DEPLOY_SMART_CONTRACT)

	const [settings, setSettings] = React.useState<Settings>()

	React.useEffect(() => {
		if (!mutDeployContract.data?.deploySmartContract) return
		setSettings(mutDeployContract.data.deploySmartContract)
	}, [mutDeployContract.loading, mutDeployContract.data])

	React.useEffect(() => {
		if (!settingsQuery.data?.settings) return
		setSettings(settingsQuery.data.settings)
	}, [settingsQuery.loading, settingsQuery.data])

	const truncate = (input: string) => {
		if (input.length > 10) return input.substring(0, 10) + "..."
		else return input
	}

	// Styling
	const [css, theme] = useStyletron()
	const transactionLink = css({
		...theme.typography.LabelSmall,
		cursor: "pointer",
		fontFamily: "monospace",
		marginLeft: "auto",
		display: "flex",
	})
	const transactionIcon = css({
		marginRight: "0.5rem",
	})

	return (
		<CenteredPage>
			<Spaced>
				<FontAwesomeIcon icon={["fas", "steak"]} size="3x" />
				<H1>Genesis</H1>
			</Spaced>

			{/* Account Balance */}
			<Card
				title={
					<Spread>
						<div>Account Balance</div>
						<div
							className={transactionLink}
							onClick={() => {
								if (!settingsQuery.data || !ethQuery.data) return
								window.open(`${settingsQuery.data.settings.etherscanHost}/address/${ethQuery.data.ethereumAccountAddress}`, "_blank")
							}}
						>
							<div className={transactionIcon}>{ethQuery.loading ? <Spinner size="15px" /> : <FontAwesomeIcon icon={["fas", "external-link"]} />}</div>
							{ethQuery.data?.ethereumAccountAddress && truncate(ethQuery.data.ethereumAccountAddress)}
						</div>
					</Spread>
				}
			>
				{ethQuery.loading && <Spinner />}
				<div>{ethQuery.data ? `${+ethQuery.data.ethereumAccountBalance / 10 ** 18} ETH` : "-"}</div>
			</Card>

			{/* Smart contract address + Deploy Contract button */}
			<Card
				overrides={CardMargin}
				title={
					<Spread>
						<div>Blockchain Contract</div>
						{settingsQuery.loading && <Spinner />}
						{settings !== undefined && settings.smartContractAddress !== "" && (
							<div
								className={transactionLink}
								onClick={() => {
									if (!settings) return
									window.open(`${settings.etherscanHost}/address/${settings.smartContractAddress}`, "_blank")
								}}
							>
								<div className={transactionIcon}>{ethQuery.loading ? <Spinner size="15px" /> : <FontAwesomeIcon icon={["fas", "external-link"]} />}</div>
								{truncate(settings.smartContractAddress)}
							</div>
						)}
						{settings !== undefined && settings.smartContractAddress === "" && (
							<Button
								kind="secondary"
								type="button"
								onClick={() => {
									if (!settingsQuery.data || mutDeployContract.loading || mutDeployContract.called) return
									deployContract()
								}}
								isLoading={mutDeployContract.loading}
							>
								Deploy Contract
							</Button>
						)}
					</Spread>
				}
			>
				{settings !== undefined && settings.smartContractAddress === "" && (
					<Notification
						kind="warning"
						overrides={{
							Body: {
								style: {
									width: "auto",
								},
							},
						}}
					>
						No blockchain smart contract address found. Perform first time set-up by clicking "Deploy Contract".
					</Notification>
				)}
				{mutDeployContract.error && <ErrorNotification message={mutDeployContract.error.message} />}
			</Card>

			<PendingTransactions />
		</CenteredPage>
	)
}

const PendingTransactions = () => {
	const tickerInfoQuery = useQuery<{ getTickerInfo: TickerInfo }>(graphql.query.TICKER_INFO)
	const { data, loading, error, refetch } = useQuery<{ pendingTransactionsCount: number }>(graphql.query.PENDING_TRANSACTIONS_COUNT, {
		fetchPolicy: "network-only",
	})
	const [flushToBlockchain, mutFlushToBlockchain] = useMutation<{ flushPendingTransactions: Boolean }>(graphql.mutation.FLUSH_PENDING_TRANSACTIONS)

	const [tickerInfo, setTickerInfo] = React.useState<TickerInfo>()
	const [pendingTotal, setPendingTotal] = React.useState<number>()

	React.useEffect(() => {
		if (!data?.pendingTransactionsCount) return
		setPendingTotal(data?.pendingTransactionsCount)
	}, [loading, data])
	React.useEffect(() => {
		if (!mutFlushToBlockchain.data || !mutFlushToBlockchain.data.flushPendingTransactions) return
		// Success - reset pendingTotal and lastTick
		setPendingTotal(0)
		if (tickerInfo) setTickerInfo({ lastTick: new Date().toUTCString(), tickInterval: tickerInfo?.tickInterval })
	}, [mutFlushToBlockchain.data])

	React.useEffect(() => {
		if (!tickerInfoQuery.data?.getTickerInfo) return
		setTickerInfo(tickerInfoQuery.data.getTickerInfo)
	}, [tickerInfoQuery.loading, tickerInfoQuery.data])

	return (
		<Card
			title={
				<Spread>
					<div>Pending Transactions</div>
					<Button
						kind="secondary"
						type="button"
						onClick={() => {
							if (mutFlushToBlockchain.loading || mutFlushToBlockchain.called) return
							flushToBlockchain()
						}}
						isLoading={mutFlushToBlockchain.loading}
						disabled={loading || !pendingTotal}
					>
						Commit to Blockchain
					</Button>
				</Spread>
			}
			overrides={CardMargin}
		>
			{loading && <Spinner />}
			{!loading && pendingTotal !== undefined && <div>{`${pendingTotal} transaction${pendingTotal == 1 ? " is" : "s are"} pending`}</div>}
			{(error || mutFlushToBlockchain.error) && <ErrorNotification message={mutFlushToBlockchain.error?.message || error?.message || ""} />}
			{tickerInfo && <TimeUntilNextTick tickerInfo={tickerInfo} />}
		</Card>
	)
}

const TimeUntilNextTick = (props: { tickerInfo: TickerInfo }) => {
	const { tickerInfo } = props

	const getNextTick = () => {
		let t = new Date(tickerInfo.lastTick)
		t.setHours(t.getHours() + tickerInfo.tickInterval)
		// if next tick is somehow in the past, just add another tickInterval (default: 24 hours)
		if (t < new Date()) t.setHours(t.getHours() + tickerInfo.tickInterval)
		return t
	}
	const [nextTick] = React.useState(getNextTick())

	const [css, theme] = useStyletron()
	const timeUntilStyle = css({
		...theme.typography.font200,
		color: "grey",
	})

	if (!tickerInfo) return null

	return (
		<div className={timeUntilStyle}>
			<span>Next automatic commit: </span>
			<TimeAgo date={nextTick} />
		</div>
	)
}

export default Dashboard
