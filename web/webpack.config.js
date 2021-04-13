const webpack = require("webpack")
const path = require("path")
const HtmlWebpackPlugin = require("html-webpack-plugin")
const CopyWebpackPlugin = require("copy-webpack-plugin")
const fs = require("fs-extra")

// remember build start time
const wdat = {
	build_start: new Date().toJSON(),
}
fs.writeFileSync("./stats.json", JSON.stringify(wdat), function(err) {
	if (err) {
		reportError("failed to write stats.json")
		process.exit(3)
	}
})

// tried clean-webpack-plugin and remove-webpack-plugin didnt work.
// after wasting all morning trying to find issue.
// k.i.s.s. delete dist folder and start over, this works!
fs.removeSync("./dist")

const config = {
	entry: {
		admin: "./src/admin.tsx",
		consumer: "./src/consumer.tsx",
	},
	output: {
		path: path.resolve(__dirname, "dist"),
		publicPath: "/",
		filename: "[name].[contenthash].js",
	},
	module: {
		rules: [
			{
				test: /\.(js|jsx)$/,
				use: "babel-loader",
				exclude: /node_modules/,
			},
			{
				test: /\.(ts|tsx)?$/,
				loader: "ts-loader",
				exclude: /node_modules/,
			},
			{
				test: /\.(sass|scss|css)$/,
				use: ["style-loader", "css-loader", "sass-loader"],
			},
			{
				test: /\.(svg|otf|eot|ttf|woff|gif)$/,
				use: "file-loader",
			},
			{
				test: /\.png$/,
				use: [
					{
						loader: "url-loader",
						options: {
							mimetype: "image/png",
						},
					},
				],
			},
		],
	},
	resolve: {
		extensions: [".js", ".jsx", ".tsx", ".ts"],
	},
	plugins: [
		new webpack.ContextReplacementPlugin(/moment[\/\\]locale$/, /en/),
		new HtmlWebpackPlugin({
			template: "src/index.html",
			filename: "admin.html",
			inject: true,
			excludeChunks: ["consumer"],
			appMountId: "app",
			favicon: "src/assets/images/favicon.ico",
		}),
		new HtmlWebpackPlugin({
			template: "src/index.html",
			filename: "consumer.html",
			inject: true,
			excludeChunks: ["admin"],
			appMountId: "app",
			favicon: "src/assets/images/favicon2.ico",
		}),
		new CopyWebpackPlugin([{ from: "./src/assets/static-files", to: "static-files" }]),
	],
	optimization: {
		runtimeChunk: "single",
		splitChunks: {
			cacheGroups: {
				vendor: {
					test: /[\\\/]node_modules[\\\/]/,
					name: "vendors",
					chunks: "all",
				},
			},
		},
	},
}

module.exports = config
