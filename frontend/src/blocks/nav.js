import React, { Component } from 'react'
import { Link } from 'react-router-dom'

import './nav.css'

export default class Nav extends Component {
    constructor(props) {
        super(props)
        this.state = {query: ""}
    }

    render() {
        return (
            <nav>
                <input className="search" type="search" placeholder="Search for commands"
                    value={this.props.query} onChange={this.props.onQueryChange} results="0" />
                <ul>
                    {Object.keys(this.props.commands).map(this.renderItem)}
                </ul>
            </nav>
        )
    }

    renderItem(name) {
        let cmd = this.props.commands[name]
        if (cmd.name.indexOf(this.props.query) === -1) {
            return null
        }

        let className = this.props.active && cmd.name === this.props.active ? "active" : ""
        return (
            <li key={cmd.name}>
                <Link to={"/cmd/" + cmd.name} className={className}>{cmd.name}</Link>
            </li>
        )
    }
}
