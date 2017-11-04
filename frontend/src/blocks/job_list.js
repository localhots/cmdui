import React, { Component } from 'react'
import { Link } from 'react-router-dom'

import Timestamp from './timestamp.js'
import './job_list.css'

export default class JobList extends Component {
    render() {
        if (this.props.jobs === undefined) {
            return null
        }

        return (
            <div className="history-page">
                <div className={"job-details short legend"}>
                    <div></div>
                    <div>ID</div>
                    <div>Command</div>
                    <div>User</div>
                    <div>Started</div>
                    <div>Took</div>
                </div>
                {this.props.jobs.map(this.renderJob)}
            </div>
        )
    }

    renderJob(job) {
        let shortID = job.id.substring(0, 8)
        return (
            <div key={job.id} className={"job-details short"}>
                <div className={"dot "+ job.state}></div>
                <div className="id"><Link to={"/cmd/"+ job.command +"/jobs/"+ job.id}>{shortID}</Link></div>
                <div className="command"><Link to={"/cmd/"+ job.command +"/jobs"}>{job.command}</Link></div>
                <div className="user"><Link to={"/users/"+ job.user_id +"/jobs"}>{job.user.name}</Link></div>
                <div className="started"><Timestamp date={job.started_at} /></div>
                <div className="took"><Timestamp from={job.started_at} until={job.finished_at} /></div>
            </div>
        )
    }
}

