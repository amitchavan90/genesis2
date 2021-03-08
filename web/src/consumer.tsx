import * as React from "react"
import * as ReactDOM from "react-dom"
import { App } from "./app"
import { Enviroment } from "./types/enums"

const mountNode = document.getElementById("app")
ReactDOM.render(<App enviroment={Enviroment.Consumer} />, mountNode)
