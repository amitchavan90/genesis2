import * as React from "react"
import { useStyletron } from "baseui"
import { ComposableMap, Geographies, Geography, ZoomableGroup, Point, Marker } from "react-simple-maps"
import Geohash from "latlon-geohash"

import geoUrl from "../assets/world-50m.json"

interface MarkerData {
	name: string
	coordinates: Point
}

interface MapPreviewProps {
	locationName?: string
	locationGeohash: string
	roundedCorner?: boolean
	zoom?: number
	largeMarker?: boolean
}

const MapPreview = (props: MapPreviewProps) => {
	const { locationName, locationGeohash, zoom, largeMarker } = props

	const latLong = Geohash.decode(locationGeohash)
	const markers: MarkerData[] = [
		{
			name: locationName || "?",
			coordinates: [latLong.lon, latLong.lat],
		},
	]

	const [css, theme] = useStyletron()
	const mapStyle = css({
		pointerEvents: "none",
		borderTopRightRadius: props.roundedCorner ? ".25em" : "unset",
	})
	const geoStyle = css({
		outline: "unset",
		fill: "#EAEAEC",
		stroke: "#D6D6DA",
		strokeWidth: "0.2",
	})

	return (
		<ComposableMap className={mapStyle} projection="geoEquirectangular">
			<ZoomableGroup zoom={zoom || 25} center={markers[0].coordinates} maxZoom={50} disablePanning disableZooming>
				<Geographies geography={geoUrl}>
					{({ geographies }) => geographies.map(geo => <Geography key={geo.rsmKey} geography={geo} className={geoStyle} />)}
				</Geographies>
				{markers.map(({ name, coordinates }) => (
					<MarkerShape key={name} name={name} coordinates={coordinates} largeMarker={largeMarker} />
				))}
			</ZoomableGroup>
		</ComposableMap>
	)
}

const MarkerShape = ({ name, coordinates, largeMarker }: MarkerData & { largeMarker?: boolean }) => {
	const [css, theme] = useStyletron()
	const textStyle = css({
		fontFamily: "Spartan Book Classified",
		fontSize: "2px",
		fill: "#5D5A6D",
		userSelect: "none",
	})

	const commaIndex = name.indexOf(",")
	const trimmedName = commaIndex === -1 ? name : name.substr(0, commaIndex)

	return (
		<Marker key={name} coordinates={coordinates}>
			<g
				fill="none"
				stroke="#FF5533"
				strokeWidth="2"
				strokeLinecap="round"
				strokeLinejoin="round"
				transform={`scale(${largeMarker ? 1 : 0.1}) translate(-12, -24)`}
			>
				<circle cx="12" cy="10" r="3" />
				<path d="M12 21.7C17.3 17 20 13 20 10a8 8 0 1 0-16 0c0 3 2.7 6.9 8 11.7z" />
			</g>
			{!largeMarker && (
				<text textAnchor="middle" y="-3" className={textStyle}>
					{trimmedName}
				</text>
			)}
		</Marker>
	)
}

export default MapPreview
