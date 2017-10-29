import React, { Component } from 'react';
import { Link } from 'react-router-dom';

import Timestamp from './timestamp.js';
import { api, httpGET, httpStreamGET } from '../http.js';
import './output.css';

export default class Output extends Component {
    constructor(props) {
        super(props);
        this.state = {
            job: null,
            xhr: null
        };
    }

    componentDidMount() {
        let jobID = this.props.jobID;
        if (jobID !== null) {
            this.loadCommandLog(jobID);
            this.loadJobDetails(jobID);
        }
    }

    componentWillUnmount() {
        if (this.state.xhr !== null) {
            this.state.xhr.abort();
        }
    }

    loadJobDetails(id) {
        if (id === null) {
            return;
        }

        httpGET(api("/jobs/" + id),
            (status, body) => {
                this.setState({job: JSON.parse(body)});
            },
            (error) => {
                console.log("Failed to load job details:", error);
            }
        );
    }

    loadCommandLog(id) {
        if (id === null || this.state.xhr !== null) {
            return;
        }

        let xhr = httpStreamGET(api("/jobs/" + id + "/log"),
            (chunk) => { // Progress
                let target = this.refs["output"];
                target.innerHTML += chunk.replace(/\n/g, "<br/>");
                this.autoScroll();
            },
            (status) => { // Complete
                // Request cancelled
                if (status === 0) {
                    return;
                }

                // Reload job details
                this.setState({xhr: null});
                this.loadJobDetails(id);
            },
            (error) => {
                let target = this.refs["output"];
                target.innerHTML = "Failed to fetch command log: "+ error;
            }
        );
        this.setState({xhr: xhr});
    }

    autoScroll() {
        // TODO: Figure out how to make it convinient
    }

    render() {
        var outputClass = "output";
        if (this.state.job) {
            outputClass += " visible";
        }

        return (
            <div>
                {this.renderJobDetails()}
                <div ref="output" className={outputClass}></div>
            </div>
        );
    }

    renderJobDetails() {
        let details = this.state.job;
        if (!details) {
            return (<div></div>);
        }

        // let shortID = details.id.substring(0, 8);
        var state = details.state;
        state = state.charAt(0).toUpperCase() + state.substr(1);

        var args;
        if (details.args !== "") {
            args = <span className="args">{details.args}</span>;
        }
        return (
            <div className="job-details full">
                <div className="item command">
                    <span className="command"><Link to={"/cmd/"+ details.command + "/jobs"}>{details.command}</Link></span>
                    {args}
                    <span className="flags">{details.flags}</span>
                </div>
                <div className="item id">
                    <div className="name">ID</div>
                    <div className="val"><Link to={"/cmd/"+ details.command + "/jobs/"+ details.id}>{details.id}</Link></div>
                </div>
                <div className="item state">
                    <div className="name">Status</div>
                    <div className={"val "+details.state}>{state}</div>
                </div>
                <div className="item user">
                    <div className="name">User</div>
                    <div className="val"><Link to={"/users/"+ details.user.id +"/jobs"}>{details.user.name}</Link></div>
                </div>
                {this.renderStarted()}
                {this.renderFinished()}
            </div>
        );
    }

    renderStarted() {
        let details = this.state.job;
        return (
            <div className="item started_at">
                <div className="name">Started</div>
                <div className="val"><Timestamp date={details.started_at} relative={details.finished_at === null} /></div>
            </div>
        );
    }

    renderFinished() {
        let details = this.state.job;
        if (details.finished_at !== null) {
            return [
                <div className="item finished_at" key="finished_at">
                    <div className="name">Finished</div>
                    <div className="val"><Timestamp date={details.finished_at} /></div>
                </div>,
                <div className="item took" key="took">
                    <div className="name">Took</div>
                    <div className="val"><Timestamp from={details.started_at} until={details.finished_at} /></div>
                </div>
            ];
        } else {
            return null;
        }
    }
}
