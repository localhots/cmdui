import React, { Component } from 'react';
import { Link } from 'react-router-dom';

import Nav from '../blocks/nav.js';
import JobList from '../blocks/job_list.js';
import User from '../blocks/user.js';
import { api, httpGET } from '../http.js';
import './history.css';

export default class History extends Component {
    constructor(props) {
        super(props);
        this.state = {};
    }

    componentDidMount() {
        var url;
        if (this.props.cmd !== undefined) {
            url = api("/commands/" + this.props.cmd + "/jobs")
        } else if (this.props.userID !== undefined) {
            url = api("/users/" + this.props.userID + "/jobs")
        } else {
            url = api("/jobs")
        }
        httpGET(url,
            (status, body) => {
                let jobs = JSON.parse(body);
                this.setState({jobs: jobs});
            },
            (error) => {
                console.log("Failed to load jobs:", error);
            },
        );
    }

    render() {
        let details;
        if (this.props.cmd !== undefined) {
            details = (
                <div className="command-details">
                    <span>Command</span>
                    <div className="cmd-name"><Link to={"/cmd/"+ this.props.cmd}>{this.props.cmd}</Link></div>
                </div>
            );
        } else if (this.props.userID !== undefined) {
            details = (
                <User id={this.props.userID} />
            );
        }
        return (
            <div className="container">
                <Nav commands={this.props.commands} active={this.props.cmd}
                    query={this.props.query} onQueryChange={this.props.onQueryChange} />
                <main>
                    {details}
                    <JobList jobs={this.state.jobs} />
                </main>
            </div>
        );
    }
}
