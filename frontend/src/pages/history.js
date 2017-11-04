import React, { Component } from 'react'
import { Link } from 'react-router-dom'

import Nav from '../blocks/nav.js'
import JobList from '../blocks/job_list.js'
import User from '../blocks/user.js'
import Alert from '../blocks/alert.js'
import { api, httpGET } from '../http.js'
import './history.css'

export default class History extends Component {
    constructor(props) {
        super(props)
        this.state = {}
    }

    componentDidMount() {
        httpGET(this.endpointFromProps(this.props),
            (status, body) => {
                let jobs = JSON.parse(body)
                this.setState({jobs: jobs})
            },
            (error) => {
                this.setState({error: "Failed to load jobs: "+ error})
            },
        )
    }

    endpointFromProps(props) {
        if (props.cmd !== undefined) {
            return api("/commands/" + props.cmd + "/jobs")
        } else if (props.userID !== undefined) {
            return api("/users/" + props.userID + "/jobs")
        } else {
            return api("/jobs")
        }
    }

    render() {
        return (
            <div className="container">
                <Nav commands={this.props.commands} active={this.props.cmd}
                    query={this.props.query} onQueryChange={this.props.onQueryChange} />
                <main>
                    {this.renderJobDetails()}
                    <Alert type="error" message={this.state.error} />
                    <JobList jobs={this.state.jobs} />
                </main>
            </div>
        )
    }

    renderJobDetails() {
        if (this.props.cmd !== undefined) {
            return (
                <div className="command-details">
                    <span>Command</span>
                    <div className="cmd-name"><Link to={"/cmd/"+ this.props.cmd}>{this.props.cmd}</Link></div>
                </div>
            )
        } else if (this.props.userID !== undefined) {
            return (
                <User id={this.props.userID} />
            )
        }
        return null
    }
}
