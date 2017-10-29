import React, { Component } from 'react';

import Nav from '../blocks/nav.js';
import Output from '../blocks/output.js';

export default class Job extends Component {
    render() {
        return (
            <div className="container">
                <Nav commands={this.props.commands} active={this.props.cmd}
                    query={this.props.query} onQueryChange={this.props.onQueryChange} />
                <main>
                    <Output cmd={this.props.cmd} jobID={this.props.jobID} />
                </main>
            </div>
        );
    }
}
