import React, { Component } from 'react'
import Nav from '../blocks/nav.js'

export default class SelectCommand extends Component {
    render() {
        return (
            <div className="container">
                <Nav commands={this.props.commands}
                    onQueryChange={this.props.onQueryChange}
                    query={this.props.query} />
                <main>
                    <div className="select-command">
                        Please select a command
                    </div>
                </main>
            </div>
        )
    }
}
